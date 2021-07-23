package agency_table

import (
	"context"

	"gke-go-sample/infra/cloudsql"

	"gorm.io/gorm"

	"github.com/pkg/errors"

	"gke-go-sample/adapter"
	"gke-go-sample/domain"
	"gke-go-sample/domain/agency_domain"
)

func NewRepo() adapter.AgencyRepo {
	return &repository{}
}

type repository struct {
}

func (r *repository) GetAll(ctx context.Context, db *gorm.DB) ([]*agency_domain.Agency, error) {
	var es []Entity
	if err := db.
		Order("created_at DESC").
		Find(&es).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*agency_domain.Agency, 0, len(es))
	for _, e := range es {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) Get(ctx context.Context, db *gorm.DB, id domain.AgencyID) (*agency_domain.Agency, error) {
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

func (r *repository) InsertMulti(ctx context.Context, db *gorm.DB, items []*agency_domain.Agency) error {
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
