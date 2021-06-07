package controllers

import (
	"net/http"

	"github.com/facs95/orderwell-backend/authentication"
	router "github.com/facs95/orderwell-backend/http"
	"github.com/facs95/orderwell-backend/service"
)

type controller struct{}

type Controller interface {
	GetTenant(response http.ResponseWriter, request *http.Request)
	GetTenantId(response http.ResponseWriter, request *http.Request)
	CreateTenant(response http.ResponseWriter, request *http.Request)
	ValidateCompanyName(response http.ResponseWriter, request *http.Request)

	GetEmployees(response http.ResponseWriter, request *http.Request)

	GetClients(response http.ResponseWriter, request *http.Request)

	GetOrders(response http.ResponseWriter, request *http.Request)

	GetServices(response http.ResponseWriter, request *http.Request)
}

var (
	authProvider authentication.Auth
	dbService    service.Service
	httpRouter   router.Router
)

func NewController(service service.Service, auth authentication.Auth, router router.Router) Controller {
	dbService = service
	authProvider = auth
	httpRouter = router
	return &controller{}
}
