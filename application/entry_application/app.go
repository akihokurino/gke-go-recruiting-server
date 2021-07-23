package entry_application

import (
	"gke-go-sample/adapter"
)

type app struct {
	tx        adapter.TX
	workRepo  adapter.WorkRepo
	entryRepo adapter.EntryRepo
}

func NewApp(
	tx adapter.TX,
	workRepo adapter.WorkRepo,
	entryRepo adapter.EntryRepo) adapter.EntryApp {
	return &app{
		tx:        tx,
		workRepo:  workRepo,
		entryRepo: entryRepo,
	}
}

func (a *app) PublicBuild() adapter.EntryAppForPublic {
	return &publicApp{
		app: a,
	}
}
