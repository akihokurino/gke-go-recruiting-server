package master_handler

import (
	"context"

	"gke-go-recruiting-server/adapter"
	"gke-go-recruiting-server/domain"
	"gke-go-recruiting-server/handler/response"
	pb "gke-go-recruiting-server/proto/go/pb"
)

func NewAgencyQuery(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	auth adapter.AdminAuthorization,
	agencyRepo adapter.AgencyRepo) pb.MasterAgencyQuery {
	return &agencyQuery{
		errorConverter: errorConverter,
		db:             db,
		auth:           auth,
		agencyRepo:     agencyRepo,
	}
}

type agencyQuery struct {
	errorConverter adapter.ErrorConverter
	db             adapter.DB
	auth           adapter.AdminAuthorization
	agencyRepo     adapter.AgencyRepo
}

func (h *agencyQuery) List(ctx context.Context, req *pb.Empty) (*pb.AgencyList, error) {
	db := h.db(ctx)

	_, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	agencies, err := h.agencyRepo.GetAll(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	resItems := make([]*pb.Agency, 0, len(agencies))
	for _, agency := range agencies {
		resItems = append(resItems, response.AgencyFrom(agency))
	}

	return &pb.AgencyList{
		Items: resItems,
	}, nil
}

func (h *agencyQuery) Get(ctx context.Context, req *pb.AgencyID) (*pb.Agency, error) {
	db := h.db(ctx)

	_, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	agencyID := domain.AgencyID(req.Id)

	agency, err := h.agencyRepo.Get(ctx, db, agencyID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return response.AgencyFrom(agency), nil
}
