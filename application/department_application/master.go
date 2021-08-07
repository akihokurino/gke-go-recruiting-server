package department_application

import (
	"context"
	"net/url"
	"time"

	"gorm.io/gorm"

	"gke-go-recruiting-server/adapter"
	"gke-go-recruiting-server/domain/department_domain"

	"gke-go-recruiting-server/domain"
	"gke-go-recruiting-server/domain/account_domain"

	"github.com/pkg/errors"
)

type masterApp struct {
	me *account_domain.Administrator
	*app
}

func (a *masterApp) Accept(
	ctx context.Context,
	db *gorm.DB,
	departmentID domain.DepartmentID) error {
	if err := a.tx(db, func(db *gorm.DB) error {
		department, err := a.departmentRepo.Get(ctx, db, departmentID)
		if err != nil {
			return err
		}

		if err := department.Accept(); err != nil {
			return err
		}

		return a.departmentRepo.Update(ctx, db, department)
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *masterApp) Deny(
	ctx context.Context,
	db *gorm.DB,
	departmentID domain.DepartmentID) error {
	if err := a.tx(db, func(db *gorm.DB) error {
		department, err := a.departmentRepo.Get(ctx, db, departmentID)
		if err != nil {
			return err
		}

		if err := department.Deny(); err != nil {
			return err
		}

		return a.departmentRepo.Update(ctx, db, department)
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *masterApp) Update(
	ctx context.Context,
	db *gorm.DB,
	departmentID domain.DepartmentID,
	departmentParams adapter.DepartmentParams,
	imageURLs []url.URL,
	lineIDs []domain.LineID,
	now time.Time) (*department_domain.Department, error) {
	if err := validateUpdate(departmentParams); err != nil {
		return nil, errors.WithStack(err)
	}

	newImages := make([]*department_domain.Image, 0, len(imageURLs))
	for _, u := range imageURLs {
		newImages = append(newImages, department_domain.NewImage(departmentID, u))
	}

	newStations := make([]*department_domain.Station, 0, len(lineIDs))
	for _, lineID := range lineIDs {
		newStations = append(newStations, department_domain.NewStation(departmentID, lineID))
	}

	var department *department_domain.Department
	var err error
	if err := a.tx(db, func(db *gorm.DB) error {
		department, err = a.departmentRepo.Get(ctx, db, departmentID)
		if err != nil {
			return err
		}

		department.Update(
			departmentParams.Name,
			departmentParams.BusinessCondition,
			departmentParams.PostalCode,
			departmentParams.PrefID,
			departmentParams.CityID,
			departmentParams.Address,
			departmentParams.Building,
			departmentParams.PhoneNumber,
			department_domain.Location{
				MAreaID:   departmentParams.MAreaID,
				SAreaID:   departmentParams.SAreaID,
				Latitude:  departmentParams.Latitude,
				Longitude: departmentParams.Longitude,
			},
			now)

		if err := a.departmentRepo.Update(ctx, db, department); err != nil {
			return err
		}

		if err := a.departmentImageRepo.DeleteByDepartment(ctx, db, department.ID); err != nil {
			return err
		}

		if err := a.departmentStationRepo.DeleteByDepartment(ctx, db, department.ID); err != nil {
			return err
		}

		if err := a.departmentImageRepo.InsertMulti(ctx, db, newImages); err != nil {
			return err
		}

		if err := a.departmentStationRepo.InsertMulti(ctx, db, newStations); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	return department, nil
}

func (a *masterApp) UpdateSales(
	ctx context.Context,
	db *gorm.DB,
	departmentID domain.DepartmentID,
	salesID domain.FirebaseUserID) error {
	if err := a.tx(db, func(db *gorm.DB) error {
		department, err := a.departmentRepo.Get(ctx, db, departmentID)
		if err != nil {
			return err
		}

		department.UpdateSales(salesID)

		return a.departmentRepo.Update(ctx, db, department)
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
