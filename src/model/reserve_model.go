package model

import (
	"time"

	"gorm.io/gorm"
)

type Reserve struct {
	gorm.Model
	Start        time.Time `gorm:"column:start;type:integer"`
	UserID       int
	User         User
	ConsultantID int
	Consultant   Consultant
}
