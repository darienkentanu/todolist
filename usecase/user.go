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
	GetUserIDFromCookies(r *http.Request) (uint, error)
	ValidationForAddTodos(input *model.Todo, userID uint) (todos *model.Todo, err error)
	AddCategory(input *model.Category) (output model.Category, err error)
	SaveTodo(input *model.Todo, userID uint) (output model.Todo, err error)
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

func (uc *userUsecase) GetUserIDFromCookies(r *http.Request) (uint, error) {
	id, err := middlewares.ExtractTokenUserID(r)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (uc *userUsecase) ValidationForAddTodos(input *model.Todo, userID uint) (todos *model.Todo, err error) {
	input.UserID = userID
	if input.Title == "" {
		return nil, errors.New("please input a title")
	}
	if input.Description == "" {
		return nil, errors.New("please input a description")
	}
	if input.Categories == nil {
		input.Categories = append(input.Categories, model.Category{Name: "uncategorized"})
	}
	return input, nil
}

func (uc *userUsecase) SaveTodo(input *model.Todo, userID uint) (output model.Todo, err error) {
	// insert data to categories table
	var categoriesID []uint
	for i := range input.Categories {
		exist := uc.repo.CheckCategory(&input.Categories[i])
		if !exist {
			category, err := uc.repo.AddCategory(&input.Categories[i])
			if err != nil {
				return output, err
			}
			categoriesID = append(categoriesID, category.ID)
		}
		categoriesID = append(categoriesID, input.Categories[i].ID)
	}
	// insert data to todo table
	output, err = uc.repo.SaveTodo(input)
	if err != nil {
		return output, err
	}
	// insert id to todo_category table
	todoID := output.ID
	// var todo_category = []struct {
	// 	todo_id     uint `gorm:"primarykey`
	// 	category_id uint `gorm:"primarykey`
	// }{}
	err = uc.repo.InsertTodoCategory(todoID, categoriesID)
	if err != nil {
		return output, err
	}
	output.Categories = append(output.Categories, input.Categories...)
	return output, nil
}

func (uc *userUsecase) AddCategory(input *model.Category) (output model.Category, err error) {
	exist := uc.repo.CheckCategory(input)
	if !exist {
		output, err = uc.repo.AddCategory(input)
		if err != nil {
			return output, err
		}
		return output, nil
	}
	return *input, nil
}
