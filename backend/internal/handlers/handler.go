package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	gdrive "github.com/xhermitx/gitpulse-tracker/backend/gdrive"
	"github.com/xhermitx/gitpulse-tracker/backend/internal/models"
	"github.com/xhermitx/gitpulse-tracker/backend/queue/rmq"
	"github.com/xhermitx/gitpulse-tracker/backend/store"
	"github.com/xhermitx/gitpulse-tracker/backend/utils"
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

	_, err = h.store.CreateJob(&Job)
	if err != nil {
		http.Error(w, "Failed to create the Job", http.StatusInternalServerError) // 400
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusCreated) // 201

	fmt.Fprintf(w, "Job ID: %d -> Created Successfully", Job.JobId)
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
	defer r.Body.Close()

	jobId, err := strconv.Atoi(string(data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err = h.store.DeleteJob(uint(jobId)); err != nil {
		http.Error(w, "Failed to Delete the job Id", http.StatusInternalServerError) // 500
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "Job Deleted Successfully")
}

func (h TaskHandler) ListJobs(w http.ResponseWriter, r *http.Request) {
	// VALIDATE THE TOKEN TO GET A RECRUITER INFORMATION
	recruiter, err := utils.Auth(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	jobs, err := h.store.ListJobs(recruiter.RecruiterId)
	if err != nil {
		http.Error(w, "failed to delete the job Id", http.StatusInternalServerError) // 500
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

func (h TaskHandler) Trigger(w http.ResponseWriter, r *http.Request) {

	// GET THE JOB ID FROM BODY
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read the body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	jobId, err := strconv.Atoi(string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// NOTIFY THE QUEUE THAT THE PROCESS IS INITITIATED
	data := models.StatusQueue{
		JobId:  uint(jobId),
		Status: true,
		Timer:  time.Now(),
	}

	mq := rmq.NewRabbitMQClient(data, models.STATUS_QUEUE)

	if err := mq.Publish(); err != nil {
		log.Printf("\nfailed to update the profiling status for Job %d: %v", data.JobId, err)
	}

	job, err := h.store.GetJob(uint(jobId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type candidates struct {
		JobID     uint
		Usernames []string
	}

	folderId, err := extractFolderID(job.DriveLink)
	if err != nil {
		http.Error(w, "invalid drive link", http.StatusBadRequest)
		return
	}

	log.Println("FolderID: ", folderId)

	usernames, err := gdrive.GetDriveDetails(folderId)
	if err != nil {
		http.Error(w, "error fetching data from Drive", http.StatusInternalServerError)
		return
	}

	payload, err := json.Marshal(candidates{job.JobId, usernames})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	requestURL := fmt.Sprintf("http://github-service%s/github", os.Getenv("GITHUB_ADDRESS"))
	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(payload))
	if err != nil {
		http.Error(w, "client: could not create request", http.StatusInternalServerError)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "client: error making request", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(res.StatusCode)
	fmt.Fprint(w, "PROFILING TRIGGERED")
}

func (h TaskHandler) TopCandidates(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	jobId, err := strconv.Atoi(values.Get("jobId"))
	if err != nil {
		log.Println("Invalid Job ID")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	topCandidates, err := h.store.ListCandidates(uint(jobId))
	if err != nil {
		log.Println("Error fetching the candidates from DB")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	payload, err := json.Marshal(topCandidates)
	if err != nil {
		log.Println("Error Marshalling the candidates data")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

func extractFolderID(link string) (string, error) {

	pattern := `https://drive\.google\.com/drive/folders/([0-9A-Za-z-_]+)`

	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(link)

	if len(matches) > 1 {
		// THE FIRST MATCH IS THE ENTIRE MATCH, AND THE SECOND IS THE CAPTURED GROUP
		return matches[1], nil
	}

	return "", fmt.Errorf("folder ID not found in link")
}
