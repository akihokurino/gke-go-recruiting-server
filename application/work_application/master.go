package work_application

import (
	"context"

	"gorm.io/gorm"

	"gke-go-recruiting-server/domain"
	"gke-go-recruiting-server/domain/account_domain"

	"github.com/pkg/errors"
)

type masterApp struct {
	me *account_domain.Administrator
	*app
}

func (a *masterApp) Accept(ctx context.Context, db *gorm.DB, workID domain.WorkID) error {
	if err := a.tx(db, func(db *gorm.DB) error {
		work, err := a.workRepo.Get(ctx, db, workID)
		if err != nil {
			return err
		}

		if err := work.Accept(); err != nil {
			return err
		}

		if err := a.workRepo.Update(ctx, db, work); err != nil {
			return err
		}

		if err := a.workIndexRepo.Save(ctx, work); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *masterApp) Deny(ctx context.Context, db *gorm.DB, workID domain.WorkID) error {
	if err := a.tx(db, func(db *gorm.DB) error {
		work, err := a.workRepo.Get(ctx, db, workID)
		if err != nil {
			return err
		}

		if err := work.Deny(); err != nil {
			return err
		}

		if err := a.workRepo.Update(ctx, db, work); err != nil {
			return err
		}

		if err := a.workIndexRepo.Save(ctx, work); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
