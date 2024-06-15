package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/xhermitx/gitpulse-tracker/auth-service/internal/models"
	"github.com/xhermitx/gitpulse-tracker/auth-service/internal/store"
)

type TaskHandler struct {
	store store.Store
}

func NewTaskHandler(s store.Store) *TaskHandler {
	return &TaskHandler{store: s}
}

func (t *TaskHandler) Register(w http.ResponseWriter, r *http.Request) {

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading the request body", http.StatusBadRequest) // 400
		return
	}
	defer r.Body.Close()

	var recruiter models.Recruiter
	if err := json.Unmarshal(data, &recruiter); err != nil {
		http.Error(w, "Error reading the request body", http.StatusBadRequest) // 400
		log.Println("Error Unmarshalling the data")
		return
	}
	defer r.Body.Close()

	if err = t.store.CreateRecruiter(&recruiter); err != nil {
		http.Error(w, "Failed to create the recruiter", http.StatusInternalServerError) // 400
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusCreated) // 201

	fmt.Fprintf(w, "Signup Successful")
}

func (t *TaskHandler) Login(w http.ResponseWriter, r *http.Request) {

	var credentials *models.Credentials

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading the request body", http.StatusBadRequest) //400
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(data, &credentials); err != nil {
		http.Error(w, "Error reading the request body", http.StatusBadRequest) //400
		log.Println("Error Unmarshalling the data")
		return
	}

	token, err := t.store.AuthenticateRecruiter(credentials)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized) //401
		return
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
}

func (t *TaskHandler) Validate(w http.ResponseWriter, r *http.Request) {

	// Extract the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "authorization header required", http.StatusUnauthorized)
		return
	}

	// Split the header to get the token part
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		http.Error(w, "authorization header format must be Bearer {token}", http.StatusUnauthorized)
		return
	}

	tokenString := headerParts[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate the expected alg:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		log.Print(err)
	}

	log.Println("Debug1")

	if claims, ok := token.Claims.(jwt.MapClaims); ok {

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			http.Error(w, "token expired", http.StatusUnauthorized) //401
			return
		}

		recruiter, err := t.store.FindRecruiter(uint(claims["id"].(float64)))
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		payload, err := json.Marshal(recruiter)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(payload))

	} else {
		log.Print("Invalid Token")
	}
}
