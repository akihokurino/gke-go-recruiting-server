package usage_statement

import (
	"database/sql"
	"time"

	"gke-go-sample/domain/contract_domain"

	"gke-go-sample/infra/cloudsql/main_contract_table"

	"gke-go-sample/domain/statement_domain"

	"gke-go-sample/domain"

	"github.com/guregu/null"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

func (e *Entity) TableName() string {
	return "usage_statements"
}

type Entity struct {
	ID             string         `gorm:"column:id;primary_key"`
	DepartmentID   string         `gorm:"column:department_id"`
	MainContractID sql.NullString `gorm:"column:main_contract_id"`
	Price          uint64         `gorm:"column:price"`
	CreatedAt      time.Time      `gorm:"column:created_at"`

	MainContract *main_contract_table.Entity `gorm:"PRELOAD:false;foreignkey:main_contract_id"`
}

func (e *Entity) toDomain() *statement_domain.Usage {
	var mainContract *contract_domain.Main
	if e.MainContract != nil {
		mainContract = e.MainContract.ToDomain()
	}

	return &statement_domain.Usage{
		ID:             domain.UsageStatementID(e.ID),
		DepartmentID:   domain.DepartmentID(e.DepartmentID),
		MainContractID: domain.MainContractID(e.MainContractID.String),
		Price:          e.Price,
		CreatedAt:      e.CreatedAt,

		With: statement_domain.UsageWith{
			MainContract: mainContract,
		},
	}
}

func entityFrom(d *statement_domain.Usage) *Entity {
	return &Entity{
		ID:           d.ID.String(),
		DepartmentID: d.DepartmentID.String(),
		MainContractID: sql.NullString{
			Valid:  d.MainContractID.String() != "",
			String: d.MainContractID.String(),
		},
		Price:     d.Price,
		CreatedAt: d.CreatedAt,
	}
}
