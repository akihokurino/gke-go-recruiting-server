package master_handler

import (
	"context"
	"gke-go-recruiting-server/util"
	"net/url"

	"github.com/hashicorp/go-multierror"

	"gke-go-recruiting-server/adapter"
	"gke-go-recruiting-server/domain"
	"gke-go-recruiting-server/handler/response"
	pb "gke-go-recruiting-server/proto/go/pb"
)

func NewDepartmentQuery(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	auth adapter.AdminAuthorization,
	departmentRepo adapter.DepartmentRepo,
	cityRepo adapter.CityRepo) pb.MasterDepartmentQuery {
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
	auth           adapter.AdminAuthorization
	departmentRepo adapter.DepartmentRepo
	cityRepo       adapter.CityRepo
}

func (h *departmentQuery) ListByFilter(ctx context.Context, req *pb.MasterDepartmentFilterParams) (*pb.DepartmentList, error) {
	db := h.db(ctx)

	_, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	pager := domain.NewPager(req.Pager.Page, req.Pager.Offset)

	departments, err := h.departmentRepo.GetByFilterWithPager(ctx, db, pager, adapter.DepartmentFilterParams{
		AgencyID:       domain.AgencyID(req.AgencyId),
		CompanyID:      domain.CompanyID(req.CompanyId),
		DepartmentID:   domain.DepartmentID(req.DepartmentId),
		DepartmentName: req.DepartmentName,
		SalesID:        domain.FirebaseUserID(req.SalesId),
		Status:         req.Status,
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

func (h *departmentQuery) CountByFilter(ctx context.Context, req *pb.MasterDepartmentCountFilterParams) (*pb.Count, error) {
	db := h.db(ctx)

	_, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	count, err := h.departmentRepo.GetCountByFilter(ctx, db, adapter.DepartmentFilterParams{
		AgencyID:       domain.AgencyID(req.AgencyId),
		CompanyID:      domain.CompanyID(req.CompanyId),
		DepartmentID:   domain.DepartmentID(req.DepartmentId),
		DepartmentName: req.DepartmentName,
		SalesID:        domain.FirebaseUserID(req.SalesId),
		Status:         req.Status,
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

	_, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	departmentID := domain.DepartmentID(req.Id)

	department, err := h.departmentRepo.Get(ctx, db, departmentID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
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
	tx adapter.TX,
	auth adapter.AdminAuthorization,
	departmentApp adapter.DepartmentApp) pb.MasterDepartmentCommand {
	return &departmentCommand{
		errorConverter: errorConverter,
		db:             db,
		tx:             tx,
		auth:           auth,
		departmentApp:  departmentApp,
	}
}

type departmentCommand struct {
	errorConverter adapter.ErrorConverter
	db             adapter.DB
	tx             adapter.TX
	auth           adapter.AdminAuthorization
	departmentApp  adapter.DepartmentApp
}

func (h *departmentCommand) Accept(ctx context.Context, req *pb.DepartmentID) (*pb.Empty, error) {
	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	departmentID := domain.DepartmentID(req.Id)

	app := h.departmentApp.MasterBuild(me)

	if err := app.Accept(ctx, db, departmentID); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

func (h *departmentCommand) Deny(ctx context.Context, req *pb.DepartmentID) (*pb.Empty, error) {
	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	departmentID := domain.DepartmentID(req.Id)

	app := h.departmentApp.MasterBuild(me)

	if err := app.Deny(ctx, db, departmentID); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

func (h *departmentCommand) Update(ctx context.Context, req *pb.MasterDepartmentUpdateParams) (*pb.Empty, error) {
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

	app := h.departmentApp.MasterBuild(me)

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

func (h *departmentCommand) validateUpdate(req *pb.MasterDepartmentUpdateParams) error {
	var result *multierror.Error
	result = multierror.Append(result, util.ValidateTextRange(req.Name, 1, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.NameKana, 1, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.PostalCode, 1, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.Address, 1, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.Building, 0, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.PhoneNumber, 1, 255))
	if req.BusinessCondition == pb.Department_BusinessCondition_Unknown {
		result = multierror.Append(result, domain.NewBadRequestErr(domain.BadRequestMsg))
	}
	return result.ErrorOrNil()
}

func (h *departmentCommand) UpdateSales(ctx context.Context, req *pb.MasterDepartmentUpdateSalesParams) (*pb.Empty, error) {
	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	departmentID := domain.DepartmentID(req.DepartmentId)
	salesID := domain.FirebaseUserID(req.SalesId)

	app := h.departmentApp.MasterBuild(me)

	if err := app.UpdateSales(ctx, db, departmentID, salesID); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}
