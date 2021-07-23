package master_handler

import (
	"context"

	"gke-go-sample/adapter"
	"gke-go-sample/domain"
	"gke-go-sample/handler/response"
	pb "gke-go-sample/proto/go/pb"
)

func NewWorkQuery(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	auth adapter.AdminAuthorization,
	workRepo adapter.WorkRepo,
	cityRepo adapter.CityRepo) pb.MasterWorkQuery {
	return &workQuery{
		errorConverter: errorConverter,
		db:             db,
		auth:           auth,
		workRepo:       workRepo,
		cityRepo:       cityRepo,
	}
}

type workQuery struct {
	errorConverter adapter.ErrorConverter
	db             adapter.DB
	auth           adapter.AdminAuthorization
	workRepo       adapter.WorkRepo
	cityRepo       adapter.CityRepo
}

func (h *workQuery) ListByFilter(ctx context.Context, req *pb.MasterWorkFilterParams) (*pb.WorkList, error) {
	db := h.db(ctx)

	_, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	now := domain.NowUTC()
	pager := domain.NewPager(req.Pager.Page, req.Pager.Offset)

	var dateFromRange *domain.DateRange
	if req.DateFromRange != nil {
		dateFromRange = domain.NewDateRangeFromString(req.DateFromRange.From, req.DateFromRange.To)
	}

	var dateToRange *domain.DateRange
	if req.DateToRange != nil {
		dateToRange = domain.NewDateRangeFromString(req.DateToRange.From, req.DateToRange.To)
	}

	works, err := h.workRepo.GetByFilterWithPager(ctx, db, pager, adapter.WorkFilterParams{
		AgencyID:          domain.AgencyID(req.AgencyId),
		CompanyID:         domain.CompanyID(req.CompanyId),
		DepartmentID:      domain.DepartmentID(req.DepartmentId),
		DepartmentName:    req.DepartmentName,
		WorkID:            domain.WorkID(req.WorkId),
		SalesID:           domain.FirebaseUserID(req.SalesId),
		BusinessCondition: req.BusinessCondition,
		WorkType:          req.WorkType,
		DateFromRange:     dateFromRange,
		DateToRange:       dateToRange,
		PrefID:            domain.PrefID(req.PrefId),
		Status:            req.Status,
	})
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

func (h *workQuery) CountByFilter(ctx context.Context, req *pb.MasterWorkCountFilterParams) (*pb.Count, error) {
	db := h.db(ctx)

	_, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	var dateFromRange *domain.DateRange
	if req.DateFromRange != nil {
		dateFromRange = domain.NewDateRangeFromString(req.DateFromRange.From, req.DateFromRange.To)
	}

	var dateToRange *domain.DateRange
	if req.DateToRange != nil {
		dateToRange = domain.NewDateRangeFromString(req.DateToRange.From, req.DateToRange.To)
	}

	count, err := h.workRepo.GetCountByFilter(ctx, db, adapter.WorkFilterParams{
		AgencyID:          domain.AgencyID(req.AgencyId),
		CompanyID:         domain.CompanyID(req.CompanyId),
		DepartmentID:      domain.DepartmentID(req.DepartmentId),
		DepartmentName:    req.DepartmentName,
		WorkID:            domain.WorkID(req.WorkId),
		SalesID:           domain.FirebaseUserID(req.SalesId),
		BusinessCondition: req.BusinessCondition,
		WorkType:          req.WorkType,
		DateFromRange:     dateFromRange,
		DateToRange:       dateToRange,
		PrefID:            domain.PrefID(req.PrefId),
		Status:            req.Status,
	})
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Count{
		Count: count,
	}, nil
}

func (h *workQuery) Get(ctx context.Context, req *pb.WorkID) (*pb.Work, error) {
	db := h.db(ctx)

	_, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

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

func NewWorkCommand(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	tx adapter.TX,
	auth adapter.AdminAuthorization,
	workApp adapter.WorkApp) pb.MasterWorkCommand {
	return &workCommand{
		errorConverter: errorConverter,
		db:             db,
		tx:             tx,
		auth:           auth,
		workApp:        workApp,
	}
}

type workCommand struct {
	errorConverter adapter.ErrorConverter
	db             adapter.DB
	tx             adapter.TX
	auth           adapter.AdminAuthorization
	workApp        adapter.WorkApp
}

func (h *workCommand) Accept(ctx context.Context, req *pb.WorkID) (*pb.Empty, error) {
	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	workID := domain.WorkID(req.Id)

	app := h.workApp.MasterBuild(me)

	if err := app.Accept(ctx, db, workID); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

func (h *workCommand) Deny(ctx context.Context, req *pb.WorkID) (*pb.Empty, error) {
	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	workID := domain.WorkID(req.Id)

	app := h.workApp.MasterBuild(me)

	if err := app.Deny(ctx, db, workID); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}
