package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todolist/constants"
	"todolist/model"
	"todolist/usecase"

	"github.com/gorilla/mux"
)

type userDelivery struct {
	uc usecase.UserUsecase
}

func NewUserDelivery(r *mux.Router, uc usecase.UserUsecase) *userDelivery {
	handler := userDelivery{uc: uc}
	r.HandleFunc("/register", handler.Register).Methods("POST")
	r.HandleFunc("/login", handler.Login).Methods("POST")
	return &handler
}

func (ud *userDelivery) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user = model.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user.Name == "" {
		http.Error(w, "please input a valid name", http.StatusBadRequest)
		return
	}
	if user.Email == "" {
		http.Error(w, "please input a valid email", http.StatusBadRequest)
		return
	}
	if user.Password == "" {
		http.Error(w, "please input a valid password", http.StatusBadRequest)
		return
	}
	u, err := ud.uc.Register(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		w.Write([]byte("internal server error"))
		return
	}
	json.NewEncoder(w).Encode(&u)
	w.Write([]byte("register success"))
}

func (ud *userDelivery) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var login = model.Login{}
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if login.Email == "" {
		http.Error(w, "please input a valid email", http.StatusBadRequest)
		return
	}
	if login.Password == "" {
		http.Error(w, "please input a valid password", http.StatusBadRequest)
		return
	}
	token, message, err := ud.uc.VerifyUserPassword(w, r, login.Email, login.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		w.Write([]byte("invalid username or password"))
		return
	}

	cookie := &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: constants.EXPIRATION_TIME,
	}

	http.SetCookie(w, cookie)
	fmt.Fprintf(w, "login success %v", message)

	// w.Write([]byte(token))
}
