package adapter

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gke-go-recruiting-server/domain/statement_domain"

	"gke-go-recruiting-server/domain/contract_domain"

	"gke-go-recruiting-server/domain/account_domain"
	"gke-go-recruiting-server/domain/agency_domain"

	pb "gke-go-recruiting-server/proto/go/pb"

	"gke-go-recruiting-server/domain/entry_domain"

	"gke-go-recruiting-server/domain/master_domain"

	"gke-go-recruiting-server/domain/work_domain"

	"gke-go-recruiting-server/domain/department_domain"

	"gke-go-recruiting-server/domain/company_domain"

	"gke-go-recruiting-server/domain"
)

type DB func(ctx context.Context) *gorm.DB

type TX func(db *gorm.DB, fn func(db *gorm.DB) error) error

type CompanyFilterParams struct {
	AgencyID    domain.AgencyID
	CompanyID   domain.CompanyID
	CompanyName string
}

type CompanyRepo interface {
	GetByFilterWithPager(ctx context.Context, db *gorm.DB, pager *domain.Pager, params CompanyFilterParams) ([]*company_domain.Company, error)
	GetCountByFilter(ctx context.Context, db *gorm.DB, params CompanyFilterParams) (uint64, error)
	Get(ctx context.Context, db *gorm.DB, id domain.CompanyID) (*company_domain.Company, error)
	Exist(ctx context.Context, db *gorm.DB, id domain.CompanyID) (bool, error)
	Insert(ctx context.Context, db *gorm.DB, item *company_domain.Company) error
	InsertMulti(ctx context.Context, db *gorm.DB, items []*company_domain.Company) error
	Update(ctx context.Context, db *gorm.DB, item *company_domain.Company) error
}

type DepartmentFilterParams struct {
	AgencyID       domain.AgencyID
	CompanyID      domain.CompanyID
	DepartmentID   domain.DepartmentID
	DepartmentName string
	SalesID        domain.FirebaseUserID
	Status         pb.Department_Status
	PhoneNumber    string
}

type DepartmentRepo interface {
	GetByFilterWithPager(ctx context.Context, db *gorm.DB, pager *domain.Pager, params DepartmentFilterParams) ([]*department_domain.Department, error)
	GetCountByFilter(ctx context.Context, db *gorm.DB, params DepartmentFilterParams) (uint64, error)
	Get(ctx context.Context, db *gorm.DB, id domain.DepartmentID) (*department_domain.Department, error)
	Exist(ctx context.Context, db *gorm.DB, id domain.DepartmentID) (bool, error)
	Insert(ctx context.Context, db *gorm.DB, item *department_domain.Department) error
	InsertMulti(ctx context.Context, db *gorm.DB, items []*department_domain.Department) error
	Update(ctx context.Context, db *gorm.DB, item *department_domain.Department) error
}

type DepartmentImageRepo interface {
	InsertMulti(ctx context.Context, db *gorm.DB, items []*department_domain.Image) error
	DeleteByDepartment(ctx context.Context, db *gorm.DB, departmentID domain.DepartmentID) error
}

type DepartmentStationRepo interface {
	GetByRail(ctx context.Context, db *gorm.DB, railID domain.RailID) ([]*department_domain.Station, error)
	GetByStation(ctx context.Context, db *gorm.DB, stationID domain.StationID) ([]*department_domain.Station, error)
	InsertMulti(ctx context.Context, db *gorm.DB, items []*department_domain.Station) error
	DeleteByDepartment(ctx context.Context, db *gorm.DB, departmentID domain.DepartmentID) error
}

type WorkFilterParams struct {
	AgencyID          domain.AgencyID
	CompanyID         domain.CompanyID
	DepartmentID      domain.DepartmentID
	DepartmentName    string
	WorkID            domain.WorkID
	SalesID           domain.FirebaseUserID
	BusinessCondition pb.Department_BusinessCondition
	WorkType          pb.Work_Type
	DateFromRange     *domain.DateRange
	DateToRange       *domain.DateRange
	PrefID            domain.PrefID
	Status            pb.Work_Status
}

type WorkRepo interface {
	GetByStatusAndSEO(ctx context.Context, db *gorm.DB, status pb.Work_Status, forSEO bool) ([]*work_domain.Work, error)
	GetByWillStart(ctx context.Context, db *gorm.DB, t time.Time) ([]*work_domain.Work, error)
	GetByWillFinish(ctx context.Context, db *gorm.DB, t time.Time) ([]*work_domain.Work, error)
	GetByFilterWithPager(ctx context.Context, db *gorm.DB, pager *domain.Pager, params WorkFilterParams) ([]*work_domain.Work, error)
	GetCountByFilter(ctx context.Context, db *gorm.DB, params WorkFilterParams) (uint64, error)
	GetMulti(ctx context.Context, db *gorm.DB, ids []domain.WorkID) ([]*work_domain.Work, error)
	Exist(ctx context.Context, db *gorm.DB, id domain.WorkID) (bool, error)
	ExistByDepartmentAndTypeAndTime(
		ctx context.Context,
		db *gorm.DB,
		departmentID domain.DepartmentID,
		workType pb.Work_Type,
		t time.Time) (bool, error)
	GetCountByActiveAndDepartment(ctx context.Context, db *gorm.DB, departmentID domain.DepartmentID) (uint64, error)
	GetCountByActiveAndDepartments(ctx context.Context, db *gorm.DB, departmentIDs []domain.DepartmentID) (uint64, error)
	GetCountByActiveAndPref(ctx context.Context, db *gorm.DB, prefID domain.PrefID) (uint64, error)
	GetCountByActiveAndCity(ctx context.Context, db *gorm.DB, cityID domain.CityID) (uint64, error)
	GetCountByActiveAndMArea(ctx context.Context, db *gorm.DB, mAreaID domain.MAreaID) (uint64, error)
	GetCountByActiveAndSArea(ctx context.Context, db *gorm.DB, sAreaID domain.SAreaID) (uint64, error)
	Get(ctx context.Context, db *gorm.DB, id domain.WorkID) (*work_domain.Work, error)
	Insert(ctx context.Context, db *gorm.DB, item *work_domain.Work) error
	Update(ctx context.Context, db *gorm.DB, item *work_domain.Work) error
	InsertMulti(ctx context.Context, db *gorm.DB, items []*work_domain.Work) error
	UpdateMulti(ctx context.Context, db *gorm.DB, items []*work_domain.Work) error
}

type WorkActivePlanRepo interface {
	GetAll(ctx context.Context, db *gorm.DB) ([]*contract_domain.ActivePlan, error)
	Upsert(
		ctx context.Context,
		db *gorm.DB,
		id domain.WorkID,
		mainContractID domain.MainContractID,
		publishedOrder int) error
	UpdateMulti(ctx context.Context, db *gorm.DB, items []*contract_domain.ActivePlan) error
	Delete(ctx context.Context, db *gorm.DB, id domain.WorkID) error
}

type WorkImageRepo interface {
	InsertMulti(ctx context.Context, db *gorm.DB, items []*work_domain.Image) error
	DeleteByWork(ctx context.Context, db *gorm.DB, workID domain.WorkID) error
}

type WorkMovieRepo interface {
	InsertMulti(ctx context.Context, db *gorm.DB, items []*work_domain.Movie) error
	DeleteByWork(ctx context.Context, db *gorm.DB, workID domain.WorkID) error
}

type WorkMeritRepo interface {
	InsertMulti(ctx context.Context, db *gorm.DB, items []*work_domain.Merit) error
	DeleteByWork(ctx context.Context, db *gorm.DB, workID domain.WorkID) error
}

type CityRepo interface {
	GetAllPrefecture(ctx context.Context, db *gorm.DB) ([]*master_domain.Prefecture, error)
	GetByPrefecture(ctx context.Context, db *gorm.DB, prefID domain.PrefID) ([]*master_domain.City, error)
	Exist(ctx context.Context, db *gorm.DB, id domain.CityID) (bool, error)
	InsertMulti(ctx context.Context, db *gorm.DB, items []*master_domain.City) error
}

type LineRepo interface {
	GetRailByPrefecture(ctx context.Context, db *gorm.DB, prefID domain.PrefID) ([]*master_domain.Rail, error)
	GetRailByCompany(ctx context.Context, db *gorm.DB, companyName string) ([]*master_domain.Rail, error)
	GetByRail(ctx context.Context, db *gorm.DB, railID domain.RailID) ([]*master_domain.Line, error)
	GetRailCompanyByPrefecture(ctx context.Context, db *gorm.DB, prefID domain.PrefID) ([]*master_domain.RailCompany, error)
	GetByDistance(
		ctx context.Context,
		db *gorm.DB,
		latitude float64,
		longitude float64,
		distanceKM uint64,
	) ([]*master_domain.Line, error)
	Exist(ctx context.Context, db *gorm.DB, id domain.LineID) (bool, error)
	InsertMulti(ctx context.Context, db *gorm.DB, items []*master_domain.Line) error
}

type RegionRepo interface {
	GetMAreaByLArea(ctx context.Context, db *gorm.DB, lArea domain.LAreaID) ([]*master_domain.MArea, error)
	GetSAreaByMArea(ctx context.Context, db *gorm.DB, mArea domain.MAreaID) ([]*master_domain.SArea, error)
}

type EntryFilterParams struct {
	AgencyID          domain.AgencyID
	CompanyID         domain.CompanyID
	DepartmentID      domain.DepartmentID
	DepartmentName    string
	WorkID            domain.WorkID
	DateRange         *domain.DateRange
	SalesID           domain.FirebaseUserID
	BusinessCondition pb.Department_BusinessCondition
	PrefID            domain.PrefID
	Status            pb.Entry_Status
}

type EntryRepo interface {
	GetByFilterWithPager(ctx context.Context, db *gorm.DB, pager *domain.Pager, params EntryFilterParams) ([]*entry_domain.Entry, error)
	GetCountByFilter(ctx context.Context, db *gorm.DB, params EntryFilterParams) (uint64, error)
	GetCountByInProgressAndDepartment(ctx context.Context, db *gorm.DB, departmentID domain.DepartmentID) (uint64, error)
	Insert(ctx context.Context, db *gorm.DB, item *entry_domain.Entry) error
}

type AgencyAccountRepo interface {
	GetAll(ctx context.Context, db *gorm.DB) ([]*account_domain.AgencyAccount, error)
	GetByAgency(ctx context.Context, db *gorm.DB, agencyID domain.AgencyID) ([]*account_domain.AgencyAccount, error)
	Get(ctx context.Context, db *gorm.DB, id domain.FirebaseUserID) (*account_domain.AgencyAccount, error)
	Insert(ctx context.Context, db *gorm.DB, item *account_domain.AgencyAccount) error
	InsertMulti(ctx context.Context, db *gorm.DB, items []*account_domain.AgencyAccount) error
	Delete(ctx context.Context, db *gorm.DB, id domain.FirebaseUserID) error
}

type AgencyRepo interface {
	GetAll(ctx context.Context, db *gorm.DB) ([]*agency_domain.Agency, error)
	Get(ctx context.Context, db *gorm.DB, id domain.AgencyID) (*agency_domain.Agency, error)
	InsertMulti(ctx context.Context, db *gorm.DB, items []*agency_domain.Agency) error
}

type AdministratorRepo interface {
	Get(ctx context.Context, db *gorm.DB, id domain.FirebaseUserID) (*account_domain.Administrator, error)
	Insert(ctx context.Context, db *gorm.DB, item *account_domain.Administrator) error
}

type MainContractFilterParams struct {
	AgencyID       domain.AgencyID
	Status         pb.MainContract_Status
	CompanyID      domain.CompanyID
	DepartmentID   domain.DepartmentID
	DepartmentName string
	DateRange      *domain.DateRange
	Plan           pb.MainProduct_Plan
	SalesID        domain.FirebaseUserID
}

type MainContractRepo interface {
	GetByFilterWithPager(
		ctx context.Context,
		db *gorm.DB,
		pager *domain.Pager,
		params MainContractFilterParams) ([]*contract_domain.Main, error)
	GetCountByFilter(
		ctx context.Context,
		db *gorm.DB,
		params MainContractFilterParams) (uint64, error)
	GetByActiveAndDepartmentAndTime(
		ctx context.Context,
		db *gorm.DB,
		departmentID domain.DepartmentID,
		t time.Time) (*contract_domain.Main, error)
	ExistByDepartmentAndTime(ctx context.Context, db *gorm.DB, departmentID domain.DepartmentID, t time.Time) (bool, error)
	Get(ctx context.Context, db *gorm.DB, id domain.MainContractID) (*contract_domain.Main, error)
	Insert(ctx context.Context, db *gorm.DB, item *contract_domain.Main) error
	InsertMulti(ctx context.Context, db *gorm.DB, items []*contract_domain.Main) error
	Update(ctx context.Context, db *gorm.DB, item *contract_domain.Main) error
}

type UsageStatementFilterParams struct {
	AgencyID         domain.AgencyID
	CompanyID        domain.CompanyID
	DepartmentID     domain.DepartmentID
	DepartmentName   string
	MainProductPlan  pb.MainProduct_Plan
	DateRange        *domain.DateRange
	ExcludeFree      bool
	UsageStatementID domain.UsageStatementID
}

type UsageStatementRepo interface {
	GetByFilterWithPager(
		ctx context.Context,
		db *gorm.DB,
		pager *domain.Pager,
		params UsageStatementFilterParams) ([]*statement_domain.Usage, error)
	GetCountByFilter(ctx context.Context, db *gorm.DB, params UsageStatementFilterParams) (uint64, error)
	Insert(ctx context.Context, db *gorm.DB, item *statement_domain.Usage) error
}
