package department_domain

import (
	"gke-go-sample/domain/master_domain"
)

type StationWith struct {
	Line *master_domain.Line
}
