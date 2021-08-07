package department_application

import (
	"gke-go-recruiting-server/adapter"
	"gke-go-recruiting-server/domain"
	pb "gke-go-recruiting-server/proto/go/pb"

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
