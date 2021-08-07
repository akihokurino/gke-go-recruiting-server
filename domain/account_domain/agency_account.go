package account_domain

import "gke-go-recruiting-server/domain"

type AgencyAccount struct {
	ID       domain.FirebaseUserID
	V1ID     string
	AgencyID domain.AgencyID
	Email    string
	Name     string
	NameKana string
	domain.Meta
}
