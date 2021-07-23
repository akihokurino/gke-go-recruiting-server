package ope_handler

import (
	"context"
	"time"

	"gke-go-sample/util"

	"github.com/hashicorp/go-multierror"

	"gke-go-sample/domain"

	pb "gke-go-sample/proto/go/pb"

	"gke-go-sample/adapter"
	"gke-go-sample/handler/response"
)

func NewAccountQuery(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	auth adapter.AgencyAuthorization,
	agencyAccountRepo adapter.AgencyAccountRepo) pb.OpeAccountQuery {
	return &accountQuery{
		errorConverter:    errorConverter,
		db:                db,
		auth:              auth,
		agencyAccountRepo: agencyAccountRepo,
	}
}

type accountQuery struct {
	errorConverter    adapter.ErrorConverter
	db                adapter.DB
	auth              adapter.AgencyAuthorization
	agencyAccountRepo adapter.AgencyAccountRepo
}

func (h *accountQuery) List(ctx context.Context, req *pb.Empty) (*pb.AgencyAccountList, error) {
	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	accounts, err := h.agencyAccountRepo.GetByAgency(ctx, db, me.AgencyID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	resItems := make([]*pb.AgencyAccount, 0, len(accounts))
	for _, account := range accounts {
		resItems = append(resItems, response.AgencyAccountFrom(account))
	}

	return &pb.AgencyAccountList{
		Items: resItems,
	}, nil
}

func (h *accountQuery) Get(ctx context.Context, req *pb.AgencyAccountID) (*pb.AgencyAccount, error) {
	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	accountID := domain.FirebaseUserID(req.Id)

	account, err := h.agencyAccountRepo.Get(ctx, db, accountID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	if account.AgencyID != me.AgencyID {
		return nil, h.errorConverter(ctx, domain.NewForbiddenErr(domain.ForbiddenAgencyMsg))
	}

	return response.AgencyAccountFrom(account), nil
}

func NewAccountCommand(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	auth adapter.AgencyAuthorization,
	accountApp adapter.AccountApp) pb.OpeAccountCommand {
	return &accountCommand{
		errorConverter: errorConverter,
		db:             db,
		auth:           auth,
		accountApp:     accountApp,
	}
}

type accountCommand struct {
	errorConverter adapter.ErrorConverter
	db             adapter.DB
	auth           adapter.AgencyAuthorization
	accountApp     adapter.AccountApp
}

func (h *accountCommand) Create(ctx context.Context, req *pb.CreateAgencyAccountParams) (*pb.AgencyAccount, error) {
	if err := h.validateCreate(req); err != nil {
		return nil, h.errorConverter(ctx, domain.NewBadRequestErr(domain.BadRequestMsg))
	}

	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	now := time.Now()

	app := h.accountApp.OperationBuild(me)

	account, err := app.CreateAgencyAccount(
		ctx,
		db,
		adapter.AgencyAccountParams{
			Email:    req.Email,
			Password: req.Password,
			Name:     req.Name,
			NameKana: req.NameKana,
		},
		now)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return response.AgencyAccountFrom(account), nil
}

func (h *accountCommand) validateCreate(req *pb.CreateAgencyAccountParams) error {
	var result *multierror.Error
	result = multierror.Append(result, util.ValidateTextRange(req.Email, 1, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.Password, 1, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.Name, 1, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.NameKana, 1, 255))
	return result.ErrorOrNil()
}

func (h *accountCommand) Delete(ctx context.Context, req *pb.AgencyAccountID) (*pb.Empty, error) {
	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	accountID := domain.FirebaseUserID(req.Id)

	app := h.accountApp.OperationBuild(me)

	if err := app.DeleteAgencyAccount(ctx, db, accountID); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}
