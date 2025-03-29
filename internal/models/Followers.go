package models

import "gorm.io/gorm"

type Followers struct {
	gorm.Model
	FollowerID  uint `json:"follower_id" gorm:"not null"`
	FollowingID uint `json:"following_id" gorm:"not null"`
}
