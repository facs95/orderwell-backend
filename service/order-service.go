package service

import "github.com/facs95/orderwell-backend/entity"

type orderService interface {
	FindAll(tenantId string, tenant *[]entity.Order) error
}

type orderServiceInstance struct{}

func (*service) Order() orderService {
	return &orderServiceInstance{}
}

func (*orderServiceInstance) FindAll(id string, orders *[]entity.Order) error {
	return repo.FindAllInTenant(id, entity.Order{}, orders)
}
