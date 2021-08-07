package contract_domain

import (
	"time"

	"github.com/google/uuid"

	"gke-go-recruiting-server/domain"
	pb "gke-go-recruiting-server/proto/go/pb"
)

type Main struct {
	ID           domain.MainContractID
	DepartmentID domain.DepartmentID
	Status       pb.MainContract_Status
	Plan         pb.MainProduct_Plan
	DateRange    domain.DateRange
	Price        uint64
	domain.Meta

	Department *DepartmentOverview
}

func NewMainContract(
	departmentID domain.DepartmentID,
	plan pb.MainProduct_Plan,
	dateRange domain.DateRange,
	price uint64,
	now time.Time) *Main {
	return &Main{
		ID:           domain.MainContractID(uuid.New().String()),
		DepartmentID: departmentID,
		Status:       pb.MainContract_Status_Review,
		Plan:         plan,
		DateRange:    dateRange,
		Price:        price,
		Meta: domain.Meta{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

func (m *Main) Accept() error {
	if m.Status != pb.MainContract_Status_Review {
		return domain.NewConflictErr(domain.ConflictStatusMsg)
	}

	m.Status = pb.MainContract_Status_OK
	return nil
}

func (m *Main) Deny() error {
	if m.Status != pb.MainContract_Status_Review {
		return domain.NewConflictErr(domain.ConflictStatusMsg)
	}

	m.Status = pb.MainContract_Status_NG
	return nil
}
