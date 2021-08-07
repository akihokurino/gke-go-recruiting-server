package account_application

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gke-go-recruiting-server/adapter"

	"gke-go-recruiting-server/domain"
	"gke-go-recruiting-server/domain/account_domain"

	"github.com/pkg/errors"
)

type masterApp struct {
	me *account_domain.Administrator
	*app
}

func (a *masterApp) CreateAdministrator(
	ctx context.Context,
	db *gorm.DB,
	accountParams adapter.AdministratorParams,
	now time.Time) (*account_domain.Administrator, error) {
	if err := validateCreateAdministrator(accountParams); err != nil {
		return nil, errors.WithStack(err)
	}

	firebaseUser, err := a.firebaseRepo.CreateEmailUser(ctx, accountParams.Email, accountParams.Password)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	newAccount := firebaseUser.NewAdministrator(accountParams.Name, now)

	if err := a.tx(db, func(db *gorm.DB) error {
		return a.administratorRepo.Insert(ctx, db, newAccount)
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	return newAccount, nil
}

func (a *masterApp) CreateAgencyAccount(
	ctx context.Context,
	db *gorm.DB,
	accountParams adapter.AgencyAccountParams,
	agencyID domain.AgencyID,
	now time.Time) (*account_domain.AgencyAccount, error) {
	if err := validateCreateAgencyAccount(accountParams); err != nil {
		return nil, errors.WithStack(err)
	}

	agency, err := a.agencyRepo.Get(ctx, db, agencyID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	firebaseUser, err := a.firebaseRepo.CreateEmailUser(ctx, accountParams.Email, accountParams.Password)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	newAccount := firebaseUser.NewAgencyAccount(agency.ID, accountParams.Name, accountParams.NameKana, now)

	if err := a.tx(db, func(db *gorm.DB) error {
		return a.agencyAccountRepo.Insert(ctx, db, newAccount)
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	return newAccount, nil
}

func (a *masterApp) DeleteAgencyAccount(
	ctx context.Context,
	db *gorm.DB,
	agencyID domain.AgencyID,
	accountID domain.FirebaseUserID) error {
	account, err := a.agencyAccountRepo.Get(ctx, db, accountID)
	if err != nil {
		return errors.WithStack(err)
	}

	if account.AgencyID != agencyID {
		return domain.NewForbiddenErr(domain.ForbiddenAgencyMsg)
	}

	if err := a.tx(db, func(db *gorm.DB) error {
		if err := a.agencyAccountRepo.Delete(ctx, db, account.ID); err != nil {
			return err
		}

		return a.firebaseRepo.Delete(ctx, account.ID)
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
