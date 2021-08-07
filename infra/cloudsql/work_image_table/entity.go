package work_image_table

import (
	"database/sql"
	"net/url"
	"time"

	"gke-go-recruiting-server/domain/work_domain"

	"github.com/guregu/null"

	"gke-go-recruiting-server/domain"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

func (e *Entity) TableName() string {
	return "work_images"
}

type Entity struct {
	ID        string `gorm:"column:id;primary_key"`
	WorkID    string `gorm:"column:work_id"`
	URL       string `gorm:"column:url"`
	ViewOrder uint64 `gorm:"column:view_order"`
	Comment   string `gorm:"column:comment"`
}

func (e *Entity) ToDomain() *work_domain.Image {
	u, _ := url.Parse(e.URL)

	return &work_domain.Image{
		ID:        e.ID,
		WorkID:    domain.WorkID(e.WorkID),
		URL:       *u,
		ViewOrder: e.ViewOrder,
		Comment:   e.Comment,
	}
}

func entityFrom(d *work_domain.Image) *Entity {
	return &Entity{
		ID:        d.ID,
		WorkID:    d.WorkID.String(),
		URL:       d.URL.String(),
		ViewOrder: d.ViewOrder,
		Comment:   d.Comment,
	}
}
