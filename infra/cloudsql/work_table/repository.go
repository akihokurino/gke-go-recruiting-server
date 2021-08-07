package work_table

import (
	"context"
	"time"

	"gke-go-recruiting-server/infra/cloudsql"

	"gorm.io/gorm"

	"github.com/pkg/errors"

	"gke-go-recruiting-server/domain/work_domain"
	pb "gke-go-recruiting-server/proto/go/pb"

	"gke-go-recruiting-server/adapter"
	"gke-go-recruiting-server/domain"
)

func NewRepo() adapter.WorkRepo {
	return &repository{}
}

type repository struct {
}

func (r *repository) GetByStatusAndSEO(ctx context.Context, db *gorm.DB, status pb.Work_Status, forSEO bool) ([]*work_domain.Work, error) {
	var es []Entity
	if err := db.
		Joins("inner join departments on departments.id = works.department_id").
		Where("works.status = ?", int32(status)).
		Where("works.for_seo = ?", forSEO).
		Order("works.created_at DESC").
		Preload("Department").
		Preload("Department.Stations").
		Preload("Department.Stations.Line").
		Preload("Department.ImmediatelyWorkableCalendars").
		Preload("Images").
		Preload("Movies").
		Preload("Merits").
		Preload("WorkActivePlan").
		Find(&es).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*work_domain.Work, 0, len(es))
	for _, e := range es {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetByWillStart(ctx context.Context, db *gorm.DB, t time.Time) ([]*work_domain.Work, error) {
	var es []Entity
	if err := db.
		Joins("inner join departments on departments.id = works.department_id").
		Where("works.date_from <= ?", t).
		Where("works.date_to > ?", t).
		Where("works.status = ?", int32(pb.Work_Status_Reserved)).
		Where("departments.status = ?", int32(pb.Department_Status_OK)).
		Preload("Department").
		Preload("Department.Stations").
		Preload("Department.Stations.Line").
		Preload("Department.ImmediatelyWorkableCalendars").
		Preload("Images").
		Preload("Movies").
		Preload("Merits").
		Preload("WorkActivePlan").
		Find(&es).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*work_domain.Work, 0, len(es))
	for _, e := range es {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetByWillFinish(ctx context.Context, db *gorm.DB, t time.Time) ([]*work_domain.Work, error) {
	var es []Entity
	if err := db.
		Joins("inner join departments on departments.id = works.department_id").
		Where("works.date_to <= ?", t).
		Where("works.status = ?", int32(pb.Work_Status_Active)).
		Preload("Department").
		Preload("Department.Stations").
		Preload("Department.Stations.Line").
		Preload("Department.ImmediatelyWorkableCalendars").
		Preload("Images").
		Preload("Movies").
		Preload("Merits").
		Preload("WorkActivePlan").
		Find(&es).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*work_domain.Work, 0, len(es))
	for _, e := range es {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) filterQuery(db *gorm.DB, params adapter.WorkFilterParams) *gorm.DB {
	query := db.Model(&Entity{}).
		Joins("inner join departments on departments.id = works.department_id")

	if params.AgencyID.String() != "" {
		query = query.Where("departments.agency_id = ?", params.AgencyID.String())
	}

	if params.CompanyID.String() != "" {
		query = query.Where("departments.company_id = ?", params.CompanyID.String())
	}

	if params.DepartmentID.String() != "" {
		query = query.Where("works.department_id = ?", params.DepartmentID.String())
	}

	if params.DepartmentName != "" {
		query = query.Where("departments.name LIKE ?", "%"+params.DepartmentName+"%")
	}

	if params.WorkID.String() != "" {
		query = query.Where("works.id = ?", params.WorkID.String())
	}

	if params.SalesID.String() != "" {
		query = query.Where("departments.sales_id = ?", params.SalesID.String())
	}

	if params.BusinessCondition != pb.Department_BusinessCondition_Unknown {
		query = query.Where("departments.business_condition = ?", int32(params.BusinessCondition))
	}

	if params.WorkType != pb.Work_Type_Unknown {
		query = query.Where("works.work_type = ?", int32(params.WorkType))
	}

	if params.DateFromRange != nil {
		query = query.
			Where("works.date_from > ?", params.DateFromRange.From).
			Where("works.date_from < ?", params.DateFromRange.To)
	}

	if params.DateToRange != nil {
		query = query.
			Where("works.date_to > ?", params.DateToRange.From).
			Where("works.date_to < ?", params.DateToRange.To)
	}

	if params.Status != pb.Work_Status_Unknown {
		query = query.Where("works.status = ?", int32(params.Status))
	}

	if params.PrefID.String() != "" {
		query = query.Where("departments.pref_id = ?", params.PrefID.String())
	}

	query = query.Order("works.created_at DESC")

	return query
}

func (r *repository) GetByFilterWithPager(
	ctx context.Context,
	db *gorm.DB,
	pager *domain.Pager,
	params adapter.WorkFilterParams) ([]*work_domain.Work, error) {
	var es []Entity

	if err := r.filterQuery(db, params).
		Offset(pager.Offset()).
		Limit(pager.Limit()).
		Preload("Department").
		Preload("Department.Stations").
		Preload("Department.Stations.Line").
		Preload("Department.ImmediatelyWorkableCalendars").
		Preload("Images").
		Preload("Movies").
		Preload("Merits").
		Preload("WorkActivePlan").
		Find(&es).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*work_domain.Work, 0, len(es))
	for _, e := range es {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetCountByFilter(
	ctx context.Context,
	db *gorm.DB,
	params adapter.WorkFilterParams) (uint64, error) {
	var count int64

	if err := r.filterQuery(db, params).Count(&count).Error; err != nil {
		return 0, errors.WithStack(err)
	}

	return uint64(count), nil
}

func (r *repository) GetMulti(
	ctx context.Context,
	db *gorm.DB,
	ids []domain.WorkID) ([]*work_domain.Work, error) {
	var es []Entity
	if err := db.
		Joins("inner join departments on departments.id = works.department_id").
		Where("works.id IN (?)", ids).
		Preload("Department").
		Preload("Department.Stations").
		Preload("Department.Stations.Line").
		Preload("Department.ImmediatelyWorkableCalendars").
		Preload("Images").
		Preload("Movies").
		Preload("Merits").
		Preload("WorkActivePlan").
		Preload("ColorAD").
		Find(&es).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*work_domain.Work, 0, len(es))
	for _, e := range es {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) Exist(ctx context.Context, db *gorm.DB, id domain.WorkID) (bool, error) {
	var count int64
	if err := db.
		Model(&Entity{}).
		Where("id = ?", id.String()).
		Count(&count).Error; err != nil {
		return false, errors.WithStack(err)
	}

	return count > 0, nil
}

func (r *repository) ExistByDepartmentAndTypeAndTime(
	ctx context.Context,
	db *gorm.DB,
	departmentID domain.DepartmentID,
	workType pb.Work_Type,
	t time.Time) (bool, error) {
	var count int64
	if err := db.
		Model(&Entity{}).
		Where("department_id = ?", departmentID.String()).
		Where("work_type = ?", int32(workType)).
		Where("date_from <= ?", t).
		Where("date_to > ?", t).
		Count(&count).Error; err != nil {
		return false, errors.WithStack(err)
	}

	return count > 0, nil
}

func (r *repository) GetCountByActiveAndDepartment(
	ctx context.Context,
	db *gorm.DB,
	departmentID domain.DepartmentID) (uint64, error) {
	var count int64

	if err := db.Model(&Entity{}).
		Where("status = ?", int32(pb.Work_Status_Active)).
		Where("department_id = ?", departmentID.String()).
		Count(&count).Error; err != nil {
		return 0, errors.WithStack(err)
	}

	return uint64(count), nil
}

func (r *repository) GetCountByActiveAndDepartments(
	ctx context.Context,
	db *gorm.DB,
	departmentIDs []domain.DepartmentID) (uint64, error) {
	var count int64

	if err := db.Model(&Entity{}).
		Where("status = ?", int32(pb.Work_Status_Active)).
		Where("department_id IN (?)", departmentIDs).
		Count(&count).Error; err != nil {
		return 0, errors.WithStack(err)
	}

	return uint64(count), nil
}

func (r *repository) GetCountByActiveAndPref(
	ctx context.Context,
	db *gorm.DB,
	prefID domain.PrefID) (uint64, error) {
	var count int64

	if err := db.Model(&Entity{}).
		Joins("inner join departments on departments.id = works.department_id").
		Where("works.status = ?", int32(pb.Work_Status_Active)).
		Where("departments.pref_id = ?", prefID.String()).
		Count(&count).Error; err != nil {
		return 0, errors.WithStack(err)
	}

	return uint64(count), nil
}

func (r *repository) GetCountByActiveAndCity(
	ctx context.Context,
	db *gorm.DB,
	cityID domain.CityID) (uint64, error) {
	var count int64

	if err := db.Model(&Entity{}).
		Joins("inner join departments on departments.id = works.department_id").
		Where("works.status = ?", int32(pb.Work_Status_Active)).
		Where("departments.city_id = ?", cityID.String()).
		Count(&count).Error; err != nil {
		return 0, errors.WithStack(err)
	}

	return uint64(count), nil
}

func (r *repository) GetCountByActiveAndMArea(
	ctx context.Context,
	db *gorm.DB,
	mAreaID domain.MAreaID) (uint64, error) {
	var count int64

	if err := db.Model(&Entity{}).
		Joins("inner join departments on departments.id = works.department_id").
		Where("works.status = ?", int32(pb.Work_Status_Active)).
		Where("departments.m_area_id = ?", mAreaID.String()).
		Count(&count).Error; err != nil {
		return 0, errors.WithStack(err)
	}

	return uint64(count), nil
}

func (r *repository) GetCountByActiveAndSArea(
	ctx context.Context,
	db *gorm.DB,
	sAreaID domain.SAreaID) (uint64, error) {
	var count int64

	if err := db.Model(&Entity{}).
		Joins("inner join departments on departments.id = works.department_id").
		Where("works.status = ?", int32(pb.Work_Status_Active)).
		Where("departments.s_area_id = ?", sAreaID.String()).
		Count(&count).Error; err != nil {
		return 0, errors.WithStack(err)
	}

	return uint64(count), nil
}

func (r *repository) Get(ctx context.Context, db *gorm.DB, id domain.WorkID) (*work_domain.Work, error) {
	var e Entity
	if err := db.
		Joins("inner join departments on departments.id = works.department_id").
		Where("works.id = ?", id.String()).
		Preload("Department").
		Preload("Department.Stations").
		Preload("Department.Stations.Line").
		Preload("Department.ImmediatelyWorkableCalendars").
		Preload("Images").
		Preload("Movies").
		Preload("Merits").
		Preload("WorkActivePlan").
		First(&e).Error; err != nil {
		if cloudsql.IsNotFoundError(err) {
			return nil, domain.NewNotFoundErr()
		}
		return nil, errors.WithStack(err)
	}

	return e.toDomain(), nil
}

func (r *repository) Insert(ctx context.Context, db *gorm.DB, item *work_domain.Work) error {
	if err := db.Create(entityFrom(item)).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repository) Update(ctx context.Context, db *gorm.DB, item *work_domain.Work) error {
	if err := db.Save(entityFrom(item)).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repository) InsertMulti(ctx context.Context, db *gorm.DB, items []*work_domain.Work) error {
	for _, item := range items {
		if err := db.Create(entityFrom(item)).Error; err != nil {
			if cloudsql.IsDuplicateError(err) {
				continue
			}
			return errors.WithStack(err)
		}
	}

	return nil
}

func (r *repository) UpdateMulti(ctx context.Context, db *gorm.DB, items []*work_domain.Work) error {
	for _, item := range items {
		if err := db.Save(entityFrom(item)).Error; err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
