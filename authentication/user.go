package authentication

import (
	"context"

	"firebase.google.com/go/auth"
	"github.com/facs95/orderwell-backend/entity"
	"github.com/facs95/orderwell-backend/logger"
)

func (*googleAuth) CreateUser(tenantOauthId string, newUser *entity.NewAuthUser) (*auth.UserRecord, error) {
	ctx := context.Background()
	authInstance, err := app.Auth(ctx)
	if err != nil {
		logger.TenantErrorLog(tenantOauthId).Println(err)
		return nil, err
	}

	client, err := authInstance.TenantManager.AuthForTenant(tenantOauthId)
	if err != nil {
		logger.TenantErrorLog(tenantOauthId).Println(err)
		return nil, err
	}

	user := (&auth.UserToCreate{}).
		Email(newUser.Email).
		Password(newUser.Password).
		DisplayName(newUser.DisplayName).
		Disabled(false).
		PhoneNumber(newUser.PhoneNumber)
	createdUser, err := client.CreateUser(ctx, user)

	if err != nil {
		logger.ErrorLogger.Println("There was an issue creating the oauth user: ", err)
		return nil, err
	}
	logger.InfoLogger.Println("User created succesfully")

	return createdUser, nil
}
