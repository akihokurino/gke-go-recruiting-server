package domain

import "fmt"

type CompanyID string

func (id CompanyID) String() string {
	return string(id)
}

type DepartmentID string

func (id DepartmentID) String() string {
	return string(id)
}

type WorkID string

func (id WorkID) String() string {
	return string(id)
}

type CityID string

func (id CityID) String() string {
	return string(id)
}

type PrefID string

func (id PrefID) String() string {
	return string(id)
}

type RailID string

func (id RailID) String() string {
	return string(id)
}

type StationID string

func (id StationID) String() string {
	return string(id)
}

type MunicipalityID string

func (id MunicipalityID) String() string {
	return string(id)
}

type LAreaID string

func (id LAreaID) String() string {
	return string(id)
}

func NewLAreaID(prefID PrefID) LAreaID {
	return LAreaID(fmt.Sprintf("L-0%s", prefID))
}

type MAreaID string

func (id MAreaID) String() string {
	return string(id)
}

type SAreaID string

func (id SAreaID) String() string {
	return string(id)
}

type EntryID string

func (id EntryID) String() string {
	return string(id)
}

type AgencyID string

func (id AgencyID) String() string {
	return string(id)
}

type FirebaseUserID string

func (id FirebaseUserID) String() string {
	return string(id)
}

type MainContractID string

func (id MainContractID) String() string {
	return string(id)
}

type UsageStatementID string

func (id UsageStatementID) String() string {
	return string(id)
}

type LineID string

func (id LineID) String() string {
	return string(id)
}

func NewLineID(railID RailID, stationID StationID) LineID {
	return LineID(fmt.Sprintf("%s-%s", railID.String(), stationID.String()))
}
