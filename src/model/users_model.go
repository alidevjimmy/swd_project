package model

import "gorm.io/gorm"

type UserModel struct {
	gorm.Model
	Phone    string `gorm:"column:phone;unique;not null;type:varchar(255)"`
	FullName string `gorm:"column:full_name;not null;type:varchar(255)"`
	Password string `gorm:"column:password;not null;type:varchar(255)"`
}
