package work_domain

import (
	"gke-go-sample/domain/contract_domain"
	"gke-go-sample/domain/department_domain"
)

type WorkWith struct {
	Department *department_domain.Department
	ActivePlan *contract_domain.ActivePlan
}
