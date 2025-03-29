package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	UserID   uint      `json:"user_id" gorm:"not null"`
	Title    string    `json:"title" gorm:"not null"`
	Body     string    `json:"body" gorm:"not null"`
	ImageURL string    `json:"image_url"`
	Comments []Comment `gorm:"foreignKey:PostID"`
	Likes    []Like    `gorm:"foreignKey:PostID"`
}
