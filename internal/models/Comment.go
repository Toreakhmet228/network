package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	UserID  uint   `json:"user_id" gorm:"not null"`
	PostID  uint   `json:"post_id" gorm:"not null"`
	Content string `json:"content" gorm:"not null"`
}
