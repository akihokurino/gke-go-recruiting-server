package ope_handler

import (
	"context"

	"gke-go-sample/adapter"
	"gke-go-sample/domain/product_domain"
	"gke-go-sample/handler/response"
	pb "gke-go-sample/proto/go/pb"
)

func NewProductQuery(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	auth adapter.AgencyAuthorization) pb.OpeProductQuery {
	return &productQuery{
		errorConverter: errorConverter,
		db:             db,
		auth:           auth,
	}
}

type productQuery struct {
	errorConverter adapter.ErrorConverter
	db             adapter.DB
	auth           adapter.AgencyAuthorization
}

func (h *productQuery) ListMain(ctx context.Context, req *pb.Empty) (*pb.MainProductList, error) {
	db := h.db(ctx)

	_, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	products := product_domain.GetMainList()

	resItems := make([]*pb.MainProduct, 0, len(products))
	for _, product := range products {
		resItems = append(resItems, response.MainProductFrom(product))
	}

	return &pb.MainProductList{
		Items: resItems,
	}, nil
}