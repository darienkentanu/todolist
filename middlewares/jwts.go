package middlewares

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
	"todolist/constants"

	jwt "github.com/dgrijalva/jwt-go"
)

var EXPIRATION_TIME = time.Now().Add(24 * time.Hour)

type M map[string]interface{}

type MyClaims struct {
	jwt.StandardClaims
	ID uint `json:"id"`
}

func CreateToken(userID uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["id"] = userID
	claims["exp"] = constants.JWT_SIGNATURE_KEY

	tokenString, err := token.SignedString([]byte(constants.JWT_SIGNATURE_KEY))
	if err != nil {
		fmt.Printf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func ExtractTokenUserID(r *http.Request) (uint, error) {

	// get our token string from Cookie
	biscuit, err := r.Cookie("token")

	var tokenString string
	if err != nil {
		tokenString = ""
	} else {
		tokenString = biscuit.Value
	}

	// abort
	if tokenString == "" {
		return 0, nil
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(constants.JWT_SIGNATURE_KEY), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userIDInterface := claims["id"]
		userID, ok := userIDInterface.(float64)
		if ok {
			return uint(userID), nil
		}
	}
	return 0, errors.New("an error has been occured in extracting token")
}

func IsLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			log.Println("unauthorized")
			http.Error(w, "unauthorized", 404)
			return
		}
		tokenString := cookie.Value
		if tokenString == "" {
			log.Println("unauthorized")
			http.Error(w, "unauthorized", 404)
			return
		}
		// check apakah tanda tangan token sudah benar
		_, err = jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(constants.JWT_SIGNATURE_KEY), nil
		})
		if err != nil {
			log.Println("unauthorized")
			http.Error(w, "unauthorized", 404)
			return
		}
		next.ServeHTTP(w, r)
	})
}
