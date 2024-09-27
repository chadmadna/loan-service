package models

import (
	"loan-service/services/auth"

	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name     string        `json:"name"`
	RoleType auth.RoleType `json:"role_type" gorm:"column:role_type"`
}

func (Role) TableName() string {
	return "roles"
}
