package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"gke-go-sample/domain"

	"gke-go-sample/adapter"
	pb "gke-go-sample/proto/go/pb"

	"github.com/twitchtv/twirp"
)

func apply(target http.Handler, handlers ...func(http.Handler) http.Handler) http.Handler {
	h := target
	for _, mw := range handlers {
		h = mw(h)
	}
	return h
}

type APIHandler func()

func NewAPIHandler(
	logger adapter.CompositeLogger,
	cros adapter.Cros,
	adminAuth adapter.AdminAuthenticate,
	cityQuery pb.CityQuery,
	lineQuery pb.LineQuery,
	regionQuery pb.RegionQuery,
	workQuery pb.WorkQuery,
	entryCommand pb.EntryCommand,
	masterAgencyQuery pb.MasterAgencyQuery,
	masterAccountQuery pb.MasterAccountQuery,
	masterAccountCommand pb.MasterAccountCommand,
	masterCompanyQuery pb.MasterCompanyQuery,
	masterDepartmentQuery pb.MasterDepartmentQuery,
	masterDepartmentCommand pb.MasterDepartmentCommand,
	masterWorkQuery pb.MasterWorkQuery,
	masterWorkCommand pb.MasterWorkCommand,
	masterEntryQuery pb.MasterEntryQuery,
	masterProductQuery pb.MasterProductQuery,
	masterContractQuery pb.MasterContractQuery,
	masterContractCommand pb.MasterContractCommand,
	masterUsageStatementQuery pb.MasterUsageStatementQuery,
	opeCompanyQuery pb.OpeCompanyQuery,
	opeCompanyCommand pb.OpeCompanyCommand,
	opeDepartmentQuery pb.OpeDepartmentQuery,
	opeDepartmentCommand pb.OpeDepartmentCommand,
	opeWorkQuery pb.OpeWorkQuery,
	opeWorkCommand pb.OpeWorkCommand,
	opeEntryQuery pb.OpeEntryQuery,
	opeAccountQuery pb.OpeAccountQuery,
	opeAccountCommand pb.OpeAccountCommand,
	opeProductQuery pb.OpeProductQuery,
	opeContractQuery pb.OpeContractQuery,
	opeContractCommand pb.OpeContractCommand,
	opeUsageStatementQuery pb.OpeUsageStatementQuery,
	opeLineQuery pb.OpeLineQuery,
	opeCityQuery pb.OpeCityQuery,
	opeRegionQuery pb.OpeRegionQuery) APIHandler {

	masterMW := func(server pb.TwirpServer) http.Handler {
		return apply(server, adminAuth, cros)
	}

	opeMW := func(server pb.TwirpServer) http.Handler {
		return apply(server, adminAuth, cros)
	}

	publicMW := func(server pb.TwirpServer) http.Handler {
		return apply(server, cros)
	}

	hooks := &twirp.ServerHooks{
		RequestReceived: func(ctx context.Context) (context.Context, error) {
			return ctx, nil
		},
		ResponseSent: func(ctx context.Context) {

		},
		Error: func(ctx context.Context, err twirp.Error) context.Context {
			return ctx
		},
	}

	return func() {
		port := os.Getenv("APP_PORT")
		if len(port) == 0 {
			port = "3000"
		}

		http.Handle(pb.CityQueryPathPrefix, publicMW(pb.NewCityQueryServer(cityQuery, hooks)))
		http.Handle(pb.LineQueryPathPrefix, publicMW(pb.NewLineQueryServer(lineQuery, hooks)))
		http.Handle(pb.RegionQueryPathPrefix, publicMW(pb.NewRegionQueryServer(regionQuery, hooks)))
		http.Handle(pb.WorkQueryPathPrefix, publicMW(pb.NewWorkQueryServer(workQuery, hooks)))
		http.Handle(pb.EntryCommandPathPrefix, publicMW(pb.NewEntryCommandServer(entryCommand, hooks)))

		http.Handle(pb.MasterAgencyQueryPathPrefix, masterMW(pb.NewMasterAgencyQueryServer(masterAgencyQuery, hooks)))
		http.Handle(pb.MasterAccountQueryPathPrefix, masterMW(pb.NewMasterAccountQueryServer(masterAccountQuery, hooks)))
		http.Handle(pb.MasterAccountCommandPathPrefix, masterMW(pb.NewMasterAccountCommandServer(masterAccountCommand, hooks)))
		http.Handle(pb.MasterCompanyQueryPathPrefix, masterMW(pb.NewMasterCompanyQueryServer(masterCompanyQuery, hooks)))
		http.Handle(pb.MasterDepartmentQueryPathPrefix, masterMW(pb.NewMasterDepartmentQueryServer(masterDepartmentQuery, hooks)))
		http.Handle(pb.MasterDepartmentCommandPathPrefix, masterMW(pb.NewMasterDepartmentCommandServer(masterDepartmentCommand, hooks)))
		http.Handle(pb.MasterWorkQueryPathPrefix, masterMW(pb.NewMasterWorkQueryServer(masterWorkQuery, hooks)))
		http.Handle(pb.MasterWorkCommandPathPrefix, masterMW(pb.NewMasterWorkCommandServer(masterWorkCommand, hooks)))
		http.Handle(pb.MasterEntryQueryPathPrefix, masterMW(pb.NewMasterEntryQueryServer(masterEntryQuery, hooks)))
		http.Handle(pb.MasterProductQueryPathPrefix, masterMW(pb.NewMasterProductQueryServer(masterProductQuery, hooks)))
		http.Handle(pb.MasterContractQueryPathPrefix, masterMW(pb.NewMasterContractQueryServer(masterContractQuery, hooks)))
		http.Handle(pb.MasterContractCommandPathPrefix, masterMW(pb.NewMasterContractCommandServer(masterContractCommand, hooks)))
		http.Handle(pb.MasterUsageStatementQueryPathPrefix, masterMW(pb.NewMasterUsageStatementQueryServer(masterUsageStatementQuery, hooks)))

		http.Handle(pb.OpeCompanyQueryPathPrefix, opeMW(pb.NewOpeCompanyQueryServer(opeCompanyQuery, hooks)))
		http.Handle(pb.OpeCompanyCommandPathPrefix, opeMW(pb.NewOpeCompanyCommandServer(opeCompanyCommand, hooks)))
		http.Handle(pb.OpeDepartmentQueryPathPrefix, opeMW(pb.NewOpeDepartmentQueryServer(opeDepartmentQuery, hooks)))
		http.Handle(pb.OpeDepartmentCommandPathPrefix, opeMW(pb.NewOpeDepartmentCommandServer(opeDepartmentCommand, hooks)))
		http.Handle(pb.OpeWorkQueryPathPrefix, opeMW(pb.NewOpeWorkQueryServer(opeWorkQuery, hooks)))
		http.Handle(pb.OpeWorkCommandPathPrefix, opeMW(pb.NewOpeWorkCommandServer(opeWorkCommand, hooks)))
		http.Handle(pb.OpeEntryQueryPathPrefix, opeMW(pb.NewOpeEntryQueryServer(opeEntryQuery, hooks)))
		http.Handle(pb.OpeAccountQueryPathPrefix, opeMW(pb.NewOpeAccountQueryServer(opeAccountQuery, hooks)))
		http.Handle(pb.OpeAccountCommandPathPrefix, opeMW(pb.NewOpeAccountCommandServer(opeAccountCommand, hooks)))
		http.Handle(pb.OpeProductQueryPathPrefix, opeMW(pb.NewOpeProductQueryServer(opeProductQuery, hooks)))
		http.Handle(pb.OpeContractQueryPathPrefix, opeMW(pb.NewOpeContractQueryServer(opeContractQuery, hooks)))
		http.Handle(pb.OpeContractCommandPathPrefix, opeMW(pb.NewOpeContractCommandServer(opeContractCommand, hooks)))
		http.Handle(pb.OpeUsageStatementQueryPathPrefix, opeMW(pb.NewOpeUsageStatementQueryServer(opeUsageStatementQuery, hooks)))
		http.Handle(pb.OpeLineQueryPathPrefix, opeMW(pb.NewOpeLineQueryServer(opeLineQuery, hooks)))
		http.Handle(pb.OpeCityQueryPathPrefix, opeMW(pb.NewOpeCityQueryServer(opeCityQuery, hooks)))
		http.Handle(pb.OpeRegionQueryPathPrefix, opeMW(pb.NewOpeRegionQueryServer(opeRegionQuery, hooks)))

		http.HandleFunc("/health_check", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("ok"))
		})

		logger.Info().Printf("running server port:%s", port)

		_ = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	}
}

type BatchHandler func(ctx context.Context, task pb.BatchTask)

func NewBatchHandler(
	logger adapter.CompositeLogger,
	db adapter.DB,
	workApp adapter.WorkApp) BatchHandler {
	return func(ctx context.Context, task pb.BatchTask) {
		db := db(ctx)

		now := domain.NowUTC()

		logger.Info().Printf("start %s batch at %#v", task, now.String())

		var err error
		switch task {
		case pb.BatchTask_ReIndexSearch:
			err = workApp.BatchBuild().ReIndex(ctx, db)
		case pb.BatchTask_ProceedWorkStatus:
			err = workApp.BatchBuild().Proceed(ctx, db, now)
		case pb.BatchTask_ReOrderWork:
			err = workApp.BatchBuild().ReOrder(ctx, db)
		}

		if err != nil {
			panic(fmt.Errorf("raise error of %s, error: %s", task.String(), err.Error()))
		}
	}
}
