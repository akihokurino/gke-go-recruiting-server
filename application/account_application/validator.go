package account_application

import (
	"gke-go-sample/adapter"
	"gke-go-sample/domain"

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
