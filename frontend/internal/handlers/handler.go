package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/xhermitx/gitpulse-tracker/frontend/internal/models"
	"github.com/xhermitx/gitpulse-tracker/frontend/internal/store"
)

type TaskHandler struct {
	store store.Store
}

func NewTaskHandler(s store.Store) *TaskHandler {
	return &TaskHandler{store: s}
}

func (h TaskHandler) CreateJob(w http.ResponseWriter, r *http.Request) {

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading the request body", http.StatusBadRequest) // 400
		return
	}
	defer r.Body.Close()

	var Job models.Job
	if err := json.Unmarshal(data, &Job); err != nil {
		http.Error(w, "Error reading the request body", http.StatusBadRequest) // 400
		log.Println("Error Unmarshalling the data")
		return
	}
	defer r.Body.Close()

	job, err := h.store.CreateJob(&Job)
	if err != nil {
		http.Error(w, "Failed to create the Job", http.StatusInternalServerError) // 400
		return
	}

	log.Println(job)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusCreated) // 201

	fmt.Fprintf(w, "Job Created Successfully")
}

func (h TaskHandler) UpdateJob(w http.ResponseWriter, r *http.Request) {
	// TO BE IMPLEMENTED
	http.Error(w, "Failed to Update the Job", http.StatusNotImplemented) //501
}

func (h TaskHandler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read the body", http.StatusBadRequest) // 400
		return
	}

	var jobId uint

	if err = json.Unmarshal(data, &jobId); err != nil {
		http.Error(w, "failed to read the body", http.StatusBadRequest) // 400
		log.Println("error unmarshalling the data")
		return
	}

	if err = h.store.DeleteJob(jobId); err != nil {
		http.Error(w, "Failed to Delete the job Id", http.StatusInternalServerError) // 500
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "Job Deleted Successfully")
}

func (h TaskHandler) ListJobs(w http.ResponseWriter, r *http.Request) {
	recruiter_id, err := strconv.Atoi(r.URL.Query().Get("recruiter_id"))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	jobs, err := h.store.ListJobs(uint(recruiter_id))
	if err != nil {
		http.Error(w, "Failed to Delete the job Id", http.StatusInternalServerError) // 500
		return
	}

	res_data, err := json.Marshal(jobs)
	if err != nil {
		http.Error(w, "failed to get the job list", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res_data)
}
