package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/xhermitx/gitpulse-tracker/auth-service/models"
	"github.com/xhermitx/gitpulse-tracker/auth-service/store"
)

var (
	ErrNotFound     = &models.APIError{StatusCode: http.StatusNotFound, Message: "Resource not found"}
	ErrBadRequest   = &models.APIError{StatusCode: http.StatusBadRequest, Message: "Bad request"}
	ErrUnauthorized = &models.APIError{StatusCode: http.StatusUnauthorized, Message: "Unauthorized"}
	ErrInternal     = &models.APIError{StatusCode: http.StatusInternalServerError, Message: "Internal Server Error"}
)

type TaskHandler struct {
	store store.Store
}

func NewTaskHandler(s store.Store) *TaskHandler {
	return &TaskHandler{store: s}
}

type Handler func(http.ResponseWriter, *http.Request) error

func Wrapper(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)

		handleError(w, err)
	}
}

func handleError(w http.ResponseWriter, err error) {

	var customErr *models.APIError

	if !errors.As(err, customErr) {
		customErr = models.NewAPIError(http.StatusInternalServerError, "Internal server error")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(customErr.StatusCode)
	json.NewEncoder(w).Encode(customErr)
}

func (t *TaskHandler) Register(w http.ResponseWriter, r *http.Request) error {

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return ErrBadRequest
	}
	defer r.Body.Close()

	var recruiter models.Recruiter
	if err := json.Unmarshal(data, &recruiter); err != nil {
		log.Println("Error Unmarshalling the data")
		return ErrBadRequest
	}
	defer r.Body.Close()

	if err = t.store.CreateRecruiter(&recruiter); err != nil {
		return ErrInternal
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated) // 201
	fmt.Fprintf(w, "Signup Successful")
	return nil
}

func (t *TaskHandler) Login(w http.ResponseWriter, r *http.Request) error {

	var credentials models.Credentials

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return ErrBadRequest
	}
	defer r.Body.Close()

	if err := json.Unmarshal(data, &credentials); err != nil {
		log.Println("Error Unmarshalling the data")
		return ErrBadRequest
	}

	token, err := t.store.AuthenticateRecruiter(&credentials)
	if err != nil {
		return ErrUnauthorized
	}

	cookie := http.Cookie{
		Name:     "Authorization",
		Value:    token,
		Secure:   false,
		HttpOnly: true,
		// SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(w, &cookie)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK) // 200
	fmt.Fprint(w, "Login Successful: ", token)
	return nil
}

func (t *TaskHandler) Validate(w http.ResponseWriter, r *http.Request) error {

	tokenString, err := getToken(r)
	if err != nil {
		return ErrBadRequest
	}

	log.Println("TOKEN: ", tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate the expected alg:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("incorrect signing method")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		log.Println("error while parsing the token: ", err)
		return ErrUnauthorized
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return ErrInternal
		}

		recruiter, err := t.store.FindRecruiter(uint(claims["id"].(float64)))
		if err != nil {
			log.Println("Invalid Token")
			return ErrUnauthorized
		}

		payload, err := json.Marshal(recruiter)
		if err != nil {
			return ErrInternal
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(payload)
		return nil

	}
	log.Println("Invalid Token")
	return ErrUnauthorized
}

func getToken(r *http.Request) (string, error) {
	// Extract the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Println("missing authorization header")
		return "", fmt.Errorf("authorization header missing")
	}

	// Split the header to get the token part
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		log.Println("incorrect authorization header")
		return "", fmt.Errorf("authorization header format must be Bearer {token}")
	}

	return headerParts[1], nil
}
