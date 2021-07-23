package ope_handler

import (
	"context"

	"gke-go-sample/adapter"
	"gke-go-sample/handler/response"
	pb "gke-go-sample/proto/go/pb"
)

func NewLineQuery(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	lineRepo adapter.LineRepo) pb.OpeLineQuery {
	return &lineQuery{
		errorConverter: errorConverter,
		db:             db,
		lineRepo:       lineRepo,
	}
}

type lineQuery struct {
	errorConverter adapter.ErrorConverter
	db             adapter.DB
	lineRepo       adapter.LineRepo
}

func (h *lineQuery) ListByDistance(ctx context.Context, req *pb.LineListByDistanceParams) (*pb.LineList, error) {
	db := h.db(ctx)

	lines, err := h.lineRepo.GetByDistance(ctx, db, req.Latitude, req.Longitude, req.DistanceKm)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	resItems := make([]*pb.Line, 0, len(lines))
	for _, line := range lines {
		resItems = append(resItems, response.LineFrom(line))
	}

	return &pb.LineList{
		Items: resItems,
	}, nil
}
