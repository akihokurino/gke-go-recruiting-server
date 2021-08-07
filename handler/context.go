package handler

import (
	"context"

	"gke-go-recruiting-server/adapter"
	"gke-go-recruiting-server/domain"
)

const (
	firebaseUserIDStoreKey = "__firebase_user_id_store_key__"
)

type contextProvider struct {
}

func NewContextProvider() adapter.ContextProvider {
	return &contextProvider{}
}

func (c *contextProvider) WithFirebaseUserID(ctx context.Context, id domain.FirebaseUserID) (context.Context, error) {
	return context.WithValue(ctx, firebaseUserIDStoreKey, id), nil
}

func (c *contextProvider) FirebaseUserID(ctx context.Context) (domain.FirebaseUserID, error) {
	id, ok := ctx.Value(firebaseUserIDStoreKey).(domain.FirebaseUserID)
	if !ok {
		return "", domain.NewInternalServerErr()
	}
	return id, nil
}
