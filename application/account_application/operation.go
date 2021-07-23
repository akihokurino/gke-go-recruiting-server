package account_application

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gke-go-sample/adapter"

	"gke-go-sample/domain"
	"gke-go-sample/domain/account_domain"

	"github.com/pkg/errors"
)

type operationApp struct {
	me *account_domain.AgencyAccount
	*app
}

func (a *operationApp) CreateAgencyAccount(
	ctx context.Context,
	db *gorm.DB,
	accountParams adapter.AgencyAccountParams,
	now time.Time) (*account_domain.AgencyAccount, error) {
	if err := validateCreateAgencyAccount(accountParams); err != nil {
		return nil, errors.WithStack(err)
	}

	agency, err := a.agencyRepo.Get(ctx, db, a.me.AgencyID)
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

func (a *operationApp) DeleteAgencyAccount(
	ctx context.Context,
	db *gorm.DB,
	accountID domain.FirebaseUserID) error {
	account, err := a.agencyAccountRepo.Get(ctx, db, accountID)
	if err != nil {
		return errors.WithStack(err)
	}

	if account.AgencyID != a.me.AgencyID {
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
