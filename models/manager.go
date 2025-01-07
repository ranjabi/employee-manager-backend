package models

import "github.com/jackc/pgx/v5/pgtype"

type Manager struct {
	Id              string      `json:"id" db:"id"`
	Email           string      `json:"email" db:"email"`
	Password        string      `json:"password" db:"password"`
	Name            pgtype.Text `json:"name" db:"name"`
	UserImageUri    pgtype.Text `json:"userImageUri" db:"user_image_uri"`
	CompanyName     pgtype.Text `json:"companyName" db:"company_name"`
	CompanyImageUri pgtype.Text `json:"companyImageUri" db:"company_image_uri"`
	Token           string      `json:"token" db:"-"`
}
