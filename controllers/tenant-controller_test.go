package controllers

import (
	"github.com/facs95/orderwell-backend/authentication"

	"firebase.google.com/go/auth"
	"github.com/stretchr/testify/mock"
)

type MockedAuthentication struct {
	mock.Mock
}

func (m *MockedAuthentication) CreateOauthTenant(companyName string) (oauthTenant *auth.Tenant, err error) {
	args := m.Called(companyName)
	return args.Get(0).(*auth.Tenant), args.Error(1)
}

func (m *MockedAuthentication) DeleteTenant(tenantId string) error {
	args := m.Called(tenantId)
	return args.Error(1)
}

func (m *MockedAuthentication) CreateUser(tenantOauthId string, newUser *authentication.NewUser) (*auth.UserRecord, error) {
	args := m.Called(tenantOauthId, newUser)
	return args.Get(0).(*auth.UserRecord), args.Error(1)
}

type MockedRepository struct {
	mock.Mock
}

func (m *MockedRepository) CreateRecord(table interface{}) error {
	args := m.Called(table)
	return args.Error(0)
}

func (m *MockedRepository) CreateSchema(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *MockedRepository) DeleteRecord(table interface{}, id string) error {
	args := m.Called(table, id)
	return args.Error(0)
}

func (m *MockedRepository) RunSchemaMigration(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *MockedRepository) InsertRecordInSchema(value interface{}, table interface{}, tenantId string) error {
	args := m.Called(value, table, tenantId)
	return args.Error(0)
}

func (m *MockedRepository) IsValidCompanyName(companyName string) (isValid bool, err error) {
	args := m.Called(companyName)
	return args.Bool(1), args.Error(0)
}

// func TestCreateTenants() {
// 	ctx := context.Background()
// 	controllers.CreateTenant(ctx)
// }
