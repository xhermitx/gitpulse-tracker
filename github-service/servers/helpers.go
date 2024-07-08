package servers

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
	"github.com/xhermitx/gitpulse-tracker/github-service/api"
	"github.com/xhermitx/gitpulse-tracker/github-service/models"
)

type Error struct {
	err        error
	msg        string
	httpStatus int
	w          http.ResponseWriter
}

func NewError(err error, msg string, httpStatus int, w http.ResponseWriter) Error {
	return Error{
		err:        err,
		msg:        msg,
		httpStatus: httpStatus,
		w:          w,
	}
}

func (e Error) HandleError() {
	log.Println(e.err)
	http.Error(e.w, e.msg, e.httpStatus)
}

func HttpServer() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/github", FetchData).Methods("POST")
	log.Fatal(http.ListenAndServe(os.Getenv("GITHUB_ADDRESS"), router))
}

func ProcessUsername(res models.Job, client *api.Client, wg *sync.WaitGroup, w http.ResponseWriter) {
	var candidate models.Candidate
	for _, u := range res.Usernames {
		profile, err := api.GetUserDetails(u) // GET EACH CANDIDATE'S DATA FROM GITHUB
		if err != nil {
			log.Println("failed to fetch data for username: ", u)
		} else {
			wg.Add(1)
			candidate = *CreateCandidate(res.JobID, profile)
			go func(candidate models.Candidate) {
				defer wg.Done()
				if err = client.Publish(GITHUB_QUEUE, candidate); err != nil {
					log.Println("error publishing the data for candidate: ", candidate.GithubId)
				}
			}(candidate)
		}
	}

	wg.Wait()
	if err := client.Publish(GITHUB_QUEUE, models.Candidate{JobId: candidate.JobId, Status: true}); err != nil {
		customErr := NewError(err, INTERNAL_ERROR, http.StatusInternalServerError, w)
		customErr.HandleError()
		return
	}
}

func CreateCandidate(jobId uint, profile models.GitResponse) *models.Candidate {
	return &models.Candidate{
		JobId:         uint(jobId),
		GithubId:      profile.Data.User.Login,
		Followers:     uint(profile.Data.User.Followers.TotalCount),
		Contributions: uint(profile.Data.User.ContributionsCollection.ContributionCalendar.TotalContributions),
		MostPopularRepo: func() string {
			if len(profile.Data.User.Repositories.Nodes) > 0 {
				return profile.Data.User.Repositories.Nodes[0].Name
			}
			return ""
		}(),
		RepoStars: func() uint {
			if len(profile.Data.User.Repositories.Nodes) > 0 {
				return uint(profile.Data.User.Repositories.Nodes[0].StargazerCount)
			}
			return 0
		}(),
	}
}
