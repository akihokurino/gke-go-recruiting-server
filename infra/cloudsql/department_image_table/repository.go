package department_image_table

import (
	"context"

	"gke-go-recruiting-server/infra/cloudsql"

	"gorm.io/gorm"

	"gke-go-recruiting-server/domain"

	"github.com/pkg/errors"

	"gke-go-recruiting-server/domain/department_domain"

	"gke-go-recruiting-server/adapter"
)

func NewRepo() adapter.DepartmentImageRepo {
	return &repository{}
}

type repository struct {
}

func (r *repository) InsertMulti(ctx context.Context, db *gorm.DB, items []*department_domain.Image) error {
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

func (r *repository) DeleteByDepartment(ctx context.Context, db *gorm.DB, departmentID domain.DepartmentID) error {
	if err := db.Where("department_id = ?", departmentID.String()).Delete(Entity{}).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
