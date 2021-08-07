package adapter

import (
	"context"
	"net/http"

	"gorm.io/gorm"

	"gke-go-recruiting-server/domain/account_domain"

	"gke-go-recruiting-server/domain"
)

type ContextProvider interface {
	WithFirebaseUserID(ctx context.Context, id domain.FirebaseUserID) (context.Context, error)
	FirebaseUserID(ctx context.Context) (domain.FirebaseUserID, error)
}

type Cros func(base http.Handler) http.Handler

type AdminAuthenticate func(base http.Handler) http.Handler

type AdminAuthorization func(ctx context.Context, db *gorm.DB) (*account_domain.Administrator, error)

type AgencyAuthorization func(ctx context.Context, db *gorm.DB) (*account_domain.AgencyAccount, error)

type ErrorConverter func(ctx context.Context, err error) error
