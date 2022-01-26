package jwts

import (
	"fmt"
	"net/http"
	"todolist/constants"

	jwt "github.com/dgrijalva/jwt-go"
)

type M map[string]interface{}

const JWT_SIGNATURE_KEY = "rahasia"

type MyClaims struct {
	jwt.StandardClaims
	ID uint `json:"id"`
}

func CreateToken(userID uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["id"] = userID
	claims["exp"] = constants.EXPIRATION_TIME

	tokenString, err := token.SignedString([]byte(JWT_SIGNATURE_KEY))
	if err != nil {
		fmt.Printf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func ExtractTokenUserID(r *http.Request) (string, error) {

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
		return "", nil
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWT_SIGNATURE_KEY), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		username := fmt.Sprintf("%v", claims["id"]) // convert to string
		// if err != nil {
		// 	return "", err
		// }
		return username, nil
	}
	return "", nil
}
