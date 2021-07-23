package firebase

import (
	"context"

	"gke-go-sample/adapter"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func NewFirebaseAppFactory(opts option.ClientOption) adapter.FirebaseAppFactory {
	return func(ctx context.Context) *firebase.App {
		app, err := firebase.NewApp(ctx, nil, opts)
		if err != nil {
			panic(err)
		}

		return app
	}
}
