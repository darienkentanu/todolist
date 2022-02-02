package delivery

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"todolist/middlewares"
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
	getR := r.Methods("GET").Subrouter()
	getR.Use(middlewares.IsLoggedIn)
	getR.HandleFunc("/logout", handler.Logout)

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
		log.Println(err)
		http.Error(w, "internal server error", http.StatusBadRequest)
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
		log.Println(err)
		http.Error(w, "an error has been occured", 500)
		return
	}
	if login.Email == "" {
		http.Error(w, "please input a valid email", 400)
		return
	}
	if login.Password == "" {
		http.Error(w, "please input a valid password", 400)
		return
	}

	userID, err := ud.uc.GetUserIDByEmail(login.Email)
	if err != nil {
		log.Println(err)
		http.Error(w, "an error has been occured", 400)
		return
	}

	err = ud.uc.VerifyUserPassword(w, r, login.Email, login.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, "invalid username or password", 400)
		return
	}

	ud.uc.RemoveCookie(w)

	token, err := ud.uc.CreateToken(userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "an error has been occured", 500)
		return
	}
	cookie := ud.uc.CreateCookie(token)

	http.SetCookie(w, cookie)
	fmt.Fprintf(w, "login success")
}

func (ud *userDelivery) Logout(w http.ResponseWriter, r *http.Request) {
	message := ud.uc.RemoveCookie(w)
	log.Printf("logout success %v", message)
	var response = struct {
		Status  int
		Message string
	}{
		200,
		"logout success",
	}
	json.NewEncoder(w).Encode(&response)
}
