package department_image_table

import (
	"database/sql"
	"net/url"
	"time"

	"gke-go-sample/domain/department_domain"

	"github.com/guregu/null"

	"gke-go-sample/domain"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

func (e *Entity) TableName() string {
	return "department_images"
}

type Entity struct {
	ID           string `gorm:"column:id;primary_key"`
	DepartmentID string `gorm:"column:department_id"`
	URL          string `gorm:"column:url"`
}

func (e *Entity) ToDomain() *department_domain.Image {
	u, _ := url.Parse(e.URL)

	return &department_domain.Image{
		ID:           e.ID,
		DepartmentID: domain.DepartmentID(e.DepartmentID),
		URL:          *u,
	}
}

func entityFrom(d *department_domain.Image) *Entity {
	return &Entity{
		ID:           d.ID,
		DepartmentID: d.DepartmentID.String(),
		URL:          d.URL.String(),
	}
}
