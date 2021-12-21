package model

import (
	"time"

	"gorm.io/gorm"
)

type Schudule struct {
	gorm.Model
	Date         time.Time `gorm:"column:date;index;type:timestamp"`
	StartHour    int       `gorm:"column:start_hour;type:integer"`
	FinishHour   int       `gorm:"column:finish_hour;type:integer"`
	ConsultantID int
	Consultant   Consultant
}
