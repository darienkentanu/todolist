package usecase

import (
	"errors"
	"net/http"
	"todolist/helper"
	"todolist/jwts"
	"todolist/model"
	"todolist/repository"
)

type UserUsecase interface {
	Register(user model.User) (model.User, error)
	VerifyUserPassword(w http.ResponseWriter, r *http.Request, email string, password string) (token, message string, err error)
}

type userUsecase struct {
	repo repository.UserRepo
}

func NewUserUsecase(repo repository.UserRepo) *userUsecase {
	return &userUsecase{repo: repo}
}

func (uc *userUsecase) Register(user model.User) (model.User, error) {
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

func (uc *userUsecase) VerifyUserPassword(w http.ResponseWriter, r *http.Request, email string, password string) (token, message string, err error) {
	exist := uc.repo.CheckUserByEmail(email)
	if !exist {
		return "", "", err
	}
	userID, err := uc.repo.GetUserIDByEmail(email)
	if err != nil {
		return "", "", err
	}
	hash, err := uc.repo.GetPasswordHash(email)
	if err != nil {
		return "", "", err
	}
	ok := helper.CheckPasswordHash(password, hash)
	if !ok {
		return "", "", errors.New("invalid password")
	}
	// token, err = jwts.ExtractTokenUserID(r)
	// if err == nil {
	// 	return token, "token exist", nil
	// }
	message = helper.RemoveCookie(w)
	token, err = jwts.CreateToken(userID)
	if err != nil {
		return "", message, err
	}
	// uc.repo.UpdateToken(email, token)
	return token, message, nil
}
