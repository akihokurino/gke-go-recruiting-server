package entry_table

import (
	"database/sql"
	"time"

	"gke-go-sample/domain"
	"gke-go-sample/domain/entry_domain"
	pb "gke-go-sample/proto/go/pb"

	"github.com/guregu/null"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

func (e *Entity) TableName() string {
	return "entries"
}

type Entity struct {
	ID                     string    `gorm:"column:id;primary_key"`
	DepartmentID           string    `gorm:"column:department_id"`
	WorkID                 string    `gorm:"column:work_id"`
	FullName               string    `gorm:"column:full_name"`
	FullNameKana           string    `gorm:"column:full_name_kana"`
	Birthdate              time.Time `gorm:"column:birthdate"`
	Gender                 int32     `gorm:"column:gender"`
	PhoneNumber            string    `gorm:"column:phone_number"`
	Email                  string    `gorm:"column:email"`
	Question               string    `gorm:"column:question"`
	Category               *int32    `gorm:"column:category"`
	PrefID                 *string   `gorm:"column:pref_id"`
	PreferredContactMethod *int32    `gorm:"column:preferred_contact_method"`
	PreferredContactTime   *string   `gorm:"column:preferred_contact_time"`
	Status                 int32     `gorm:"column:status"`
	CreatedAt              time.Time `gorm:"column:created_at"`
	UpdatedAt              time.Time `gorm:"column:updated_at"`
}

func (e *Entity) toDomain() *entry_domain.Entry {
	var category *pb.User_Category
	if e.Category != nil {
		tmp := pb.User_Category(*e.Category)
		category = &tmp
	}

	var prefID *domain.PrefID
	if e.PrefID != nil {
		tmp := domain.PrefID(*e.PrefID)
		prefID = &tmp
	}

	var method *pb.Entry_PreferredContactMethod
	if e.PreferredContactMethod != nil {
		tmp := pb.Entry_PreferredContactMethod(*e.PreferredContactMethod)
		method = &tmp
	}

	return &entry_domain.Entry{
		ID:                     domain.EntryID(e.ID),
		DepartmentID:           domain.DepartmentID(e.DepartmentID),
		WorkID:                 domain.WorkID(e.WorkID),
		FullName:               e.FullName,
		FullNameKana:           e.FullNameKana,
		Birthdate:              e.Birthdate,
		Gender:                 pb.User_Gender(e.Gender),
		PhoneNumber:            e.PhoneNumber,
		Email:                  e.Email,
		Question:               e.Question,
		Category:               category,
		PrefID:                 prefID,
		PreferredContactMethod: method,
		PreferredContactTime:   e.PreferredContactTime,
		Status:                 pb.Entry_Status(e.Status),
		Meta: domain.Meta{
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		},
	}
}

func entityFrom(d *entry_domain.Entry) *Entity {
	var category *int32
	if d.Category != nil {
		tmp := int32(*d.Category)
		category = &tmp
	}

	var prefID *string
	if d.PrefID != nil {
		tmp := (*d.PrefID).String()
		prefID = &tmp
	}

	var method *int32
	if d.PreferredContactMethod != nil {
		tmp := int32(*d.PreferredContactMethod)
		method = &tmp
	}

	return &Entity{
		ID:                     d.ID.String(),
		DepartmentID:           d.DepartmentID.String(),
		WorkID:                 d.WorkID.String(),
		FullName:               d.FullName,
		FullNameKana:           d.FullNameKana,
		Birthdate:              d.Birthdate,
		Gender:                 int32(d.Gender),
		PhoneNumber:            d.PhoneNumber,
		Email:                  d.Email,
		Question:               d.Question,
		Category:               category,
		PrefID:                 prefID,
		PreferredContactMethod: method,
		PreferredContactTime:   d.PreferredContactTime,
		Status:                 int32(d.Status),
		CreatedAt:              d.Meta.CreatedAt,
		UpdatedAt:              d.Meta.UpdatedAt,
	}
}
