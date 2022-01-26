package config

import (
	"sync"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	DB_USERNAME string
	DB_PASSWORD string
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string
}

var lock = &sync.Mutex{}
var appConfig *AppConfig

func GetConfig() *AppConfig {
	lock.Lock()
	defer lock.Unlock()

	if appConfig == nil {
		appConfig = initConfig()
	}

	return appConfig
}

func initConfig() *AppConfig {
	var defaultConfig AppConfig
	defaultConfig.DB_USERNAME = "root"
	defaultConfig.DB_PASSWORD = "password"
	defaultConfig.DB_HOST = "localhost"
	defaultConfig.DB_PORT = "3306"
	defaultConfig.DB_NAME = "todolist_test"

	config, err := godotenv.Read()

	if err != nil {
		return &defaultConfig
	}

	var finalConfig AppConfig
	finalConfig.DB_USERNAME = config["DB_USERNAME"]
	finalConfig.DB_PASSWORD = config["DB_PASSWORD"]
	finalConfig.DB_HOST = config["DB_HOST"]
	finalConfig.DB_PORT = config["DB_PORT"]
	finalConfig.DB_NAME = config["DB_NAME"]

	return &finalConfig
}
