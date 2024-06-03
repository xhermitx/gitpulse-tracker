package servers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/xhermitx/gitpulse-tracker/github-service/API"
	"github.com/xhermitx/gitpulse-tracker/github-service/models"
)

// type myServer struct {
// 	gitfetch.UnimplementedGithubServer
// }

// func (s *myServer) FetchData(ctx context.Context, in *gitfetch.Profile) (*gitfetch.Response, error) {

// 	if len(in.Usernames) == 0 {
// 		return nil, fmt.Errorf("error processing the requests")
// 	}

// 	for _, userID := range in.Usernames {
// 		user, err := API.GetUserDetails(userID)
// 		if err != nil {
// 			log.Printf("Error fetching the user %s : %v", user.Data.User.Login, err)
// 		}
// 	}

// 	return &gitfetch.Response{Candidate: []*gitfetch.User{}, Status: true}, nil
// }

// func GrpcServer() {
// 	lis, err := net.Listen("tcp", ":8080")
// 	if err != nil {
// 		log.Fatalf("cannot create listener: %v", err)
// 	}

// 	server := grpc.NewServer()

// 	gitfetch.RegisterGithubServer(server, &myServer{})
// 	log.Printf("gRPC server is listening at %v", lis.Addr())
// 	if err := server.Serve(lis); err != nil {
// 		log.Fatalf("failed to server: %v", err)
// 	}
// }

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
		http.Error(w, "Error processing the request", http.StatusBadRequest)
	}

	wg := sync.WaitGroup{}
	wg.Add(len(res.Usernames))
	// GET EACH CANDIDATE'S DATA FROM GITHUB
	for i, u := range res.Usernames {
		// profile, err := API.GetUserDetails(u)
		// if err != nil {
		// 	log.Println(err)
		// } else {
		// 	candidate := models.Candidate{
		// 		JobID:         res.JobID,
		// 		Username:      profile.Data.User.Login,
		// 		Followers:     profile.Data.User.Followers.TotalCount,
		// 		Contributions: profile.Data.User.ContributionsCollection.ContributionCalendar.TotalContributions,
		// 		MostPopularRepo: func() string {
		// 			if len(profile.Data.User.Repositories.Nodes) > 0 {
		// 				return profile.Data.User.Repositories.Nodes[0].Name
		// 			}
		// 			return ""
		// 		}(),
		// 		RepoStars: func() int {
		// 			if len(profile.Data.User.Repositories.Nodes) > 0 {
		// 				return profile.Data.User.Repositories.Nodes[0].StargazerCount
		// 			}
		// 			return 0
		// 		}(),
		// 	}

		candidate := models.Candidate{
			JobID:           2,
			Username:        u,
			Followers:       4 + i, // Added a variable for different scores on redis
			Contributions:   20,
			MostPopularRepo: "Test",
			RepoStars:       200,
			Status:          false,
		}

		// CREATE A GO ROUTINE FOR EACH PUBLISH ON THE QUEUE
		go func(candidate models.Candidate) {
			defer wg.Done()
			if err = API.Publish(candidate); err != nil {
				fmt.Print(err)
			}
		}(candidate)
	}

	wg.Wait()
	if err = API.Publish(models.Candidate{JobID: 2, Status: true}); err != nil {
		log.Print(err)
	}
}

func HttpServer() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", FetchData).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
