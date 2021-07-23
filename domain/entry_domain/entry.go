package entry_domain

import (
	"time"

	"github.com/google/uuid"

	pb "gke-go-sample/proto/go/pb"

	"gke-go-sample/domain"
)

type Entry struct {
	ID                     domain.EntryID
	DepartmentID           domain.DepartmentID
	WorkID                 domain.WorkID
	FullName               string
	FullNameKana           string
	Birthdate              time.Time
	Gender                 pb.User_Gender
	PhoneNumber            string
	Email                  string
	Question               string
	Category               *pb.User_Category
	PrefID                 *domain.PrefID
	PreferredContactMethod *pb.Entry_PreferredContactMethod
	PreferredContactTime   *string
	Status                 pb.Entry_Status
	domain.Meta
}

func New(
	departmentID domain.DepartmentID,
	workID domain.WorkID,
	fullName string,
	fullNameKana string,
	birthdate time.Time,
	gender pb.User_Gender,
	phoneNumber string,
	email string,
	question string,
	category *pb.User_Category,
	prefID *domain.PrefID,
	preferredContactMethod *pb.Entry_PreferredContactMethod,
	preferredContactTime *string,
	now time.Time) *Entry {
	return &Entry{
		ID:                     domain.EntryID(uuid.New().String()),
		DepartmentID:           departmentID,
		WorkID:                 workID,
		FullName:               fullName,
		FullNameKana:           fullNameKana,
		Birthdate:              birthdate,
		Gender:                 gender,
		PhoneNumber:            phoneNumber,
		Email:                  email,
		Question:               question,
		Category:               category,
		PrefID:                 prefID,
		PreferredContactMethod: preferredContactMethod,
		PreferredContactTime:   preferredContactTime,
		Status:                 pb.Entry_Status_InProgress,
		Meta: domain.Meta{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}
