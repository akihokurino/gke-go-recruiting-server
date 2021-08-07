package master_domain

import "gke-go-recruiting-server/domain"

type Line struct {
	ID              domain.LineID
	Rail            Rail
	PrefID          domain.PrefID
	StationID       domain.StationID
	StopOrder       uint64
	StationName     string
	StationNameKana string
	Latitude        float64
	Longitude       float64
}

type Rail struct {
	ID        domain.RailID
	Name1     string
	NameKana1 string
	Name2     string
	NameKana2 string
	Company   RailCompany
}

type RailCompany struct {
	Name1     string
	NameKana1 string
	Name2     string
	Kind      uint64
	KindName  string
}
