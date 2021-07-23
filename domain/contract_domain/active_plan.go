package contract_domain

import (
	"gke-go-sample/domain"
	pb "gke-go-sample/proto/go/pb"
)

type ActivePlan struct {
	WorkID         domain.WorkID
	MainContractID domain.MainContractID
	PublishedOrder int

	MainContract *Main
}

func (a *ActivePlan) ReOrder() {
	a.PublishedOrder = PublishedOrderFrom(a.MainContract.Plan)
}

func PublishedOrderFrom(plan pb.MainProduct_Plan) int {
	return int(plan)
}
