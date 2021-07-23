package company_application

import (
	"gke-go-sample/adapter"
	"gke-go-sample/domain"
	pb "gke-go-sample/proto/go/pb"

	"github.com/pkg/errors"
)

func validateCreate(params adapter.CompanyParams) error {
	if params.Rank == pb.Company_Rank_Unknown ||
		params.RankType == pb.Company_RankType_Unknown ||
		params.Name == "" ||
		params.NameKana == "" ||
		params.PostalCode == "" ||
		params.PrefID == "" ||
		params.Address == "" ||
		params.PhoneNumber == "" {
		return errors.WithStack(domain.NewBadRequestErr(domain.BadRequestMsg))
	}
	return nil
}
