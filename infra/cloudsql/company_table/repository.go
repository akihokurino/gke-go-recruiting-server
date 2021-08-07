package company_table

import (
	"context"

	"gke-go-recruiting-server/infra/cloudsql"

	"gorm.io/gorm"

	"github.com/pkg/errors"

	"gke-go-recruiting-server/domain"
	"gke-go-recruiting-server/domain/company_domain"

	"gke-go-recruiting-server/adapter"
)

func NewRepo() adapter.CompanyRepo {
	return &repository{}
}

type repository struct {
}

func (r *repository) filterQuery(db *gorm.DB, params adapter.CompanyFilterParams) *gorm.DB {
	query := db.Model(&Entity{})

	if params.CompanyID.String() != "" {
		query = query.Where("id = ?", params.CompanyID.String())
	}

	if params.CompanyName != "" {
		query = query.Where("name LIKE ?", "%"+params.CompanyName+"%")
	}

	if params.AgencyID.String() != "" {
		query = query.Where("agency_id = ?", params.AgencyID.String())
	}

	query = query.Order("created_at DESC")

	return query
}

func (r *repository) GetByFilterWithPager(
	ctx context.Context,
	db *gorm.DB,
	pager *domain.Pager,
	params adapter.CompanyFilterParams) ([]*company_domain.Company, error) {
	var es []Entity

	if err := r.filterQuery(db, params).
		Offset(pager.Offset()).
		Limit(pager.Limit()).
		Find(&es).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*company_domain.Company, 0, len(es))
	for _, e := range es {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetCountByFilter(
	ctx context.Context,
	db *gorm.DB,
	params adapter.CompanyFilterParams) (uint64, error) {
	var count int64

	if err := r.filterQuery(db, params).Count(&count).Error; err != nil {
		return 0, errors.WithStack(err)
	}

	return uint64(count), nil
}

func (r *repository) Get(ctx context.Context, db *gorm.DB, id domain.CompanyID) (*company_domain.Company, error) {
	var e Entity
	if err := db.Where("id = ?", id.String()).
		First(&e).Error; err != nil {
		if cloudsql.IsNotFoundError(err) {
			return nil, domain.NewNotFoundErr()
		}
		return nil, errors.WithStack(err)
	}

	return e.toDomain(), nil
}

func (r *repository) Exist(ctx context.Context, db *gorm.DB, id domain.CompanyID) (bool, error) {
	var count int64
	if err := db.
		Model(&Entity{}).
		Where("id = ?", id.String()).
		Count(&count).Error; err != nil {
		return false, errors.WithStack(err)
	}

	return count > 0, nil
}

func (r *repository) Insert(ctx context.Context, db *gorm.DB, item *company_domain.Company) error {
	if err := db.Create(entityFrom(item)).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repository) InsertMulti(ctx context.Context, db *gorm.DB, items []*company_domain.Company) error {
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

func (r *repository) Update(ctx context.Context, db *gorm.DB, item *company_domain.Company) error {
	if err := db.Save(entityFrom(item)).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
