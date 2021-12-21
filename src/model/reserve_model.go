package model

import (
	"gorm.io/gorm"
)

type Reserve struct {
	gorm.Model
	Hour         int `gorm:"column:hour;type:integer"`
	UserID       int
	User         User
	SchuduleID   int
	Schudule     Schudule
}
