package model

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	ImageUrl string `gorm:"column:image_url;not null;type:varchar(255)"`
	Title    string `gorm:"column:title;not null;type:varchar(255)"`
	Abstract string `gorm:"column:abstract;not null;type:text"`
	Body     string `gorm:"column:body;not null;type:text"`
}
