package config

import (
	"fmt"

	"todolist/model"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetConfig() (config map[string]string) {
	conf, err := godotenv.Read()
	if err != nil {
		panic(err)
	}
	return conf
}

func InitDB() *gorm.DB {
	conf := GetConfig()
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		conf["DB_USERNAME"],
		conf["DB_PASSWORD"],
		conf["DB_HOST"],
		conf["DB_PORT"],
		conf["DB_NAME"],
	)

	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	initMigration(db)
	return db
}

func initMigration(db *gorm.DB) {
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Todo{})
}
