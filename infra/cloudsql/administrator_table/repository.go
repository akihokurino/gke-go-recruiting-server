package administrator_table

import (
	"context"

	"gke-go-sample/infra/cloudsql"

	"gorm.io/gorm"

	"gke-go-sample/adapter"
	"gke-go-sample/domain"
	"gke-go-sample/domain/account_domain"

	"github.com/pkg/errors"
)

func NewRepo() adapter.AdministratorRepo {
	return &repository{}
}

type repository struct {
}

func (r *repository) Get(ctx context.Context, db *gorm.DB, id domain.FirebaseUserID) (*account_domain.Administrator, error) {
	var e Entity
	if err := db.Where("id = ?", id.String()).
		First(&e).Error; err != nil {
		if cloudsql.IsNotFoundError(err) {
			return nil, domain.NewNotFoundErr()
		}
		return nil, errors.WithStack(err)
	}

	return e.toDomain(), nil
}

func (r *repository) Insert(ctx context.Context, db *gorm.DB, item *account_domain.Administrator) error {
	if err := db.Create(entityFrom(item)).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
