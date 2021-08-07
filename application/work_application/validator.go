package work_application

import (
	"gke-go-recruiting-server/adapter"
	"gke-go-recruiting-server/domain"
	pb "gke-go-recruiting-server/proto/go/pb"

	"github.com/pkg/errors"
)

func validateCreate(params adapter.WorkParams) error {
	if params.WorkType == pb.Work_Type_Unknown || params.JobCode == pb.Work_Job_Unknown || params.Title == "" {
		return errors.WithStack(domain.NewBadRequestErr(domain.BadRequestMsg))
	}
	return nil
}

func validateUpdate(params adapter.WorkParams) error {
	return validateCreate(params)
}
