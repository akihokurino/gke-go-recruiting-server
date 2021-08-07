package master_handler

import (
	"context"
	"time"

	"gke-go-recruiting-server/util"

	"github.com/hashicorp/go-multierror"

	"gke-go-recruiting-server/handler/response"

	"gke-go-recruiting-server/domain"

	"gke-go-recruiting-server/adapter"
	pb "gke-go-recruiting-server/proto/go/pb"
)

func NewAccountQuery(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	auth adapter.AdminAuthorization,
	agencyAccountRepo adapter.AgencyAccountRepo) pb.MasterAccountQuery {
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
	auth              adapter.AdminAuthorization
	agencyAccountRepo adapter.AgencyAccountRepo
}

func (h *accountQuery) ListAgencyAccount(ctx context.Context, req *pb.AgencyID) (*pb.AgencyAccountList, error) {
	db := h.db(ctx)

	_, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	agencyID := domain.AgencyID(req.Id)

	accounts, err := h.agencyAccountRepo.GetByAgency(ctx, db, agencyID)
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

func (h *accountQuery) GetAgencyAccount(ctx context.Context, req *pb.AgencyAccountID) (*pb.AgencyAccount, error) {
	db := h.db(ctx)

	_, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	accountID := domain.FirebaseUserID(req.Id)

	account, err := h.agencyAccountRepo.Get(ctx, db, accountID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return response.AgencyAccountFrom(account), nil
}

func (h *accountQuery) ListAllAgencyAccount(ctx context.Context, req *pb.Empty) (*pb.AgencyAccountList, error) {
	db := h.db(ctx)

	_, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	accounts, err := h.agencyAccountRepo.GetAll(ctx, db)
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

func NewAccountCommand(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	auth adapter.AdminAuthorization,
	accountApp adapter.AccountApp) pb.MasterAccountCommand {
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
	auth           adapter.AdminAuthorization
	accountApp     adapter.AccountApp
}

func (h *accountCommand) CreateAdministrator(ctx context.Context, req *pb.CreateAdministratorParams) (*pb.Administrator, error) {
	if err := h.validateCreateAdministrator(req); err != nil {
		return nil, h.errorConverter(ctx, domain.NewBadRequestErr(domain.BadRequestMsg))
	}

	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	now := time.Now()

	app := h.accountApp.MasterBuild(me)

	account, err := app.CreateAdministrator(ctx, db, adapter.AdministratorParams{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}, now)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return response.AdministratorFrom(account), nil
}

func (h *accountCommand) validateCreateAdministrator(req *pb.CreateAdministratorParams) error {
	var result *multierror.Error
	result = multierror.Append(result, util.ValidateTextRange(req.Email, 1, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.Password, 1, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.Name, 1, 255))
	return result.ErrorOrNil()
}

func (h *accountCommand) CreateAgency(ctx context.Context, req *pb.CreateAgencyAccountByMasterParams) (*pb.AgencyAccount, error) {
	if err := h.validateCreateAgency(req); err != nil {
		return nil, h.errorConverter(ctx, domain.NewBadRequestErr(domain.BadRequestMsg))
	}

	db := h.db(ctx)

	me, err := h.auth(ctx, db)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	now := time.Now()

	app := h.accountApp.MasterBuild(me)

	agencyID := domain.AgencyID(req.AgencyId)

	account, err := app.CreateAgencyAccount(
		ctx,
		db,
		adapter.AgencyAccountParams{
			Email:    req.Email,
			Password: req.Password,
			Name:     req.Name,
			NameKana: req.NameKana,
		},
		agencyID,
		now)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return response.AgencyAccountFrom(account), nil
}

func (h *accountCommand) validateCreateAgency(req *pb.CreateAgencyAccountByMasterParams) error {
	var result *multierror.Error
	result = multierror.Append(result, util.ValidateTextRange(req.Email, 1, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.Password, 1, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.Name, 1, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.NameKana, 1, 255))
	return result.ErrorOrNil()
}
