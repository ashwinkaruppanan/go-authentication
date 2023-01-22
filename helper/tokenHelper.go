package helper

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Details struct {
	Name  string
	Email string
	jwt.StandardClaims
}

var secretKey = os.Getenv("SECRET_KEY")

func GenerateToken(name string, email string) (string, error) {
	Claims := &Details{
		Name:  name,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims).SignedString([]byte(secretKey))
	if err != nil {
		log.Panic(err)
	}

	return token, err
}

func ValidateToken(singnedToken string) (claims *Details, msg string) {
	token, err := jwt.ParseWithClaims(singnedToken, &Details{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*Details)
	if !ok {
		msg = fmt.Sprint("the token is invalid")
		return
	}
	if claims.ExpiresAt < time.Now().Unix() {
		msg = fmt.Sprint("token is expired")
		return
	}

	return claims, msg

}
