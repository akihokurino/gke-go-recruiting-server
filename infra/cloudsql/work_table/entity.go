package work_table

import (
	"database/sql"
	"time"

	"gke-go-recruiting-server/domain/contract_domain"

	"gke-go-recruiting-server/domain/department_domain"

	"gke-go-recruiting-server/infra/cloudsql/work_active_plan_table"

	"gke-go-recruiting-server/infra/cloudsql/work_merit_table"

	"gke-go-recruiting-server/domain/work_domain"

	"gke-go-recruiting-server/infra/cloudsql/work_movie_table"

	"gke-go-recruiting-server/infra/cloudsql/work_image_table"

	"gke-go-recruiting-server/infra/cloudsql/department_table"

	"github.com/guregu/null"

	"gke-go-recruiting-server/domain"
	pb "gke-go-recruiting-server/proto/go/pb"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

func (e *Entity) TableName() string {
	return "works"
}

type Entity struct {
	ID           string    `gorm:"column:id;primary_key"`
	DepartmentID string    `gorm:"column:department_id"`
	Status       int32     `gorm:"column:status"`
	WorkType     int32     `gorm:"column:work_type"`
	JobCode      int32     `gorm:"column:job_code"`
	Title        string    `gorm:"column:title"`
	Content      string    `gorm:"column:content"`
	DateFrom     time.Time `gorm:"column:date_from"`
	DateTo       time.Time `gorm:"column:date_to"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`

	Department     *department_table.Entity       `gorm:"PRELOAD:false;foreignkey:department_id"`
	Images         []*work_image_table.Entity     `gorm:"PRELOAD:false;foreignkey:work_id"`
	Movies         []*work_movie_table.Entity     `gorm:"PRELOAD:false;foreignkey:work_id"`
	Merits         []*work_merit_table.Entity     `gorm:"PRELOAD:false;foreignkey:work_id"`
	WorkActivePlan *work_active_plan_table.Entity `gorm:"PRELOAD:false;foreignkey:work_id"`
}

func (e *Entity) toDomain() *work_domain.Work {
	var department *department_domain.Department
	if e.Department != nil {
		department = e.Department.ToDomain()
	}

	images := make([]*work_domain.Image, 0, len(e.Images))
	for _, image := range e.Images {
		images = append(images, image.ToDomain())
	}

	movies := make([]*work_domain.Movie, 0, len(e.Movies))
	for _, movie := range e.Movies {
		movies = append(movies, movie.ToDomain())
	}

	merits := make([]*work_domain.Merit, 0, len(e.Merits))
	for _, merit := range e.Merits {
		merits = append(merits, merit.ToDomain())
	}

	var activePlan *contract_domain.ActivePlan
	if e.WorkActivePlan != nil {
		activePlan = e.WorkActivePlan.ToDomain()
	}

	return &work_domain.Work{
		ID:           domain.WorkID(e.ID),
		DepartmentID: domain.DepartmentID(e.DepartmentID),
		Status:       pb.Work_Status(e.Status),
		WorkType:     pb.Work_Type(e.WorkType),
		JobCode:      pb.Work_Job(e.JobCode),
		Title:        e.Title,
		Content:      e.Content,
		DateRange: domain.DateRange{
			From: e.DateFrom,
			To:   e.DateTo,
		},
		Meta: domain.Meta{
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		},

		Images: images,
		Movies: movies,
		Merits: merits,
		With: work_domain.WorkWith{
			Department: department,
			ActivePlan: activePlan,
		},
	}
}

func entityFrom(d *work_domain.Work) *Entity {
	return &Entity{
		ID:           d.ID.String(),
		DepartmentID: d.DepartmentID.String(),
		Status:       int32(d.Status),
		WorkType:     int32(d.WorkType),
		JobCode:      int32(d.JobCode),
		Title:        d.Title,
		Content:      d.Content,
		DateFrom:     d.DateRange.From,
		DateTo:       d.DateRange.To,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}
