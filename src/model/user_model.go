package model

import (
	"time"

	"gorm.io/gorm"
)

type UserStatus string

const (
	Green UserStatus = "GREEN"
	Red   UserStatus = "RED"
)

type User struct {
	gorm.Model
	Phone        string     `gorm:"column:phone;unique;not null;type:varchar(255)"`
	Name         string     `gorm:"column:name;not null;type:varchar(255)"`
	Family       string     `gorm:"column:family;not null;type:varchar(255)"`
	NationalCode string     `gorm:"column:national_code;type:varchar(255)"`
	Password     string     `gorm:"column:password;not null;type:varchar(255)"`
	Birth        time.Time  `gorm:"column:birth;not null;type:timestamp"`
	Status       UserStatus `gorm:"column:status;not null;type:varchar(255);default:GREEN"`
}
