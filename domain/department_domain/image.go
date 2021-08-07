package department_domain

import (
	"net/url"

	"gke-go-recruiting-server/domain"

	"github.com/google/uuid"
)

type Image struct {
	ID           string
	DepartmentID domain.DepartmentID
	URL          url.URL
}

func NewImage(departmentID domain.DepartmentID, url url.URL) *Image {
	return &Image{
		ID:           uuid.New().String(),
		DepartmentID: departmentID,
		URL:          url,
	}
}
