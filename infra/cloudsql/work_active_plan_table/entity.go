package work_active_plan_table

import (
	"database/sql"
	"time"

	"gke-go-sample/domain/contract_domain"

	"gke-go-sample/infra/cloudsql/main_contract_table"

	"gke-go-sample/domain"

	"github.com/guregu/null"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

func (e *Entity) TableName() string {
	return "work_active_plans"
}

type Entity struct {
	WorkID         string         `gorm:"column:work_id;primary_key"`
	MainContractID sql.NullString `gorm:"column:main_contract_id"`
	PublishedOrder int            `gorm:"column:published_order"`

	MainContract *main_contract_table.Entity `gorm:"PRELOAD:false;foreignkey:main_contract_id"`
}

func (e *Entity) ToDomain() *contract_domain.ActivePlan {
	var contract *contract_domain.Main
	if e.MainContract != nil {
		contract = e.MainContract.ToDomain()
	}

	return &contract_domain.ActivePlan{
		WorkID:         domain.WorkID(e.WorkID),
		MainContractID: domain.MainContractID(e.MainContractID.String),
		PublishedOrder: e.PublishedOrder,
		MainContract:   contract,
	}
}

func entityFrom(d *contract_domain.ActivePlan) *Entity {
	return &Entity{
		WorkID: d.WorkID.String(),
		MainContractID: sql.NullString{
			String: d.MainContractID.String(),
			Valid:  d.MainContractID.String() != "",
		},
		PublishedOrder: d.PublishedOrder,
	}
}

func onlyID(workID domain.WorkID) *Entity {
	return &Entity{
		WorkID: workID.String(),
	}
}
