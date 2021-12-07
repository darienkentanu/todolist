package jwts

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type M map[string]interface{}

const APPLICATION_NAME = "Todolist"
const LOGIN_EXPIRATION_DURATION = time.Duration(time.Hour * 2)
const JWT_SIGNATURE_KEY = "rahasia"

var JWT_SIGNING_METHOD = jwt.SigningMethodHS256

type MyClaims struct {
	jwt.StandardClaims
	ID uint `json:"id"`
}

func CreateToken(userID uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix() // Token expires after 2 hour

	tokenString, err := token.SignedString([]byte(JWT_SIGNATURE_KEY))
	if err != nil {
		fmt.Printf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}
