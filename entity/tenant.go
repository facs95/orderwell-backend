package entity

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

//Tenant handles new organizatios
type Tenant struct {
	ID          string    `json:"id" gorm:"primary_key"`
	Email       string    `json:"email"`
	CompanyName string    `json:"companyName" gorm:"unique;size:20"`
	OauthId     string    `json:"oauthId" gorm:"unique"`
	CreatedAt   time.Time `json:"createdAt"`
}

type TenantWithUser struct {
	Tenant
	Password    string `json:"password"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhoneNumber string `json:"phoneNumber"`
}

//Validate tenant body params
func (t TenantWithUser) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.CompanyName, validation.Required, validation.Length(4, 20)), //To be able to comply with google
		validation.Field(&t.Email, validation.Required, is.Email),
		validation.Field(&t.Password, validation.Required, validation.Length(6, 0)),
		validation.Field(&t.FirstName, validation.Required),
		validation.Field(&t.LastName, validation.Required),
		validation.Field(&t.PhoneNumber, validation.Required, is.E164),
	)
}
