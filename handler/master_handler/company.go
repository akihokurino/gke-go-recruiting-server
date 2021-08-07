package master_handler

import (
	"context"

	"gke-go-recruiting-server/adapter"
	"gke-go-recruiting-server/domain"
	"gke-go-recruiting-server/handler/response"
	pb "gke-go-recruiting-server/proto/go/pb"
)

func NewCompanyQuery(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	auth adapter.AdminAuthorization,
	companyRepo adapter.CompanyRepo) pb.MasterCompanyQuery {
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
	auth           adapter.AdminAuthorization
	companyRepo    adapter.CompanyRepo
}

func (h *companyQuery) ListByFilter(ctx context.Context, req *pb.MasterCompanyFilterParams) (*pb.CompanyList, error) {
	db := h.db(ctx)

	_, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	pager := domain.NewPager(req.Pager.Page, req.Pager.Offset)

	companies, err := h.companyRepo.GetByFilterWithPager(ctx, db, pager, adapter.CompanyFilterParams{
		CompanyID:   domain.CompanyID(req.CompanyId),
		CompanyName: req.CompanyName,
		AgencyID:    domain.AgencyID(req.AgencyId),
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

func (h *companyQuery) CountByFilter(ctx context.Context, req *pb.MasterCompanyCountFilterParams) (*pb.Count, error) {
	db := h.db(ctx)

	_, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	count, err := h.companyRepo.GetCountByFilter(ctx, db, adapter.CompanyFilterParams{
		CompanyID:   domain.CompanyID(req.CompanyId),
		CompanyName: req.CompanyName,
		AgencyID:    domain.AgencyID(req.AgencyId),
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

	_, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	companyID := domain.CompanyID(req.Id)

	company, err := h.companyRepo.Get(ctx, db, companyID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return response.CompanyFrom(company), nil
}
