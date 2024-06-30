package servers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/gorilla/mux"
	"github.com/xhermitx/gitpulse-tracker/github-service/api"
	"github.com/xhermitx/gitpulse-tracker/github-service/models"
)

func HttpServer() {
	router := mux.NewRouter().StrictSlash(true)

	router.Use(authMW)

	router.HandleFunc("/github", FetchData).Methods("POST")

	log.Fatal(http.ListenAndServe(os.Getenv("GITHUB_ADDRESS"), router))
}

func FetchData(w http.ResponseWriter, r *http.Request) {

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error processing the request", http.StatusBadRequest)
	}
	defer r.Body.Close()

	var res models.Job

	if err = json.Unmarshal(reqBody, &res); err != nil {
		http.Error(w, "Error processing the request", http.StatusBadRequest)
	}

	if len(res.Usernames) == 0 {
		log.Println("No usernames found")
		http.Error(w, "Error processing the request", http.StatusBadRequest)
	}

	wg := sync.WaitGroup{}

	var candidate models.Candidate
	// GET EACH CANDIDATE'S DATA FROM GITHUB
	for _, u := range res.Usernames {
		profile, err := api.GetUserDetails(u)
		if err != nil {
			log.Println(err)
		} else {
			wg.Add(1)
			candidate = models.Candidate{
				JobId:         res.JobID,
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

			// CREATE A GO ROUTINE FOR EACH PUBLISH ON THE QUEUE
			go func(candidate models.Candidate) {
				defer wg.Done()
				if err = api.Publish(candidate); err != nil {
					fmt.Print(err)
				}
			}(candidate)

		}
	}

	wg.Wait()
	if err = api.Publish(models.Candidate{JobId: candidate.JobId, Status: true}); err != nil {
		log.Print(err)
	}

	// SEND A RESPONSE STATING SUCCESSFUL TRIGGER
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "PROFILING TRIGGER")
}

func authMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r.URL = &url.URL{
			Path: fmt.Sprintf("http://auth-service%s/auth/validate", os.Getenv("AUTH_ADDRESS")),
		}

		client := &http.Client{}
		resp, err := client.Do(r)
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
