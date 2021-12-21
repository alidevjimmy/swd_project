package model

import "gorm.io/gorm"

type TokenModel struct {
	gorm.Model
	Token  string    `gorm:"column:token;not null;type:varchar(255)"`
	UserID uint      `gorm:"column:user_id"`
	User   UserModel `gorm:"foreignKey:UserID"`
}
