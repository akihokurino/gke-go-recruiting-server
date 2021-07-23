package department_domain

import (
	"gke-go-sample/domain"

	"github.com/google/uuid"
)

type Station struct {
	ID           string
	DepartmentID domain.DepartmentID
	LineID       domain.LineID

	With StationWith
}

func NewStation(departmentID domain.DepartmentID, lineID domain.LineID) *Station {
	return &Station{
		ID:           uuid.New().String(),
		DepartmentID: departmentID,
		LineID:       lineID,
	}
}
