package authentication

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/facs95/orderwell-backend/entity"
	"github.com/facs95/orderwell-backend/logger"
	"github.com/facs95/orderwell-backend/service"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

var (
	dbService service.Service
	app       *firebase.App
	err       error
)

type googleAuth struct{}

type contextKey string

const contextUserKey contextKey = "user"

type Auth interface {
	CreateOauthTenant(companyName string) (oauthTenant *auth.Tenant, err error)
	DeleteTenant(tenantOauthId string) error
	AuthMiddleware(next http.Handler) http.Handler
	UserFromContext(ctx context.Context) entity.User
	CreateUser(tenantOauthId string, newUser *entity.NewAuthUser) (*auth.UserRecord, error)
}

//InitGoogleSdk Inits connection with google SDK
func NewGoogleAuth(service service.Service) Auth {
	dbService = service
	ctx := context.Background()
	app, err = firebase.NewApp(ctx, nil)
	if err != nil {
		logger.ErrorLogger.Fatal("error initializing google SDK: ", err)
	}
	logger.InfoLogger.Println("Google SDK Initialized With Success")
	return &googleAuth{}
}

func (*googleAuth) UserFromContext(ctx context.Context) entity.User {
	return ctx.Value(contextUserKey).(entity.User)
}

//Creates a new tenant on google oauth with an independant user pool
func (*googleAuth) CreateOauthTenant(companyName string) (oauthTenant *auth.Tenant, err error) {
	ctx := context.Background()
	authInstance, err := app.Auth(ctx)
	if err != nil {
		logger.ErrorLogger.Println("We encountered an issue creating the auth: ", err)
		return nil, err
	}

	formattedCompanyName := standardizeCompanyName(companyName) //Removes spaces and join with hyphen

	newCreatedTenant := auth.TenantToCreate{}
	newCreatedTenant.DisplayName(formattedCompanyName)
	newCreatedTenant.AllowPasswordSignUp(true)
	newCreatedTenant.EnableEmailLinkSignIn(true)
	tenant, err := authInstance.TenantManager.CreateTenant(ctx, &newCreatedTenant)
	if err != nil {
		logger.ErrorLogger.Println("We encountered a problem creating the tenant: ", err)
		return nil, err
	}

	logger.InfoLogger.Print("Oauth Tenant succesfully created: ", tenant.ID)
	return tenant, nil
}

//Delete Oauth tenant by the ID
func (*googleAuth) DeleteTenant(tenantOauthId string) error {
	ctx := context.Background()
	authInstance, err := app.Auth(ctx)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}
	if err := authInstance.TenantManager.DeleteTenant(ctx, tenantOauthId); err != nil {
		message := fmt.Sprintf("There was an issue removing the tenant'%v': %v", tenantOauthId, err)
		logger.ErrorLogger.Println(message)
		return err
	}
	message := fmt.Sprintf("Deleted tenant '%v' succesfully", tenantOauthId)
	logger.InfoLogger.Println(message)
	return nil
}

//Handles HTTP authentication
func (*googleAuth) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Create parent context with a timeout of 60 seconds
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(60*time.Second))
		defer cancel()
		//Need to get subdomain from the request
		subDomain, err := getSubDomain(r) //Gets the first subdomain. There should only be one
		if err != nil {
			logger.ErrorLogger.Println("No subdomain found: ", err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		//Get the ID of the tenant based on that subdomaain
		tenant := entity.Tenant{}

		err = dbService.Tenant().FindBySubDomain(subDomain, &tenant)

		if err != nil {
			logger.TenantErrorLog(tenant.ID).Println("Error querying the tenant: ", err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if tenant.ID == "" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		//Use that tenant Id to authenticate the reequest
		reqToken := r.Header.Get("Authorization")
		if reqToken == "" {
			unauthorizedResponse(w)
			return
		}

		authInstance, err := app.Auth(ctx)
		if err != nil {
			logger.TenantErrorLog(tenant.ID).Println(err)
			unauthorizedResponse(w)
			return
		}

		client, err := authInstance.TenantManager.AuthForTenant(tenant.OauthId)
		if err != nil {
			logger.TenantErrorLog(tenant.ID).Println(err)
			unauthorizedResponse(w)
			return
		}

		// Need to remove the bearer text
		idToken := strings.Split(reqToken, "Bearer ")

		if len(idToken) != 2 {
			unauthorizedResponse(w)
			return
		}

		token, err := client.VerifyIDToken(ctx, idToken[1])
		if err != nil {
			unauthorizedResponse(w)
			return
		}

		ctx = context.WithValue(r.Context(), contextUserKey, entity.User{
			ID:       token.UID,
			TenantId: tenant.ID,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
