package firebase

import (
	"context"

	"gke-go-sample/domain"

	"github.com/pkg/errors"

	"gke-go-sample/domain/account_domain"

	"firebase.google.com/go/auth"

	"gke-go-sample/adapter"
)

func NewRepo(fc adapter.FirebaseAppFactory) adapter.FirebaseRepo {
	return &repository{
		fc: fc,
	}
}

type repository struct {
	fc adapter.FirebaseAppFactory
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*account_domain.FirebaseUser, error) {
	authClient, err := r.fc(ctx).Auth(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	userRecord, err := authClient.GetUserByEmail(ctx, email)
	if err != nil {
		if auth.IsUserNotFound(err) {
			return nil, domain.NewNotFoundErr()
		}
		return nil, errors.WithStack(err)
	}

	return &account_domain.FirebaseUser{
		ID:    domain.FirebaseUserID(userRecord.UID),
		Email: userRecord.Email,
	}, nil
}

func (r *repository) CreateEmailUser(ctx context.Context, email string, password string) (*account_domain.FirebaseUser, error) {
	authClient, err := r.fc(ctx).Auth(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password).
		EmailVerified(false).
		Disabled(false)

	userRecord, err := authClient.CreateUser(ctx, params)
	if err != nil {
		if auth.IsEmailAlreadyExists(err) {
			return nil, domain.NewConflictErr("そのメールアドレスはすでに存在します")
		}
		if auth.IsInvalidEmail(err) {
			return nil, domain.NewBadRequestErr("そのメールアドレスは不正です")
		}
		return nil, errors.WithStack(err)
	}

	return &account_domain.FirebaseUser{
		ID:    domain.FirebaseUserID(userRecord.UID),
		Email: userRecord.Email,
	}, nil
}

func (r *repository) Delete(ctx context.Context, id domain.FirebaseUserID) error {
	authClient, err := r.fc(ctx).Auth(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := authClient.DeleteUser(ctx, id.String()); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
