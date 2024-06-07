package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/xhermitx/gitpulse-tracker/auth-service/internal/app/handlers"
	msql "github.com/xhermitx/gitpulse-tracker/auth-service/internal/store/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DEFINE THE HOME PAGE
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the home Page")
	fmt.Println("Endpoint hit: homepage")
}

// HANDLE THE ROUTES
func handleRequests(handler *handlers.TaskHandler) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homePage)
	router.HandleFunc("/auth/register", handler.Register).Methods("POST")
	router.HandleFunc("/auth/login", handler.Login).Methods("POST")
	router.HandleFunc("/auth/validate", handler.Validate).Methods("GET")

	log.Fatal(http.ListenAndServe(os.Getenv("ADDRESS"), router))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panic(err)
	}

	db, err := gorm.Open(mysql.Open(os.Getenv("DB_SERVER")), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}
	mysqlDB := msql.NewMySQLStore(db)
	taskHandler := handlers.NewTaskHandler(mysqlDB)

	handleRequests(taskHandler)
}
