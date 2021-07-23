package agency_account_table

import (
	"database/sql"
	"time"

	"gke-go-sample/domain/department_domain"

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
	return "agency_accounts"
}

type Entity struct {
	ID        string         `gorm:"column:id;primary_key"`
	V1ID      sql.NullString `gorm:"column:v1_id"`
	AgencyID  string         `gorm:"column:agency_id"`
	Email     string         `gorm:"column:email"`
	Name      string         `gorm:"column:name"`
	NameKana  string         `gorm:"column:name_kana"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
}

func (e *Entity) toDomain() *account_domain.AgencyAccount {
	return &account_domain.AgencyAccount{
		ID:       domain.FirebaseUserID(e.ID),
		V1ID:     e.V1ID.String,
		AgencyID: domain.AgencyID(e.AgencyID),
		Email:    e.Email,
		Name:     e.Name,
		NameKana: e.NameKana,
		Meta: domain.Meta{
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		},
	}
}

func (e *Entity) ToDepartmentDomainOverview() *department_domain.SalesOverview {
	return &department_domain.SalesOverview{
		ID:   domain.FirebaseUserID(e.ID),
		Name: e.Name,
	}
}

func entityFrom(d *account_domain.AgencyAccount) *Entity {
	return &Entity{
		ID: d.ID.String(),
		V1ID: sql.NullString{
			String: d.V1ID,
			Valid:  d.V1ID != "",
		},
		AgencyID:  d.AgencyID.String(),
		Email:     d.Email,
		Name:      d.Name,
		NameKana:  d.NameKana,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func onlyID(d domain.FirebaseUserID) *Entity {
	return &Entity{
		ID: d.String(),
	}
}
