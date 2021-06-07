package repository

import (
	"errors"
	"fmt"
	"os"

	"github.com/facs95/orderwell-backend/entity"
	"github.com/facs95/orderwell-backend/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type postgressRepo struct{}

var (
	host      = os.Getenv("DB_HOST")
	port      = os.Getenv("DB_PORT")
	user      = os.Getenv("DB_USER")
	password  = os.Getenv("DB_PW")
	dbname    = os.Getenv("DB_NAME")
	dbSslMode = os.Getenv("DB_SSL")
)

var db *gorm.DB

func InitPostgressRepo() Repository {
	err := gormConnectDatabase()
	if err != nil {
		logger.ErrorLogger.Fatal(err)
	}
	return &postgressRepo{}
}

//GormConnectDatabase to postgress database through gorm
func gormConnectDatabase() error {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, dbSslMode)
	db, err = gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		logger.ErrorLogger.Fatal(err)
	}
	db.AutoMigrate(&entity.Tenant{})
	logger.InfoLogger.Print("Database Connection Sucessfull")
	return nil
}

func (*postgressRepo) CreateTenant(tenant *entity.Tenant) error {
	//Create record
	if err := createRecord(&tenant); err != nil {
		return err
	}

	//Create the schema for the tenant - Independant tables
	schemaName := tenant.ID
	if err := createSchema(schemaName); err != nil {
		deleteRecord(&entity.Tenant{}, tenant.ID)
		return err
	}

	//Run the migration for that schema
	if err := runSchemaMigration(schemaName); err != nil {
		deleteRecord(&entity.Tenant{}, tenant.ID)
		return err
	}

	return nil
}

func (*postgressRepo) DeleteTenant(id string) error {
	if err := deleteRecord(&entity.Tenant{}, id); err != nil {
		return err
	}
	if err := deleteSchema(id); err != nil {
		return err
	}
	return nil
}

//TODO do not like having to pass 3 values here
//InesertRecord in Table
func (*postgressRepo) SaveInTenant(tenantId string, table interface{}, value interface{}) error {
	name := getSchemaTableName(table, tenantId)
	result := db.Table(name).Create(value)
	if result.Error != nil {
		logger.ErrorLogger.Println("Following issue creating the record: ", result.Error)
		return result.Error
	}
	logger.TenantInfoLog(tenantId).Println("Employee created successfully")
	return nil
}

//Get Tenant by subdomain
func (*postgressRepo) FindTenantBySubDomain(subDomain string, value interface{}) error {
	//The subdomain is going to be the same as the oauth_id for now
	tx := db.First(value, "oauth_id = ?", subDomain)
	return tx.Error
}

//Get Record in table by ID
func (*postgressRepo) FindTenant(id string, value interface{}) error {
	tx := db.First(value, "id = ?", id)
	return tx.Error
}

//Checks if company name already exists
func (*postgressRepo) ValidateTenantCompanyName(companyName string) (isValid bool, err error) {
	tenant := entity.Tenant{}
	result := db.Where(&entity.Tenant{CompanyName: companyName}, "company_name").First(&tenant)
	if result.Error != nil {
		if isNotFound := errors.Is(result.Error, gorm.ErrRecordNotFound); isNotFound == true {
			return true, nil
		}
		logger.ErrorLogger.Println("There was an issue accessing database: ", err)
		return false, result.Error
	}
	return false, nil
}

//Get specific table in schema.
//Results: Should be a pointer
func (*postgressRepo) FindAllInTenant(tenantId string, table interface{}, results interface{}) error {
	name := getSchemaTableName(table, tenantId)
	tx := db.Table(name).Find(results)
	return tx.Error
}

//CreateRecord creates new record in table
func createRecord(value interface{}) error {
	result := db.Create(value)
	if result.Error != nil {
		logger.ErrorLogger.Println("Following issue creating the record: ", result.Error)
		return result.Error
	}
	logger.InfoLogger.Println("Record created succesfully")
	return nil
}

//CreateSchema on postgress database
func createSchema(name string) error {
	//We need to turn the name into a hash for schema name requirements
	if err := db.Exec("CREATE SCHEMA " + name).Error; err != nil {
		logger.ErrorLogger.Println("There was a problem creating the schema ", err)
		return err
	}
	logger.InfoLogger.Println("Schema created succesfully ", name)
	return nil
}

//CreateSchema on postgress database
func deleteSchema(name string) error {
	//We need to turn the name into a hash for schema name requirements
	cmd := fmt.Sprintf("DROP SCHEMA %v CASCADE", name)
	if err := db.Exec(cmd).Error; err != nil {
		logger.ErrorLogger.Println("There was a problem deleting the schema ", err)
		return err
	}
	logger.InfoLogger.Println("Schema deleted succesfully ", name)
	return nil
}

func deleteRecord(table interface{}, id string) error {
	result := db.Where("id = ?", id).Delete(table)
	if result.Error != nil {
		logger.ErrorLogger.Println("Following issue deleting the record: ", result.Error)
		return result.Error
	}
	logger.InfoLogger.Println("Record deleted succesfully ", id)
	return nil
}

//Get Record in table by ID
func getRecordById(table interface{}, id string) (tx *gorm.DB) {
	return db.First(table, "id = ?", id)
}

//TODO do not like having to pass 3 values here
//InesertRecord in Table
func insertRecordInSchema(value interface{}, table interface{}, tenantId string) error {
	name := getSchemaTableName(table, tenantId)
	result := db.Table(name).Create(value)
	if result.Error != nil {
		logger.ErrorLogger.Println("Following issue creating the record: ", result.Error)
		return result.Error
	}
	logger.TenantInfoLog(tenantId).Println("Employee created successfully")
	return nil
}

//Get specific table in schema.
//Results: Should be a pointer
func getTableWithinSchema(results interface{}, table interface{}, tenantId string) (tx *gorm.DB) {
	name := getSchemaTableName(table, tenantId)
	return db.Table(name).Find(results)
}

//runSchemaMigration runs migrations needed for new schema
func runSchemaMigration(name string) error {
	//Need to generate new connection to configure names
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, dbSslMode)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   name + ".",
			SingularTable: true,
		},
	})
	if err != nil {
		logger.InfoLogger.Println("Error connecting to the database", err)
		return err
	}

	DB, err := db.DB()
	if err != nil {
		logger.InfoLogger.Println("Error Initializing the database instance", err)
		return err
	}
	//Close database connection specific to run tenant migrations
	defer DB.Close()

	err = db.AutoMigrate(&entity.EmployeeRole{}, &entity.Employee{}, &entity.Service{}, &entity.Order{}, &entity.Client{})
	if err != nil {
		message := fmt.Sprintf("Error running schema migration for schema:%v with:%v ", name, err.Error())
		logger.ErrorLogger.Println(message)
		return err
	}

	err = populateRolesTable(name)
	if err != nil {
		message := fmt.Sprintf("Error generating roles for schema:%v with:%v ", name, err.Error())
		logger.ErrorLogger.Println(message)
		return err
	}
	logger.InfoLogger.Println("Schema migration ran succesfully for ID -> ", name)
	return nil
}

func getEmployeeRoles() [3]entity.EmployeeRole {
	return [3]entity.EmployeeRole{{Role: "owner"}, {Role: "admin"}, {Role: "employee"}}
}

func populateRolesTable(tenantId string) error {
	roles := getEmployeeRoles()
	return insertRecordInSchema(&roles, entity.EmployeeRole{}, tenantId)
}
