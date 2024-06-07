package handlers

import (
	"fmt"
	"net/http"

	"github.com/xhermitx/gitpulse-tracker/auth-service/internal/store"
)

type TaskHandler struct {
	store store.Store
}

func NewTaskHandler(s store.Store) *TaskHandler {
	return &TaskHandler{
		store: s,
	}
}

func (t *TaskHandler) Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Register Recruiter")
}

func (t *TaskHandler) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Login Recruiter")
}

func (t *TaskHandler) Validate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Validate User")
}
