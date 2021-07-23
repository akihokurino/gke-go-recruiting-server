package account_domain

import "gke-go-sample/domain"

type Administrator struct {
	ID    domain.FirebaseUserID
	Email string
	Name  string
	domain.Meta
}
