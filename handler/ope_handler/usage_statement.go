package ope_handler

import (
	"context"

	pb "gke-go-recruiting-server/proto/go/pb"

	"gke-go-recruiting-server/adapter"
	"gke-go-recruiting-server/domain"
	"gke-go-recruiting-server/handler/response"
)

func NewUsageStatementQuery(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	auth adapter.AgencyAuthorization,
	usageStatementRepo adapter.UsageStatementRepo) pb.OpeUsageStatementQuery {
	return &usageStatementQuery{
		errorConverter:     errorConverter,
		db:                 db,
		auth:               auth,
		usageStatementRepo: usageStatementRepo,
	}
}

type usageStatementQuery struct {
	errorConverter     adapter.ErrorConverter
	db                 adapter.DB
	auth               adapter.AgencyAuthorization
	usageStatementRepo adapter.UsageStatementRepo
}

func (h *usageStatementQuery) ListByFilter(
	ctx context.Context,
	req *pb.OpeUsageStatementFilterParams) (*pb.UsageStatementList, error) {
	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	pager := domain.NewPager(req.Pager.Page, req.Pager.Offset)

	statements, err := h.usageStatementRepo.GetByFilterWithPager(ctx, db, pager, adapter.UsageStatementFilterParams{
		AgencyID:         me.AgencyID,
		CompanyID:        domain.CompanyID(req.CompanyId),
		DepartmentID:     domain.DepartmentID(req.DepartmentId),
		DepartmentName:   req.DepartmentName,
		UsageStatementID: domain.UsageStatementID(req.UsageStatementId),
		DateRange:        domain.NewDateRangeFromString(req.DateFrom, req.DateTo),
	})
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	resItems := make([]*pb.UsageStatement, 0, len(statements))
	for _, statement := range statements {
		resItems = append(resItems, response.UsageStatementFrom(statement))
	}

	return &pb.UsageStatementList{
		Items: resItems,
	}, nil
}

func (h *usageStatementQuery) CountByFilter(
	ctx context.Context,
	req *pb.OpeUsageStatementCountFilterParams) (*pb.Count, error) {
	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	count, err := h.usageStatementRepo.GetCountByFilter(ctx, db, adapter.UsageStatementFilterParams{
		AgencyID:         me.AgencyID,
		CompanyID:        domain.CompanyID(req.CompanyId),
		DepartmentID:     domain.DepartmentID(req.DepartmentId),
		DepartmentName:   req.DepartmentName,
		UsageStatementID: domain.UsageStatementID(req.UsageStatementId),
		DateRange:        domain.NewDateRangeFromString(req.DateFrom, req.DateTo),
	})
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Count{
		Count: count,
	}, nil
}
