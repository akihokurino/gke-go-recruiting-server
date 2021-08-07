package account_domain

import "gke-go-recruiting-server/domain"

type Administrator struct {
	ID    domain.FirebaseUserID
	Email string
	Name  string
	domain.Meta
}
