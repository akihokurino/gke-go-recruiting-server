package work_domain

import (
	"gke-go-recruiting-server/domain/contract_domain"
	"gke-go-recruiting-server/domain/department_domain"
)

type WorkWith struct {
	Department *department_domain.Department
	ActivePlan *contract_domain.ActivePlan
}
