package work_domain

import (
	"gke-go-recruiting-server/domain"
	pb "gke-go-recruiting-server/proto/go/pb"

	"github.com/google/uuid"
)

type Merit struct {
	ID     string
	WorkID domain.WorkID
	Value  pb.Work_Merit
}

func NewMerit(workID domain.WorkID, value pb.Work_Merit) *Merit {
	return &Merit{
		ID:     uuid.New().String(),
		WorkID: workID,
		Value:  value,
	}
}
