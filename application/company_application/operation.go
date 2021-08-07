package company_application

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gke-go-recruiting-server/adapter"
	"gke-go-recruiting-server/domain/account_domain"
	"gke-go-recruiting-server/domain/company_domain"

	"github.com/pkg/errors"
)

type operationApp struct {
	me *account_domain.AgencyAccount
	*app
}

func (a *operationApp) Create(
	ctx context.Context,
	db *gorm.DB,
	companyParams adapter.CompanyParams,
	now time.Time) (*company_domain.Company, error) {
	if err := validateCreate(companyParams); err != nil {
		return nil, errors.WithStack(err)
	}

	newCompany := company_domain.New(
		a.me.AgencyID,
		companyParams.RankType,
		companyParams.Rank,
		companyParams.Name,
		companyParams.NameKana,
		companyParams.PostalCode,
		companyParams.PrefID,
		companyParams.Address,
		companyParams.Building,
		companyParams.PhoneNumber,
		now)

	if err := a.tx(db, func(db *gorm.DB) error {
		return a.companyRepo.Update(ctx, db, newCompany)
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	return newCompany, nil
}
