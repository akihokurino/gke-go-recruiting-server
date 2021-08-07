package master_domain

import "gke-go-recruiting-server/domain"

type City struct {
	Prefecture   *Prefecture
	ID           domain.CityID
	Name         string
	NameKana     string
	AreaName     string
	AreaNameKana string
}

func NewCity(
	prefID domain.PrefID,
	prefName string,
	cityID domain.CityID,
	cityName string,
	cityNameKana string,
	areaName string,
	areaNameKana string) *City {
	return &City{
		Prefecture: &Prefecture{
			ID:   prefID,
			Name: prefName,
		},
		ID:           cityID,
		Name:         cityName,
		NameKana:     cityNameKana,
		AreaName:     areaName,
		AreaNameKana: areaNameKana,
	}
}
