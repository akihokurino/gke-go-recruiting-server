package ope_handler

import (
	"context"

	"gke-go-sample/util"

	"github.com/hashicorp/go-multierror"

	"gke-go-sample/adapter"
	"gke-go-sample/domain"
	"gke-go-sample/handler/response"
	pb "gke-go-sample/proto/go/pb"
)

func NewCompanyQuery(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	auth adapter.AgencyAuthorization,
	companyRepo adapter.CompanyRepo) pb.OpeCompanyQuery {
	return &companyQuery{
		errorConverter: errorConverter,
		db:             db,
		auth:           auth,
		companyRepo:    companyRepo,
	}
}

type companyQuery struct {
	errorConverter adapter.ErrorConverter
	db             adapter.DB
	auth           adapter.AgencyAuthorization
	companyRepo    adapter.CompanyRepo
}

func (h *companyQuery) ListByFilter(ctx context.Context, req *pb.OpeCompanyFilterParams) (*pb.CompanyList, error) {
	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	pager := domain.NewPager(req.Pager.Page, req.Pager.Offset)

	companies, err := h.companyRepo.GetByFilterWithPager(ctx, db, pager, adapter.CompanyFilterParams{
		AgencyID:    me.AgencyID,
		CompanyID:   domain.CompanyID(req.CompanyId),
		CompanyName: req.CompanyName,
	})
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	resItems := make([]*pb.Company, 0, len(companies))
	for _, company := range companies {
		resItems = append(resItems, response.CompanyFrom(company))
	}

	return &pb.CompanyList{
		Items: resItems,
	}, nil
}

func (h *companyQuery) CountByFilter(ctx context.Context, req *pb.OpeCompanyCountFilterParams) (*pb.Count, error) {
	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	count, err := h.companyRepo.GetCountByFilter(ctx, db, adapter.CompanyFilterParams{
		AgencyID:    me.AgencyID,
		CompanyID:   domain.CompanyID(req.CompanyId),
		CompanyName: req.CompanyName,
	})
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Count{
		Count: count,
	}, nil
}

func (h *companyQuery) Get(ctx context.Context, req *pb.CompanyID) (*pb.Company, error) {
	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	companyID := domain.CompanyID(req.Id)

	company, err := h.companyRepo.Get(ctx, db, companyID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	if company.AgencyID != me.AgencyID {
		return nil, h.errorConverter(ctx, domain.NewForbiddenErr(domain.ForbiddenAgencyMsg))
	}

	return response.CompanyFrom(company), nil
}

func NewCompanyCommand(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	auth adapter.AgencyAuthorization,
	companyApp adapter.CompanyApp) pb.OpeCompanyCommand {
	return &companyCommand{
		errorConverter: errorConverter,
		db:             db,
		auth:           auth,
		companyApp:     companyApp,
	}
}

type companyCommand struct {
	errorConverter adapter.ErrorConverter
	db             adapter.DB
	auth           adapter.AgencyAuthorization
	companyApp     adapter.CompanyApp
}

func (h *companyCommand) Create(ctx context.Context, req *pb.OpeCompanyCreateParams) (*pb.CompanyID, error) {
	if err := h.validateCreate(req); err != nil {
		return nil, h.errorConverter(ctx, domain.NewBadRequestErr(domain.BadRequestMsg))
	}

	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	companyParams := adapter.CompanyParams{
		Rank:        req.Rank,
		RankType:    req.RankType,
		Name:        req.Name,
		NameKana:    req.NameKana,
		PostalCode:  req.PostalCode,
		PrefID:      domain.PrefID(req.PrefId),
		Address:     req.Address,
		Building:    req.Building,
		PhoneNumber: req.PhoneNumber,
	}

	now := domain.NowUTC()

	app := h.companyApp.OperationBuild(me)

	company, err := app.Create(ctx, db, companyParams, now)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.CompanyID{
		Id: company.ID.String(),
	}, nil
}

func (h *companyCommand) validateCreate(req *pb.OpeCompanyCreateParams) error {
	var result *multierror.Error
	result = multierror.Append(result, util.ValidateTextRange(req.Name, 1, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.NameKana, 1, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.PostalCode, 1, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.Address, 1, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.Building, 0, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.PhoneNumber, 1, 255))
	return result.ErrorOrNil()
}
