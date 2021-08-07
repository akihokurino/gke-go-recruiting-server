package account_application

import (
	"gke-go-recruiting-server/adapter"
	"gke-go-recruiting-server/domain"

	"github.com/pkg/errors"
)

func validateCreateAdministrator(params adapter.AdministratorParams) error {
	if params.Email == "" ||
		params.Password == "" ||
		params.Name == "" {
		return errors.WithStack(domain.NewBadRequestErr(domain.BadRequestMsg))
	}
	return nil
}

func validateCreateAgencyAccount(params adapter.AgencyAccountParams) error {
	if params.Email == "" ||
		params.Password == "" ||
		params.Name == "" ||
		params.NameKana == "" {
		return errors.WithStack(domain.NewBadRequestErr(domain.BadRequestMsg))
	}
	return nil
}
