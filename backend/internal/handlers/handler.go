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

func (h TaskHandler) CreateJob(w http.ResponseWriter, r *http.Request) error {

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return ErrBadRequest
	}
	defer r.Body.Close()

	var Job models.Job
	if err := json.Unmarshal(data, &Job); err != nil {
		log.Println("Error Unmarshalling the data")
		return ErrBadRequest
	}
	defer r.Body.Close()

	_, err = h.store.CreateJob(&Job)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated) // 201
	fmt.Fprintf(w, "Job ID: %d -> Created Successfully", Job.JobId)

	return nil
}

func (h TaskHandler) DeleteJob(w http.ResponseWriter, r *http.Request) error {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return ErrBadRequest
	}
	defer r.Body.Close()

	jobId, err := strconv.Atoi(string(data))
	if err != nil {
		return ErrBadRequest
	}

	if err = h.store.DeleteJob(uint(jobId)); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Job Deleted Successfully")

	return nil
}

func (h TaskHandler) ListJobs(w http.ResponseWriter, r *http.Request) error {
	// VALIDATE THE TOKEN TO GET A RECRUITER INFORMATION
	recruiter, err := utils.Auth(r)
	if err != nil {
		return ErrUnauthorized
	}

	jobs, err := h.store.ListJobs(recruiter.RecruiterId)
	if err != nil {
		return err
	}

	res_data, err := json.Marshal(jobs)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res_data)

	return nil
}

func (h TaskHandler) Trigger(w http.ResponseWriter, r *http.Request) error {

	jobId, err := getID(r)
	if err != nil {
		return ErrBadRequest
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
		return err
	}

	type candidates struct {
		JobID     uint
		Usernames []string
	}

	folderId, err := extractFolderID(job.DriveLink)
	if err != nil {
		return ErrBadRequest
	}

	usernames, err := gdrive.GetDriveDetails(folderId)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(candidates{job.JobId, usernames})
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf("http://github-service%s/github", os.Getenv("GITHUB_ADDRESS"))
	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(res.StatusCode)
	fmt.Fprint(w, "PROFILING TRIGGERED")

	return nil
}

func (h TaskHandler) TopCandidates(w http.ResponseWriter, r *http.Request) error {
	values := r.URL.Query()
	jobId, err := strconv.Atoi(values.Get("jobId"))
	if err != nil {
		log.Println("Invalid Job ID")
		return ErrBadRequest
	}

	topCandidates, err := h.store.ListCandidates(uint(jobId))
	if err != nil {
		log.Println("Error fetching the candidates from DB")
		return ErrNotFound
	}

	payload, err := json.Marshal(topCandidates)
	if err != nil {
		log.Println("Error Marshalling the candidates data")
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)

	return nil
}

func (h TaskHandler) UpdateJob(w http.ResponseWriter, r *http.Request) error {
	// TO BE IMPLEMENTED
	return models.NewAPIError(http.StatusNotImplemented, "Endpoint not available")
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

func getID(r *http.Request) (int, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return 0, err
	}
	defer r.Body.Close()

	return strconv.Atoi(string(body))
}
