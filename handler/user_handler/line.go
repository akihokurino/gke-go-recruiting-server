package user_handler

import (
	"context"

	"gke-go-sample/handler/response"

	"gke-go-sample/adapter"
	"gke-go-sample/domain"
	pb "gke-go-sample/proto/go/pb"
)

func NewLineQuery(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	lineRepo adapter.LineRepo,
	workRepo adapter.WorkRepo,
	departmentStationRepo adapter.DepartmentStationRepo) pb.LineQuery {
	return &lineQuery{
		errorConverter:        errorConverter,
		db:                    db,
		lineRepo:              lineRepo,
		workRepo:              workRepo,
		departmentStationRepo: departmentStationRepo,
	}
}

type lineQuery struct {
	errorConverter        adapter.ErrorConverter
	db                    adapter.DB
	lineRepo              adapter.LineRepo
	workRepo              adapter.WorkRepo
	departmentStationRepo adapter.DepartmentStationRepo
}

func (h *lineQuery) RailListByPrefecture(ctx context.Context, req *pb.PrefectureID) (*pb.RailList, error) {
	db := h.db(ctx)

	prefID := domain.PrefID(req.Id)

	rails, err := h.lineRepo.GetRailByPrefecture(ctx, db, prefID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	resItems := make([]*pb.Rail, 0, len(rails))
	for _, rail := range rails {
		departmentStations, err := h.departmentStationRepo.GetByRail(ctx, db, rail.ID)
		if err != nil {
			return nil, h.errorConverter(ctx, err)
		}

		departmentIDs := make([]domain.DepartmentID, 0, len(departmentStations))
		for _, station := range departmentStations {
			departmentIDs = append(departmentIDs, station.DepartmentID)
		}

		count, err := h.workRepo.GetCountByActiveAndDepartments(ctx, db, departmentIDs)
		if err != nil {
			return nil, h.errorConverter(ctx, err)
		}

		resItems = append(resItems, response.RailFrom(rail, count))
	}

	return &pb.RailList{
		Items: resItems,
	}, nil
}

func (h *lineQuery) RailListByCompany(ctx context.Context, req *pb.RailListByCompanyParams) (*pb.RailList, error) {
	db := h.db(ctx)

	companyName := req.CompanyName

	rails, err := h.lineRepo.GetRailByCompany(ctx, db, companyName)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	resItems := make([]*pb.Rail, 0, len(rails))
	for _, rail := range rails {
		departmentStations, err := h.departmentStationRepo.GetByRail(ctx, db, rail.ID)
		if err != nil {
			return nil, h.errorConverter(ctx, err)
		}

		departmentIDs := make([]domain.DepartmentID, 0, len(departmentStations))
		for _, station := range departmentStations {
			departmentIDs = append(departmentIDs, station.DepartmentID)
		}

		count, err := h.workRepo.GetCountByActiveAndDepartments(ctx, db, departmentIDs)
		if err != nil {
			return nil, h.errorConverter(ctx, err)
		}

		resItems = append(resItems, response.RailFrom(rail, count))
	}

	return &pb.RailList{
		Items: resItems,
	}, nil
}

func (h *lineQuery) StationListByRail(ctx context.Context, req *pb.RailID) (*pb.StationList, error) {
	db := h.db(ctx)

	railID := domain.RailID(req.Id)

	stations, err := h.lineRepo.GetByRail(ctx, db, railID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	resItems := make([]*pb.Station, 0, len(stations))
	for _, station := range stations {
		departmentStations, err := h.departmentStationRepo.GetByStation(ctx, db, station.StationID)
		if err != nil {
			return nil, h.errorConverter(ctx, err)
		}

		departmentIDs := make([]domain.DepartmentID, 0, len(stations))
		for _, station := range departmentStations {
			departmentIDs = append(departmentIDs, station.DepartmentID)
		}

		count, err := h.workRepo.GetCountByActiveAndDepartments(ctx, db, departmentIDs)
		if err != nil {
			return nil, h.errorConverter(ctx, err)
		}

		resItems = append(resItems, response.StationFrom(station, count))
	}

	return &pb.StationList{
		Items: resItems,
	}, nil
}

func (h *lineQuery) RailCompanyListByPrefecture(ctx context.Context, req *pb.PrefectureID) (*pb.RailCompanyList, error) {
	db := h.db(ctx)

	prefID := domain.PrefID(req.Id)

	companies, err := h.lineRepo.GetRailCompanyByPrefecture(ctx, db, prefID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	resItems := make([]*pb.RailCompany, 0, len(companies))
	for _, line := range companies {
		resItems = append(resItems, response.RailCompanyFrom(line))
	}

	return &pb.RailCompanyList{
		Items: resItems,
	}, nil
}
