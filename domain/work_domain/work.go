package work_domain

import (
	"time"

	"gke-go-sample/domain/contract_domain"

	"github.com/google/uuid"

	"gke-go-sample/domain"
	pb "gke-go-sample/proto/go/pb"
)

type Work struct {
	ID           domain.WorkID
	DepartmentID domain.DepartmentID
	Status       pb.Work_Status
	WorkType     pb.Work_Type
	JobCode      pb.Work_Job
	Title        string
	Content      string
	DateRange    domain.DateRange
	domain.Meta

	Images []*Image
	Movies []*Movie
	Merits []*Merit
	With   WorkWith
}

func New(
	departmentID domain.DepartmentID,
	workType pb.Work_Type,
	jobCode pb.Work_Job,
	title string,
	content string,
	dateRange domain.DateRange,
	now time.Time) *Work {
	status := pb.Work_Status_Review

	return &Work{
		ID:           domain.WorkID(uuid.New().String()),
		DepartmentID: departmentID,
		Status:       status,
		WorkType:     workType,
		JobCode:      jobCode,
		Title:        title,
		Content:      content,
		DateRange:    dateRange,
		Meta: domain.Meta{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

func (w *Work) Update(
	workType pb.Work_Type,
	jobCode pb.Work_Job,
	title string,
	content string,
	dateRange domain.DateRange,
	now time.Time) error {
	if w.IsFinish() {
		return domain.NewForbiddenErr("終了した求人は変更できません")
	}

	w.WorkType = workType
	w.JobCode = jobCode
	w.Title = title
	w.Content = content
	w.DateRange = dateRange
	w.Meta.UpdatedAt = now

	return nil
}

func (w *Work) Accept() error {
	if w.Status != pb.Work_Status_Review {
		return domain.NewConflictErr(domain.ConflictStatusMsg)
	}

	w.Status = pb.Work_Status_Reserved
	return nil
}

func (w *Work) Deny() error {
	if w.Status != pb.Work_Status_Review {
		return domain.NewConflictErr(domain.ConflictStatusMsg)
	}

	w.Status = pb.Work_Status_NG
	return nil
}

func (w *Work) Active() error {
	if w.Status != pb.Work_Status_Reserved {
		return domain.NewConflictErr(domain.ConflictStatusMsg)
	}

	w.Status = pb.Work_Status_Active
	return nil
}

func (w *Work) Finish() error {
	if w.Status != pb.Work_Status_Active {
		return domain.NewConflictErr(domain.ConflictStatusMsg)
	}

	w.Status = pb.Work_Status_Finish
	return nil
}

func (w *Work) EarlyFinish() error {
	if w.Status != pb.Work_Status_Active {
		return domain.NewConflictErr(domain.ConflictStatusMsg)
	}

	w.Status = pb.Work_Status_Early_Finish
	return nil
}

func (w *Work) Continue(contract *contract_domain.Main) {
	w.Status = pb.Work_Status_Active
	w.DateRange.To = contract.DateRange.To
}

func (w *Work) IsNew(now time.Time) bool {
	diff := now.Sub(domain.UTC(w.CreatedAt))
	return diff <= 7*24*time.Hour
}

func (w *Work) IsFinish() bool {
	return w.Status == pb.Work_Status_Early_Finish || w.Status == pb.Work_Status_Finish
}
