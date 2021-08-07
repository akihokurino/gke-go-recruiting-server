package company_domain

import (
	"time"

	"github.com/google/uuid"

	"gke-go-recruiting-server/domain"
	pb "gke-go-recruiting-server/proto/go/pb"
)

type Company struct {
	ID          domain.CompanyID
	AgencyID    domain.AgencyID
	Status      pb.Company_Status
	RankType    pb.Company_RankType
	Rank        pb.Company_Rank
	Name        string
	NameKana    string
	PostalCode  string
	PrefID      domain.PrefID
	Address     string
	Building    string
	PhoneNumber string
	domain.Meta
}

func New(
	agencyID domain.AgencyID,
	rankType pb.Company_RankType,
	rank pb.Company_Rank,
	name string,
	nameKana string,
	postalCode string,
	prefID domain.PrefID,
	address string,
	building string,
	phoneNumber string,
	now time.Time) *Company {
	return &Company{
		ID:          domain.CompanyID(uuid.New().String()),
		AgencyID:    agencyID,
		Status:      pb.Company_Status_OK,
		RankType:    rankType,
		Rank:        rank,
		Name:        name,
		NameKana:    nameKana,
		PostalCode:  postalCode,
		PrefID:      prefID,
		Address:     address,
		Building:    building,
		PhoneNumber: phoneNumber,
		Meta: domain.Meta{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}
