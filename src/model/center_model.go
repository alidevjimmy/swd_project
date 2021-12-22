package model

import (
	"gorm.io/gorm"
)

type Center struct {
	gorm.Model
	Name      string `gorm:"column:name;not null;type:varchar(255)"`
	Website   string `gorm:"column:website;not null;type:varchar(255)"`
	Telephone string `gorm:"column:telephone;not null;type:varchar(255)"`
	Sms       string `gorm:"column:sms;not null;type:text"`
}
