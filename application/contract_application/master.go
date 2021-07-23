package contract_application

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gke-go-sample/domain"
	"gke-go-sample/domain/account_domain"
	"gke-go-sample/domain/statement_domain"

	"github.com/pkg/errors"
)

type masterApp struct {
	me *account_domain.Administrator
	*app
}

func (a *masterApp) AcceptMainContract(
	ctx context.Context,
	db *gorm.DB,
	contractID domain.MainContractID,
	now time.Time) error {
	if err := a.tx(db, func(db *gorm.DB) error {
		contract, err := a.mainContractRepo.Get(ctx, db, contractID)
		if err != nil {
			return err
		}

		if err := contract.Accept(); err != nil {
			return nil
		}

		if err := a.mainContractRepo.Update(ctx, db, contract); err != nil {
			return err
		}

		statement := statement_domain.NewUsageFromMain(
			contract.DepartmentID,
			contract.ID,
			contract.Price,
			now,
		)

		if err := a.usageStatementRepo.Insert(ctx, db, statement); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *masterApp) DenyMainContract(
	ctx context.Context,
	db *gorm.DB,
	contractID domain.MainContractID) error {
	if err := a.tx(db, func(db *gorm.DB) error {
		contract, err := a.mainContractRepo.Get(ctx, db, contractID)
		if err != nil {
			return err
		}

		if err := contract.Deny(); err != nil {
			return err
		}

		return a.mainContractRepo.Update(ctx, db, contract)
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
