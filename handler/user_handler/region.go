package user_handler

import (
	"context"

	"gke-go-sample/handler/response"

	"gke-go-sample/adapter"
	"gke-go-sample/domain"
	pb "gke-go-sample/proto/go/pb"
)

func NewRegionQuery(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	regionRepo adapter.RegionRepo,
	workRepo adapter.WorkRepo) pb.RegionQuery {
	return &regionQuery{
		errorConverter: errorConverter,
		db:             db,
		regionRepo:     regionRepo,
		workRepo:       workRepo,
	}
}

type regionQuery struct {
	errorConverter adapter.ErrorConverter
	db             adapter.DB
	regionRepo     adapter.RegionRepo
	workRepo       adapter.WorkRepo
}

func (h *regionQuery) MAreaListByPrefecture(ctx context.Context, req *pb.PrefectureID) (*pb.AreaList, error) {
	db := h.db(ctx)

	prefID := domain.PrefID(req.Id)

	regions, err := h.regionRepo.GetMAreaByLArea(ctx, db, domain.NewLAreaID(prefID))
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	resItems := make([]*pb.Area, 0, len(regions))
	for _, region := range regions {
		count, err := h.workRepo.GetCountByActiveAndMArea(ctx, db, region.ID)
		if err != nil {
			return nil, h.errorConverter(ctx, err)
		}

		resItems = append(resItems, response.MAreaFrom(region, count))
	}

	return &pb.AreaList{
		Items: resItems,
	}, nil
}

func (h *regionQuery) SAreaListByMArea(ctx context.Context, req *pb.MAreaID) (*pb.AreaList, error) {
	db := h.db(ctx)

	mAreaID := domain.MAreaID(req.Id)

	regions, err := h.regionRepo.GetSAreaByMArea(ctx, db, mAreaID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	resItems := make([]*pb.Area, 0, len(regions))
	for _, region := range regions {
		count, err := h.workRepo.GetCountByActiveAndSArea(ctx, db, region.ID)
		if err != nil {
			return nil, h.errorConverter(ctx, err)
		}

		resItems = append(resItems, response.SAreaFrom(region, count))
	}

	return &pb.AreaList{
		Items: resItems,
	}, nil
}
