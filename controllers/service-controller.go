package controllers

import (
	"net/http"

	"github.com/facs95/orderwell-backend/entity"
)

//GetTenantServices returns all the services under that tenant
func (*controller) GetServices(w http.ResponseWriter, r *http.Request) {
	user := authProvider.UserFromContext(r.Context())
	services := []entity.Service{}
	err := dbService.Service().FindAll(user.TenantId, &services)
	if err != nil {
		errorHttpResponse(w, http.StatusInternalServerError)
		return
	}
	successHttpResponseEncoder(w, services)
	return
}
