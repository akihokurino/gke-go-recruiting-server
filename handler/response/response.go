package response

import (
	"gke-go-recruiting-server/domain"
	"gke-go-recruiting-server/domain/account_domain"
	"gke-go-recruiting-server/domain/agency_domain"

	"gke-go-recruiting-server/domain/company_domain"
	"gke-go-recruiting-server/domain/contract_domain"
	"gke-go-recruiting-server/domain/department_domain"
	"gke-go-recruiting-server/domain/entry_domain"
	"gke-go-recruiting-server/domain/master_domain"
	"gke-go-recruiting-server/domain/product_domain"
	"gke-go-recruiting-server/domain/statement_domain"
	"gke-go-recruiting-server/domain/work_domain"
	pb "gke-go-recruiting-server/proto/go/pb"
)

func PrefectureFrom(d *master_domain.Prefecture, workCount uint64) *pb.Prefecture {
	return &pb.Prefecture{
		Id:        d.ID.String(),
		Name:      d.Name,
		WorkCount: workCount,
	}
}

func CityFrom(d *master_domain.City, workCount uint64) *pb.City {
	return &pb.City{
		Id:           d.ID.String(),
		Name:         d.Name,
		NameKana:     d.NameKana,
		AreaName:     d.AreaName,
		AreaNameKana: d.AreaNameKana,
		WorkCount:    workCount,
	}
}

func RailFrom(d *master_domain.Rail, workCount uint64) *pb.Rail {
	return &pb.Rail{
		Id:                d.ID.String(),
		Name_1:            d.Name1,
		NameKana_1:        d.NameKana1,
		Name_2:            d.Name2,
		NameKana_2:        d.NameKana2,
		KindName:          d.Company.KindName,
		CompanyName_1:     d.Company.Name1,
		CompanyNameKana_1: d.Company.NameKana1,
		CompanyName_2:     d.Company.Name2,
		WorkCount:         workCount,
	}
}

func RailCompanyFrom(d *master_domain.RailCompany) *pb.RailCompany {
	return &pb.RailCompany{
		Name_1:     d.Name1,
		NameKana_1: d.NameKana1,
		Name_2:     d.Name2,
		KindName:   d.KindName,
	}
}

func StationFrom(d *master_domain.Line, workCount uint64) *pb.Station {
	return &pb.Station{
		Id:        d.StationID.String(),
		PrefId:    d.PrefID.String(),
		Name:      d.StationName,
		NameKana:  d.StationNameKana,
		StopOrder: d.StopOrder,
		WorkCount: workCount,
	}
}

func LineFrom(d *master_domain.Line) *pb.Line {
	return &pb.Line{
		Id:          d.ID.String(),
		RailName:    d.Rail.Name1,
		StationName: d.StationName,
		Latitude:    d.Latitude,
		Longitude:   d.Longitude,
	}
}

func MAreaFrom(d *master_domain.MArea, workCount uint64) *pb.Area {
	return &pb.Area{
		Id:        d.ID.String(),
		Name:      d.Name,
		WorkCount: workCount,
	}
}

func SAreaFrom(d *master_domain.SArea, workCount uint64) *pb.Area {
	return &pb.Area{
		Id:        d.ID.String(),
		Name:      d.Name,
		WorkCount: workCount,
	}
}

func WorkFrom(d *work_domain.Work, isNew bool, prefs []*master_domain.Prefecture) *pb.Work {
	stations := make([]*pb.DepartmentStation, 0, len(d.With.Department.Stations))
	for _, station := range d.With.Department.Stations {
		stations = append(stations, DepartmentStationFrom(station))
	}

	images := make([]*pb.WorkImage, 0, len(d.Images))
	for _, image := range d.Images {
		images = append(images, WorkImageFrom(image))
	}

	movies := make([]*pb.WorkMovie, 0, len(d.Movies))
	for _, movie := range d.Movies {
		movies = append(movies, WorkMovieFrom(movie))
	}

	merits := make([]pb.Work_Merit, 0, len(d.Merits))
	for _, merit := range d.Merits {
		merits = append(merits, merit.Value)
	}

	return &pb.Work{
		Id:           d.ID.String(),
		DepartmentId: d.DepartmentID.String(),
		Status:       d.Status,
		WorkType:     d.WorkType,
		JobCode:      d.JobCode,
		Title:        d.Title,
		Content:      d.Content,
		DateFrom:     domain.JSTStringFromDateTime(d.DateRange.From),
		DateTo:       domain.JSTStringFromDateTime(d.DateRange.To),
		CreatedAt:    domain.JSTStringFromDateTime(d.CreatedAt),
		UpdatedAt:    domain.JSTStringFromDateTime(d.UpdatedAt),
		Department:   DepartmentFrom(d.With.Department, prefs),
		Images:       images,
		Movies:       movies,
		Merits:       merits,
		IsNew:        isNew,
	}
}

func WorkImageFrom(d *work_domain.Image) *pb.WorkImage {
	return &pb.WorkImage{
		Id:        d.ID,
		WorkId:    d.WorkID.String(),
		Url:       d.URL.String(),
		ViewOrder: d.ViewOrder,
		Comment:   d.Comment,
	}
}

func WorkMovieFrom(d *work_domain.Movie) *pb.WorkMovie {
	return &pb.WorkMovie{
		Id:     d.ID,
		WorkId: d.WorkID.String(),
		Url:    d.URL.String(),
	}
}

func DepartmentFrom(d *department_domain.Department, prefs []*master_domain.Prefecture) *pb.Department {
	images := make([]*pb.DepartmentImage, 0, len(d.Images))
	for _, image := range d.Images {
		images = append(images, DepartmentImageFrom(image))
	}

	stations := make([]*pb.DepartmentStation, 0, len(d.Stations))
	for _, station := range d.Stations {
		stations = append(stations, DepartmentStationFrom(station))
	}

	prefName := ""
	for _, pref := range prefs {
		if pref.ID == d.PrefID {
			prefName = pref.Name
			break
		}
	}

	return &pb.Department{
		Id:                d.ID.String(),
		AgencyId:          d.AgencyID.String(),
		CompanyId:         d.CompanyID.String(),
		SalesId:           d.SalesID.String(),
		SalesName:         d.Sales.Name,
		Status:            d.Status,
		Name:              d.Name,
		BusinessCondition: d.BusinessCondition,
		PostalCode:        d.PostalCode,
		PrefId:            d.PrefID.String(),
		PrefName:          prefName,
		CityId:            d.CityID.String(),
		Address:           d.Address,
		Building:          d.Building,
		PhoneNumber:       d.PhoneNumber,
		MAreaId:           d.Location.MAreaID.String(),
		SAreaId:           d.Location.SAreaID.String(),
		Latitude:          d.Location.Latitude,
		Longitude:         d.Location.Longitude,
		Images:            images,
		Stations:          stations,
	}
}

func DepartmentImageFrom(d *department_domain.Image) *pb.DepartmentImage {
	return &pb.DepartmentImage{
		Id:           d.ID,
		DepartmentId: d.DepartmentID.String(),
		Url:          d.URL.String(),
	}
}

func DepartmentStationFrom(d *department_domain.Station) *pb.DepartmentStation {
	return &pb.DepartmentStation{
		Id:           d.ID,
		DepartmentId: d.DepartmentID.String(),
		LineId:       d.LineID.String(),
		RailId:       d.With.Line.Rail.ID.String(),
		RailName:     d.With.Line.Rail.Name1,
		StationId:    d.With.Line.StationID.String(),
		StationName:  d.With.Line.StationName,
	}
}

func AgencyAccountFrom(d *account_domain.AgencyAccount) *pb.AgencyAccount {
	return &pb.AgencyAccount{
		Id:       d.ID.String(),
		AgencyId: d.AgencyID.String(),
		Email:    d.Email,
		Name:     d.Name,
	}
}

func AdministratorFrom(d *account_domain.Administrator) *pb.Administrator {
	return &pb.Administrator{
		Id:    d.ID.String(),
		Email: d.Email,
		Name:  d.Name,
	}
}

func CompanyFrom(d *company_domain.Company) *pb.Company {
	return &pb.Company{
		Id:          d.ID.String(),
		AgencyId:    d.AgencyID.String(),
		Status:      d.Status,
		RankType:    d.RankType,
		Rank:        d.Rank,
		Name:        d.Name,
		NameKana:    d.NameKana,
		PostalCode:  d.PostalCode,
		PrefId:      d.PrefID.String(),
		Address:     d.Address,
		Building:    d.Building,
		PhoneNumber: d.PhoneNumber,
		CreatedAt:   domain.JSTStringFromDateTime(d.CreatedAt),
		UpdatedAt:   domain.JSTStringFromDateTime(d.UpdatedAt),
	}
}

func EntryFrom(d *entry_domain.Entry) *pb.Entry {
	category := pb.User_Category_Unknown
	if d.Category != nil {
		category = *d.Category
	}

	prefID := ""
	if d.PrefID != nil {
		tmp := *d.PrefID
		prefID = tmp.String()
	}

	preferredContactMethod := pb.Entry_PreferredContactMethod_Unknown
	if d.PreferredContactMethod != nil {
		preferredContactMethod = *d.PreferredContactMethod
	}

	preferredContactTime := ""
	if d.PreferredContactTime != nil {
		preferredContactTime = *d.PreferredContactTime
	}

	return &pb.Entry{
		Id:                     d.ID.String(),
		DepartmentId:           d.DepartmentID.String(),
		WorkId:                 d.WorkID.String(),
		Fullname:               d.FullName,
		FullnameKana:           d.FullNameKana,
		Birthdate:              domain.JSTStringFromDateTime(d.Birthdate),
		Gender:                 d.Gender,
		PhoneNumber:            d.PhoneNumber,
		Email:                  d.Email,
		Question:               d.Question,
		Category:               category,
		PrefId:                 prefID,
		PreferredContactMethod: preferredContactMethod,
		PreferredContactTime:   preferredContactTime,
		Status:                 d.Status,
		CreatedAt:              domain.JSTStringFromDateTime(d.CreatedAt),
		UpdatedAt:              domain.JSTStringFromDateTime(d.UpdatedAt),
	}
}

func AgencyFrom(d *agency_domain.Agency) *pb.Agency {
	return &pb.Agency{
		Id:         d.ID.String(),
		Name:       d.Name,
		NameKana:   d.NameKana,
		PostalCode: d.PostalCode,
		PrefId:     d.PrefID.String(),
		Address:    d.Address,
		CreatedAt:  domain.JSTStringFromDateTime(d.CreatedAt),
		UpdatedAt:  domain.JSTStringFromDateTime(d.UpdatedAt),
	}
}

func MainProductFrom(d *product_domain.Main) *pb.MainProduct {
	return &pb.MainProduct{
		Plan:  d.Plan,
		Price: d.Price,
	}
}

func MainContractFrom(d *contract_domain.Main) *pb.MainContract {
	return &pb.MainContract{
		Id:           d.ID.String(),
		DepartmentId: d.DepartmentID.String(),
		Status:       d.Status,
		Plan:         d.Plan,
		DateFrom:     domain.JSTStringFromDateTime(d.DateRange.From),
		DateTo:       domain.JSTStringFromDateTime(d.DateRange.To),
		Price:        d.Price,
		CreatedAt:    domain.JSTStringFromDateTime(d.CreatedAt),
		UpdatedAt:    domain.JSTStringFromDateTime(d.UpdatedAt),
	}
}

func UsageStatementFrom(d *statement_domain.Usage) *pb.UsageStatement {
	return &pb.UsageStatement{
		Id:           d.ID.String(),
		DepartmentId: d.DepartmentID.String(),
		Main:         MainContractFrom(d.With.MainContract),
		Price:        d.Price,
		CreatedAt:    domain.JSTStringFromDateTime(d.CreatedAt),
	}
}
