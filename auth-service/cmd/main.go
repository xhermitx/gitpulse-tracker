package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/xhermitx/gitpulse-tracker/auth-service/handlers"
	msql "github.com/xhermitx/gitpulse-tracker/auth-service/store/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// HANDLE THE ROUTES
func handleRequests(handler *handlers.TaskHandler) {
	router := mux.NewRouter().StrictSlash(true)

	authRouter := router.PathPrefix("/auth").Subrouter()

	authRouter.HandleFunc("/register", handlers.Wrapper(handler.Register)).Methods("POST")
	authRouter.HandleFunc("/login", handlers.Wrapper(handler.Login)).Methods("POST")
	authRouter.HandleFunc("/validate", handlers.Wrapper(handler.Validate)).Methods("POST")

	log.Fatal(http.ListenAndServe(os.Getenv("AUTH_ADDRESS"), authRouter))
}

func main() {

	fmt.Println("DB SERVER ADDRESS: ", os.Getenv("DB_SERVER"))

	db, err := gorm.Open(mysql.Open(os.Getenv("DB_SERVER")), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to Connect to DB")
	}
	mysqlDB := msql.NewMySQLStore(db)
	taskHandler := handlers.NewTaskHandler(mysqlDB)

	handleRequests(taskHandler)
}
