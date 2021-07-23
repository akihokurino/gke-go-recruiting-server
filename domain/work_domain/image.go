package work_domain

import (
	"net/url"

	"gke-go-sample/domain"

	"github.com/google/uuid"
)

type Image struct {
	ID        string
	WorkID    domain.WorkID
	URL       url.URL
	ViewOrder uint64
	Comment   string
}

func NewImage(workID domain.WorkID, url url.URL, viewOrder uint64, comment string) *Image {
	return &Image{
		ID:        uuid.New().String(),
		WorkID:    workID,
		URL:       url,
		ViewOrder: viewOrder,
		Comment:   comment,
	}
}
