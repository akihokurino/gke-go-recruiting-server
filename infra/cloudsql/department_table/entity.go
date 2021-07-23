package department_table

import (
	"database/sql"
	"time"

	"gke-go-sample/infra/cloudsql/agency_account_table"

	"gke-go-sample/domain/contract_domain"

	"gke-go-sample/domain/account_domain"

	"gke-go-sample/infra/cloudsql/department_station_table"

	"gke-go-sample/domain/department_domain"

	"gke-go-sample/infra/cloudsql/department_image_table"

	"gke-go-sample/domain"
	pb "gke-go-sample/proto/go/pb"

	"github.com/guregu/null"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

func (e *Entity) TableName() string {
	return "departments"
}

type Entity struct {
	ID                string         `gorm:"column:id;primary_key"`
	AgencyID          string         `gorm:"column:agency_id"`
	CompanyID         string         `gorm:"column:company_id"`
	SalesID           string         `gorm:"column:sales_id"`
	Status            int32          `gorm:"column:status"`
	Name              string         `gorm:"column:name"`
	BusinessCondition int32          `gorm:"column:business_condition"`
	PostalCode        string         `gorm:"column:postal_code"`
	PrefID            string         `gorm:"column:pref_id"`
	CityID            sql.NullString `gorm:"column:city_id"`
	Address           string         `gorm:"column:address"`
	Building          string         `gorm:"column:building"`
	PhoneNumber       string         `gorm:"column:phone_number"`
	MAreaID           sql.NullString `gorm:"column:m_area_id"`
	SAreaID           sql.NullString `gorm:"column:s_area_id"`
	Latitude          float64        `gorm:"column:latitude"`
	Longitude         float64        `gorm:"column:longitude"`
	CreatedAt         time.Time      `gorm:"column:created_at"`
	UpdatedAt         time.Time      `gorm:"column:updated_at"`

	Images        []*department_image_table.Entity   `gorm:"PRELOAD:false;foreignkey:department_id"`
	Stations      []*department_station_table.Entity `gorm:"PRELOAD:false;foreignkey:department_id"`
	AgencyAccount *agency_account_table.Entity       `gorm:"PRELOAD:false;foreignkey:sales_id"`
}

func (e *Entity) ToDomain() *department_domain.Department {
	images := make([]*department_domain.Image, 0, len(e.Images))
	for _, image := range e.Images {
		images = append(images, image.ToDomain())
	}

	stations := make([]*department_domain.Station, 0, len(e.Stations))
	for _, station := range e.Stations {
		d, err := station.ToDomain()
		if err != nil {
			continue
		}
		stations = append(stations, d)
	}

	sales := &department_domain.SalesOverview{
		ID:   "",
		Name: "",
	}
	if e.AgencyAccount != nil {
		sales = e.AgencyAccount.ToDepartmentDomainOverview()
	}

	return &department_domain.Department{
		ID:                domain.DepartmentID(e.ID),
		AgencyID:          domain.AgencyID(e.AgencyID),
		CompanyID:         domain.CompanyID(e.CompanyID),
		SalesID:           domain.FirebaseUserID(e.SalesID),
		Status:            pb.Department_Status(e.Status),
		Name:              e.Name,
		BusinessCondition: pb.Department_BusinessCondition(e.BusinessCondition),
		PostalCode:        e.PostalCode,
		PrefID:            domain.PrefID(e.PrefID),
		CityID:            domain.CityID(e.CityID.String),
		Address:           e.Address,
		Building:          e.Building,
		PhoneNumber:       e.PhoneNumber,
		Location: department_domain.Location{
			MAreaID:   domain.MAreaID(e.MAreaID.String),
			SAreaID:   domain.SAreaID(e.SAreaID.String),
			Latitude:  e.Latitude,
			Longitude: e.Longitude,
		},
		Meta: domain.Meta{
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		},

		Images:   images,
		Stations: stations,
		Sales:    sales,
	}
}

func (e *Entity) ToAccountDomainOverview() *account_domain.DepartmentOverview {
	return &account_domain.DepartmentOverview{
		ID:       domain.DepartmentID(e.ID),
		AgencyID: domain.AgencyID(e.AgencyID),
		Name:     e.Name,
	}
}

func (e *Entity) ToContractDomainOverview() *contract_domain.DepartmentOverview {
	return &contract_domain.DepartmentOverview{
		ID:       domain.DepartmentID(e.ID),
		AgencyID: domain.AgencyID(e.AgencyID),
	}
}

func entityFrom(d *department_domain.Department) *Entity {
	return &Entity{
		ID:                d.ID.String(),
		AgencyID:          d.AgencyID.String(),
		CompanyID:         d.CompanyID.String(),
		SalesID:           d.SalesID.String(),
		Status:            int32(d.Status),
		Name:              d.Name,
		BusinessCondition: int32(d.BusinessCondition),
		PostalCode:        d.PostalCode,
		PrefID:            d.PrefID.String(),
		CityID: sql.NullString{
			String: d.CityID.String(),
			Valid:  d.CityID.String() != "",
		},
		Address:     d.Address,
		Building:    d.Building,
		PhoneNumber: d.PhoneNumber,
		MAreaID: sql.NullString{
			String: d.Location.MAreaID.String(),
			Valid:  d.Location.MAreaID.String() != "",
		},
		SAreaID: sql.NullString{
			String: d.Location.SAreaID.String(),
			Valid:  d.Location.SAreaID.String() != "",
		},
		Latitude:  d.Location.Latitude,
		Longitude: d.Location.Longitude,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}
