package adapter

import (
	"context"

	"gke-go-sample/domain"

	"gke-go-sample/domain/account_domain"

	firebase "firebase.google.com/go"
)

type FirebaseAppFactory func(ctx context.Context) *firebase.App

type FirebaseRepo interface {
	GetByEmail(ctx context.Context, email string) (*account_domain.FirebaseUser, error)
	CreateEmailUser(ctx context.Context, email string, password string) (*account_domain.FirebaseUser, error)
	Delete(ctx context.Context, id domain.FirebaseUserID) error
}
