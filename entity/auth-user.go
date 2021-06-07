package entity

type User struct {
	ID       string
	TenantId string
}

type NewAuthUser struct {
	Email       string
	Password    string
	DisplayName string
	PhoneNumber string
}
