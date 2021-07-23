package administrator_table

import (
	"database/sql"
	"time"

	"gke-go-sample/domain"
	"gke-go-sample/domain/account_domain"

	"github.com/guregu/null"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

func (e *Entity) TableName() string {
	return "administrators"
}

type Entity struct {
	ID        string    `gorm:"column:id;primary_key"`
	Email     string    `gorm:"column:email"`
	Name      string    `gorm:"column:name"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (e *Entity) toDomain() *account_domain.Administrator {
	return &account_domain.Administrator{
		ID:    domain.FirebaseUserID(e.ID),
		Email: e.Email,
		Name:  e.Name,
		Meta: domain.Meta{
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		},
	}
}

func entityFrom(d *account_domain.Administrator) *Entity {
	return &Entity{
		ID:        d.ID.String(),
		Email:     d.Email,
		Name:      d.Name,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}
