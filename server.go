package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/facs95/orderwell-backend/authentication"
	"github.com/facs95/orderwell-backend/controllers"
	"github.com/facs95/orderwell-backend/logger"
	"github.com/facs95/orderwell-backend/repository"

	router "github.com/facs95/orderwell-backend/http"

	"github.com/facs95/orderwell-backend/service"
)

var (
	httpRouter   router.Router          = router.NewChiRouter()
	dbRepository repository.Repository  = repository.InitPostgressRepo()
	dbService    service.Service        = service.NewDbService(dbRepository)
	authProvider authentication.Auth    = authentication.NewGoogleAuth(dbService)
	controller   controllers.Controller = controllers.NewController(dbService, authProvider, httpRouter)
)

func main() {
	httpRouter.UseCors()

	httpRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, welcome to Orderwell API 👋!")
	})

	httpRouter.Post("/create-tenant", controller.CreateTenant)
	httpRouter.Get("/tenant/valid", controller.ValidateCompanyName)
	httpRouter.Get("/tenant/id", controller.GetTenantId)

	//From here on all the routes will be private routes
	httpRouter.PrivateGet("/tenant", authProvider.AuthMiddleware, controller.GetTenant)
	httpRouter.PrivateGet("/employees", authProvider.AuthMiddleware, controller.GetEmployees)
	httpRouter.PrivateGet("/clients", authProvider.AuthMiddleware, controller.GetClients)
	httpRouter.PrivateGet("/orders", authProvider.AuthMiddleware, controller.GetOrders)
	httpRouter.PrivateGet("/services", authProvider.AuthMiddleware, controller.GetServices)

	err := httpRouter.Serve(fmt.Sprintf(":%v", os.Getenv("PORT")))
	if err != nil {
		logger.ErrorLogger.Panic("Failed to start server", err)
	}
}
