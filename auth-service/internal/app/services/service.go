package services

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// GENERATE TOKENS
func JwtAuth(id int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		log.Fatal(err)
	}

	return tokenString
}
