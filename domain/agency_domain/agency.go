package agency_domain

import "gke-go-sample/domain"

type Agency struct {
	ID         domain.AgencyID
	Name       string
	NameKana   string
	PostalCode string
	PrefID     domain.PrefID
	Address    string
	domain.Meta
}
