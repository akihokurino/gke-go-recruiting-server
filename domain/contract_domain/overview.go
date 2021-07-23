package contract_domain

import "gke-go-sample/domain"

type DepartmentOverview struct {
	ID       domain.DepartmentID
	AgencyID domain.AgencyID
}
