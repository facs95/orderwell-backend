package controllers

import (
	"net/http"

	"github.com/facs95/orderwell-backend/entity"
)

//GetTenantEmployees returns all the employees under that tenant
func (*controller) GetEmployees(w http.ResponseWriter, r *http.Request) {
	user := authProvider.UserFromContext(r.Context())
	employees := []entity.Employee{}
	err := dbService.Employee().FindAll(user.TenantId, &employees)
	if err != nil {
		errorHttpResponse(w, http.StatusInternalServerError)
		return
	}
	successHttpResponseEncoder(w, employees)
	return
}
