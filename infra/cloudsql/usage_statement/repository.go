package usage_statement

import (
	"context"

	"gorm.io/gorm"

	pb "gke-go-recruiting-server/proto/go/pb"

	"github.com/pkg/errors"

	"gke-go-recruiting-server/domain/statement_domain"

	"gke-go-recruiting-server/adapter"
	"gke-go-recruiting-server/domain"
)

func NewRepo() adapter.UsageStatementRepo {
	return &repository{}
}

type repository struct {
}

func (r *repository) filterQuery(db *gorm.DB, params adapter.UsageStatementFilterParams) *gorm.DB {
	query := db.Model(&Entity{}).
		Joins("inner join departments on departments.id = usage_statements.department_id")

	if params.AgencyID.String() != "" {
		query = query.Where("departments.agency_id = ?", params.AgencyID.String())
	}

	if params.CompanyID.String() != "" {
		query = query.Where("departments.company_id = ?", params.CompanyID.String())
	}

	if params.DepartmentID.String() != "" {
		query = query.Where("usage_statements.department_id = ?", params.DepartmentID.String())
	}

	if params.DepartmentName != "" {
		query = query.Where("departments.name LIKE ?", "%"+params.DepartmentName+"%")
	}

	if params.MainProductPlan != pb.MainProduct_Plan_Unknown {
		query = query.
			Joins("inner join main_contracts on main_contracts.id = usage_statements.main_contract_id").
			Where("main_contracts.plan = ?", int32(params.MainProductPlan))
	}

	if params.DateRange != nil {
		query = query.
			Where("usage_statements.created_at > ?", params.DateRange.From).
			Where("usage_statements.created_at < ?", params.DateRange.To)
	}

	if params.ExcludeFree {
		query = query.Where("usage_statements.price != 0")
	}

	if params.UsageStatementID.String() != "" {
		query = query.Where("usage_statements.id = ?", params.UsageStatementID.String())
	}

	query = query.Order("usage_statements.created_at DESC")

	return query
}

func (r *repository) GetByFilterWithPager(
	ctx context.Context,
	db *gorm.DB,
	pager *domain.Pager,
	params adapter.UsageStatementFilterParams) ([]*statement_domain.Usage, error) {
	var es []Entity

	if err := r.filterQuery(db, params).
		Offset(pager.Offset()).
		Limit(pager.Limit()).
		Preload("MainContract").
		Preload("OptionContract").
		Find(&es).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*statement_domain.Usage, 0, len(es))
	for _, e := range es {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetCountByFilter(
	ctx context.Context,
	db *gorm.DB,
	params adapter.UsageStatementFilterParams) (uint64, error) {
	var count int64

	if err := r.filterQuery(db, params).Count(&count).Error; err != nil {
		return 0, errors.WithStack(err)
	}

	return uint64(count), nil
}

func (r *repository) Insert(ctx context.Context, db *gorm.DB, item *statement_domain.Usage) error {
	if err := db.Create(entityFrom(item)).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
