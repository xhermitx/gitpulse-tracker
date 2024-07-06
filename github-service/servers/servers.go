package servers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/xhermitx/gitpulse-tracker/github-service/api"
	"github.com/xhermitx/gitpulse-tracker/github-service/models"
)

const (
	INTERNAL_ERROR = "internal server error"
	BAD_REQUEST    = "error processing the request"
	GITHUB_QUEUE   = "github_data_queue"
)

func FetchData(w http.ResponseWriter, r *http.Request) {

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		customErr := NewError(err, BAD_REQUEST, http.StatusBadRequest, w)
		customErr.HandleError()
		return
	}
	defer r.Body.Close()

	var res models.Job
	err = json.Unmarshal(reqBody, &res)
	if err != nil || len(res.Usernames) == 0 {
		customErr := NewError(err, BAD_REQUEST, http.StatusBadRequest, w)
		customErr.HandleError()
		return
	}

	wg := &sync.WaitGroup{}
	conn, err := api.Connect()
	if err != nil {
		customErr := NewError(err, INTERNAL_ERROR, http.StatusInternalServerError, w)
		customErr.HandleError()
		return
	}
	defer conn.Close()
	client := api.NewQueueConnection(conn)

	ProcessUsername(res, client, wg, w)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "STARTED PROFILING...")
}
