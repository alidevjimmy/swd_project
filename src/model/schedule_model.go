package model

import (
	"time"

	"gorm.io/gorm"
)

type Schedule struct {
	gorm.Model
	Start        time.Time `gorm:"column:start;index;type:timestamp"`
	End          time.Time `gorm:"column:end;index;type:timestamp"`
	Each         int       `gorm:"column:each;type:integer"`
	ConsultantID int       `gorm:"column:consultant_id;not null;type:integer"`
	Consultant   Consultant
}
