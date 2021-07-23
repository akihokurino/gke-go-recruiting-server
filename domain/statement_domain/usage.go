package statement_domain

import (
	"time"

	"github.com/google/uuid"

	"gke-go-sample/domain"
)

type Usage struct {
	ID             domain.UsageStatementID
	DepartmentID   domain.DepartmentID
	MainContractID domain.MainContractID
	Price          uint64
	CreatedAt      time.Time

	With UsageWith
}

func NewUsageFromMain(
	departmentID domain.DepartmentID,
	contractID domain.MainContractID,
	price uint64,
	now time.Time) *Usage {
	return &Usage{
		ID:             domain.UsageStatementID(uuid.New().String()),
		DepartmentID:   departmentID,
		MainContractID: contractID,
		Price:          price,
		CreatedAt:      now,
	}
}
