package service

import "github.com/facs95/orderwell-backend/entity"

type clientService interface {
	FindAll(tenantId string, tenant *[]entity.Client) error
}

type clientServiceInstance struct{}

func (*service) Client() clientService {
	return &clientServiceInstance{}
}

func (*clientServiceInstance) FindAll(id string, clients *[]entity.Client) error {
	return repo.FindAllInTenant(id, entity.Client{}, clients)
}
