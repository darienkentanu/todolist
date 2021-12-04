package usecase

import (
	"todolist/model"
	"todolist/repository"
)

type UserImplementation interface {
	Register(user model.User) (model.User, error)
}

type UserUsecase struct {
	repo repository.UserDB
}

func NewUserUsecase(repo *repository.UserDB) *UserUsecase {
	return &UserUsecase{repo: *repo}
}

func (uc *UserUsecase) Register(user model.User) (model.User, error) {
	u, err := uc.repo.Register(user)
	if err != nil {
		return model.User{}, err
	}
	return u, nil
}
