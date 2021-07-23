package work_application

import (
	"gke-go-sample/adapter"
	"gke-go-sample/domain/account_domain"
)

type app struct {
	logger             adapter.CompositeLogger
	tx                 adapter.TX
	workRepo           adapter.WorkRepo
	workIndexRepo      adapter.WorkIndexRepo
	workImageRepo      adapter.WorkImageRepo
	workMovieRepo      adapter.WorkMovieRepo
	workMeritRepo      adapter.WorkMeritRepo
	mainContractRepo   adapter.MainContractRepo
	departmentRepo     adapter.DepartmentRepo
	workActivePlanRepo adapter.WorkActivePlanRepo
}

func NewApp(
	logger adapter.CompositeLogger,
	tx adapter.TX,
	workRepo adapter.WorkRepo,
	workIndexRepo adapter.WorkIndexRepo,
	workImageRepo adapter.WorkImageRepo,
	workMovieRepo adapter.WorkMovieRepo,
	workMeritRepo adapter.WorkMeritRepo,
	mainContractRepo adapter.MainContractRepo,
	departmentRepo adapter.DepartmentRepo,
	workActivePlanRepo adapter.WorkActivePlanRepo) adapter.WorkApp {
	return &app{
		logger:             logger,
		tx:                 tx,
		workRepo:           workRepo,
		workIndexRepo:      workIndexRepo,
		workImageRepo:      workImageRepo,
		workMovieRepo:      workMovieRepo,
		workMeritRepo:      workMeritRepo,
		mainContractRepo:   mainContractRepo,
		departmentRepo:     departmentRepo,
		workActivePlanRepo: workActivePlanRepo,
	}
}

func (a *app) MasterBuild(me *account_domain.Administrator) adapter.WorkAppForMaster {
	return &masterApp{
		me:  me,
		app: a,
	}
}

func (a *app) OperationBuild(me *account_domain.AgencyAccount) adapter.WorkAppForOperation {
	return &operationApp{
		me:  me,
		app: a,
	}
}

func (a *app) BatchBuild() adapter.WorkAppForBatch {
	return &batchApp{
		app: a,
	}
}
