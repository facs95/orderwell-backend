package entity

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

//Tenant handles new organizatios
type Client struct {
	ID           string    `json:"id" gorm:"primary_key"`
	PrimaryEmail string    `json:"primaryEmail"`
	PhoneNumber  string    `json:"phoneNumber"`
	CreatedAt    time.Time `json:"createdAt"`
}

//Validate Client body params
func (t Client) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.PhoneNumber, validation.Required, is.E164),
		validation.Field(&t.PrimaryEmail, validation.Required, is.Email),
	)
}
