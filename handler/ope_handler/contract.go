package ope_handler

import (
	"context"

	"github.com/hashicorp/go-multierror"

	"gke-go-recruiting-server/domain"

	"gke-go-recruiting-server/adapter"
	"gke-go-recruiting-server/handler/response"
	pb "gke-go-recruiting-server/proto/go/pb"
)

func NewContractQuery(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	auth adapter.AgencyAuthorization,
	mainContractRepo adapter.MainContractRepo) pb.OpeContractQuery {
	return &contractQuery{
		errorConverter:   errorConverter,
		db:               db,
		auth:             auth,
		mainContractRepo: mainContractRepo,
	}
}

type contractQuery struct {
	errorConverter   adapter.ErrorConverter
	db               adapter.DB
	auth             adapter.AgencyAuthorization
	mainContractRepo adapter.MainContractRepo
}

func (h *contractQuery) ListMainByFilter(ctx context.Context, req *pb.OpeMainContractFilterParams) (*pb.MainContractList, error) {
	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	pager := domain.NewPager(req.Pager.Page, req.Pager.Offset)

	contracts, err := h.mainContractRepo.GetByFilterWithPager(ctx, db, pager, adapter.MainContractFilterParams{
		AgencyID:       me.AgencyID,
		CompanyID:      domain.CompanyID(req.CompanyId),
		DepartmentID:   domain.DepartmentID(req.DepartmentId),
		DepartmentName: req.DepartmentName,
		DateRange:      domain.NewDateRangeFromString(req.DateFrom, req.DateTo),
		Plan:           req.Plan,
		Status:         req.Status,
		SalesID:        domain.FirebaseUserID(req.SalesId),
	})
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	resItems := make([]*pb.MainContract, 0, len(contracts))
	for _, contract := range contracts {
		resItems = append(resItems, response.MainContractFrom(contract))
	}

	return &pb.MainContractList{
		Items: resItems,
	}, nil
}

func (h *contractQuery) CountMainByFilter(ctx context.Context, req *pb.OpeMainContractCountFilterParams) (*pb.Count, error) {
	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	count, err := h.mainContractRepo.GetCountByFilter(ctx, db, adapter.MainContractFilterParams{
		AgencyID:       me.AgencyID,
		CompanyID:      domain.CompanyID(req.CompanyId),
		DepartmentID:   domain.DepartmentID(req.DepartmentId),
		DepartmentName: req.DepartmentName,
		DateRange:      domain.NewDateRangeFromString(req.DateFrom, req.DateTo),
		Plan:           req.Plan,
		Status:         req.Status,
		SalesID:        domain.FirebaseUserID(req.SalesId),
	})
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Count{
		Count: count,
	}, nil
}

func (h *contractQuery) GetMain(ctx context.Context, req *pb.MainContractID) (*pb.MainContract, error) {
	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	contractID := domain.MainContractID(req.Id)

	contract, err := h.mainContractRepo.Get(ctx, db, contractID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	if contract.Department.AgencyID != me.AgencyID {
		return nil, h.errorConverter(ctx, domain.NewForbiddenErr(domain.ForbiddenAgencyMsg))
	}

	return response.MainContractFrom(contract), nil
}

func NewContractCommand(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	auth adapter.AgencyAuthorization,
	contractApp adapter.ContractApp) pb.OpeContractCommand {
	return &contractCommand{
		errorConverter: errorConverter,
		db:             db,
		auth:           auth,
		contractApp:    contractApp,
	}
}

type contractCommand struct {
	errorConverter adapter.ErrorConverter
	db             adapter.DB
	auth           adapter.AgencyAuthorization
	contractApp    adapter.ContractApp
}

func (h *contractCommand) CreateMain(ctx context.Context, req *pb.OpeMainContractCreateParams) (*pb.Empty, error) {
	if err := h.validateCreateMain(req); err != nil {
		return nil, h.errorConverter(ctx, domain.NewBadRequestErr(domain.BadRequestMsg))
	}

	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	departmentID := domain.DepartmentID(req.DepartmentId)

	dateFrom, err := domain.UTCFrom(req.DateFrom)
	if err != nil {
		return nil, h.errorConverter(ctx, domain.NewBadRequestErr(domain.BadRequestDateMsg))
	}

	now := domain.NowUTC()

	app := h.contractApp.OperationBuild(me)

	if _, err := app.CreateMainContract(
		ctx,
		db,
		departmentID,
		adapter.MainContractParams{
			Plan:     req.Plan,
			DateFrom: dateFrom,
		},
		now); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

func (h *contractCommand) validateCreateMain(req *pb.OpeMainContractCreateParams) error {
	var result *multierror.Error
	result = multierror.Append(result, nil)
	if req.Plan == pb.MainProduct_Plan_Unknown {
		result = multierror.Append(result, domain.NewBadRequestErr(domain.BadRequestMsg))
	}
	return result.ErrorOrNil()
}
