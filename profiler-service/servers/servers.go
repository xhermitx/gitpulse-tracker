package servers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func FetchData(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write([]byte("Profiler-service is running...")); err != nil {
		log.Println(err)
	}
}

func HttpServer() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", FetchData).Methods("GET")

	log.Fatal(http.ListenAndServe(":8081", router))
}
