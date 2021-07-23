package ope_handler

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
	auth adapter.AgencyAuthorization,
	regionRepo adapter.RegionRepo) pb.OpeRegionQuery {
	return &regionQuery{
		errorConverter: errorConverter,
		db:             db,
		auth:           auth,
		regionRepo:     regionRepo,
	}
}

type regionQuery struct {
	errorConverter adapter.ErrorConverter
	db             adapter.DB
	auth           adapter.AgencyAuthorization
	regionRepo     adapter.RegionRepo
}

func (h *regionQuery) MAreaListByPrefecture(ctx context.Context, req *pb.PrefectureID) (*pb.AreaList, error) {
	db := h.db(ctx)

	_, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	prefID := domain.PrefID(req.Id)

	regions, err := h.regionRepo.GetMAreaByLArea(ctx, db, domain.NewLAreaID(prefID))
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	resItems := make([]*pb.Area, 0, len(regions))
	for _, region := range regions {
		resItems = append(resItems, response.MAreaFrom(region, 0))
	}

	return &pb.AreaList{
		Items: resItems,
	}, nil
}

func (h *regionQuery) SAreaListByMArea(ctx context.Context, req *pb.MAreaID) (*pb.AreaList, error) {
	db := h.db(ctx)

	_, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	mAreaID := domain.MAreaID(req.Id)

	regions, err := h.regionRepo.GetSAreaByMArea(ctx, db, mAreaID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	resItems := make([]*pb.Area, 0, len(regions))
	for _, region := range regions {
		resItems = append(resItems, response.SAreaFrom(region, 0))
	}

	return &pb.AreaList{
		Items: resItems,
	}, nil
}
