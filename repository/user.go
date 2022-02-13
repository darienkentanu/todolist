package repository

import (
	"todolist/model"

	"gorm.io/gorm"
)

type UserRepo interface {
	Register(user model.User) (model.User, error)
	// Login(user model.User) (model.User, error)
	CheckUserByEmail(email string) bool
	GetUserIDByEmail(email string) (uint, error)
	GetPasswordHash(email string) (string, error)
	// UpdateToken(email, token string) error
	CheckCategory(input *model.Category) bool
	AddCategory(input *model.Category) (output model.Category, err error)
	AddTodo(input *model.Todo) (output model.Todo, err error)
	SaveTodo(input *model.Todo) (output model.Todo, err error)
	InsertTodoCategory(todoID uint, categoriesID []uint) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *userRepo {
	return &userRepo{db: db}
}

func (udb *userRepo) Register(user model.User) (model.User, error) {
	if err := udb.db.Create(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (udb *userRepo) CheckUserByEmail(email string) bool {
	if row := udb.db.Where("email = ?", email).Find(&model.User{}).RowsAffected; row == 1 {
		return true
	}
	return false
}

func (udb *userRepo) GetUserIDByEmail(email string) (uint, error) {
	var user = model.User{}

	if err := udb.db.Where("email = ?", email).First(&user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (udb *userRepo) GetPasswordHash(email string) (string, error) {
	var user = model.User{}

	if err := udb.db.Where("email = ?", email).First(&user).Error; err != nil {
		return "", err
	}

	return user.Password, nil
}

func (udb *userRepo) CheckCategory(input *model.Category) bool {
	var output model.Category
	if rows := udb.db.Where("name=?", input.Name).First(&output).RowsAffected; rows != 0 {
		return true
	}
	return false
}

func (udb *userRepo) AddCategory(input *model.Category) (output model.Category, err error) {
	if err := udb.db.Save(input).Error; err != nil {
		return model.Category{}, err
	}
	if err = udb.db.Where("id=?", input.ID).First(&output).Error; err != nil {
		return model.Category{}, err
	}
	return output, nil
}

func (udb *userRepo) AddTodo(input *model.Todo) (output model.Todo, err error) {
	if err := udb.db.Save(input).Error; err != nil {
		return model.Todo{}, err
	}
	if err := udb.db.Where("id=?", input.ID).First(&output).Error; err != nil {
		return model.Todo{}, err
	}
	return output, nil
}

func (udb *userRepo) SaveTodo(input *model.Todo) (output model.Todo, err error) {
	if err := udb.db.Save(input).Error; err != nil {
		return output, err
	}
	if err := udb.db.Where("id=?", input.ID).First(&output).Error; err != nil {
		return output, err
	}
	return output, nil
}

// type todo_category struct {
// 	todo_id     uint `gorm:"primarykey"`
// 	category_id uint `gorm:"primarykey"`
// }

func (udb *userRepo) InsertTodoCategory(todoID uint, categoriesID []uint) error {
	for i := range categoriesID {
		// var input = todo_category{
		// 	todo_id:     todoID,
		// 	category_id: categoriesID[i],
		// }
		if err := udb.db.Raw("insert into todo_category (todo_id, category_id) values (?,?)", todoID, categoriesID[i]).Error; err != nil {
			return err
		}
	}
	return nil
}

// func (udb *UserDB) UpdateToken(email, token string) error {
// 	var user = model.User{}

// 	if err := udb.db.Where("email = ?", email).First(&user).Error; err != nil {
// 		return err
// 	}
// 	if err := udb.db.Model(&user).Update("token", token).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }
