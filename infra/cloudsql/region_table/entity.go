package region_table

import (
	"database/sql"
	"time"

	"gke-go-sample/domain/master_domain"

	"gke-go-sample/domain"

	"github.com/guregu/null"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

func (e *Entity) TableName() string {
	return "regions"
}

type Entity struct {
	Geocode1          string         `gorm:"column:geocode_1"`
	Geocode2          string         `gorm:"column:geocode_2"`
	Zipcode           string         `gorm:"column:zipcode"`
	Address           string         `gorm:"column:address"`
	LArea             string         `gorm:"column:l_area"`
	LAreaName         string         `gorm:"column:l_area_name"`
	OriginalMArea     string         `gorm:"column:original_m_area"`
	OriginalMAreaName string         `gorm:"column:original_m_area_name"`
	OriginalSArea     string         `gorm:"column:original_s_area"`
	OriginalSAreaName string         `gorm:"column:original_s_area_name"`
	MArea             string         `gorm:"column:m_area"`
	MAreaName         string         `gorm:"column:m_area_name"`
	SArea             sql.NullString `gorm:"column:s_area"`
	SAreaName         sql.NullString `gorm:"column:s_area_name"`
}

func (e *Entity) toDomain() *master_domain.Region {
	return &master_domain.Region{
		Geocode1: e.Geocode1,
		Geocode2: e.Geocode2,
		Zipcode:  e.Zipcode,
		Address:  e.Address,
		LArea: master_domain.LArea{
			ID:   domain.LAreaID(e.LArea),
			Name: e.LAreaName,
		},
		OriginalMArea:     e.OriginalMArea,
		OriginalMAreaName: e.OriginalMAreaName,
		OriginalSArea:     e.OriginalSArea,
		OriginalSAreaName: e.OriginalSAreaName,
		MArea:             *e.toMAreaDomain(),
		SArea:             *e.toSAreaDomain(),
	}
}

func (e *Entity) toMAreaDomain() *master_domain.MArea {
	return &master_domain.MArea{
		ID:   domain.MAreaID(e.MArea),
		Name: e.MAreaName,
	}
}

func (e *Entity) toSAreaDomain() *master_domain.SArea {
	return &master_domain.SArea{
		ID:   domain.SAreaID(e.SArea.String),
		Name: e.SAreaName.String,
	}
}
