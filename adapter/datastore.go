package adapter

import (
	"context"

	w "go.mercari.io/datastore"
)

type DataStoreFactory func(ctx context.Context) w.Client
