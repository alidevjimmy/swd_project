package model

import (
	"time"

	"gorm.io/gorm"
)

type Report struct {
	gorm.Model
	Subject string `gorm:"column:subject;type:varchar(255)"`
	Long    float64    `gorm:"column:long;type:float"`
	Lat     float64    `gorm:"column:lat;type:float"`
	Address string `gorm:"column:address;type:varchar(255)"`
	Active  bool   `gorm:"column:active;type:bool"`
	UserID  int
	User    User
	Until   time.Time `gorm:"column:until;type:timestamp;default:"`
}
