package entity

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

//Employee within an organization
type Employee struct {
	ID          string    `json:"id" gorm:"primary_key"`
	Email       string    `json:"email" gorm:"unique"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	PhoneNumber string    `json:"phoneNumber"`
	CreatedAt   time.Time `json:"createdAt"`
	Role        string    `json:"role"`
	// Services    []Service `json:"services" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

//Validate Employee body params
func (t *Employee) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.FirstName, validation.Required),
		validation.Field(&t.LastName, validation.Required),
		validation.Field(&t.Email, validation.Required, is.Email),
		validation.Field(&t.PhoneNumber, validation.Required, is.E164),
	)
}
