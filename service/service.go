package service

import "github.com/facs95/orderwell-backend/repository"

type Service interface {
	Tenant() tenantService
	Employee() employeeService
	Client() clientService
	Order() orderService
	Service() serviceService
}

type service struct{}

var (
	repo repository.Repository
)

func NewDbService(repository repository.Repository) Service {
	repo = repository
	return &service{}
}
