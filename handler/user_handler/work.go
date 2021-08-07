package user_handler

import (
	"context"
	"gke-go-recruiting-server/handler/response"

	"gke-go-recruiting-server/domain"

	"gke-go-recruiting-server/adapter"
	pb "gke-go-recruiting-server/proto/go/pb"
)

func NewWorkQuery(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	workRepo adapter.WorkRepo,
	workIndexRepo adapter.WorkIndexRepo,
	departmentRepo adapter.DepartmentRepo,
	departmentImageRepo adapter.DepartmentImageRepo,
	cityRepo adapter.CityRepo,
) pb.WorkQuery {
	return &workQuery{
		errorConverter:      errorConverter,
		db:                  db,
		workRepo:            workRepo,
		workIndexRepo:       workIndexRepo,
		departmentRepo:      departmentRepo,
		departmentImageRepo: departmentImageRepo,
		cityRepo:            cityRepo,
	}
}

type workQuery struct {
	errorConverter      adapter.ErrorConverter
	db                  adapter.DB
	workRepo            adapter.WorkRepo
	workIndexRepo       adapter.WorkIndexRepo
	departmentRepo      adapter.DepartmentRepo
	departmentImageRepo adapter.DepartmentImageRepo
	cityRepo            adapter.CityRepo
}

func (h *workQuery) SearchEach(ctx context.Context, req *pb.SearchWorkParams) (*pb.WorkList, error) {
	db := h.db(ctx)

	now := domain.NowUTC()
	pager := domain.NewPager(req.Pager.Page, req.Pager.Offset)

	ids, err := h.workIndexRepo.SearchWithPager(
		ctx,
		req.Q,
		req.BusinessCondition,
		domain.PrefID(req.PrefId),
		domain.MAreaID(req.MAreaId),
		domain.SAreaID(req.SAreaId),
		domain.RailID(req.RailId),
		domain.StationID(req.StationId),
		req.Merit,
		pager,
		req.Order)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	works, err := h.workRepo.GetMulti(ctx, db, ids)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	prefs, err := h.cityRepo.GetAllPrefecture(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	resItems := make([]*pb.Work, 0, len(works))
	for _, work := range works {
		resItems = append(resItems, response.WorkFrom(work, work.IsNew(now), prefs))
	}

	return &pb.WorkList{
		Items: resItems,
	}, nil
}

func (h *workQuery) SearchCount(ctx context.Context, req *pb.SearchWorkCountParams) (*pb.Count, error) {
	count, err := h.workIndexRepo.SearchCount(
		ctx,
		req.Q,
		req.BusinessCondition,
		domain.PrefID(req.PrefId),
		domain.MAreaID(req.MAreaId),
		domain.SAreaID(req.SAreaId),
		domain.RailID(req.RailId),
		domain.StationID(req.StationId),
		req.Merit,
		req.Order)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Count{
		Count: count,
	}, nil
}

func (h *workQuery) GetMulti(ctx context.Context, req *pb.WorkIDList) (*pb.WorkList, error) {
	db := h.db(ctx)

	now := domain.NowUTC()

	workIDs := make([]domain.WorkID, 0, len(req.Ids))
	for _, id := range req.Ids {
		workIDs = append(workIDs, domain.WorkID(id))
	}

	works, err := h.workRepo.GetMulti(ctx, db, workIDs)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	prefs, err := h.cityRepo.GetAllPrefecture(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	resItems := make([]*pb.Work, 0, len(works))
	for _, work := range works {
		resItems = append(resItems, response.WorkFrom(work, work.IsNew(now), prefs))
	}

	return &pb.WorkList{
		Items: resItems,
	}, nil
}

func (h *workQuery) Get(ctx context.Context, req *pb.WorkID) (*pb.Work, error) {
	db := h.db(ctx)

	now := domain.NowUTC()
	workID := domain.WorkID(req.Id)

	work, err := h.workRepo.Get(ctx, db, workID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	prefs, err := h.cityRepo.GetAllPrefecture(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return response.WorkFrom(work, work.IsNew(now), prefs), nil
}
