package region_table

import (
	"context"

	"gorm.io/gorm"

	"gke-go-recruiting-server/domain/master_domain"

	"github.com/pkg/errors"

	"gke-go-recruiting-server/adapter"
	"gke-go-recruiting-server/domain"
)

func NewRepo() adapter.RegionRepo {
	return &repository{}
}

type repository struct {
}

func (r *repository) GetMAreaByLArea(ctx context.Context, db *gorm.DB, lArea domain.LAreaID) ([]*master_domain.MArea, error) {
	var es []Entity
	if err := db.
		Select("DISTINCT m_area, m_area_name").
		Where("l_area = ?", lArea.String()).
		Find(&es).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*master_domain.MArea, 0, len(es))
	for _, e := range es {
		items = append(items, e.toMAreaDomain())
	}

	return items, nil
}

func (r *repository) GetSAreaByMArea(ctx context.Context, db *gorm.DB, mArea domain.MAreaID) ([]*master_domain.SArea, error) {
	var es []Entity
	if err := db.
		Select("DISTINCT s_area, s_area_name").
		Where("m_area = ?", mArea.String()).
		Find(&es).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*master_domain.SArea, 0, len(es))
	for _, e := range es {
		items = append(items, e.toSAreaDomain())
	}

	return items, nil
}
