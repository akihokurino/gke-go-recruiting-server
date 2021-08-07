package work_movie_table

import (
	"database/sql"
	"net/url"
	"time"

	"gke-go-recruiting-server/domain/work_domain"

	"gke-go-recruiting-server/domain"

	"github.com/guregu/null"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

func (e *Entity) TableName() string {
	return "work_movies"
}

type Entity struct {
	ID     string `gorm:"column:id;primary_key"`
	WorkID string `gorm:"column:work_id"`
	URL    string `gorm:"column:url"`
}

func (e *Entity) ToDomain() *work_domain.Movie {
	u, _ := url.Parse(e.URL)

	return &work_domain.Movie{
		ID:     e.ID,
		WorkID: domain.WorkID(e.WorkID),
		URL:    *u,
	}
}

func entityFrom(d *work_domain.Movie) *Entity {
	return &Entity{
		ID:     d.ID,
		WorkID: d.WorkID.String(),
		URL:    d.URL.String(),
	}
}
