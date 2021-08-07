package work_merit_table

import (
	"context"

	"gke-go-recruiting-server/infra/cloudsql"

	"gorm.io/gorm"

	"gke-go-recruiting-server/domain"

	"github.com/pkg/errors"

	"gke-go-recruiting-server/adapter"
	"gke-go-recruiting-server/domain/work_domain"
)

func NewRepo() adapter.WorkMeritRepo {
	return &repository{}
}

type repository struct {
}

func (r *repository) InsertMulti(ctx context.Context, db *gorm.DB, items []*work_domain.Merit) error {
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

func (r *repository) DeleteByWork(ctx context.Context, db *gorm.DB, workID domain.WorkID) error {
	if err := db.Where("work_id = ?", workID.String()).Delete(Entity{}).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
