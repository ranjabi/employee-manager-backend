package models

import (
	"time"
)

type Department struct {
	Id        string    `json:"departmentId" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	ManagerId string	`json:"-" db:"manager_id"`
}
