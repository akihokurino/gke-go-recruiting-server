package contract_application

import (
	"gke-go-recruiting-server/adapter"
	"gke-go-recruiting-server/domain"
	pb "gke-go-recruiting-server/proto/go/pb"

	"github.com/pkg/errors"
)

func validateCreateMainContract(params adapter.MainContractParams) error {
	if params.Plan == pb.MainProduct_Plan_Unknown || params.DateFrom.IsZero() {
		return errors.WithStack(domain.NewBadRequestErr(domain.BadRequestMsg))
	}
	return nil
}
