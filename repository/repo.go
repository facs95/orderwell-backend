package repository

import "github.com/facs95/orderwell-backend/entity"

type Repository interface {
	CreateTenant(tenant *entity.Tenant) error //Creation of tenants is independant from the controller
	DeleteTenant(id string) error
	FindTenant(id string, value interface{}) error
	FindTenantBySubDomain(id string, value interface{}) error
	ValidateTenantCompanyName(companyName string) (isValid bool, err error)

	SaveInTenant(tenantId string, table interface{}, value interface{}) error
	FindAllInTenant(tenantId string, table interface{}, value interface{}) error
}

// type TenantRepo interface {
// 	Save(tenant entity.Tenant) entity.Tenant // Creation of tenants is independant from the controller
// 	Delete(id string) error                  //This will have to delete everything related with that tenant
// }

// type EmployeeRepo interface {
// 	Save(employee entity.Employee) (entity.Employee, error)
// 	Delete(id string) error
// }
