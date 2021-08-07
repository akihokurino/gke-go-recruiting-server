package ope_handler

import (
	"context"
	"net/url"

	"gke-go-recruiting-server/util"

	"github.com/hashicorp/go-multierror"

	"gke-go-recruiting-server/adapter"
	"gke-go-recruiting-server/domain"
	"gke-go-recruiting-server/handler/response"
	pb "gke-go-recruiting-server/proto/go/pb"
)

func NewWorkQuery(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	auth adapter.AgencyAuthorization,
	workRepo adapter.WorkRepo,
	cityRepo adapter.CityRepo) pb.OpeWorkQuery {
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
	auth           adapter.AgencyAuthorization
	workRepo       adapter.WorkRepo
	cityRepo       adapter.CityRepo
}

func (h *workQuery) ListByFilter(ctx context.Context, req *pb.OpeWorkFilterParams) (*pb.WorkList, error) {
	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	now := domain.NowUTC()
	pager := domain.NewPager(req.Pager.Page, req.Pager.Offset)

	works, err := h.workRepo.GetByFilterWithPager(ctx, db, pager, adapter.WorkFilterParams{
		AgencyID:          me.AgencyID,
		WorkID:            domain.WorkID(req.WorkId),
		DepartmentID:      domain.DepartmentID(req.DepartmentId),
		DepartmentName:    req.DepartmentName,
		SalesID:           domain.FirebaseUserID(req.SalesId),
		BusinessCondition: req.BusinessCondition,
		WorkType:          req.WorkType,
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

func (h *workQuery) CountByFilter(ctx context.Context, req *pb.OpeWorkCountFilterParams) (*pb.Count, error) {
	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	count, err := h.workRepo.GetCountByFilter(ctx, db, adapter.WorkFilterParams{
		AgencyID:          me.AgencyID,
		WorkID:            domain.WorkID(req.WorkId),
		DepartmentID:      domain.DepartmentID(req.DepartmentId),
		DepartmentName:    req.DepartmentName,
		SalesID:           domain.FirebaseUserID(req.SalesId),
		BusinessCondition: req.BusinessCondition,
		WorkType:          req.WorkType,
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

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	now := domain.NowUTC()
	workID := domain.WorkID(req.Id)

	work, err := h.workRepo.Get(ctx, db, workID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	if work.With.Department.AgencyID != me.AgencyID {
		return nil, h.errorConverter(ctx, domain.NewForbiddenErr(domain.ForbiddenAgencyMsg))
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
	auth adapter.AgencyAuthorization,
	workApp adapter.WorkApp) pb.OpeWorkCommand {
	return &workCommand{
		errorConverter: errorConverter,
		db:             db,
		auth:           auth,
		workApp:        workApp,
	}
}

type workCommand struct {
	errorConverter adapter.ErrorConverter
	db             adapter.DB
	auth           adapter.AgencyAuthorization
	workApp        adapter.WorkApp
}

func (h *workCommand) Create(ctx context.Context, req *pb.OpeWorkCreateParams) (*pb.WorkID, error) {
	if err := h.validateCreate(req); err != nil {
		return nil, h.errorConverter(ctx, domain.NewBadRequestErr(domain.BadRequestMsg))
	}

	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	dateFrom, err := domain.UTCFrom(req.DateFrom)
	if err != nil {
		return nil, h.errorConverter(ctx, domain.NewBadRequestErr(domain.BadRequestDateMsg))
	}

	dateTo, err := domain.UTCFrom(req.DateTo)
	if err != nil {
		return nil, h.errorConverter(ctx, domain.NewBadRequestErr(domain.BadRequestDateMsg))
	}

	images := make([]adapter.WorkImageParams, 0, len(req.Images))
	for _, image := range req.Images {
		u, err := url.Parse(image.Url)
		if err != nil {
			continue
		}
		images = append(images, adapter.WorkImageParams{
			URL:       *u,
			ViewOrder: image.ViewOrder,
			Comment:   image.Comment,
		})
	}

	movieURLs := make([]url.URL, 0, len(req.MovieUrls))
	for _, urlString := range req.MovieUrls {
		u, err := url.Parse(urlString)
		if err != nil {
			continue
		}
		movieURLs = append(movieURLs, *u)
	}

	workParams := adapter.WorkParams{
		WorkType: req.WorkType,
		JobCode:  req.JobCode,
		Title:    req.Title,
		Content:  req.Content,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	departmentID := domain.DepartmentID(req.DepartmentId)

	now := domain.NowUTC()

	app := h.workApp.OperationBuild(me)

	work, err := app.Create(
		ctx,
		db,
		departmentID,
		workParams,
		images,
		movieURLs,
		req.Merits,
		now)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.WorkID{
		Id: work.ID.String(),
	}, nil
}

func (h *workCommand) validateCreate(req *pb.OpeWorkCreateParams) error {
	var result *multierror.Error
	result = multierror.Append(result, util.ValidateTextRange(req.Title, 1, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.Content, 1, 1000))
	if req.WorkType == pb.Work_Type_Unknown {
		result = multierror.Append(result, domain.NewBadRequestErr(domain.BadRequestMsg))
	}
	if req.JobCode == pb.Work_Job_Unknown {
		result = multierror.Append(result, domain.NewBadRequestErr(domain.BadRequestMsg))
	}
	return result.ErrorOrNil()
}

func (h *workCommand) Update(ctx context.Context, req *pb.OpeWorkUpdateParams) (*pb.Empty, error) {
	if err := h.validateUpdate(req); err != nil {
		return nil, h.errorConverter(ctx, domain.NewBadRequestErr(domain.BadRequestMsg))
	}

	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	dateFrom, err := domain.UTCFrom(req.DateFrom)
	if err != nil {
		return nil, h.errorConverter(ctx, domain.NewBadRequestErr(domain.BadRequestDateMsg))
	}

	dateTo, err := domain.UTCFrom(req.DateTo)
	if err != nil {
		return nil, h.errorConverter(ctx, domain.NewBadRequestErr(domain.BadRequestDateMsg))
	}

	images := make([]adapter.WorkImageParams, 0, len(req.Images))
	for _, image := range req.Images {
		u, err := url.Parse(image.Url)
		if err != nil {
			continue
		}
		images = append(images, adapter.WorkImageParams{
			URL:       *u,
			ViewOrder: image.ViewOrder,
			Comment:   image.Comment,
		})
	}

	movieURLs := make([]url.URL, 0, len(req.MovieUrls))
	for _, urlString := range req.MovieUrls {
		u, err := url.Parse(urlString)
		if err != nil {
			continue
		}
		movieURLs = append(movieURLs, *u)
	}

	workParams := adapter.WorkParams{
		WorkType: req.WorkType,
		JobCode:  req.JobCode,
		Title:    req.Title,
		Content:  req.Content,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	workID := domain.WorkID(req.WorkId)
	now := domain.NowUTC()

	app := h.workApp.OperationBuild(me)

	if _, err := app.Update(
		ctx,
		db,
		workID,
		workParams,
		images,
		movieURLs,
		req.Merits,
		now); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

func (h *workCommand) validateUpdate(req *pb.OpeWorkUpdateParams) error {
	var result *multierror.Error
	result = multierror.Append(result, util.ValidateTextRange(req.Title, 1, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.Content, 1, 1000))
	if req.WorkType == pb.Work_Type_Unknown {
		result = multierror.Append(result, domain.NewBadRequestErr(domain.BadRequestMsg))
	}
	if req.JobCode == pb.Work_Job_Unknown {
		result = multierror.Append(result, domain.NewBadRequestErr(domain.BadRequestMsg))
	}
	return result.ErrorOrNil()
}

func (h *workCommand) EarlyFinish(ctx context.Context, req *pb.WorkID) (*pb.Empty, error) {
	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	workID := domain.WorkID(req.Id)

	app := h.workApp.OperationBuild(me)

	if err := app.EarlyFinish(ctx, db, workID); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}
