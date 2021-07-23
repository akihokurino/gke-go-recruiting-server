package contract_application

import (
	"gke-go-sample/adapter"
	"gke-go-sample/domain/account_domain"
)

type app struct {
	tx                 adapter.TX
	mainContractRepo   adapter.MainContractRepo
	departmentRepo     adapter.DepartmentRepo
	usageStatementRepo adapter.UsageStatementRepo
	workActivePlanRepo adapter.WorkActivePlanRepo
	workRepo           adapter.WorkRepo
}

func NewApp(
	tx adapter.TX,
	mainContractRepo adapter.MainContractRepo,
	departmentRepo adapter.DepartmentRepo,
	usageStatementRepo adapter.UsageStatementRepo,
	workActivePlanRepo adapter.WorkActivePlanRepo,
	workRepo adapter.WorkRepo) adapter.ContractApp {
	return &app{
		tx:                 tx,
		mainContractRepo:   mainContractRepo,
		departmentRepo:     departmentRepo,
		usageStatementRepo: usageStatementRepo,
		workActivePlanRepo: workActivePlanRepo,
		workRepo:           workRepo,
	}
}

func (a *app) MasterBuild(me *account_domain.Administrator) adapter.ContractAppForMaster {
	return &masterApp{
		me:  me,
		app: a,
	}
}

func (a *app) OperationBuild(me *account_domain.AgencyAccount) adapter.ContractAppForOperation {
	return &operationApp{
		me:  me,
		app: a,
	}
}
