package ope_handler

import (
	"context"

	"gke-go-recruiting-server/handler/response"

	"gke-go-recruiting-server/adapter"
	"gke-go-recruiting-server/domain"
	pb "gke-go-recruiting-server/proto/go/pb"
)

func NewCityQuery(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	auth adapter.AgencyAuthorization,
	cityRepo adapter.CityRepo) pb.OpeCityQuery {
	return &cityQuery{
		errorConverter: errorConverter,
		db:             db,
		auth:           auth,
		cityRepo:       cityRepo,
	}
}

type cityQuery struct {
	errorConverter adapter.ErrorConverter
	db             adapter.DB
	auth           adapter.AgencyAuthorization
	cityRepo       adapter.CityRepo
}

func (h *cityQuery) PrefectureList(ctx context.Context, req *pb.Empty) (*pb.PrefectureList, error) {
	db := h.db(ctx)

	_, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	prefectures, err := h.cityRepo.GetAllPrefecture(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	resItems := make([]*pb.Prefecture, 0, len(prefectures))
	for _, prefecture := range prefectures {
		resItems = append(resItems, response.PrefectureFrom(prefecture, 0))
	}

	return &pb.PrefectureList{
		Items: resItems,
	}, nil
}

func (h *cityQuery) CityListByPrefecture(ctx context.Context, req *pb.PrefectureID) (*pb.CityList, error) {
	db := h.db(ctx)

	_, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	prefID := domain.PrefID(req.Id)

	cities, err := h.cityRepo.GetByPrefecture(ctx, db, prefID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	resItems := make([]*pb.City, 0, len(cities))
	for _, city := range cities {
		resItems = append(resItems, response.CityFrom(city, 0))
	}

	return &pb.CityList{
		Items: resItems,
	}, nil
}
