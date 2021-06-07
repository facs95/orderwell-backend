package entity

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"gorm.io/datatypes"
)

//Service within an organization
type Service struct {
	ID         string         `json:"id" gorm:"primary_key"`
	Name       string         `json:"name" gorm:"unique"`
	CreatedAt  time.Time      `json:"createdAt"`
	Data       datatypes.JSON `json:"data"`
	Orders     []Order        `json:"orders" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	EmployeeId string         `json:"employeeId"`
}

//Validate Service body params
func (t *Service) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Name, validation.Required),
	)
}
