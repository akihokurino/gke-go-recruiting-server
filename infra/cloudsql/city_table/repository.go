package city_table

import (
	"context"

	"gke-go-sample/infra/cloudsql"

	"gorm.io/gorm"

	"github.com/pkg/errors"

	"gke-go-sample/domain/master_domain"

	"gke-go-sample/adapter"
	"gke-go-sample/domain"
)

func NewRepo() adapter.CityRepo {
	return &repository{}
}

type repository struct {
}

func (r *repository) GetAllPrefecture(ctx context.Context, db *gorm.DB) ([]*master_domain.Prefecture, error) {
	var es []Entity
	if err := db.Select("DISTINCT pref_id, pref_name").Find(&es).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*master_domain.Prefecture, 0, len(es))
	for _, e := range es {
		items = append(items, e.toPrefectureDomain())
	}

	return items, nil
}

func (r *repository) GetByPrefecture(ctx context.Context, db *gorm.DB, prefID domain.PrefID) ([]*master_domain.City, error) {
	var es []Entity
	if err := db.Where("pref_id = ?", prefID.String()).Find(&es).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*master_domain.City, 0, len(es))
	for _, e := range es {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) Exist(ctx context.Context, db *gorm.DB, id domain.CityID) (bool, error) {
	var count int64
	if err := db.
		Model(&Entity{}).
		Where("city_id = ?", id.String()).
		Count(&count).Error; err != nil {
		return false, errors.WithStack(err)
	}

	return count > 0, nil
}

func (r *repository) InsertMulti(ctx context.Context, db *gorm.DB, items []*master_domain.City) error {
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
