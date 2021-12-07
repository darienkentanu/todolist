package usecase

import (
	"errors"
	"todolist/helper"
	"todolist/jwts"
	"todolist/model"
	"todolist/repository"
)

type UserImplementation interface {
	Register(user model.User) (model.User, error)
	VerifyUserPassword(email string, password string) (token string, err error)
}

type UserUsecase struct {
	repo repository.UserDB
}

func NewUserUsecase(repo *repository.UserDB) *UserUsecase {
	return &UserUsecase{repo: *repo}
}

func (uc *UserUsecase) Register(user model.User) (model.User, error) {
	pwd, err := helper.GenerateHashPassword(user.Password)
	if err != nil {
		return model.User{}, err
	}
	user.Password = pwd
	u, err := uc.repo.Register(user)
	if err != nil {
		return model.User{}, err
	}
	return u, nil
}

func (uc *UserUsecase) VerifyUserPassword(email string, password string) (token string, err error) {
	exist := uc.repo.CheckUserByEmail(email)
	if !exist {
		return "", err
	}
	userID, err := uc.repo.GetUserIDByEmail(email)
	if err != nil {
		return "", err
	}
	hash, err := uc.repo.GetPasswordHash(email)
	if err != nil {
		return "", err
	}
	ok := helper.CheckPasswordHash(password, hash)
	if !ok {
		return "", errors.New("invalid password")
	}
	token, err = jwts.CreateToken(userID)
	if err != nil {
		return "", err
	}
	uc.repo.UpdateToken(email, token)
	return token, nil
}
