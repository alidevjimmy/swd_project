package model

import "gorm.io/gorm"

type CommandModel struct {
	gorm.Model
	Title string `gorm:"column:title;unique;not null;type:varchar(255)"`
	Auto  bool   `gorm:"column:auto;not null;type:varchar(255)"`
}
