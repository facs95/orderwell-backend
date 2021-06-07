package entity

//Employee within an organization
type EmployeeRole struct {
	Role      string     `json:"role" gorm:"primary_key"`
	Employees []Employee `json:"employees" gorm:"foreignKey:Role;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
