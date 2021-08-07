package entry_application

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gke-go-recruiting-server/adapter"

	pb "gke-go-recruiting-server/proto/go/pb"

	"github.com/pkg/errors"

	"gke-go-recruiting-server/domain"
	"gke-go-recruiting-server/domain/entry_domain"
)

type publicApp struct {
	*app
}

func (a *publicApp) Entry(
	ctx context.Context,
	db *gorm.DB,
	workID domain.WorkID,
	entryParams adapter.EntryParams,
	now time.Time) error {
	if err := validateEntry(entryParams); err != nil {
		return errors.WithStack(err)
	}

	work, err := a.workRepo.Get(ctx, db, workID)
	if err != nil {
		return errors.WithStack(err)
	}

	var nilOrCategory *pb.User_Category
	if entryParams.Category != pb.User_Category_Unknown {
		nilOrCategory = &entryParams.Category
	}

	var nilOrPrefID *domain.PrefID
	if entryParams.PrefID.String() != "" {
		nilOrPrefID = &entryParams.PrefID
	}

	var nilOrPreferredContactMethod *pb.Entry_PreferredContactMethod
	if entryParams.PreferredContactMethod != pb.Entry_PreferredContactMethod_Unknown {
		nilOrPreferredContactMethod = &entryParams.PreferredContactMethod
	}

	var nilOrPreferredContactTime *string
	if entryParams.PreferredContactTime != "" {
		nilOrPreferredContactTime = &entryParams.PreferredContactTime
	}

	newEntry := entry_domain.New(
		work.DepartmentID,
		work.ID,
		entryParams.FullName,
		entryParams.FullNameKana,
		entryParams.Birthdate,
		entryParams.Gender,
		entryParams.PhoneNumber,
		entryParams.Email,
		entryParams.Question,
		nilOrCategory,
		nilOrPrefID,
		nilOrPreferredContactMethod,
		nilOrPreferredContactTime,
		now)

	if err := a.tx(db, func(db *gorm.DB) error {
		if err := a.entryRepo.Insert(ctx, db, newEntry); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
