package types

type UpdateManagerProfilePayload = struct {
	Email           *string `json:"email,omitempty" db:"email" validate:"omitempty,email"`
	Name            *string `json:"name,omitempty" db:"name" validate:"omitempty,min=4,max=52"`
	UserImageUri    *string `json:"userImageUri,omitempty" db:"user_image_uri" validate:"omitempty,uri"`
	CompanyName     *string `json:"companyName,omitempty" db:"company_name" validate:"omitempty,min=4,max=52"`
	CompanyImageUri *string `json:"companyImageUri,omitempty" db:"company_image_uri" validate:"omitempty,uri"`
}