package entry_table

import (
	"context"

	"gorm.io/gorm"

	pb "gke-go-recruiting-server/proto/go/pb"

	"github.com/pkg/errors"

	"gke-go-recruiting-server/domain"
	"gke-go-recruiting-server/domain/entry_domain"

	"gke-go-recruiting-server/adapter"
)

func NewRepo() adapter.EntryRepo {
	return &repository{}
}

type repository struct {
}

func (r *repository) filterQuery(db *gorm.DB, params adapter.EntryFilterParams) *gorm.DB {
	query := db.Model(&Entity{}).
		Joins("inner join departments on departments.id = entries.department_id")

	if params.AgencyID.String() != "" {
		query = query.Where("departments.agency_id = ?", params.AgencyID.String())
	}

	if params.CompanyID.String() != "" {
		query = query.Where("departments.company_id = ?", params.CompanyID.String())
	}

	if params.DepartmentID.String() != "" {
		query = query.Where("entries.department_id = ?", params.DepartmentID.String())
	}

	if params.DepartmentName != "" {
		query = query.Where("departments.name LIKE ?", "%"+params.DepartmentName+"%")
	}

	if params.WorkID.String() != "" {
		query = query.Where("entries.work_id = ?", params.WorkID.String())
	}

	if params.DateRange != nil {
		query = query.
			Where("entries.created_at > ?", params.DateRange.From).
			Where("entries.created_at < ?", params.DateRange.To)
	}

	if params.SalesID.String() != "" {
		query = query.Where("departments.sales_id = ?", params.SalesID.String())
	}

	if params.BusinessCondition != pb.Department_BusinessCondition_Unknown {
		query = query.Where("departments.business_condition = ?", int32(params.BusinessCondition))
	}

	if params.PrefID.String() != "" {
		query = query.Where("departments.pref_id = ?", params.PrefID.String())
	}

	if params.Status != pb.Entry_Status_Unknown {
		query = query.Where("entries.status = ?", int32(params.Status))
	}

	query = query.Order("entries.created_at DESC")

	return query
}

func (r *repository) GetByFilterWithPager(
	ctx context.Context,
	db *gorm.DB,
	pager *domain.Pager,
	params adapter.EntryFilterParams) ([]*entry_domain.Entry, error) {
	var es []Entity

	if err := r.filterQuery(db, params).
		Offset(pager.Offset()).
		Limit(pager.Limit()).
		Find(&es).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*entry_domain.Entry, 0, len(es))
	for _, e := range es {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetCountByFilter(
	ctx context.Context,
	db *gorm.DB,
	params adapter.EntryFilterParams) (uint64, error) {
	var count int64

	if err := r.filterQuery(db, params).Count(&count).Error; err != nil {
		return 0, errors.WithStack(err)
	}

	return uint64(count), nil
}

func (r *repository) GetCountByInProgressAndDepartment(ctx context.Context, db *gorm.DB, departmentID domain.DepartmentID) (uint64, error) {
	var count int64

	if err := db.Model(&Entity{}).
		Where("department_id = ?", departmentID.String()).
		Where("status = ?", int32(pb.Entry_Status_InProgress)).
		Count(&count).Error; err != nil {
		return 0, errors.WithStack(err)
	}

	return uint64(count), nil
}

func (r *repository) Insert(ctx context.Context, db *gorm.DB, item *entry_domain.Entry) error {
	if err := db.Create(entityFrom(item)).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
