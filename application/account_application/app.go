package account_application

import (
	"gke-go-recruiting-server/adapter"
	"gke-go-recruiting-server/domain/account_domain"
)

type app struct {
	tx                adapter.TX
	firebaseRepo      adapter.FirebaseRepo
	administratorRepo adapter.AdministratorRepo
	agencyRepo        adapter.AgencyRepo
	agencyAccountRepo adapter.AgencyAccountRepo
}

func NewApp(
	tx adapter.TX,
	firebaseRepo adapter.FirebaseRepo,
	administratorRepo adapter.AdministratorRepo,
	agencyRepo adapter.AgencyRepo,
	agencyAccountRepo adapter.AgencyAccountRepo) adapter.AccountApp {
	return &app{
		tx:                tx,
		firebaseRepo:      firebaseRepo,
		administratorRepo: administratorRepo,
		agencyRepo:        agencyRepo,
		agencyAccountRepo: agencyAccountRepo,
	}
}

func (a *app) MasterBuild(me *account_domain.Administrator) adapter.AccountAppForMaster {
	return &masterApp{
		me:  me,
		app: a,
	}
}

func (a *app) OperationBuild(me *account_domain.AgencyAccount) adapter.AccountAppForOperation {
	return &operationApp{
		me:  me,
		app: a,
	}
}
