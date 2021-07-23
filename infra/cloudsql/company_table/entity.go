package company_table

import (
	"database/sql"
	"time"

	"gke-go-sample/domain/company_domain"

	"github.com/guregu/null"

	"gke-go-sample/domain"
	pb "gke-go-sample/proto/go/pb"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

func (e *Entity) TableName() string {
	return "companies"
}

type Entity struct {
	ID          string    `gorm:"column:id;primary_key"`
	AgencyID    string    `gorm:"column:agency_id"`
	Status      int32     `gorm:"column:status"`
	RankType    int32     `gorm:"column:rank_type"`
	Rank        int32     `gorm:"column:rank"`
	Name        string    `gorm:"column:name"`
	NameKana    string    `gorm:"column:name_kana"`
	PostalCode  string    `gorm:"column:postal_code"`
	PrefID      string    `gorm:"column:pref_id"`
	Address     string    `gorm:"column:address"`
	Building    string    `gorm:"column:building"`
	PhoneNumber string    `gorm:"column:phone_number"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (e *Entity) toDomain() *company_domain.Company {
	return &company_domain.Company{
		ID:          domain.CompanyID(e.ID),
		AgencyID:    domain.AgencyID(e.AgencyID),
		Status:      pb.Company_Status(e.Status),
		RankType:    pb.Company_RankType(e.RankType),
		Rank:        pb.Company_Rank(e.Rank),
		Name:        e.Name,
		NameKana:    e.NameKana,
		PostalCode:  e.PostalCode,
		PrefID:      domain.PrefID(e.PrefID),
		Address:     e.Address,
		Building:    e.Building,
		PhoneNumber: e.PhoneNumber,
		Meta: domain.Meta{
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		},
	}
}

func entityFrom(d *company_domain.Company) *Entity {
	return &Entity{
		ID:          d.ID.String(),
		AgencyID:    d.AgencyID.String(),
		Status:      int32(d.Status),
		RankType:    int32(d.RankType),
		Rank:        int32(d.Rank),
		Name:        d.Name,
		NameKana:    d.NameKana,
		PostalCode:  d.PostalCode,
		PrefID:      d.PrefID.String(),
		Address:     d.Address,
		Building:    d.Building,
		PhoneNumber: d.PhoneNumber,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}
