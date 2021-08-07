package ope_handler

import (
	"context"
	"gke-go-recruiting-server/util"
	"net/url"

	"github.com/hashicorp/go-multierror"

	pb "gke-go-recruiting-server/proto/go/pb"

	"gke-go-recruiting-server/adapter"
	"gke-go-recruiting-server/domain"
	"gke-go-recruiting-server/handler/response"
)

func NewDepartmentQuery(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	auth adapter.AgencyAuthorization,
	departmentRepo adapter.DepartmentRepo,
	cityRepo adapter.CityRepo) pb.OpeDepartmentQuery {
	return &departmentQuery{
		errorConverter: errorConverter,
		db:             db,
		auth:           auth,
		departmentRepo: departmentRepo,
		cityRepo:       cityRepo,
	}
}

type departmentQuery struct {
	errorConverter adapter.ErrorConverter
	db             adapter.DB
	auth           adapter.AgencyAuthorization
	departmentRepo adapter.DepartmentRepo
	cityRepo       adapter.CityRepo
}

func (h *departmentQuery) ListByFilter(ctx context.Context, req *pb.OpeDepartmentFilterParams) (*pb.DepartmentList, error) {
	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	pager := domain.NewPager(req.Pager.Page, req.Pager.Offset)

	departments, err := h.departmentRepo.GetByFilterWithPager(ctx, db, pager, adapter.DepartmentFilterParams{
		AgencyID:       me.AgencyID,
		CompanyID:      domain.CompanyID(req.CompanyId),
		DepartmentID:   domain.DepartmentID(req.DepartmentId),
		DepartmentName: req.DepartmentName,
		SalesID:        domain.FirebaseUserID(req.SalesId),
		Status:         req.Status,
		PhoneNumber:    req.PhoneNumber,
	})
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	prefs, err := h.cityRepo.GetAllPrefecture(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	resItems := make([]*pb.Department, 0, len(departments))
	for _, department := range departments {
		resItems = append(resItems, response.DepartmentFrom(department, prefs))
	}

	return &pb.DepartmentList{
		Items: resItems,
	}, nil
}

func (h *departmentQuery) CountByFilter(ctx context.Context, req *pb.OpeDepartmentCountFilterParams) (*pb.Count, error) {
	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	count, err := h.departmentRepo.GetCountByFilter(ctx, db, adapter.DepartmentFilterParams{
		AgencyID:       me.AgencyID,
		CompanyID:      domain.CompanyID(req.CompanyId),
		DepartmentID:   domain.DepartmentID(req.DepartmentId),
		DepartmentName: req.DepartmentName,
		SalesID:        domain.FirebaseUserID(req.SalesId),
		Status:         req.Status,
		PhoneNumber:    req.PhoneNumber,
	})
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Count{
		Count: count,
	}, nil
}

func (h *departmentQuery) Get(ctx context.Context, req *pb.DepartmentID) (*pb.Department, error) {
	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	departmentID := domain.DepartmentID(req.Id)

	department, err := h.departmentRepo.Get(ctx, db, departmentID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	if department.AgencyID != me.AgencyID {
		return nil, h.errorConverter(ctx, domain.NewForbiddenErr(domain.ForbiddenAgencyMsg))
	}

	prefs, err := h.cityRepo.GetAllPrefecture(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return response.DepartmentFrom(department, prefs), nil
}

func NewDepartmentCommand(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	auth adapter.AgencyAuthorization,
	departmentApp adapter.DepartmentApp) pb.OpeDepartmentCommand {
	return &departmentCommand{
		errorConverter: errorConverter,
		db:             db,
		auth:           auth,
		departmentApp:  departmentApp,
	}
}

type departmentCommand struct {
	errorConverter adapter.ErrorConverter
	db             adapter.DB
	auth           adapter.AgencyAuthorization
	departmentApp  adapter.DepartmentApp
}

func (h *departmentCommand) Create(ctx context.Context, req *pb.OpeDepartmentCreateParams) (*pb.DepartmentID, error) {
	if err := h.validateCreate(req); err != nil {
		return nil, h.errorConverter(ctx, domain.NewBadRequestErr(domain.BadRequestMsg))
	}

	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	departmentParams := adapter.DepartmentParams{
		Name:              req.Name,
		BusinessCondition: req.BusinessCondition,
		PostalCode:        req.PostalCode,
		PrefID:            domain.PrefID(req.PrefId),
		CityID:            domain.CityID(req.CityId),
		Address:           req.Address,
		Building:          req.Building,
		PhoneNumber:       req.PhoneNumber,
		MAreaID:           domain.MAreaID(req.MAreaId),
		SAreaID:           domain.SAreaID(req.SAreaId),
		Latitude:          req.Latitude,
		Longitude:         req.Longitude,
	}

	imageURLs := make([]url.URL, 0, len(req.ImageUrls))
	for _, urlString := range req.ImageUrls {
		u, err := url.Parse(urlString)
		if err != nil {
			continue
		}
		imageURLs = append(imageURLs, *u)
	}

	lineIDs := make([]domain.LineID, 0, len(req.LineIds))
	for _, id := range req.LineIds {
		lineIDs = append(lineIDs, domain.LineID(id))
	}

	companyID := domain.CompanyID(req.CompanyId)
	salesID := domain.FirebaseUserID(req.SalesId)

	now := domain.NowUTC()

	app := h.departmentApp.OperationBuild(me)

	department, err := app.Create(
		ctx,
		db,
		companyID,
		salesID,
		departmentParams,
		imageURLs,
		lineIDs,
		now)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.DepartmentID{
		Id: department.ID.String(),
	}, nil
}

func (h *departmentCommand) validateCreate(req *pb.OpeDepartmentCreateParams) error {
	var result *multierror.Error
	result = multierror.Append(result, util.ValidateTextRange(req.Name, 1, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.PostalCode, 1, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.Address, 1, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.Building, 0, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.PhoneNumber, 1, 255))
	if req.BusinessCondition == pb.Department_BusinessCondition_Unknown {
		result = multierror.Append(result, domain.NewBadRequestErr(domain.BadRequestMsg))
	}
	return result.ErrorOrNil()
}

func (h *departmentCommand) Update(ctx context.Context, req *pb.OpeDepartmentUpdateParams) (*pb.Empty, error) {
	if err := h.validateUpdate(req); err != nil {
		return nil, h.errorConverter(ctx, domain.NewBadRequestErr(domain.BadRequestMsg))
	}

	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	departmentID := domain.DepartmentID(req.Id)

	departmentParams := adapter.DepartmentParams{
		Name:              req.Name,
		BusinessCondition: req.BusinessCondition,
		PostalCode:        req.PostalCode,
		PrefID:            domain.PrefID(req.PrefId),
		CityID:            domain.CityID(req.CityId),
		Address:           req.Address,
		Building:          req.Building,
		PhoneNumber:       req.PhoneNumber,
		MAreaID:           domain.MAreaID(req.MAreaId),
		SAreaID:           domain.SAreaID(req.SAreaId),
		Latitude:          req.Latitude,
		Longitude:         req.Longitude,
	}

	imageURLs := make([]url.URL, 0, len(req.ImageUrls))
	for _, urlString := range req.ImageUrls {
		u, err := url.Parse(urlString)
		if err != nil {
			continue
		}
		imageURLs = append(imageURLs, *u)
	}

	lineIDs := make([]domain.LineID, 0, len(req.LineIds))
	for _, id := range req.LineIds {
		lineIDs = append(lineIDs, domain.LineID(id))
	}

	now := domain.NowUTC()

	app := h.departmentApp.OperationBuild(me)

	if _, err := app.Update(
		ctx,
		db,
		departmentID,
		departmentParams,
		imageURLs,
		lineIDs,
		now); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

func (h *departmentCommand) validateUpdate(req *pb.OpeDepartmentUpdateParams) error {
	var result *multierror.Error
	result = multierror.Append(result, util.ValidateTextRange(req.Name, 1, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.PostalCode, 1, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.Address, 1, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.Building, 0, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.PhoneNumber, 1, 255))
	if req.BusinessCondition == pb.Department_BusinessCondition_Unknown {
		result = multierror.Append(result, domain.NewBadRequestErr(domain.BadRequestMsg))
	}
	return result.ErrorOrNil()
}

func (h *departmentCommand) UpdateSales(ctx context.Context, req *pb.OpeDepartmentUpdateSalesParams) (*pb.Empty, error) {
	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	departmentID := domain.DepartmentID(req.DepartmentId)
	salesID := domain.FirebaseUserID(req.SalesId)

	app := h.departmentApp.OperationBuild(me)

	if err := app.UpdateSales(ctx, db, departmentID, salesID); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}
