package department_application

import (
	"gke-go-recruiting-server/adapter"
	"gke-go-recruiting-server/domain/account_domain"
)

type app struct {
	tx                    adapter.TX
	departmentRepo        adapter.DepartmentRepo
	departmentImageRepo   adapter.DepartmentImageRepo
	departmentStationRepo adapter.DepartmentStationRepo
	agencyAccountRepo     adapter.AgencyAccountRepo
}

func NewApp(
	tx adapter.TX,
	departmentRepo adapter.DepartmentRepo,
	departmentImageRepo adapter.DepartmentImageRepo,
	departmentStationRepo adapter.DepartmentStationRepo,
	agencyAccountRepo adapter.AgencyAccountRepo) adapter.DepartmentApp {
	return &app{
		tx:                    tx,
		departmentRepo:        departmentRepo,
		departmentImageRepo:   departmentImageRepo,
		departmentStationRepo: departmentStationRepo,
		agencyAccountRepo:     agencyAccountRepo,
	}
}

func (a *app) MasterBuild(me *account_domain.Administrator) adapter.DepartmentAppForMaster {
	return &masterApp{
		me:  me,
		app: a,
	}
}

func (a *app) OperationBuild(me *account_domain.AgencyAccount) adapter.DepartmentAppForOperation {
	return &operationApp{
		me:  me,
		app: a,
	}
}
