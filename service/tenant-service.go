package service

import "github.com/facs95/orderwell-backend/entity"

type tenantService interface {
	Create(tenant *entity.Tenant) error
	Delete(id string) error
	Find(id string, tenant *entity.Tenant) error
	FindBySubDomain(subdomain string, tenant *entity.Tenant) error
	ValidateCompanyName(companyName string) (isValid bool, err error)
}

type tenantServiceInstance struct{}

func (*service) Tenant() tenantService {
	return &tenantServiceInstance{}
}

func (*tenantServiceInstance) Create(tenant *entity.Tenant) error {
	tenant.ID = generateId()
	return repo.CreateTenant(tenant)
}

func (*tenantServiceInstance) Delete(id string) error {
	return repo.DeleteTenant(id)
}

func (*tenantServiceInstance) Find(id string, tenant *entity.Tenant) error {
	return repo.FindTenant(id, tenant)
}

func (*tenantServiceInstance) FindBySubDomain(subdomain string, tenant *entity.Tenant) error {
	return repo.FindTenantBySubDomain(subdomain, tenant)
}
func (*tenantServiceInstance) ValidateCompanyName(companyName string) (isValid bool, err error) {
	return repo.ValidateTenantCompanyName(companyName)
}
