package department_station_table

import (
	"context"

	"gke-go-sample/infra/cloudsql"

	"gorm.io/gorm"

	"gke-go-sample/domain"

	"github.com/pkg/errors"

	"gke-go-sample/adapter"
	"gke-go-sample/domain/department_domain"
)

func NewRepo() adapter.DepartmentStationRepo {
	return &repository{}
}

type repository struct {
}

func (r *repository) GetByRail(ctx context.Context, db *gorm.DB, railID domain.RailID) ([]*department_domain.Station, error) {
	var es []Entity
	if err := db.
		Joins("inner join `lines` on department_stations.line_id = `lines`.id").
		Where("`lines`.rail_id = ?", railID.String()).
		Preload("Line").
		Find(&es).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*department_domain.Station, 0, len(es))
	for _, e := range es {
		d, err := e.ToDomain()
		if err != nil {
			continue
		}
		items = append(items, d)
	}

	return items, nil
}

func (r *repository) GetByStation(ctx context.Context, db *gorm.DB, stationID domain.StationID) ([]*department_domain.Station, error) {
	var es []Entity
	if err := db.
		Joins("inner join `lines` on department_stations.line_id = `lines`.id").
		Where("`lines`.station_id = ?", stationID.String()).
		Preload("Line").
		Find(&es).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*department_domain.Station, 0, len(es))
	for _, e := range es {
		d, err := e.ToDomain()
		if err != nil {
			continue
		}
		items = append(items, d)
	}

	return items, nil
}

func (r *repository) InsertMulti(ctx context.Context, db *gorm.DB, items []*department_domain.Station) error {
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

func (r *repository) DeleteByDepartment(ctx context.Context, db *gorm.DB, departmentID domain.DepartmentID) error {
	if err := db.Where("department_id = ?", departmentID.String()).Delete(Entity{}).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
