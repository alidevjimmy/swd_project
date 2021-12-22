package model

import (
	"gorm.io/gorm"
)

type Consultant struct {
	gorm.Model
	Phone       string `gorm:"column:phone;unique;not null;type:varchar(255)"`
	Name        string `gorm:"column:name;not null;type:varchar(255)"`
	Family      string `gorm:"column:family;not null;type:varchar(255)"`
	Description string `gorm:"column:description;not null;type:text"`
}
