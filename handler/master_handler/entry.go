package master_handler

import (
	"context"

	"gke-go-recruiting-server/adapter"
	"gke-go-recruiting-server/domain"
	"gke-go-recruiting-server/handler/response"
	pb "gke-go-recruiting-server/proto/go/pb"
)

func NewEntryQuery(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	auth adapter.AdminAuthorization,
	entryRepo adapter.EntryRepo) pb.MasterEntryQuery {
	return &entryQuery{
		errorConverter: errorConverter,
		db:             db,
		auth:           auth,
		entryRepo:      entryRepo,
	}
}

type entryQuery struct {
	errorConverter adapter.ErrorConverter
	db             adapter.DB
	auth           adapter.AdminAuthorization
	entryRepo      adapter.EntryRepo
}

func (h *entryQuery) ListByFilter(ctx context.Context, req *pb.MasterEntryFilterParams) (*pb.EntryList, error) {
	db := h.db(ctx)

	_, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	pager := domain.NewPager(req.Pager.Page, req.Pager.Offset)

	entries, err := h.entryRepo.GetByFilterWithPager(ctx, db, pager, adapter.EntryFilterParams{
		AgencyID:          domain.AgencyID(req.AgencyId),
		CompanyID:         domain.CompanyID(req.CompanyId),
		DepartmentID:      domain.DepartmentID(req.DepartmentId),
		DepartmentName:    req.DepartmentName,
		WorkID:            domain.WorkID(req.WorkId),
		DateRange:         domain.NewDateRangeFromString(req.DateFrom, req.DateTo),
		SalesID:           domain.FirebaseUserID(req.SalesId),
		BusinessCondition: req.BusinessCondition,
		PrefID:            domain.PrefID(req.PrefId),
	})
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	resItems := make([]*pb.Entry, 0, len(entries))
	for _, entry := range entries {
		resItems = append(resItems, response.EntryFrom(entry))
	}

	return &pb.EntryList{
		Items: resItems,
	}, nil
}

func (h *entryQuery) CountByFilter(ctx context.Context, req *pb.MasterEntryCountFilterParams) (*pb.Count, error) {
	db := h.db(ctx)

	_, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	count, err := h.entryRepo.GetCountByFilter(ctx, db, adapter.EntryFilterParams{
		AgencyID:          domain.AgencyID(req.AgencyId),
		CompanyID:         domain.CompanyID(req.CompanyId),
		DepartmentID:      domain.DepartmentID(req.DepartmentId),
		DepartmentName:    req.DepartmentName,
		WorkID:            domain.WorkID(req.WorkId),
		DateRange:         domain.NewDateRangeFromString(req.DateFrom, req.DateTo),
		SalesID:           domain.FirebaseUserID(req.SalesId),
		BusinessCondition: req.BusinessCondition,
		PrefID:            domain.PrefID(req.PrefId),
	})
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Count{
		Count: count,
	}, nil
}
