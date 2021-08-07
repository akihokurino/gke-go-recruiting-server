package work_domain

import (
	"net/url"

	"gke-go-recruiting-server/domain"

	"github.com/google/uuid"
)

type Movie struct {
	ID     string
	WorkID domain.WorkID
	URL    url.URL
}

func NewMovie(workID domain.WorkID, url url.URL) *Movie {
	return &Movie{
		ID:     uuid.New().String(),
		WorkID: workID,
		URL:    url,
	}
}
