package work_active_plan_table

import (
	"context"
	"database/sql"
	"fmt"

	"gorm.io/gorm"

	"gke-go-sample/adapter"
	"gke-go-sample/domain"
	"gke-go-sample/domain/contract_domain"

	"github.com/pkg/errors"
)

func NewRepo() adapter.WorkActivePlanRepo {
	return &repository{}
}

type repository struct {
}

func (r *repository) GetAll(ctx context.Context, db *gorm.DB) ([]*contract_domain.ActivePlan, error) {
	var es []Entity
	if err := db.Preload("MainContract").Find(&es).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*contract_domain.ActivePlan, 0, len(es))
	for _, e := range es {
		items = append(items, e.ToDomain())
	}

	return items, nil
}

func (r *repository) Upsert(
	ctx context.Context,
	db *gorm.DB,
	id domain.WorkID,
	mainContractID domain.MainContractID,
	publishedOrder int) error {
	q := fmt.Sprintf(`
insert into %s (work_id, main_contract_id, published_order) 
values (?,?,?) 
on duplicate key update work_id=?, main_contract_id=?, published_order=?`,
		(&Entity{}).TableName())

	if err := db.Exec(
		q,
		id.String(),
		sql.NullString{
			String: mainContractID.String(),
			Valid:  mainContractID.String() != "",
		},
		publishedOrder,
		id.String(),
		sql.NullString{
			String: mainContractID.String(),
			Valid:  mainContractID.String() != "",
		},
		publishedOrder).Error; err != nil {

		return errors.WithStack(err)
	}

	return nil
}

func (r *repository) UpdateMulti(ctx context.Context, db *gorm.DB, items []*contract_domain.ActivePlan) error {
	for _, item := range items {
		if err := db.Save(entityFrom(item)).Error; err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, db *gorm.DB, id domain.WorkID) error {
	if err := db.Delete(onlyID(id)).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
