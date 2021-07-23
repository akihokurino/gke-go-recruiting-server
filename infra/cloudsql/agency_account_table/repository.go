package agency_account_table

import (
	"context"

	"gke-go-sample/infra/cloudsql"

	"gorm.io/gorm"

	"github.com/pkg/errors"

	"gke-go-sample/domain"
	"gke-go-sample/domain/account_domain"

	"gke-go-sample/adapter"
)

func NewRepo() adapter.AgencyAccountRepo {
	return &repository{}
}

type repository struct {
}

func (r *repository) GetAll(ctx context.Context, db *gorm.DB) ([]*account_domain.AgencyAccount, error) {
	var es []Entity
	if err := db.
		Order("created_at DESC").
		Find(&es).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*account_domain.AgencyAccount, 0, len(es))
	for _, e := range es {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetByAgency(ctx context.Context, db *gorm.DB, agencyID domain.AgencyID) ([]*account_domain.AgencyAccount, error) {
	var es []Entity
	if err := db.
		Where("agency_id = ?", agencyID.String()).
		Order("created_at DESC").
		Find(&es).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*account_domain.AgencyAccount, 0, len(es))
	for _, e := range es {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) Get(ctx context.Context, db *gorm.DB, id domain.FirebaseUserID) (*account_domain.AgencyAccount, error) {
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

func (r *repository) Insert(ctx context.Context, db *gorm.DB, item *account_domain.AgencyAccount) error {
	if err := db.Create(entityFrom(item)).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repository) InsertMulti(ctx context.Context, db *gorm.DB, items []*account_domain.AgencyAccount) error {
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

func (r *repository) Delete(ctx context.Context, db *gorm.DB, id domain.FirebaseUserID) error {
	if err := db.Delete(onlyID(id)).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
