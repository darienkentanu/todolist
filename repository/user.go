package repository

import (
	"todolist/model"

	"gorm.io/gorm"
)

type UserRepo interface {
	Register(user model.User) (model.User, error)
	Login(user model.User) (model.User, error)
	CheckUserByEmail(email string) bool
	GetUserIDByEmail(email string) (uint, error)
	UpdateToken(email, token string) error
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

func (udb *UserDB) CheckUserByEmail(email string) bool {
	if row := udb.db.Where("email = ?", email).Find(&model.User{}).RowsAffected; row == 1 {
		return true
	}
	return false
}

func (udb *UserDB) GetUserIDByEmail(email string) (uint, error) {
	var user = model.User{}

	if err := udb.db.Where("email = ?", email).First(&user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (udb *UserDB) GetPasswordHash(email string) (string, error) {
	var user = model.User{}

	if err := udb.db.Where("email = ?", email).First(&user).Error; err != nil {
		return "", err
	}

	return user.Password, nil
}

func (udb *UserDB) UpdateToken(email, token string) error {
	var user = model.User{}

	if err := udb.db.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}
	if err := udb.db.Model(&user).Update("token", token).Error; err != nil {
		return err
	}
	return nil
}
