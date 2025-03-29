package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Id        int32       `gorm:"primaryKey;autoIncrement"`
	FirstName string      `json:"first_name" gorm:"not null"`
	LastName  string      `json:"last_name" gorm:"not null"`
	Email     string      `json:"email" gorm:"unique;not null"`
	Password  string      `json:"password" gorm:"not null"`
	Posts     []Post      `gorm:"foreignKey:UserID"`
	Comments  []Comment   `gorm:"foreignKey:UserID"`
	Likes     []Like      `gorm:"foreignKey:UserID"`
	Followers []Followers `gorm:"foreignKey:FollowerID"`
	Following []Followers `gorm:"foreignKey:FollowingID"`
}
