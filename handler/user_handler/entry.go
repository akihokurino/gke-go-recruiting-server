package user_handler

import (
	"context"

	"gke-go-recruiting-server/util"

	"github.com/hashicorp/go-multierror"

	"gke-go-recruiting-server/domain"

	"gke-go-recruiting-server/adapter"
	pb "gke-go-recruiting-server/proto/go/pb"
)

func NewEntryCommand(
	errorConverter adapter.ErrorConverter,
	db adapter.DB,
	entryApp adapter.EntryApp,
) pb.EntryCommand {
	return &entryCommand{
		errorConverter: errorConverter,
		db:             db,
		entryApp:       entryApp,
	}
}

type entryCommand struct {
	errorConverter adapter.ErrorConverter
	db             adapter.DB
	entryApp       adapter.EntryApp
}

func (h *entryCommand) Create(ctx context.Context, req *pb.EntryParams) (*pb.Empty, error) {
	if err := h.validateCreate(req); err != nil {
		return nil, h.errorConverter(ctx, domain.NewBadRequestErr(domain.BadRequestMsg))
	}

	db := h.db(ctx)

	now := domain.NowUTC()

	birthdate, err := domain.DateFrom(req.Birthdate)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	app := h.entryApp.PublicBuild()

	if err := app.Entry(
		ctx,
		db,
		domain.WorkID(req.WorkId),
		adapter.EntryParams{
			FullName:               req.Fullname,
			FullNameKana:           req.FullnameKana,
			Birthdate:              birthdate,
			Gender:                 req.Gender,
			PhoneNumber:            req.PhoneNumber,
			Email:                  req.Email,
			Question:               req.Question,
			Category:               req.Category,
			PrefID:                 domain.PrefID(req.PrefId),
			PreferredContactMethod: req.PreferredContactMethod,
			PreferredContactTime:   req.PreferredContactTime,
		},
		now); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

func (h *entryCommand) validateCreate(req *pb.EntryParams) error {
	var result *multierror.Error
	result = multierror.Append(result, util.ValidateTextRange(req.Fullname, 1, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.FullnameKana, 1, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.PhoneNumber, 1, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.Email, 1, 255))
	result = multierror.Append(result, util.ValidateTextRange(req.Question, 0, 255))
	if req.Gender == pb.User_Gender_Unknown {
		result = multierror.Append(result, domain.NewBadRequestErr(domain.BadRequestMsg))
	}
	if req.Category == pb.User_Category_Unknown {
		result = multierror.Append(result, domain.NewBadRequestErr(domain.BadRequestMsg))
	}
	return result.ErrorOrNil()
}
