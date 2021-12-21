package model

import (
	"time"

	"gorm.io/gorm"
)

type Report struct {
	gorm.Model
	Subject string `gorm:"column:subject;type:varchar(255)"`
	Long    int    `gorm:"column:long;type:integer"`
	Lat     int    `gorm:"column:lat;type:integer"`
	Address string `gorm:"column:address;type:varchar(255)"`
	Active  bool   `gorm:"column:active;type:varchar(255)"`
	UserID  int
	User    User
	Until   time.Time `gorm:"column:until;type:timestamp;default:"`
}
