package product_domain

import (
	"time"

	"gke-go-sample/domain"
	pb "gke-go-sample/proto/go/pb"
)

type Main struct {
	Plan  pb.MainProduct_Plan
	Price uint64
}

func (m *Main) CalcDateRange(dateFrom time.Time) (domain.DateRange, error) {
	switch m.Plan {
	case pb.MainProduct_Plan_A:
		return domain.NewDateRange(dateFrom, dateFrom.Add(time.Duration(4*7*24*time.Hour))), nil
	case pb.MainProduct_Plan_B:
		return domain.NewDateRange(dateFrom, dateFrom.Add(time.Duration(8*7*24*time.Hour))), nil
	}

	return domain.DateRange{}, domain.NewInternalServerErr()
}

func GetMainList() []*Main {
	return []*Main{
		{
			Plan:  pb.MainProduct_Plan_A,
			Price: 10000,
		},
		{
			Plan:  pb.MainProduct_Plan_B,
			Price: 50000,
		},
	}
}
