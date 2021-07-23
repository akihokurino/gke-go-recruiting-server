package department_station_table

import (
	"database/sql"
	"time"

	"gke-go-sample/infra/cloudsql/line_table"

	"gke-go-sample/domain"
	"gke-go-sample/domain/department_domain"

	"github.com/guregu/null"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

func (e *Entity) TableName() string {
	return "department_stations"
}

type Entity struct {
	ID           string `gorm:"column:id;primary_key"`
	DepartmentID string `gorm:"column:department_id"`
	LineID       string `gorm:"column:line_id"`

	Line *line_table.Entity `gorm:"PRELOAD:false;foreignkey:line_id"`
}

func (e *Entity) ToDomain() (*department_domain.Station, error) {
	if e.Line == nil {
		return nil, domain.NewNotFoundErr()
	}

	return &department_domain.Station{
		ID:           e.ID,
		DepartmentID: domain.DepartmentID(e.DepartmentID),
		LineID:       domain.LineID(e.LineID),
		With: department_domain.StationWith{
			Line: e.Line.ToDomain(),
		},
	}, nil
}

func entityFrom(d *department_domain.Station) *Entity {
	return &Entity{
		ID:           d.ID,
		DepartmentID: d.DepartmentID.String(),
		LineID:       d.LineID.String(),
	}
}
