package entry_application

import (
	"gke-go-sample/adapter"
	"gke-go-sample/domain"
	pb "gke-go-sample/proto/go/pb"

	"github.com/pkg/errors"
)

func validateEntry(params adapter.EntryParams) error {
	if params.FullName == "" ||
		params.FullNameKana == "" ||
		params.Birthdate.IsZero() ||
		params.Gender == pb.User_Gender_Unknown ||
		params.PhoneNumber == "" ||
		params.Email == "" {
		return errors.WithStack(domain.NewBadRequestErr(domain.BadRequestMsg))
	}
	return nil
}
