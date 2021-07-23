package work_index

import (
	"gke-go-sample/domain"
	"gke-go-sample/domain/work_domain"
)

const newerIndexName = "work_newer"
const paymentHigherIndexName = "work_payment_higher"

type index struct {
	ObjectID          string   `json:"objectID"`
	Title             string   `json:"title"`
	Content           string   `json:"content"`
	Status            int32    `json:"status"`
	DepartmentName    string   `json:"departmentName"`
	BusinessCondition int32    `json:"businessCondition"`
	PrefID            string   `json:"prefId"`
	Address           string   `json:"address"`
	MAreaID           string   `json:"mAreaId"`
	SAreaID           string   `json:"sAreaId"`
	Merits            []int32  `json:"merits"`
	RailIDs           []string `json:"railIds"`
	StationIDs        []string `json:"stationIds"`
	DateFrom          int64    `json:"dateFrom"`
	DateTo            int64    `json:"dateTo"`
	PublishedOrder    int      `json:"publishedOrder"`
}

func (i *index) id() domain.WorkID {
	return domain.WorkID(i.ObjectID)
}

func entityFrom(d *work_domain.Work) *index {
	merits := make([]int32, 0, len(d.Merits))
	for _, merit := range d.Merits {
		merits = append(merits, int32(merit.Value))
	}

	railIDs := make([]string, 0, len(d.With.Department.Stations))
	stationIDs := make([]string, 0, len(d.With.Department.Stations))
	for _, station := range d.With.Department.Stations {
		railIDs = append(railIDs, station.With.Line.Rail.ID.String())
		stationIDs = append(stationIDs, station.With.Line.StationID.String())
	}

	publishedOrder := 0
	if d.With.ActivePlan != nil {
		publishedOrder = d.With.ActivePlan.PublishedOrder
	}

	return &index{
		ObjectID:          d.ID.String(),
		Title:             d.Title,
		Content:           d.Content,
		Status:            int32(d.Status),
		DepartmentName:    d.With.Department.Name,
		BusinessCondition: int32(d.With.Department.BusinessCondition),
		PrefID:            d.With.Department.PrefID.String(),
		Address:           d.With.Department.Address,
		MAreaID:           d.With.Department.Location.MAreaID.String(),
		SAreaID:           d.With.Department.Location.SAreaID.String(),
		Merits:            merits,
		RailIDs:           railIDs,
		StationIDs:        stationIDs,
		DateFrom:          d.DateRange.From.Unix(),
		DateTo:            d.DateRange.To.Unix(),
		PublishedOrder:    publishedOrder,
	}
}
