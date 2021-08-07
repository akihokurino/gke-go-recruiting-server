package department_domain

import "gke-go-recruiting-server/domain"

type Location struct {
	MAreaID   domain.MAreaID
	SAreaID   domain.SAreaID
	Latitude  float64
	Longitude float64
}

func NewLocation(
	mAreaID domain.MAreaID,
	sAreaID domain.SAreaID,
	latitude float64,
	longitude float64) *Location {
	return &Location{
		MAreaID:   mAreaID,
		SAreaID:   sAreaID,
		Latitude:  latitude,
		Longitude: longitude,
	}
}
