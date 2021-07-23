package account_domain

import (
	"time"

	"gke-go-sample/domain"
)

type FirebaseUser struct {
	ID    domain.FirebaseUserID
	Email string
}

func (u *FirebaseUser) NewAdministrator(name string, now time.Time) *Administrator {
	return &Administrator{
		ID:    u.ID,
		Email: u.Email,
		Name:  name,
		Meta: domain.Meta{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

func (u *FirebaseUser) NewAgencyAccount(
	agencyID domain.AgencyID,
	name string,
	nameKana string,
	now time.Time) *AgencyAccount {
	return &AgencyAccount{
		ID:       u.ID,
		AgencyID: agencyID,
		Email:    u.Email,
		Name:     name,
		NameKana: nameKana,
		Meta: domain.Meta{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}
