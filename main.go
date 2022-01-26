package main

import (
	"fmt"
	"net/http"
	"todolist/constants"
	"todolist/delivery"
	"todolist/repository"
	"todolist/usecase"
	"todolist/util"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	db := util.InitDB()
	repo := repository.NewUserRepo(db)
	usecase := usecase.NewUserUsecase(repo)
	_ = delivery.NewUserDelivery(r, usecase)
	fmt.Println("starting server at localhost :8080")
	http.ListenAndServe(constants.PORT, r)
}
