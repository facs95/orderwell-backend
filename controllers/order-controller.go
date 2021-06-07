package controllers

import (
	"net/http"

	"github.com/facs95/orderwell-backend/entity"
)

//GetTenantOrders returns all the orders under that tenant
func (*controller) GetOrders(w http.ResponseWriter, r *http.Request) {
	user := authProvider.UserFromContext(r.Context())
	orders := []entity.Order{}
	err := dbService.Order().FindAll(user.TenantId, &orders)
	if err != nil {
		errorHttpResponse(w, http.StatusInternalServerError)
		return
	}
	successHttpResponseEncoder(w, orders)
	return
}
