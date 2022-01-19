package model

import (
	"time"

	"gorm.io/gorm"
)

type Reserve struct {
	gorm.Model
	Start        time.Time `gorm:"column:start;not null;type:timestamp"`
	UserID       int
	User         User
	ConsultantID int
	Consultant   Consultant
}
