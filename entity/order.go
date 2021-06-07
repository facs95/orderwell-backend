package entity

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

//Order within an organization
type Order struct {
	ID        string    `json:"id" gorm:"primary_key"`
	ServiceId string    `json:"serviceId"`
	CreatedAt time.Time `json:"createdAt"`
}

//Validate Order body params
func (t *Order) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.ServiceId, validation.Required),
	)
}
