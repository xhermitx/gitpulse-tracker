package servers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/xhermitx/gitpulse-tracker/backend/internal/handlers"
	msql "github.com/xhermitx/gitpulse-tracker/backend/store/mysql"
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

	//ADD MIDDLEWARE TO HANDLE AUTHENTICATION
	router.Use(authMW)

	router.HandleFunc("/job", homePage)
	router.HandleFunc("/job/create", handler.CreateJob).Methods("POST")
	router.HandleFunc("/job/delete", handler.DeleteJob).Methods("POST")
	router.HandleFunc("/job/update", handler.UpdateJob).Methods("POST")
	router.HandleFunc("/job/list", handler.ListJobs).Methods("GET")
	router.HandleFunc("/job/trigger", handler.Trigger).Methods("POST")
	router.HandleFunc("/job/candidates", handler.TopCandidates).Methods("GET") // ?jobId=x

	// START A SERVER
	log.Fatal(http.ListenAndServe(os.Getenv("FRONTEND_ADDRESS"), router))
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
