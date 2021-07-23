package department_table

import (
	"context"

	"gke-go-sample/infra/cloudsql"

	"gorm.io/gorm"

	pb "gke-go-sample/proto/go/pb"

	"github.com/pkg/errors"

	"gke-go-sample/adapter"
	"gke-go-sample/domain"
	"gke-go-sample/domain/department_domain"
)

func NewRepo() adapter.DepartmentRepo {
	return &repository{}
}

type repository struct {
}

func (r *repository) filterQuery(db *gorm.DB, params adapter.DepartmentFilterParams) *gorm.DB {
	query := db.Model(&Entity{})

	if params.AgencyID.String() != "" {
		query = query.Where("agency_id = ?", params.AgencyID.String())
	}

	if params.CompanyID.String() != "" {
		query = query.Where("company_id = ?", params.CompanyID.String())
	}

	if params.DepartmentID.String() != "" {
		query = query.Where("id = ?", params.DepartmentID.String())
	}

	if params.DepartmentName != "" {
		query = query.Where("name LIKE ?", "%"+params.DepartmentName+"%")
	}

	if params.SalesID.String() != "" {
		query = query.Where("sales_id = ?", params.SalesID.String())
	}

	if params.Status != pb.Department_Status_Unknown {
		query = query.Where("status = ?", int32(params.Status))
	}

	if params.PhoneNumber != "" {
		query = query.Where("phone_number = ?", params.PhoneNumber)
	}

	query = query.Order("created_at DESC")

	return query
}

func (r *repository) GetByFilterWithPager(
	ctx context.Context,
	db *gorm.DB,
	pager *domain.Pager,
	params adapter.DepartmentFilterParams) ([]*department_domain.Department, error) {
	var es []Entity

	if err := r.filterQuery(db, params).
		Offset(pager.Offset()).
		Limit(pager.Limit()).
		Preload("Images").
		Preload("Stations").
		Preload("Stations.Line").
		Preload("ImmediatelyWorkableCalendars").
		Preload("AgencyAccount").
		Find(&es).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*department_domain.Department, 0, len(es))
	for _, e := range es {
		items = append(items, e.ToDomain())
	}

	return items, nil
}

func (r *repository) GetCountByFilter(
	ctx context.Context,
	db *gorm.DB,
	params adapter.DepartmentFilterParams) (uint64, error) {
	var count int64

	if err := r.filterQuery(db, params).Count(&count).Error; err != nil {
		return 0, errors.WithStack(err)
	}

	return uint64(count), nil
}

func (r *repository) Get(ctx context.Context, db *gorm.DB, id domain.DepartmentID) (*department_domain.Department, error) {
	var e Entity
	if err := db.Where("id = ?", id.String()).
		Preload("Images").
		Preload("Stations").
		Preload("Stations.Line").
		Preload("ImmediatelyWorkableCalendars").
		Preload("AgencyAccount").
		First(&e).Error; err != nil {
		if cloudsql.IsNotFoundError(err) {
			return nil, domain.NewNotFoundErr()
		}
		return nil, errors.WithStack(err)
	}

	return e.ToDomain(), nil
}

func (r *repository) Exist(ctx context.Context, db *gorm.DB, id domain.DepartmentID) (bool, error) {
	var count int64
	if err := db.
		Model(&Entity{}).
		Where("id = ?", id.String()).
		Count(&count).Error; err != nil {
		return false, errors.WithStack(err)
	}

	return count > 0, nil
}

func (r *repository) Insert(ctx context.Context, db *gorm.DB, item *department_domain.Department) error {
	if err := db.Create(entityFrom(item)).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repository) InsertMulti(ctx context.Context, db *gorm.DB, items []*department_domain.Department) error {
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

func (r *repository) Update(ctx context.Context, db *gorm.DB, item *department_domain.Department) error {
	if err := db.Save(entityFrom(item)).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
