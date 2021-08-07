package adapter

import (
	"context"
	"net/url"
	"time"

	"gorm.io/gorm"

	"gke-go-recruiting-server/domain/contract_domain"

	"gke-go-recruiting-server/domain/work_domain"

	"gke-go-recruiting-server/domain/department_domain"

	"gke-go-recruiting-server/domain/company_domain"

	"gke-go-recruiting-server/domain/account_domain"

	"gke-go-recruiting-server/domain"
	pb "gke-go-recruiting-server/proto/go/pb"
)

type AccountApp interface {
	MasterBuild(me *account_domain.Administrator) AccountAppForMaster
	OperationBuild(me *account_domain.AgencyAccount) AccountAppForOperation
}

type AdministratorParams struct {
	Email    string
	Password string
	Name     string
}

type AgencyAccountParams struct {
	Email    string
	Password string
	Name     string
	NameKana string
}

type AccountAppForMaster interface {
	CreateAdministrator(
		ctx context.Context,
		db *gorm.DB,
		accountParams AdministratorParams,
		now time.Time) (*account_domain.Administrator, error)
	CreateAgencyAccount(
		ctx context.Context,
		db *gorm.DB,
		accountParams AgencyAccountParams,
		agencyID domain.AgencyID,
		now time.Time) (*account_domain.AgencyAccount, error)
	DeleteAgencyAccount(
		ctx context.Context,
		db *gorm.DB,
		agencyID domain.AgencyID,
		accountID domain.FirebaseUserID) error
}

type AccountAppForOperation interface {
	CreateAgencyAccount(
		ctx context.Context,
		db *gorm.DB,
		accountParams AgencyAccountParams,
		now time.Time) (*account_domain.AgencyAccount, error)
	DeleteAgencyAccount(
		ctx context.Context,
		db *gorm.DB,
		accountID domain.FirebaseUserID) error
}

type EntryApp interface {
	PublicBuild() EntryAppForPublic
}

type EntryParams struct {
	FullName               string
	FullNameKana           string
	Birthdate              time.Time
	Gender                 pb.User_Gender
	PhoneNumber            string
	Email                  string
	Question               string
	Category               pb.User_Category
	PrefID                 domain.PrefID
	PreferredContactMethod pb.Entry_PreferredContactMethod
	PreferredContactTime   string
}

type EntryAppForPublic interface {
	Entry(
		ctx context.Context,
		db *gorm.DB,
		workID domain.WorkID,
		entryParams EntryParams,
		now time.Time) error
}

type WorkImageParams struct {
	URL       url.URL
	ViewOrder uint64
	Comment   string
}

type WorkParams struct {
	WorkType pb.Work_Type
	JobCode  pb.Work_Job
	Title    string
	Content  string
	DateFrom time.Time
	DateTo   time.Time
}

type WorkApp interface {
	MasterBuild(me *account_domain.Administrator) WorkAppForMaster
	OperationBuild(me *account_domain.AgencyAccount) WorkAppForOperation
	BatchBuild() WorkAppForBatch
}

type WorkAppForMaster interface {
	Accept(ctx context.Context, db *gorm.DB, workID domain.WorkID) error
	Deny(ctx context.Context, db *gorm.DB, workID domain.WorkID) error
}

type WorkAppForOperation interface {
	Create(
		ctx context.Context,
		db *gorm.DB,
		departmentID domain.DepartmentID,
		workParams WorkParams,
		imageParams []WorkImageParams,
		movieURLs []url.URL,
		merits []pb.Work_Merit,
		now time.Time) (*work_domain.Work, error)
	Update(
		ctx context.Context,
		db *gorm.DB,
		workID domain.WorkID,
		workParams WorkParams,
		imageParams []WorkImageParams,
		movieURLs []url.URL,
		merits []pb.Work_Merit,
		now time.Time) (*work_domain.Work, error)
	EarlyFinish(ctx context.Context, db *gorm.DB, workID domain.WorkID) error
}

type WorkAppForBatch interface {
	Proceed(ctx context.Context, db *gorm.DB, now time.Time) error
	ReOrder(ctx context.Context, db *gorm.DB) error
	ReIndex(ctx context.Context, db *gorm.DB) error
}

type CompanyParams struct {
	Rank        pb.Company_Rank
	RankType    pb.Company_RankType
	Name        string
	NameKana    string
	PostalCode  string
	PrefID      domain.PrefID
	Address     string
	Building    string
	PhoneNumber string
}

type CompanyApp interface {
	OperationBuild(me *account_domain.AgencyAccount) CompanyAppForOperation
}

type CompanyAppForOperation interface {
	Create(
		ctx context.Context,
		db *gorm.DB,
		companyParams CompanyParams,
		now time.Time) (*company_domain.Company, error)
}

type DepartmentParams struct {
	Name              string
	BusinessCondition pb.Department_BusinessCondition
	PostalCode        string
	PrefID            domain.PrefID
	CityID            domain.CityID
	Address           string
	Building          string
	PhoneNumber       string
	MAreaID           domain.MAreaID
	SAreaID           domain.SAreaID
	Latitude          float64
	Longitude         float64
}

type DepartmentApp interface {
	MasterBuild(me *account_domain.Administrator) DepartmentAppForMaster
	OperationBuild(me *account_domain.AgencyAccount) DepartmentAppForOperation
}

type DepartmentAppForMaster interface {
	Accept(ctx context.Context, db *gorm.DB, departmentID domain.DepartmentID) error
	Deny(ctx context.Context, db *gorm.DB, departmentID domain.DepartmentID) error
	Update(
		ctx context.Context,
		db *gorm.DB,
		departmentID domain.DepartmentID,
		departmentParams DepartmentParams,
		imageURLs []url.URL,
		lineIDs []domain.LineID,
		now time.Time) (*department_domain.Department, error)
	UpdateSales(
		ctx context.Context,
		db *gorm.DB,
		departmentID domain.DepartmentID,
		salesID domain.FirebaseUserID) error
}

type DepartmentAppForOperation interface {
	Create(
		ctx context.Context,
		db *gorm.DB,
		companyID domain.CompanyID,
		salesID domain.FirebaseUserID,
		departmentParams DepartmentParams,
		imageURLs []url.URL,
		lineIDs []domain.LineID,
		now time.Time) (*department_domain.Department, error)
	Update(
		ctx context.Context,
		db *gorm.DB,
		departmentID domain.DepartmentID,
		departmentParams DepartmentParams,
		imageURLs []url.URL,
		lineIDs []domain.LineID,
		now time.Time) (*department_domain.Department, error)
	UpdateSales(
		ctx context.Context,
		db *gorm.DB,
		departmentID domain.DepartmentID,
		salesID domain.FirebaseUserID) error
}

type ContractApp interface {
	MasterBuild(me *account_domain.Administrator) ContractAppForMaster
	OperationBuild(me *account_domain.AgencyAccount) ContractAppForOperation
}

type ContractAppForMaster interface {
	AcceptMainContract(ctx context.Context, db *gorm.DB, contractID domain.MainContractID, now time.Time) error
	DenyMainContract(ctx context.Context, db *gorm.DB, contractID domain.MainContractID) error
}

type MainContractParams struct {
	Plan     pb.MainProduct_Plan
	DateFrom time.Time
}

type ContractAppForOperation interface {
	CreateMainContract(
		ctx context.Context,
		db *gorm.DB,
		departmentID domain.DepartmentID,
		contractParams MainContractParams,
		now time.Time) (*contract_domain.Main, error)
}
