package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Name  string `json:"name" gorm:"unique;not null"`
	Users []User `gorm:"many2many:group_users"`
}
