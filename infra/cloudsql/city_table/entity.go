package city_table

import (
	"database/sql"
	"fmt"
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
	return "cities"
}

type Entity struct {
	ID           string `gorm:"column:id;primary_key"`
	PrefID       string `gorm:"column:pref_id"`
	PrefName     string `gorm:"column:pref_name"`
	CityID       string `gorm:"column:city_id"`
	CityName     string `gorm:"column:city_name"`
	CityNameKana string `gorm:"column:city_name_kana"`
	AreaName     string `gorm:"column:area_name"`
	AreaNameKana string `gorm:"column:area_name_kana"`
}

func (e *Entity) toDomain() *master_domain.City {
	return &master_domain.City{
		Prefecture: &master_domain.Prefecture{
			ID:   domain.PrefID(e.PrefID),
			Name: e.PrefName,
		},
		ID:           domain.CityID(e.CityID),
		Name:         e.CityName,
		NameKana:     e.CityNameKana,
		AreaName:     e.AreaName,
		AreaNameKana: e.AreaNameKana,
	}
}

func (e *Entity) toPrefectureDomain() *master_domain.Prefecture {
	return &master_domain.Prefecture{
		ID:   domain.PrefID(e.PrefID),
		Name: e.PrefName,
	}
}

func entityFrom(d *master_domain.City) *Entity {
	return &Entity{
		ID:           fmt.Sprintf("%s_%s", d.Prefecture.ID, d.ID),
		PrefID:       d.Prefecture.ID.String(),
		PrefName:     d.Prefecture.Name,
		CityID:       d.ID.String(),
		CityName:     d.Name,
		CityNameKana: d.NameKana,
		AreaName:     d.AreaName,
		AreaNameKana: d.AreaNameKana,
	}
}
