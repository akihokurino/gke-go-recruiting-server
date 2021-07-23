package main_contract_table

import (
	"database/sql"
	"time"

	"gke-go-sample/infra/cloudsql/department_table"

	pb "gke-go-sample/proto/go/pb"

	"gke-go-sample/domain/contract_domain"

	"gke-go-sample/domain"

	"github.com/guregu/null"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

func (e *Entity) TableName() string {
	return "main_contracts"
}

type Entity struct {
	ID           string    `gorm:"column:id;primary_key"`
	DepartmentID string    `gorm:"column:department_id"`
	Status       int32     `gorm:"column:status"`
	Plan         int32     `gorm:"column:plan"`
	DateFrom     time.Time `gorm:"column:date_from"`
	DateTo       time.Time `gorm:"column:date_to"`
	Price        uint64    `gorm:"column:price"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`

	Department *department_table.Entity `gorm:"PRELOAD:false;foreignkey:department_id"`
}

func (e *Entity) ToDomain() *contract_domain.Main {
	var department *contract_domain.DepartmentOverview
	if e.Department != nil {
		department = e.Department.ToContractDomainOverview()
	}

	return &contract_domain.Main{
		ID:           domain.MainContractID(e.ID),
		DepartmentID: domain.DepartmentID(e.DepartmentID),
		Status:       pb.MainContract_Status(e.Status),
		Plan:         pb.MainProduct_Plan(e.Plan),
		DateRange: domain.DateRange{
			From: e.DateFrom,
			To:   e.DateTo,
		},
		Price: e.Price,
		Meta: domain.Meta{
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		},

		Department: department,
	}
}

func entityFrom(d *contract_domain.Main) *Entity {
	return &Entity{
		ID:           d.ID.String(),
		DepartmentID: d.DepartmentID.String(),
		Status:       int32(d.Status),
		Plan:         int32(d.Plan),
		DateFrom:     d.DateRange.From,
		DateTo:       d.DateRange.To,
		Price:        d.Price,
		CreatedAt:    d.Meta.CreatedAt,
		UpdatedAt:    d.Meta.UpdatedAt,
	}
}
