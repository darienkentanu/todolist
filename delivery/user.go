package delivery

import (
	"encoding/json"
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
}

func (ud *UserDelivery) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user = model.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.Write([]byte("bad request"))
	}
	u, err := ud.uc.Register(user)
	if err != nil {
		w.Write([]byte("internal server error"))
	}
	json.NewEncoder(w).Encode(&u)
}
