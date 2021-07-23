package adapter

import (
	"context"
	"gke-go-sample/domain"

	pb "gke-go-sample/proto/go/pb"

	"gke-go-sample/domain/work_domain"
)

type WorkIndexSortSetting struct {
}

type WorkIndexRepo interface {
	SearchWithPager(
		ctx context.Context,
		q string,
		businessCondition pb.Department_BusinessCondition,
		prefID domain.PrefID,
		mAreaID domain.MAreaID,
		sAreaID domain.SAreaID,
		railID domain.RailID,
		stationID domain.StationID,
		merit pb.Work_Merit,
		pager *domain.Pager,
		order pb.SearchWorkOrder) ([]domain.WorkID, error)
	SearchCount(
		ctx context.Context,
		q string,
		businessCondition pb.Department_BusinessCondition,
		prefID domain.PrefID,
		mAreaID domain.MAreaID,
		sAreaID domain.SAreaID,
		railID domain.RailID,
		stationID domain.StationID,
		merit pb.Work_Merit,
		order pb.SearchWorkOrder) (uint64, error)
	Save(ctx context.Context, work *work_domain.Work) error
	SaveMulti(ctx context.Context, works []*work_domain.Work) error
	Delete(ctx context.Context, workID domain.WorkID) error
	DeleteMulti(ctx context.Context, workIDs []domain.WorkID) error
}
