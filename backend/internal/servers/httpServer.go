package servers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/xhermitx/gitpulse-tracker/backend/internal/handlers"
	"github.com/xhermitx/gitpulse-tracker/backend/internal/models"
	msql "github.com/xhermitx/gitpulse-tracker/backend/store/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Handler func(w http.ResponseWriter, r *http.Request) error

func wrapper(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		handleError(w, err)
	}
}

func handleError(w http.ResponseWriter, err error) {
	var apiErr *models.APIError

	if !errors.As(err, apiErr) {
		apiErr = handlers.ErrInternal
		log.Printf("\nInternal server error %v", err)
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(apiErr.StatusCode)

	json.NewEncoder(w).Encode(apiErr)
}

// HANDLE THE ROUTES
func handleRequests(handler *handlers.TaskHandler) {
	router := mux.NewRouter().StrictSlash(true)

	//ADD MIDDLEWARE TO HANDLE AUTHENTICATION
	router.Use(authMW)

	jobRouter := router.PathPrefix("/job").Subrouter()

	jobRouter.HandleFunc("/create", wrapper(handler.CreateJob)).Methods("POST")
	jobRouter.HandleFunc("/delete", wrapper(handler.DeleteJob)).Methods("POST")
	jobRouter.HandleFunc("/update", wrapper(handler.UpdateJob)).Methods("POST")
	jobRouter.HandleFunc("/list", wrapper(handler.ListJobs)).Methods("GET")
	jobRouter.HandleFunc("/trigger", wrapper(handler.Trigger)).Methods("POST")
	jobRouter.HandleFunc("/candidates", wrapper(handler.TopCandidates)).Methods("GET")

	// START A SERVER
	log.Fatal(http.ListenAndServe(os.Getenv("FRONTEND_ADDRESS"), jobRouter))
}

func HttpServer() {

	fmt.Println("DSN: ", os.Getenv("DB_SERVER"))

	db, err := gorm.Open(mysql.Open(os.Getenv("DB_SERVER")), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to Connect to DB")
	}
	log.Println("CONNECTED TO DB")
	mysqlDB := msql.NewMySQLStore(db)
	taskHandler := handlers.NewTaskHandler(mysqlDB)

	handleRequests(taskHandler)
}

func authMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := fmt.Sprintf("http://auth-service%s/auth/validate", os.Getenv("AUTH_ADDRESS"))

		req, err := http.NewRequest("POST", path, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Set("Authorization", r.Header.Get("Authorization"))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		if resp.StatusCode == http.StatusUnauthorized {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
