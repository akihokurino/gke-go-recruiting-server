package ope_handler

import (
	"context"

	pb "gke-go-recruiting-server/proto/go/pb"

	"gke-go-recruiting-server/adapter"
	"gke-go-recruiting-server/domain"
	"gke-go-recruiting-server/handler/response"
)

func NewEntryQuery(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	auth adapter.AgencyAuthorization,
	entryRepo adapter.EntryRepo) pb.OpeEntryQuery {
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
	auth           adapter.AgencyAuthorization
	entryRepo      adapter.EntryRepo
}

func (h *entryQuery) ListByFilter(ctx context.Context, req *pb.OpeEntryFilterParams) (*pb.EntryList, error) {
	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	pager := domain.NewPager(req.Pager.Page, req.Pager.Offset)

	entries, err := h.entryRepo.GetByFilterWithPager(ctx, db, pager, adapter.EntryFilterParams{
		AgencyID:       me.AgencyID,
		DateRange:      domain.NewDateRangeFromString(req.DateFrom, req.DateTo),
		CompanyID:      domain.CompanyID(req.CompanyId),
		DepartmentID:   domain.DepartmentID(req.DepartmentId),
		DepartmentName: req.DepartmentName,
		SalesID:        domain.FirebaseUserID(req.SalesId),
		Status:         req.Status,
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

func (h *entryQuery) CountByFilter(ctx context.Context, req *pb.OpeEntryCountFilterParams) (*pb.Count, error) {
	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	count, err := h.entryRepo.GetCountByFilter(ctx, db, adapter.EntryFilterParams{
		AgencyID:       me.AgencyID,
		DateRange:      domain.NewDateRangeFromString(req.DateFrom, req.DateTo),
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
