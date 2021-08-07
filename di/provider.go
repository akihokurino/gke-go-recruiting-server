// +build wireinject

package di

import (
	"fmt"
	"os"

	"gke-go-recruiting-server/infra/logger"

	"gke-go-recruiting-server/infra/datastore"

	"gke-go-recruiting-server/infra/cloudsql/department_station_table"

	"gke-go-recruiting-server/infra/cloudsql/work_active_plan_table"

	"gke-go-recruiting-server/infra/cloudsql/usage_statement"

	"gke-go-recruiting-server/application/contract_application"

	"gke-go-recruiting-server/handler/ope_handler"
	"gke-go-recruiting-server/infra/cloudsql/main_contract_table"

	"gke-go-recruiting-server/application/department_application"

	"gke-go-recruiting-server/application/company_application"

	"gke-go-recruiting-server/infra/cloudsql/administrator_table"

	"gke-go-recruiting-server/application/account_application"
	"gke-go-recruiting-server/infra/cloudsql/agency_account_table"
	"gke-go-recruiting-server/infra/cloudsql/agency_table"

	"gke-go-recruiting-server/handler/master_handler"

	"gke-go-recruiting-server/application/work_application"

	"gke-go-recruiting-server/infra/cloudsql/work_merit_table"
	"gke-go-recruiting-server/infra/firebase"

	"gke-go-recruiting-server/infra/algolia"
	"gke-go-recruiting-server/infra/algolia/work_index"

	"gke-go-recruiting-server/application/entry_application"

	"gke-go-recruiting-server/infra/cloudsql/entry_table"

	"gke-go-recruiting-server/handler/user_handler"
	"gke-go-recruiting-server/infra/cloudsql/work_movie_table"

	"gke-go-recruiting-server/adapter"
	"gke-go-recruiting-server/handler"
	"gke-go-recruiting-server/infra/cloudsql"
	"gke-go-recruiting-server/infra/cloudsql/city_table"
	"gke-go-recruiting-server/infra/cloudsql/company_table"
	"gke-go-recruiting-server/infra/cloudsql/department_image_table"
	"gke-go-recruiting-server/infra/cloudsql/department_table"
	"gke-go-recruiting-server/infra/cloudsql/line_table"
	"gke-go-recruiting-server/infra/cloudsql/region_table"
	"gke-go-recruiting-server/infra/cloudsql/work_image_table"
	"gke-go-recruiting-server/infra/cloudsql/work_table"

	"github.com/google/wire"
	"google.golang.org/api/option"
)

var providerSet = wire.NewSet(
	provideDB,
	provideFirebaseAppFactory,
	provideDataStoreFactory,
	provideAlgoliaClient,
	provideLogger,
	cloudsql.NewTX,
	city_table.NewRepo,
	line_table.NewRepo,
	region_table.NewRepo,
	company_table.NewRepo,
	department_table.NewRepo,
	department_image_table.NewRepo,
	department_station_table.NewRepo,
	work_table.NewRepo,
	work_index.NewRepo,
	work_image_table.NewRepo,
	work_movie_table.NewRepo,
	work_merit_table.NewRepo,
	work_active_plan_table.NewRepo,
	entry_table.NewRepo,
	agency_table.NewRepo,
	agency_account_table.NewRepo,
	firebase.NewRepo,
	administrator_table.NewRepo,
	main_contract_table.NewRepo,
	usage_statement.NewRepo,

	entry_application.NewApp,
	work_application.NewApp,
	account_application.NewApp,
	company_application.NewApp,
	department_application.NewApp,
	contract_application.NewApp,

	user_handler.NewCityQuery,
	user_handler.NewLineQuery,
	user_handler.NewRegionQuery,
	user_handler.NewWorkQuery,
	user_handler.NewEntryCommand,
	master_handler.NewAgencyQuery,
	master_handler.NewAccountQuery,
	master_handler.NewAccountCommand,
	master_handler.NewCompanyQuery,
	master_handler.NewDepartmentQuery,
	master_handler.NewDepartmentCommand,
	master_handler.NewWorkQuery,
	master_handler.NewWorkCommand,
	master_handler.NewEntryQuery,
	master_handler.NewProductQuery,
	master_handler.NewContractQuery,
	master_handler.NewContractCommand,
	master_handler.NewUsageStatementQuery,
	ope_handler.NewCompanyQuery,
	ope_handler.NewCompanyCommand,
	ope_handler.NewDepartmentQuery,
	ope_handler.NewDepartmentCommand,
	ope_handler.NewWorkQuery,
	ope_handler.NewWorkCommand,
	ope_handler.NewEntryQuery,
	ope_handler.NewAccountQuery,
	ope_handler.NewAccountCommand,
	ope_handler.NewProductQuery,
	ope_handler.NewContractQuery,
	ope_handler.NewContractCommand,
	ope_handler.NewUsageStatementQuery,
	ope_handler.NewLineQuery,
	ope_handler.NewCityQuery,
	ope_handler.NewRegionQuery,
	handler.NewAPIHandler,
	handler.NewBatchHandler,
	handler.NewContextProvider,
	handler.NewCros,
	handler.NewAdminAuthenticate,
	handler.NewAdminAuthorization,
	handler.NewAgencyAuthorization,
	handler.NewErrorConverter,
)

func provideLogger() adapter.CompositeLogger {
	return logger.NewLoggerWithMinLevel(adapter.LogLevelDebug)
}

func provideDB() adapter.DB {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/gke-go-recruiting-server?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)
	return cloudsql.NewDB(url)
}

func provideFirebaseAppFactory() adapter.FirebaseAppFactory {
	opt := option.WithCredentialsFile(os.Getenv("FIREBASE_CREDENTIALS"))
	return firebase.NewFirebaseAppFactory(opt)
}

func provideDataStoreFactory() adapter.DataStoreFactory {
	return datastore.NewDataStoreFactory(os.Getenv("PROJECT_ID"))
}

func provideAlgoliaClient() *algolia.Client {
	return algolia.NewClient(
		os.Getenv("ALGOLIA_APP_ID"),
		os.Getenv("ALGOLIA_API_KEY"),
		os.Getenv("APP_ENV"))
}

func ResolveAPIHandler() handler.APIHandler {
	panic(wire.Build(providerSet))
}

func ResolveBatchHandler() handler.BatchHandler {
	panic(wire.Build(providerSet))
}
