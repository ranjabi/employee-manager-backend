package types

type UpdateManagerProfilePayload struct {
	Email           *string `json:"email,omitempty" db:"email" validate:"omitempty,email"`
	Name            *string `json:"name,omitempty" db:"name" validate:"omitempty,min=4,max=52"`
	UserImageUri    *string `json:"userImageUri,omitempty" db:"user_image_uri" validate:"omitempty,uri"`
	CompanyName     *string `json:"companyName,omitempty" db:"company_name" validate:"omitempty,min=4,max=52"`
	CompanyImageUri *string `json:"companyImageUri,omitempty" db:"company_image_uri" validate:"omitempty,uri"`
}

type UpdateDepartmentProfilePayload struct {
	Name *string `json:"name,omitempty" db:"name" validate:"omitempty,min=4,max=33"`
}

type UpdateEmployeePayload struct {
	IdentityNumberNew *string `json:"identityNumber,omitempty" db:"identity_number" validate:"omitempty,min=5,max=33"`
	Name              *string `json:"name,omitempty" db:"name" validate:"omitempty,min=4,max=33"`
	EmployeeImageUri  *string `json:"employeeImageUri,omitempty" db:"employee_image_uri" validate:"omitempty,uri"`
	Gender            *string `json:"gender,omitempty" db:"gender" validate:"omitempty,oneof=male female"`
	DepartmentId      *string `json:"departmentId,omitempty" db:"department_id" validate:"omitempty"`
}
