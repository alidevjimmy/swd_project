package model

import (
	"time"

	"gorm.io/gorm"
)

type Report struct {
	gorm.Model
	Subject string `gorm:"column:subject;not null;type:varchar(255)"`
	Long    *float64    `gorm:"column:long;type:float"`
	Lat     *float64    `gorm:"column:lat;type:float"`
	Address string `gorm:"column:address;type:varchar(255)"`
	Active  bool   `gorm:"column:active;not null;type:bool"`
	UserID  int    `gorm:"column:user_id;not null;type:integer"`
	User    User
	Until   time.Time `gorm:"column:until;type:timestamp;default:"`
}
