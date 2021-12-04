package repository

import (
	"todolist/model"

	"gorm.io/gorm"
)

type UserRepo interface {
	Register(user model.User) (model.User, error)
}

type UserDB struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserDB {
	return &UserDB{db: db}
}

func (udb *UserDB) Register(user model.User) (model.User, error) {
	if err := udb.db.Create(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
