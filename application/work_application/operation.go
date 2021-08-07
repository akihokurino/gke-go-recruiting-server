package work_application

import (
	"context"
	"net/url"
	"time"

	"gorm.io/gorm"

	"gke-go-recruiting-server/adapter"
	"gke-go-recruiting-server/domain"
	"gke-go-recruiting-server/domain/account_domain"
	"gke-go-recruiting-server/domain/work_domain"
	pb "gke-go-recruiting-server/proto/go/pb"

	"github.com/pkg/errors"
)

type operationApp struct {
	me *account_domain.AgencyAccount
	*app
}

func (a *operationApp) Create(
	ctx context.Context,
	db *gorm.DB,
	departmentID domain.DepartmentID,
	workParams adapter.WorkParams,
	imageParams []adapter.WorkImageParams,
	movieURLs []url.URL,
	merits []pb.Work_Merit,
	now time.Time) (*work_domain.Work, error) {
	if err := validateCreate(workParams); err != nil {
		return nil, errors.WithStack(err)
	}

	mainContract, err := a.mainContractRepo.GetByActiveAndDepartmentAndTime(ctx, db, departmentID, workParams.DateFrom)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	dateRange := domain.NewDateRangeWithCap(workParams.DateFrom, workParams.DateTo, mainContract.DateRange.To)

	newWork := work_domain.New(
		departmentID,
		workParams.WorkType,
		workParams.JobCode,
		workParams.Title,
		workParams.Content,
		dateRange,
		now)

	newImages := make([]*work_domain.Image, 0, len(imageParams))
	for _, image := range imageParams {
		newImages = append(newImages, work_domain.NewImage(newWork.ID, image.URL, image.ViewOrder, image.Comment))
	}

	newMovies := make([]*work_domain.Movie, 0, len(movieURLs))
	for _, u := range movieURLs {
		newMovies = append(newMovies, work_domain.NewMovie(newWork.ID, u))
	}

	newMerits := make([]*work_domain.Merit, 0, len(merits))
	for _, merit := range merits {
		newMerits = append(newMerits, work_domain.NewMerit(newWork.ID, merit))
	}

	if err := a.tx(db, func(db *gorm.DB) error {
		if err := a.checkDuplication(ctx, db, newWork); err != nil {
			return err
		}

		if err := a.workRepo.Insert(ctx, db, newWork); err != nil {
			return err
		}

		if err := a.workImageRepo.InsertMulti(ctx, db, newImages); err != nil {
			return err
		}

		if err := a.workMovieRepo.InsertMulti(ctx, db, newMovies); err != nil {
			return err
		}

		if err := a.workMeritRepo.InsertMulti(ctx, db, newMerits); err != nil {
			return err
		}

		department, err := a.departmentRepo.Get(ctx, db, newWork.DepartmentID)
		if err != nil {
			return err
		}

		newWork.With.Department = department
		newWork.Images = newImages
		newWork.Movies = newMovies
		newWork.Merits = newMerits

		if err := a.workIndexRepo.Save(ctx, newWork); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	return newWork, nil
}

func (a *operationApp) Update(
	ctx context.Context,
	db *gorm.DB,
	workID domain.WorkID,
	workParams adapter.WorkParams,
	imageParams []adapter.WorkImageParams,
	movieURLs []url.URL,
	merits []pb.Work_Merit,
	now time.Time) (*work_domain.Work, error) {
	if err := validateUpdate(workParams); err != nil {
		return nil, errors.WithStack(err)
	}

	newImages := make([]*work_domain.Image, 0, len(imageParams))
	for _, image := range imageParams {
		newImages = append(newImages, work_domain.NewImage(workID, image.URL, image.ViewOrder, image.Comment))
	}

	newMovies := make([]*work_domain.Movie, 0, len(movieURLs))
	for _, u := range movieURLs {
		newMovies = append(newMovies, work_domain.NewMovie(workID, u))
	}

	newMerits := make([]*work_domain.Merit, 0, len(merits))
	for _, merit := range merits {
		newMerits = append(newMerits, work_domain.NewMerit(workID, merit))
	}

	var work *work_domain.Work
	var err error
	if err := a.tx(db, func(db *gorm.DB) error {
		work, err = a.workRepo.Get(ctx, db, workID)
		if err != nil {
			return err
		}

		if work.With.Department.AgencyID != a.me.AgencyID {
			return domain.NewForbiddenErr(domain.ForbiddenAgencyMsg)
		}

		mainContract, err := a.mainContractRepo.GetByActiveAndDepartmentAndTime(ctx, db, work.DepartmentID, workParams.DateFrom)
		if err != nil {
			return err
		}

		dateRange := domain.NewDateRangeWithCap(workParams.DateFrom, workParams.DateTo, mainContract.DateRange.To)

		if err := work.Update(
			workParams.WorkType,
			workParams.JobCode,
			workParams.Title,
			workParams.Content,
			dateRange,
			now); err != nil {
			return err
		}

		if err := a.workRepo.Update(ctx, db, work); err != nil {
			return err
		}

		if err := a.workImageRepo.DeleteByWork(ctx, db, work.ID); err != nil {
			return err
		}

		if err := a.workMovieRepo.DeleteByWork(ctx, db, work.ID); err != nil {
			return err
		}

		if err := a.workMeritRepo.DeleteByWork(ctx, db, work.ID); err != nil {
			return err
		}

		if err := a.workImageRepo.InsertMulti(ctx, db, newImages); err != nil {
			return err
		}

		if err := a.workMovieRepo.InsertMulti(ctx, db, newMovies); err != nil {
			return err
		}

		if err := a.workMeritRepo.InsertMulti(ctx, db, newMerits); err != nil {
			return err
		}

		work.Images = newImages
		work.Movies = newMovies
		work.Merits = newMerits

		if err := a.workIndexRepo.Save(ctx, work); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	return work, nil
}

func (a *operationApp) EarlyFinish(ctx context.Context, db *gorm.DB, workID domain.WorkID) error {
	if err := a.tx(db, func(db *gorm.DB) error {
		work, err := a.workRepo.Get(ctx, db, workID)
		if err != nil {
			return err
		}

		if work.With.Department.AgencyID != a.me.AgencyID {
			return domain.NewForbiddenErr(domain.ForbiddenAgencyMsg)
		}

		if err := work.EarlyFinish(); err != nil {
			return err
		}

		if err := a.workRepo.Update(ctx, db, work); err != nil {
			return err
		}

		if err := a.workActivePlanRepo.Delete(ctx, db, work.ID); err != nil {
			return err
		}

		if err := a.workIndexRepo.Delete(ctx, work.ID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *operationApp) checkDuplication(ctx context.Context, db *gorm.DB, work *work_domain.Work) error {
	isExist, err := a.workRepo.ExistByDepartmentAndTypeAndTime(
		ctx,
		db,
		work.DepartmentID,
		work.WorkType,
		work.DateRange.From)
	if err != nil {
		return errors.WithStack(err)
	}

	if isExist {
		return domain.NewConflictErr("その期間の同お仕事はすでに存在します")
	}

	isExist, err = a.workRepo.ExistByDepartmentAndTypeAndTime(
		ctx,
		db,
		work.DepartmentID,
		work.WorkType,
		work.DateRange.To)
	if err != nil {
		return errors.WithStack(err)
	}

	if isExist {
		return domain.NewConflictErr("その期間の同お仕事はすでに存在します")
	}

	return nil
}
