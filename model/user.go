package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(100);not null" json:"name" form:"name"`
	Email    string `gorm:"type:varchar(100);unique;not null" json:"email" form:"email"`
	Password string `gorm:"type:varchar(255);not null" json:"password" form:"password"`
	Token    string `gorm:"type:varchar(255)" json:"token" form:"token"`
	Todos    []Todo `json:"-"`
}

type Login struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}
