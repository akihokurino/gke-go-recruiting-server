package handler

import (
	"context"

	"gorm.io/gorm"

	"github.com/pkg/errors"

	"gke-go-sample/domain"

	"gke-go-sample/adapter"
	"gke-go-sample/domain/account_domain"
)

func NewAdminAuthorization(cp adapter.ContextProvider, administratorRepo adapter.AdministratorRepo) adapter.AdminAuthorization {
	return func(ctx context.Context, db *gorm.DB) (*account_domain.Administrator, error) {
		id, err := cp.FirebaseUserID(ctx)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		me, err := administratorRepo.Get(ctx, db, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return nil, domain.NewUnAuthorizedErr("マスター権限の認証に失敗しました")
			}
			return nil, errors.WithStack(err)
		}

		return me, nil
	}
}

func NewAgencyAuthorization(cp adapter.ContextProvider, agencyAccountRepo adapter.AgencyAccountRepo) adapter.AgencyAuthorization {
	return func(ctx context.Context, db *gorm.DB) (*account_domain.AgencyAccount, error) {
		id, err := cp.FirebaseUserID(ctx)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		me, err := agencyAccountRepo.Get(ctx, db, id)
		if err != nil {
			if domain.IsNotFound(err) {
				return nil, domain.NewUnAuthorizedErr("代理店権限の認証に失敗しました")
			}
			return nil, errors.WithStack(err)
		}

		return me, nil
	}
}
