package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todolist/model"
	"todolist/usecase"

	"github.com/gorilla/mux"
)

type UserDelivery struct {
	uc usecase.UserUsecase
}

func NewUserDelivery(r *mux.Router, uc *usecase.UserUsecase) {
	handler := UserDelivery{uc: *uc}
	r.HandleFunc("/register", handler.Register).Methods("POST")
	r.HandleFunc("/login", handler.Login).Methods("POST")
}

func (ud *UserDelivery) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user = model.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	u, err := ud.uc.Register(user)
	if err != nil {
		w.Write([]byte("internal server error"))
	}
	json.NewEncoder(w).Encode(&u)
	w.Write([]byte("register success"))
}

func (ud *UserDelivery) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var login = model.Login{}
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	token, err := ud.uc.VerifyUserPassword(login.Email, login.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	fmt.Fprintln(w, "login-success-token=")

	w.Write([]byte(token))
}
