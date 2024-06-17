package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/xhermitx/gitpulse-tracker/frontend/internal/models"
)

func Auth(r *http.Request) (*models.Recruiter, error) {

	tokenString, err := GetToken(r)
	if err != nil {
		return nil, err
	}

	fmt.Println("AUTH ADDRESS: ", os.Getenv("AUTH_ADDRESS"))

	requestURL := fmt.Sprintf("http://auth-service%s/auth/validate", os.Getenv("AUTH_ADDRESS"))
	req, err := http.NewRequest(http.MethodPost, requestURL, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokenString))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return nil, err
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		return nil, err
	}
	defer res.Body.Close()

	var recruiter models.Recruiter

	if err = json.Unmarshal(resBody, &recruiter); err != nil {
		log.Println("error unmarshaling the response")
		return nil, err
	}
	return &recruiter, nil
}

func GetToken(r *http.Request) (string, error) {
	// EXTRACT THE AUTHORIZATION HEADER
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Println("empty auth header")
		return "", errors.New("authorization header required")
	}

	// SPLIT THE HEADER TO GET THE TOKEN
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("authorization header format must be Bearer {token}")
	}

	// RETURN TOKEN
	return headerParts[1], nil
}
