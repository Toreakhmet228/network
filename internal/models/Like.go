package models

import "gorm.io/gorm"

type Like struct {
	gorm.Model
	UserID uint `json:"user_id" gorm:"not null"`
	PostID uint `json:"post_id" gorm:"not null"`
}
