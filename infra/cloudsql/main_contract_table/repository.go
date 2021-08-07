package main_contract_table

import (
	"context"
	"time"

	"gke-go-recruiting-server/infra/cloudsql"

	"gorm.io/gorm"

	"github.com/pkg/errors"

	pb "gke-go-recruiting-server/proto/go/pb"

	"gke-go-recruiting-server/adapter"
	"gke-go-recruiting-server/domain"
	"gke-go-recruiting-server/domain/contract_domain"
)

func NewRepo() adapter.MainContractRepo {
	return &repository{}
}

type repository struct {
}

func (r *repository) filterQuery(db *gorm.DB, params adapter.MainContractFilterParams) *gorm.DB {
	query := db.Model(&Entity{}).
		Joins("inner join departments on departments.id = main_contracts.department_id")

	if params.AgencyID.String() != "" {
		query = query.Where("departments.agency_id = ?", params.AgencyID.String())
	}

	if params.CompanyID.String() != "" {
		query = query.Where("departments.company_id = ?", params.CompanyID.String())
	}

	if params.DepartmentID.String() != "" {
		query = query.Where("main_contracts.department_id = ?", params.DepartmentID.String())
	}

	if params.DepartmentName != "" {
		query = query.Where("departments.name LIKE ?", "%"+params.DepartmentName+"%")
	}

	if params.Status != pb.MainContract_Status_Unknown {
		query = query.Where("main_contracts.status = ?", int32(params.Status))
	}

	if params.DateRange != nil {
		query = query.
			Where("main_contracts.date_from > ?", params.DateRange.From).
			Where("main_contracts.date_to < ?", params.DateRange.To)
	}

	if params.Plan != pb.MainProduct_Plan_Unknown {
		query = query.Where("main_contracts.plan = ?", int32(params.Plan))
	}

	if params.SalesID.String() != "" {
		query = query.Where("departments.sales_id = ?", params.SalesID.String())
	}

	query = query.Order("main_contracts.created_at DESC")

	return query
}

func (r *repository) GetByFilterWithPager(
	ctx context.Context,
	db *gorm.DB,
	pager *domain.Pager,
	params adapter.MainContractFilterParams) ([]*contract_domain.Main, error) {
	var es []Entity

	if err := r.filterQuery(db, params).
		Offset(pager.Offset()).
		Limit(pager.Limit()).
		Preload("Department").
		Find(&es).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*contract_domain.Main, 0, len(es))
	for _, e := range es {
		items = append(items, e.ToDomain())
	}

	return items, nil
}

func (r *repository) GetCountByFilter(
	ctx context.Context,
	db *gorm.DB,
	params adapter.MainContractFilterParams) (uint64, error) {
	var count int64

	if err := r.filterQuery(db, params).Count(&count).Error; err != nil {
		return 0, errors.WithStack(err)
	}

	return uint64(count), nil
}

func (r *repository) GetByActiveAndDepartmentAndTime(
	ctx context.Context,
	db *gorm.DB,
	departmentID domain.DepartmentID,
	t time.Time) (*contract_domain.Main, error) {
	var e Entity
	if err := db.
		Joins("inner join departments on departments.id = main_contracts.department_id").
		Where("main_contracts.department_id = ?", departmentID.String()).
		Where("main_contracts.date_from <= ?", t).
		Where("main_contracts.date_to > ?", t).
		Where("main_contracts.status = ?", int32(pb.MainContract_Status_OK)).
		Preload("Department").
		First(&e).Error; err != nil {
		if cloudsql.IsNotFoundError(err) {
			return nil, domain.NewNotFoundErr()
		}
		return nil, errors.WithStack(err)
	}

	return e.ToDomain(), nil
}

func (r *repository) ExistByDepartmentAndTime(ctx context.Context, db *gorm.DB, departmentID domain.DepartmentID, t time.Time) (bool, error) {
	var count int64
	if err := db.
		Model(&Entity{}).
		Where("date_from <= ?", t).
		Where("date_to > ?", t).
		Where("department_id = ?", departmentID.String()).
		Where("status = ? OR status = ?", int32(pb.MainContract_Status_Review), int32(pb.MainContract_Status_OK)).
		Count(&count).Error; err != nil {
		return false, errors.WithStack(err)
	}

	return count > 0, nil
}

func (r *repository) Get(ctx context.Context, db *gorm.DB, id domain.MainContractID) (*contract_domain.Main, error) {
	var e Entity
	if err := db.
		Joins("inner join departments on departments.id = main_contracts.department_id").
		Where("main_contracts.id = ?", id.String()).
		Preload("Department").
		First(&e).Error; err != nil {
		if cloudsql.IsNotFoundError(err) {
			return nil, domain.NewNotFoundErr()
		}
		return nil, errors.WithStack(err)
	}

	return e.ToDomain(), nil
}

func (r *repository) Insert(ctx context.Context, db *gorm.DB, item *contract_domain.Main) error {
	if err := db.Create(entityFrom(item)).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repository) InsertMulti(ctx context.Context, db *gorm.DB, items []*contract_domain.Main) error {
	for _, item := range items {
		if err := db.Create(entityFrom(item)).Error; err != nil {
			if cloudsql.IsDuplicateError(err) {
				continue
			}
			return errors.WithStack(err)
		}
	}

	return nil
}

func (r *repository) Update(ctx context.Context, db *gorm.DB, item *contract_domain.Main) error {
	if err := db.Save(entityFrom(item)).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
