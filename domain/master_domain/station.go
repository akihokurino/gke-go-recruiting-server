package master_domain

import "gke-go-sample/domain"

type Station struct {
	ID        domain.StationID
	PrefID    domain.PrefID
	Name      string
	Latitude  float64
	Longitude float64
}

func NewStation(
	id domain.StationID,
	prefID domain.PrefID,
	name string,
	latitude float64,
	longitude float64) *Station {
	return &Station{
		ID:        id,
		PrefID:    prefID,
		Name:      name,
		Latitude:  latitude,
		Longitude: longitude,
	}
}
