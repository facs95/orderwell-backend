package service

import "github.com/facs95/orderwell-backend/entity"

type serviceService interface {
	FindAll(tenantId string, tenant *[]entity.Service) error
}

type serviceServiceInstance struct{}

func (*service) Service() serviceService {
	return &serviceServiceInstance{}
}

func (*serviceServiceInstance) FindAll(id string, services *[]entity.Service) error {
	return repo.FindAllInTenant(id, entity.Service{}, services)
}
