package controllers

import (
	"net/http"

	"github.com/facs95/orderwell-backend/entity"
)

//GetTenantEmployees returns all the employees under that tenant
func (*controller) GetClients(w http.ResponseWriter, r *http.Request) {
	user := authProvider.UserFromContext(r.Context())
	clients := []entity.Client{}
	err := dbService.Client().FindAll(user.TenantId, &clients)
	if err != nil {
		errorHttpResponse(w, http.StatusInternalServerError)
		return
	}
	successHttpResponseEncoder(w, clients)
	return
}
