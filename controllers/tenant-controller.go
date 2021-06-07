package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/facs95/orderwell-backend/entity"
)

// func (t TenantWithPassword) Validate

//CreateTenant function
func (*controller) CreateTenant(w http.ResponseWriter, r *http.Request) {
	tenantWithUser := entity.TenantWithUser{}
	err := json.NewDecoder(r.Body).Decode(&tenantWithUser)
	if err != nil {
		code := http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
		return
	}

	if err := tenantWithUser.Validate(); err != nil {
		code := http.StatusBadRequest
		http.Error(w, http.StatusText(code), code)
		return
	}

	tenant := entity.Tenant{
		Email:       tenantWithUser.Email,
		CompanyName: tenantWithUser.CompanyName,
	}

	//make company email lowercase for ease
	tenant.CompanyName = strings.ToLower(tenant.CompanyName)
	tenantWithUser.FirstName = strings.ToLower(tenantWithUser.FirstName)
	tenantWithUser.LastName = strings.ToLower(tenantWithUser.LastName)

	//Meed to create a oauth tenant
	oauthTenant, err := authProvider.CreateOauthTenant(tenant.CompanyName)
	if err != nil {
		errorHttpResponse(w, http.StatusInternalServerError)
		return
	}
	tenant.OauthId = oauthTenant.ID

	//Create a user on the tenants oauth client
	newUser := entity.NewAuthUser{
		Email:       tenantWithUser.Email,
		Password:    tenantWithUser.Password,
		PhoneNumber: tenantWithUser.PhoneNumber,
		DisplayName: fmt.Sprintf("%s %s", tenantWithUser.FirstName, tenantWithUser.LastName),
	}

	employeeUser, err := authProvider.CreateUser(tenant.OauthId, &newUser)
	if err != nil {
		dbService.Tenant().Delete(tenant.OauthId) //Clean oauth tenant if record was not created
		errorHttpResponse(w, http.StatusInternalServerError)
		return
	}

	//Once the tenant has been created create the tenant record on database
	if err := dbService.Tenant().Create(&tenant); err != nil {
		authProvider.DeleteTenant(tenant.OauthId) //Clean oauth tenant if record was not created
		errorHttpResponse(w, http.StatusInternalServerError)
		return
	}

	//Input the Owner as an employee with respective role
	employee := entity.Employee{
		ID:          employeeUser.UID, //We are going to mantain a sync between user created and user ID
		Email:       tenantWithUser.Email,
		FirstName:   tenantWithUser.FirstName,
		LastName:    tenantWithUser.LastName,
		PhoneNumber: tenantWithUser.PhoneNumber,
		Role:        "owner",
	}

	err = dbService.Employee().Create(tenant.ID, &employee)
	if err != nil {
		handleTenantFailure(tenant)
		errorHttpResponse(w, http.StatusInternalServerError)
		return
	}

	successHttpResponseEncoder(w, tenant)
	return
}

func (*controller) GetTenantId(w http.ResponseWriter, r *http.Request) {
	tenant := entity.Tenant{}
	subDomain, err := getSubDomain(r)
	if err != nil {
		errorHttpResponse(w, http.StatusBadRequest)
		return
	}

	err = dbService.Tenant().FindBySubDomain(subDomain, &tenant)
	if err != nil {
		errorHttpResponse(w, http.StatusBadRequest)
		return
	}

	response := make(map[string]string)
	response["oauthId"] = tenant.OauthId
	response["companyName"] = tenant.CompanyName
	successHttpResponseEncoder(w, response)
	return
}

func (*controller) ValidateCompanyName(w http.ResponseWriter, r *http.Request) {
	companyName := r.URL.Query().Get("companyName")
	companyName = strings.ToLower(companyName)
	if len(companyName) > 20 || len(companyName) < 4 {
		response := make(map[string]interface{})
		response["companyName"] = companyName
		response["isValid"] = false
		successHttpResponseEncoder(w, response)
		return
	}

	exists, err := dbService.Tenant().ValidateCompanyName(companyName)
	if err != nil {
		errorHttpResponse(w, http.StatusInternalServerError)
		return
	}

	response := make(map[string]interface{})
	response["companyName"] = companyName
	response["isValid"] = exists
	successHttpResponseEncoder(w, response)
	return
}

//ALL CALLS FROM HERE SHOULD BE PRIVATE

//GetTenant returns a tenant by an ID
func (*controller) GetTenant(w http.ResponseWriter, r *http.Request) {
	user := authProvider.UserFromContext(r.Context())
	tenant := entity.Tenant{}
	err := dbService.Tenant().Find(user.TenantId, &tenant)
	if err != nil {
		errorHttpResponse(w, http.StatusNotFound)
	}
	successHttpResponseEncoder(w, tenant)
}
