package department_application

import (
	"context"
	"net/url"
	"time"

	"gorm.io/gorm"

	"github.com/pkg/errors"

	"gke-go-sample/adapter"
	"gke-go-sample/domain"
	"gke-go-sample/domain/account_domain"
	"gke-go-sample/domain/department_domain"
)

type operationApp struct {
	me *account_domain.AgencyAccount
	*app
}

func (a *operationApp) Create(
	ctx context.Context,
	db *gorm.DB,
	companyID domain.CompanyID,
	salesID domain.FirebaseUserID,
	departmentParams adapter.DepartmentParams,
	imageURLs []url.URL,
	lineIDs []domain.LineID,
	now time.Time) (*department_domain.Department, error) {
	if err := validateCreate(departmentParams); err != nil {
		return nil, errors.WithStack(err)
	}

	agencyAccount, err := a.agencyAccountRepo.Get(ctx, db, salesID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if agencyAccount.AgencyID != a.me.AgencyID {
		return nil, domain.NewForbiddenErr(domain.ForbiddenAgencyMsg)
	}

	newDepartment := department_domain.New(
		a.me.AgencyID,
		companyID,
		salesID,
		departmentParams.Name,
		departmentParams.BusinessCondition,
		departmentParams.PostalCode,
		departmentParams.PrefID,
		departmentParams.CityID,
		departmentParams.Address,
		departmentParams.Building,
		departmentParams.PhoneNumber,
		department_domain.Location{
			MAreaID:   departmentParams.MAreaID,
			SAreaID:   departmentParams.SAreaID,
			Latitude:  departmentParams.Latitude,
			Longitude: departmentParams.Longitude,
		},
		now)

	newImages := make([]*department_domain.Image, 0, len(imageURLs))
	for _, u := range imageURLs {
		newImages = append(newImages, department_domain.NewImage(newDepartment.ID, u))
	}

	newStations := make([]*department_domain.Station, 0, len(lineIDs))
	for _, lineID := range lineIDs {
		newStations = append(newStations, department_domain.NewStation(newDepartment.ID, lineID))
	}

	if err := a.tx(db, func(db *gorm.DB) error {
		if err := a.departmentRepo.Insert(ctx, db, newDepartment); err != nil {
			return err
		}

		if err := a.departmentImageRepo.InsertMulti(ctx, db, newImages); err != nil {
			return err
		}

		if err := a.departmentStationRepo.InsertMulti(ctx, db, newStations); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	newDepartment.Images = newImages

	return newDepartment, nil
}

func (a *operationApp) Update(
	ctx context.Context,
	db *gorm.DB,
	departmentID domain.DepartmentID,
	departmentParams adapter.DepartmentParams,
	imageURLs []url.URL,
	lineIDs []domain.LineID,
	now time.Time) (*department_domain.Department, error) {
	if err := validateUpdate(departmentParams); err != nil {
		return nil, errors.WithStack(err)
	}

	newImages := make([]*department_domain.Image, 0, len(imageURLs))
	for _, u := range imageURLs {
		newImages = append(newImages, department_domain.NewImage(departmentID, u))
	}

	newStations := make([]*department_domain.Station, 0, len(lineIDs))
	for _, lineID := range lineIDs {
		newStations = append(newStations, department_domain.NewStation(departmentID, lineID))
	}

	var department *department_domain.Department
	var err error
	if err := a.tx(db, func(db *gorm.DB) error {
		department, err = a.departmentRepo.Get(ctx, db, departmentID)
		if err != nil {
			return err
		}

		if department.AgencyID != a.me.AgencyID {
			return domain.NewForbiddenErr(domain.ForbiddenAgencyMsg)
		}

		department.Update(
			departmentParams.Name,
			departmentParams.BusinessCondition,
			departmentParams.PostalCode,
			departmentParams.PrefID,
			departmentParams.CityID,
			departmentParams.Address,
			departmentParams.Building,
			departmentParams.PhoneNumber,
			department_domain.Location{
				MAreaID:   departmentParams.MAreaID,
				SAreaID:   departmentParams.SAreaID,
				Latitude:  departmentParams.Latitude,
				Longitude: departmentParams.Longitude,
			},
			now)

		if err := a.departmentRepo.Update(ctx, db, department); err != nil {
			return err
		}

		if err := a.departmentImageRepo.DeleteByDepartment(ctx, db, department.ID); err != nil {
			return err
		}

		if err := a.departmentStationRepo.DeleteByDepartment(ctx, db, department.ID); err != nil {
			return err
		}

		if err := a.departmentImageRepo.InsertMulti(ctx, db, newImages); err != nil {
			return err
		}

		if err := a.departmentStationRepo.InsertMulti(ctx, db, newStations); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	return department, nil
}

func (a *operationApp) UpdateSales(
	ctx context.Context,
	db *gorm.DB,
	departmentID domain.DepartmentID,
	salesID domain.FirebaseUserID) error {
	agencyAccount, err := a.agencyAccountRepo.Get(ctx, db, salesID)
	if err != nil {
		return errors.WithStack(err)
	}

	if agencyAccount.AgencyID != a.me.AgencyID {
		return domain.NewForbiddenErr(domain.ForbiddenAgencyMsg)
	}

	if err := a.tx(db, func(db *gorm.DB) error {
		department, err := a.departmentRepo.Get(ctx, db, departmentID)
		if err != nil {
			return err
		}

		department.UpdateSales(salesID)

		return a.departmentRepo.Update(ctx, db, department)
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
