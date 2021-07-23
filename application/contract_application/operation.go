package contract_application

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gke-go-sample/adapter"

	"gke-go-sample/domain"
	"gke-go-sample/domain/account_domain"
	"gke-go-sample/domain/contract_domain"
	"gke-go-sample/domain/product_domain"

	"github.com/pkg/errors"
)

type operationApp struct {
	me *account_domain.AgencyAccount
	*app
}

func (a *operationApp) CreateMainContract(
	ctx context.Context,
	db *gorm.DB,
	departmentID domain.DepartmentID,
	contractParams adapter.MainContractParams,
	now time.Time) (*contract_domain.Main, error) {
	if err := validateCreateMainContract(contractParams); err != nil {
		return nil, errors.WithStack(err)
	}

	var thisProduct *product_domain.Main
	for _, product := range product_domain.GetMainList() {
		if product.Plan == contractParams.Plan {
			thisProduct = product
			break
		}
	}

	if thisProduct == nil {
		return nil, domain.NewBadRequestErr("不正なプランです")
	}

	department, err := a.departmentRepo.Get(ctx, db, departmentID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if department.AgencyID != a.me.AgencyID {
		return nil, domain.NewForbiddenErr(domain.ForbiddenAgencyMsg)
	}

	dateRange, err := thisProduct.CalcDateRange(contractParams.DateFrom)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	newContract := contract_domain.NewMainContract(
		departmentID,
		thisProduct.Plan,
		dateRange,
		thisProduct.Price,
		now)

	if err := a.tx(db, func(db *gorm.DB) error {
		isExist, err := a.mainContractRepo.ExistByDepartmentAndTime(ctx, db, departmentID, dateRange.From)
		if err != nil {
			return err
		}

		if isExist {
			return domain.NewConflictErr("すでにその期間の同契約は存在します")
		}

		isExist, err = a.mainContractRepo.ExistByDepartmentAndTime(ctx, db, departmentID, dateRange.To)
		if err != nil {
			return err
		}

		if isExist {
			return domain.NewConflictErr("すでにその期間の同契約は存在します")
		}

		return a.mainContractRepo.Insert(ctx, db, newContract)
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	return newContract, nil
}
