package usecase

import (
	"errors"
	"net/http"
	"todolist/helper"
	"todolist/middlewares"
	"todolist/model"
	"todolist/repository"
)

type UserUsecase interface {
	Register(user model.User) (model.User, error)
	VerifyUserPassword(w http.ResponseWriter, r *http.Request, email string, password string) (err error)
	CreateCookie(token string) *http.Cookie
	RemoveCookie(w http.ResponseWriter) string
	CreateToken(userID uint) (string, error)
	GetUserIDByEmail(email string) (uint, error)
}

type userUsecase struct {
	repo repository.UserRepo
}

func NewUserUsecase(repo repository.UserRepo) *userUsecase {
	return &userUsecase{repo: repo}
}

func (uc *userUsecase) GetUserIDByEmail(email string) (uint, error) {
	userID, err := uc.repo.GetUserIDByEmail(email)
	if err != nil {
		return 0, err
	}
	return userID, nil
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

func (uc *userUsecase) VerifyUserPassword(w http.ResponseWriter, r *http.Request, email string, password string) (err error) {
	exist := uc.repo.CheckUserByEmail(email)
	if !exist {
		return err
	}
	hash, err := uc.repo.GetPasswordHash(email)
	if err != nil {
		return err
	}
	ok := helper.CheckPasswordHash(password, hash)
	if !ok {
		return errors.New("invalid password")
	}
	return nil
}

func (uc *userUsecase) CreateToken(userID uint) (string, error) {
	token, err := middlewares.CreateToken(userID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (uc *userUsecase) CreateCookie(token string) *http.Cookie {
	cookie := &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: middlewares.EXPIRATION_TIME,
	}
	return cookie
}

func (uc *userUsecase) RemoveCookie(w http.ResponseWriter) string {
	c := http.Cookie{
		Name:   "token",
		MaxAge: -1}
	http.SetCookie(w, &c)
	return "cookie deleted"
}
