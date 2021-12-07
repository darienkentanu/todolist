package main

import (
	"fmt"
	"net/http"
	"todolist/config"
	"todolist/delivery"
	"todolist/repository"
	"todolist/usecase"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	db := config.InitDB()
	repo := repository.NewUserRepo(db)
	usecase := usecase.NewUserUsecase(repo)
	delivery.NewUserDelivery(r, usecase)
	fmt.Println("starting server at localhost :8080")
	http.ListenAndServe(config.PORT, r)
}
