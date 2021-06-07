package service

import "github.com/facs95/orderwell-backend/entity"

type employeeService interface {
	Create(tenantId string, tenant *entity.Employee) error
	FindAll(tenantId string, tenant *[]entity.Employee) error
}

type employeeServiceInstance struct{}

func (*service) Employee() employeeService {
	return &employeeServiceInstance{}
}

//TODO have a constant for table names
func (*employeeServiceInstance) Create(tenantId string, employee *entity.Employee) error {
	employee.ID = generateId()
	return repo.SaveInTenant(tenantId, entity.Employee{}, employee)
}

func (*employeeServiceInstance) FindAll(id string, employees *[]entity.Employee) error {
	return repo.FindAllInTenant(id, entity.Employee{}, employees)
}
