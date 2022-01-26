package delivery

import (
	"todolist/util"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type M map[string]interface{}

func InitGorilla() (*mux.Router, *gorm.DB) {
	r := mux.NewRouter()
	db_test := util.InitDB()
	return r, db_test
}
