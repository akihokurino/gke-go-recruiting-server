package department_application

import (
	"gke-go-sample/adapter"
	"gke-go-sample/domain"
	pb "gke-go-sample/proto/go/pb"

	"github.com/pkg/errors"
)

func validateCreate(departmentParams adapter.DepartmentParams) error {
	if departmentParams.Name == "" ||
		departmentParams.BusinessCondition == pb.Department_BusinessCondition_Unknown ||
		departmentParams.PostalCode == "" ||
		departmentParams.PrefID == "" ||
		departmentParams.Address == "" ||
		departmentParams.PhoneNumber == "" {
		return errors.WithStack(domain.NewBadRequestErr(domain.BadRequestMsg))
	}
	return nil
}

func validateUpdate(params adapter.DepartmentParams) error {
	return validateCreate(params)
}
