package user_handler

import (
	"context"

	"gke-go-sample/handler/response"

	"gke-go-sample/adapter"
	"gke-go-sample/domain"
	pb "gke-go-sample/proto/go/pb"
)

func NewCityQuery(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	cityRepo adapter.CityRepo,
	workRepo adapter.WorkRepo) pb.CityQuery {
	return &cityQuery{
		errorConverter: errorConverter,
		db:             db,
		cityRepo:       cityRepo,
		workRepo:       workRepo,
	}
}

type cityQuery struct {
	errorConverter adapter.ErrorConverter
	db             adapter.DB
	cityRepo       adapter.CityRepo
	workRepo       adapter.WorkRepo
}

func (h *cityQuery) PrefectureList(ctx context.Context, req *pb.Empty) (*pb.PrefectureList, error) {
	db := h.db(ctx)

	prefectures, err := h.cityRepo.GetAllPrefecture(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	resItems := make([]*pb.Prefecture, 0, len(prefectures))
	for _, prefecture := range prefectures {
		count, err := h.workRepo.GetCountByActiveAndPref(ctx, db, prefecture.ID)
		if err != nil {
			return nil, h.errorConverter(ctx, err)
		}

		resItems = append(resItems, response.PrefectureFrom(prefecture, count))
	}

	return &pb.PrefectureList{
		Items: resItems,
	}, nil
}

func (h *cityQuery) CityListByPrefecture(ctx context.Context, req *pb.PrefectureID) (*pb.CityList, error) {
	db := h.db(ctx)

	prefID := domain.PrefID(req.Id)

	cities, err := h.cityRepo.GetByPrefecture(ctx, db, prefID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	resItems := make([]*pb.City, 0, len(cities))
	for _, city := range cities {
		count, err := h.workRepo.GetCountByActiveAndCity(ctx, db, city.ID)
		if err != nil {
			return nil, h.errorConverter(ctx, err)
		}

		resItems = append(resItems, response.CityFrom(city, count))
	}

	return &pb.CityList{
		Items: resItems,
	}, nil
}
