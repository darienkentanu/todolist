package model

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	Title       string `gorm:"type:varchar(100);not null" json:"title" form:"title"`
	Description string `gorm:"type:varchar(255);not null" json:"description" form:"description"`
	UserID      uint
	Categories  []Category `gorm:"many2many:todo_category;"`
}
