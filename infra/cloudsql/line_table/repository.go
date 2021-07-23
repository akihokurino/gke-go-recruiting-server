package line_table

import (
	"context"
	"fmt"

	"gke-go-sample/infra/cloudsql"

	"gorm.io/gorm"

	"github.com/pkg/errors"

	"gke-go-sample/domain/master_domain"

	"gke-go-sample/adapter"
	"gke-go-sample/domain"
)

func NewRepo() adapter.LineRepo {
	return &repository{}
}

type repository struct {
}

func (r *repository) GetRailByPrefecture(ctx context.Context, db *gorm.DB, prefID domain.PrefID) ([]*master_domain.Rail, error) {
	var es []Entity
	if err := db.
		Select("DISTINCT rail_id, rail_name1, rail_name_kana1, rail_name2, rail_name_kana2, rail_company_name rail_company_kana, rail_company_name2").
		Where("pref_id = ?", prefID.String()).
		Find(&es).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*master_domain.Rail, 0, len(es))
	for _, e := range es {
		items = append(items, e.toRailDomain())
	}

	return items, nil
}

func (r *repository) GetRailByCompany(ctx context.Context, db *gorm.DB, companyName string) ([]*master_domain.Rail, error) {
	var es []Entity
	if err := db.
		Select("DISTINCT rail_id, rail_name1, rail_name_kana1, rail_name2, rail_name_kana2, rail_company_name rail_company_kana, rail_company_name2").
		Where("rail_company_name = ?", companyName).
		Find(&es).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*master_domain.Rail, 0, len(es))
	for _, e := range es {
		items = append(items, e.toRailDomain())
	}

	return items, nil
}

func (r *repository) GetByRail(ctx context.Context, db *gorm.DB, railID domain.RailID) ([]*master_domain.Line, error) {
	var es []Entity
	if err := db.
		Where("rail_id = ?", railID.String()).
		Order("stop_order ASC").
		Find(&es).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*master_domain.Line, 0, len(es))
	for _, e := range es {
		items = append(items, e.ToDomain())
	}

	return items, nil
}

func (r *repository) GetRailCompanyByPrefecture(ctx context.Context, db *gorm.DB, prefID domain.PrefID) ([]*master_domain.RailCompany, error) {
	var es []Entity
	if err := db.
		Select("DISTINCT rail_company_name, rail_company_kana, rail_company_name2, rail_kind, rail_kind_name").
		Where("pref_id = ?", prefID.String()).
		Find(&es).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*master_domain.RailCompany, 0, len(es))
	for _, e := range es {
		items = append(items, e.toRailCompanyDomain())
	}

	return items, nil
}

func (r *repository) GetByDistance(
	ctx context.Context,
	db *gorm.DB,
	latitude float64,
	longitude float64,
	distanceKM uint64,
) ([]*master_domain.Line, error) {
	var es []Entity

	queryTemplate := `select *, 
(6371 * acos(cos(radians(?)) * cos(radians(latitude)) * cos(radians(longitude) - radians(?)) + sin(radians(?)) * sin(radians(latitude)))) AS distance 
from %s having distance <= ? order by distance`

	q := fmt.Sprintf(queryTemplate, fmt.Sprintf("`%s`", (&Entity{}).TableName()))

	if err := db.Raw(q, latitude, longitude, latitude, distanceKM).Scan(&es).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*master_domain.Line, 0, len(es))
	for _, e := range es {
		items = append(items, e.ToDomain())
	}

	return items, nil
}

func (r *repository) Exist(ctx context.Context, db *gorm.DB, id domain.LineID) (bool, error) {
	var count int64
	if err := db.
		Model(&Entity{}).
		Where("id = ?", id.String()).
		Count(&count).Error; err != nil {
		return false, errors.WithStack(err)
	}

	return count > 0, nil
}

func (r *repository) InsertMulti(ctx context.Context, db *gorm.DB, items []*master_domain.Line) error {
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
