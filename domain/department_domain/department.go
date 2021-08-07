package department_domain

import (
	"time"

	"github.com/google/uuid"

	"gke-go-recruiting-server/domain"

	pb "gke-go-recruiting-server/proto/go/pb"
)

type Department struct {
	ID                domain.DepartmentID
	AgencyID          domain.AgencyID
	CompanyID         domain.CompanyID
	SalesID           domain.FirebaseUserID
	Status            pb.Department_Status
	Name              string
	BusinessCondition pb.Department_BusinessCondition
	PostalCode        string
	PrefID            domain.PrefID
	CityID            domain.CityID
	Address           string
	Building          string
	PhoneNumber       string
	Location          Location
	domain.Meta

	Images   []*Image
	Stations []*Station
	Sales    *SalesOverview
}

func New(
	agencyID domain.AgencyID,
	companyID domain.CompanyID,
	salesID domain.FirebaseUserID,
	name string,
	businessCondition pb.Department_BusinessCondition,
	postalCode string,
	prefID domain.PrefID,
	cityID domain.CityID,
	address string,
	building string,
	phoneNumber string,
	location Location,
	now time.Time) *Department {
	return &Department{
		ID:                domain.DepartmentID(uuid.New().String()),
		AgencyID:          agencyID,
		CompanyID:         companyID,
		SalesID:           salesID,
		Status:            pb.Department_Status_REVIEW,
		Name:              name,
		BusinessCondition: businessCondition,
		PostalCode:        postalCode,
		PrefID:            prefID,
		CityID:            cityID,
		Address:           address,
		Building:          building,
		PhoneNumber:       phoneNumber,
		Location:          location,
		Meta: domain.Meta{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

func (d *Department) Update(
	name string,
	businessCondition pb.Department_BusinessCondition,
	postalCode string,
	prefID domain.PrefID,
	cityID domain.CityID,
	address string,
	building string,
	phoneNumber string,
	location Location,
	now time.Time) {
	d.Name = name
	d.BusinessCondition = businessCondition
	d.PostalCode = postalCode
	d.PrefID = prefID
	d.CityID = cityID
	d.Address = address
	d.Building = building
	d.PhoneNumber = phoneNumber
	d.Location = location
	d.UpdatedAt = now
}

func (d *Department) Accept() error {
	if d.Status != pb.Department_Status_REVIEW {
		return domain.NewConflictErr(domain.ConflictStatusMsg)
	}

	d.Status = pb.Department_Status_OK
	return nil
}

func (d *Department) Deny() error {
	if d.Status != pb.Department_Status_REVIEW {
		return domain.NewConflictErr(domain.ConflictStatusMsg)
	}

	d.Status = pb.Department_Status_NG
	return nil
}

func (d *Department) UpdateSales(id domain.FirebaseUserID) {
	d.SalesID = id
}
