package line_table

import (
	"database/sql"
	"time"

	"gke-go-recruiting-server/domain/master_domain"

	"gke-go-recruiting-server/domain"

	"github.com/guregu/null"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

func (e *Entity) TableName() string {
	return "lines"
}

type Entity struct {
	ID                  string  `gorm:"column:id;primary_key"`
	RailID              string  `gorm:"column:rail_id"`
	StationID           string  `gorm:"column:station_id"`
	StopOrder           uint64  `gorm:"column:stop_order"`
	RailCompanyName     string  `gorm:"column:rail_company_name"`
	RailCompanyNameKana string  `gorm:"column:rail_company_kana"`
	RailCompanyName2    string  `gorm:"column:rail_company_name2"`
	RailName1           string  `gorm:"column:rail_name1"`
	RailNameKana1       string  `gorm:"column:rail_name_kana1"`
	RailName2           string  `gorm:"column:rail_name2"`
	RailNameKana2       string  `gorm:"column:rail_name_kana2"`
	StationName         string  `gorm:"column:station_name"`
	StationNameKana     string  `gorm:"column:station_name_kana"`
	PrefID              string  `gorm:"column:pref_id"`
	Latitude            float64 `gorm:"column:latitude"`
	Longitude           float64 `gorm:"column:longitude"`
	RailKind            uint64  `gorm:"column:rail_kind"`
	RailKindName        string  `gorm:"column:rail_kind_name"`
}

func (e *Entity) ToDomain() *master_domain.Line {
	return &master_domain.Line{
		ID:              domain.LineID(e.ID),
		Rail:            *e.toRailDomain(),
		StationID:       domain.StationID(e.StationID),
		StopOrder:       e.StopOrder,
		StationName:     e.StationName,
		StationNameKana: e.StationNameKana,
	}
}

func (e *Entity) toRailDomain() *master_domain.Rail {
	return &master_domain.Rail{
		ID:        domain.RailID(e.RailID),
		Name1:     e.RailName1,
		NameKana1: e.RailNameKana1,
		Name2:     e.RailName2,
		NameKana2: e.RailNameKana2,
		Company:   *e.toRailCompanyDomain(),
	}
}

func (e *Entity) toRailCompanyDomain() *master_domain.RailCompany {
	return &master_domain.RailCompany{
		Name1:     e.RailCompanyName,
		NameKana1: e.RailCompanyNameKana,
		Name2:     e.RailCompanyName2,
		Kind:      e.RailKind,
		KindName:  e.RailKindName,
	}
}

func entityFrom(d *master_domain.Line) *Entity {
	return &Entity{
		ID:                  d.ID.String(),
		RailID:              d.Rail.ID.String(),
		StationID:           d.StationID.String(),
		StopOrder:           d.StopOrder,
		RailCompanyName:     d.Rail.Company.Name1,
		RailCompanyNameKana: d.Rail.Company.NameKana1,
		RailCompanyName2:    d.Rail.Company.Name2,
		RailName1:           d.Rail.Name1,
		RailNameKana1:       d.Rail.NameKana1,
		RailName2:           d.Rail.Name2,
		RailNameKana2:       d.Rail.NameKana2,
		StationName:         d.StationName,
		StationNameKana:     d.StationNameKana,
		PrefID:              d.PrefID.String(),
		Latitude:            d.Latitude,
		Longitude:           d.Longitude,
		RailKind:            d.Rail.Company.Kind,
		RailKindName:        d.Rail.Company.KindName,
	}
}
