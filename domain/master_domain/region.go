package master_domain

import "gke-go-recruiting-server/domain"

type Region struct {
	Geocode1          string
	Geocode2          string
	Zipcode           string
	Address           string
	LArea             LArea
	OriginalMArea     string
	OriginalMAreaName string
	OriginalSArea     string
	OriginalSAreaName string
	MArea             MArea
	SArea             SArea
}

type LArea struct {
	ID   domain.LAreaID
	Name string
}

type MArea struct {
	ID   domain.MAreaID
	Name string
}

type SArea struct {
	ID   domain.SAreaID
	Name string
}
