package delivery_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	. "todolist/delivery"
	"todolist/model"
	"todolist/repository"
	"todolist/usecase"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupUser(db *gorm.DB) {
	err := db.Exec("drop table todo_category").Error
	if err != nil {
		fmt.Println(err)
	}
	db.Migrator().DropTable(&model.Todo{})
	db.Migrator().DropTable(&model.Category{})
	db.Migrator().DropTable(&model.User{})

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Todo{})
}

// func insertDataUser(db *gorm.DB) error {
// 	register := model.User{
// 		Name:     "darien kentanu",
// 		Email:    "darienkentanu@gmail.com",
// 		Password: "password",
// 	}

// 	hashPassword, err := helper.GenerateHashPassword(register.Password)
// 	if err != nil {
// 		return err
// 	}

// 	user := model.User{
// 		Name:     register.Name,
// 		Email:    register.Email,
// 		Password: hashPassword,
// 	}

// 	if err = db.Save(&user).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

func TestRegisterUser(t *testing.T) {

	var testCases = []struct {
		name       string
		path       string
		expectCode int
		response   string
		reqBody    M
	}{
		{
			name:       "RegisterUser",
			path:       "/register",
			expectCode: http.StatusOK,
			response:   "darien kentanu",
			reqBody: M{
				"name":     "darien kentanu",
				"email":    "darienkentanu@gmail.com",
				"password": "password",
			},
		},
		{
			name:       "RegisterUserError",
			path:       "/register",
			expectCode: http.StatusBadRequest,
			response:   "",
			reqBody: M{
				"name":     "darien kentanu",
				"email":    "darienkentanu@gmail.com",
				"password": "password",
			},
		},
	}

	r, db := InitGorilla()
	setupUser(db)
	repo := repository.NewUserRepo(db)
	uc := usecase.NewUserUsecase(repo)
	handler := NewUserDelivery(r, uc)

	for _, testCase := range testCases {
		register, err := json.Marshal(testCase.reqBody)
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(register))

		t.Run(testCase.name, func(t *testing.T) {

			handler.Register(w, r)
			var response model.User
			json.NewDecoder(w.Body).Decode(&response)
			// if testCase.response != response.Name {
			// 	t.Errorf("salah seharusnya %v", testCase.response)
			// }
			assert.Equal(t, testCase.response, response.Name, "salah seharusnya %v", testCase.response)
			// if testCase.expectCode != w.Code {
			// 	t.Errorf("salah seharusnya %v", testCase.expectCode)
			// }
			assert.Equal(t, testCase.expectCode, w.Code, "salah seharusnya %v", w.Code)
		})
	}
}
