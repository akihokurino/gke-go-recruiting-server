package work_merit_table

import (
	"database/sql"
	"time"

	pb "gke-go-recruiting-server/proto/go/pb"

	"gke-go-recruiting-server/domain"
	"gke-go-recruiting-server/domain/work_domain"

	"github.com/guregu/null"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

func (e *Entity) TableName() string {
	return "work_merits"
}

type Entity struct {
	ID     string `gorm:"column:id;primary_key"`
	WorkID string `gorm:"column:work_id"`
	Value  int32  `gorm:"column:value"`
}

func (e *Entity) ToDomain() *work_domain.Merit {
	return &work_domain.Merit{
		ID:     e.ID,
		WorkID: domain.WorkID(e.WorkID),
		Value:  pb.Work_Merit(e.Value),
	}
}

func entityFrom(d *work_domain.Merit) *Entity {
	return &Entity{
		ID:     d.ID,
		WorkID: d.WorkID.String(),
		Value:  int32(d.Value),
	}
}
