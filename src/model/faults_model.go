package model

import "gorm.io/gorm"

type FaultModel struct {
	gorm.Model
	Title    string         `gorm:"column:title;unique;not null;type:varchar(255)"`
	Commands []CommandModel `gorm:"many2many:fault_commands;"`
}
