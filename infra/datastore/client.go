package datastore

import (
	"context"

	"gke-go-sample/adapter"

	"cloud.google.com/go/datastore"

	w "go.mercari.io/datastore"

	"go.mercari.io/datastore/clouddatastore"
)

func NewDataStoreFactory(projectID string) adapter.DataStoreFactory {
	return func(ctx context.Context) w.Client {
		dataClient, err := datastore.NewClient(ctx, projectID)
		if err != nil {
			panic(err)
		}

		client, err := clouddatastore.FromClient(ctx, dataClient)
		if err != nil {
			panic(err)
		}

		return client
	}
}
