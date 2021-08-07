package master_handler

import (
	"context"

	"gke-go-recruiting-server/adapter"
	"gke-go-recruiting-server/domain"
	"gke-go-recruiting-server/handler/response"
	pb "gke-go-recruiting-server/proto/go/pb"
)

func NewContractQuery(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	auth adapter.AdminAuthorization,
	mainContractRepo adapter.MainContractRepo) pb.MasterContractQuery {
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
	auth             adapter.AdminAuthorization
	mainContractRepo adapter.MainContractRepo
}

func (h *contractQuery) ListMainByFilter(ctx context.Context, req *pb.MasterMainContractFilterParams) (*pb.MainContractList, error) {
	db := h.db(ctx)

	_, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	pager := domain.NewPager(req.Pager.Page, req.Pager.Offset)

	contracts, err := h.mainContractRepo.GetByFilterWithPager(ctx, db, pager, adapter.MainContractFilterParams{
		CompanyID:      domain.CompanyID(req.CompanyId),
		DepartmentID:   domain.DepartmentID(req.DepartmentId),
		DepartmentName: req.DepartmentName,
		DateRange:      domain.NewDateRangeFromString(req.DateFrom, req.DateTo),
		Plan:           req.Plan,
		Status:         req.Status,
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

func (h *contractQuery) CountMainByFilter(ctx context.Context, req *pb.MasterMainContractCountFilterParams) (*pb.Count, error) {
	db := h.db(ctx)

	_, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	count, err := h.mainContractRepo.GetCountByFilter(ctx, db, adapter.MainContractFilterParams{
		CompanyID:      domain.CompanyID(req.CompanyId),
		DepartmentID:   domain.DepartmentID(req.DepartmentId),
		DepartmentName: req.DepartmentName,
		DateRange:      domain.NewDateRangeFromString(req.DateFrom, req.DateTo),
		Plan:           req.Plan,
		Status:         req.Status,
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

	_, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	contractID := domain.MainContractID(req.Id)

	contract, err := h.mainContractRepo.Get(ctx, db, contractID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return response.MainContractFrom(contract), nil
}

func NewContractCommand(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	auth adapter.AdminAuthorization,
	contractApp adapter.ContractApp) pb.MasterContractCommand {
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
	auth           adapter.AdminAuthorization
	contractApp    adapter.ContractApp
}

func (h *contractCommand) AcceptMain(ctx context.Context, req *pb.MainContractID) (*pb.Empty, error) {
	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	contractID := domain.MainContractID(req.Id)
	now := domain.NowUTC()

	app := h.contractApp.MasterBuild(me)

	if err := app.AcceptMainContract(ctx, db, contractID, now); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

func (h *contractCommand) DenyMain(ctx context.Context, req *pb.MainContractID) (*pb.Empty, error) {
	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	contractID := domain.MainContractID(req.Id)

	app := h.contractApp.MasterBuild(me)

	if err := app.DenyMainContract(ctx, db, contractID); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}
