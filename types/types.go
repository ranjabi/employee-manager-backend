package types

import "encoding/json"

type UpdateManagerProfilePayload struct {
	EmailRaw           json.RawMessage `json:"email,omitempty"`
	NameRaw            json.RawMessage `json:"name,omitempty"`
	UserImageUriRaw    json.RawMessage `json:"userImageUri,omitempty"`
	CompanyNameRaw     json.RawMessage `json:"companyName,omitempty"`
	CompanyImageUriRaw json.RawMessage `json:"companyImageUri,omitempty"`
	Email              *string         `db:"email" validate:"omitempty,email"`
	Name               *string         `db:"name" validate:"omitempty,min=4,max=52"`
	UserImageUri       *string         `db:"user_image_uri" validate:"omitempty,url"`
	CompanyName        *string         `db:"company_name" validate:"omitempty,min=4,max=52"`
	CompanyImageUri    *string         `db:"company_image_uri" validate:"omitempty,uri"`
}

type UpdateDepartmentProfilePayload struct {
	NameRaw json.RawMessage `json:"name,omitempty"`
	Name    *string         `db:"name" validate:"omitempty,min=4,max=33"`
}

type UpdateEmployeePayload struct {
	IdentityNumberNewRaw json.RawMessage `json:"identityNumber,omitempty"`
	NameRaw              json.RawMessage `json:"name,omitempty"`
	EmployeeImageUriRaw  json.RawMessage `json:"employeeImageUri,omitempty"`
	GenderRaw            json.RawMessage `json:"gender,omitempty"`
	DepartmentIdRaw      json.RawMessage `json:"departmentId,omitempty"`
	IdentityNumberNew    *string         `db:"identity_number" validate:"omitempty,min=5,max=33"`
	Name                 *string         `db:"name" validate:"omitempty,min=4,max=33"`
	EmployeeImageUri     *string         `db:"employee_image_uri" validate:"omitempty,uri"`
	Gender               *string         `db:"gender" validate:"omitempty,oneof=male female"`
	DepartmentId         *string         `db:"department_id" validate:"omitempty"`
}
