package models

import "time"

type Employee struct {
	Id               string    `json:"-" db:"id"`
	IdentityNumber   string    `json:"identityNumber" db:"identity_number"`
	Name             string    `json:"name" db:"name"`
	EmployeeImageUri string    `json:"employeeImageUri" db:"employee_image_uri"`
	Gender           string    `json:"gender" db:"gender"`
	DepartmentId     string    `json:"departmentId" db:"department_id"`
	CreatedAt        time.Time `json:"-" db:"created_at"`
}
