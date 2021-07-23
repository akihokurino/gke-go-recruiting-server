package agency_table

import (
	"database/sql"
	"time"

	"gke-go-sample/domain/agency_domain"

	"gke-go-sample/domain"

	"github.com/guregu/null"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

func (e *Entity) TableName() string {
	return "agencies"
}

type Entity struct {
	ID         string    `gorm:"column:id;primary_key"`
	Name       string    `gorm:"column:name"`
	NameKana   string    `gorm:"column:name_kana"`
	PostalCode string    `gorm:"column:postal_code"`
	PrefID     string    `gorm:"column:pref_id"`
	Address    string    `gorm:"column:address"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (e *Entity) toDomain() *agency_domain.Agency {
	return &agency_domain.Agency{
		ID:         domain.AgencyID(e.ID),
		Name:       e.Name,
		NameKana:   e.NameKana,
		PostalCode: e.PostalCode,
		PrefID:     domain.PrefID(e.PrefID),
		Address:    e.Address,
		Meta: domain.Meta{
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		},
	}
}

func entityFrom(d *agency_domain.Agency) *Entity {
	return &Entity{
		ID:         d.ID.String(),
		Name:       d.Name,
		NameKana:   d.NameKana,
		PostalCode: d.PostalCode,
		PrefID:     d.PrefID.String(),
		Address:    d.Address,
		CreatedAt:  d.CreatedAt,
		UpdatedAt:  d.UpdatedAt,
	}
}
