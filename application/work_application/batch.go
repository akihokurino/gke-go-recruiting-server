package work_application

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gke-go-sample/domain/contract_domain"

	"gke-go-sample/domain"
	pb "gke-go-sample/proto/go/pb"

	"github.com/pkg/errors"
)

type batchApp struct {
	*app
}

func (a *batchApp) Proceed(ctx context.Context, db *gorm.DB, now time.Time) error {
	// 終了する求人の処理
	finishWorks, err := a.workRepo.GetByWillFinish(ctx, db, now)
	if err != nil {
		return errors.WithStack(err)
	}

	for i := range finishWorks {
		a.logger.Info().With(ctx).Printf("finish work id: %s", finishWorks[i].ID)
		if err := finishWorks[i].Finish(); err != nil {
			return errors.WithStack(err)
		}
	}

	if err := a.workRepo.UpdateMulti(ctx, db, finishWorks); err != nil {
		return errors.WithStack(err)
	}

	for _, work := range finishWorks {
		if err := a.workActivePlanRepo.Delete(ctx, db, work.ID); err != nil {
			return errors.WithStack(err)
		}
	}

	finishWorkIDs := make([]domain.WorkID, 0, len(finishWorks))
	for _, work := range finishWorks {
		finishWorkIDs = append(finishWorkIDs, work.ID)
	}

	if err := a.workIndexRepo.DeleteMulti(ctx, finishWorkIDs); err != nil {
		return errors.WithStack(err)
	}

	// 再公開する求人の処理
	for _, work := range finishWorks {
		if err := a.tx(db, func(db *gorm.DB) error {
			if work.With.ActivePlan == nil {
				return nil
			}

			currentContractID := work.With.ActivePlan.MainContractID

			nextContract, err := a.mainContractRepo.GetByActiveAndDepartmentAndTime(
				ctx,
				db,
				work.DepartmentID,
				now.Add(1*time.Hour))
			if err != nil && !domain.IsNotFound(err) {
				return err
			}

			if err != nil && domain.IsNotFound(err) {
				return nil
			}

			if currentContractID == nextContract.ID {
				return nil
			}

			work.Continue(nextContract)

			if err := a.workRepo.Update(ctx, db, work); err != nil {
				return err
			}

			return nil
		}); err != nil {
			return errors.WithStack(err)
		}
	}

	// 新規公開する求人の処理
	startWorks, err := a.workRepo.GetByWillStart(ctx, db, now)
	if err != nil {
		return errors.WithStack(err)
	}

	for i := range startWorks {
		a.logger.Info().With(ctx).Printf("open work id: %s", startWorks[i].ID)
		if err := startWorks[i].Active(); err != nil {
			return errors.WithStack(err)
		}
	}

	if err := a.workRepo.UpdateMulti(ctx, db, startWorks); err != nil {
		return errors.WithStack(err)
	}

	// 求人の掲載順設定の更新
	activeWorks, err := a.workRepo.GetByStatusAndSEO(ctx, db, pb.Work_Status_Active, false)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, work := range activeWorks {
		if err := a.tx(db, func(db *gorm.DB) error {
			contract, err := a.mainContractRepo.GetByActiveAndDepartmentAndTime(ctx, db, work.DepartmentID, now)
			if err != nil && !domain.IsNotFound(err) {
				return err
			}

			if err != nil && domain.IsNotFound(err) {
				if err := work.Finish(); err != nil {
					return err
				}

				a.logger.Info().With(ctx).Printf("finish work id: %s", work.ID)

				if err := a.workRepo.Update(ctx, db, work); err != nil {
					return err
				}

				if err := a.workActivePlanRepo.Delete(ctx, db, work.ID); err != nil {
					a.logger.Error().With(ctx).Printf("error in delete work active plan, err=%+v", err)
				}

				if err := a.workIndexRepo.Delete(ctx, work.ID); err != nil {
					return err
				}

				return nil
			}

			return a.workActivePlanRepo.Upsert(
				ctx,
				db,
				work.ID,
				contract.ID,
				contract_domain.PublishedOrderFrom(contract.Plan))
		}); err != nil {
			return errors.WithStack(err)
		}
	}

	// 検索インデックスの更新
	activeWorks, err = a.workRepo.GetByStatusAndSEO(ctx, db, pb.Work_Status_Active, false)
	if err != nil {
		return errors.WithStack(err)
	}

	a.logger.Info().With(ctx).Printf("new work index num: %d", len(activeWorks))

	if err := a.workIndexRepo.SaveMulti(ctx, activeWorks); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *batchApp) ReOrder(ctx context.Context, db *gorm.DB) error {
	plans, err := a.workActivePlanRepo.GetAll(ctx, db)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, plan := range plans {
		plan.ReOrder()
	}

	if err := a.workActivePlanRepo.UpdateMulti(ctx, db, plans); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *batchApp) ReIndex(ctx context.Context, db *gorm.DB) error {
	activeWorks, err := a.workRepo.GetByStatusAndSEO(ctx, db, pb.Work_Status_Active, false)
	if err != nil {
		return errors.WithStack(err)
	}

	activeSEOWorks, err := a.workRepo.GetByStatusAndSEO(ctx, db, pb.Work_Status_Active, true)
	if err != nil {
		return errors.WithStack(err)
	}

	activeWorks = append(activeWorks, activeSEOWorks...)

	if err := a.workIndexRepo.SaveMulti(ctx, activeWorks); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
