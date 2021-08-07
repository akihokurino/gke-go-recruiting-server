package department_domain

import (
	"gke-go-recruiting-server/domain/master_domain"
)

type StationWith struct {
	Line *master_domain.Line
}
