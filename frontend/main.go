package main

import (
	"github.com/xhermitx/gitpulse-tracker/frontend/internal/servers"
)

// MAKE REQUEST USING GRPC
// func grpcRequest(userIDs []string) {
// 	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		log.Fatalf("failed to connect to gRPC server at localhost:8080: %v", err)
// 	}
// 	defer conn.Close()
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
// 	client := gitfetch.NewGithubClient(conn)
// 	res, err := client.FetchData(ctx, &gitfetch.Profile{UserID: 1, JobID: 2, Usernames: userIDs})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Print("Successful Response: ", res)
// }

type Candidate struct {
	JobID     uint
	Usernames []string
}

// func httpRequest(candidate Candidate) {
// 	body, err := json.Marshal(candidate)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	postURL := "http://localhost:3000/"

// 	req, err := http.NewRequest("POST", postURL, bytes.NewBuffer(body))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	req.Header.Add("Content-Type", "application/json")

// 	client := &http.Client{}

// 	res, err := client.Do(req)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println(res.StatusCode)
// }

func main() {

	servers.HttpServer()

	// candidate := Candidate{JobID: 1, Usernames: []string{"xhermitx", "khalidfarooq", "jhasuraj020", "test1", "test2", "test3"}}

	// httpRequest(candidate)

	// PERFORMANCE CHECKS
	// t := time.Now()
	// defer func() {

	// 	f, err := os.OpenFile("Performance.md", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// 	if err != nil {
	// 		log.Print(err)
	// 	}
	// 	defer f.Close()
	// 	_, _ = fmt.Fprintln(f, "# PERFORMANCE WITH CONCURRENTLY READING DRIVE DATA, PDFs AND FETCHING FROM GITHUB")
	// 	_, _ = fmt.Fprintln(f, "TOTAL TIME TAKEN: ", time.Since(t).Seconds())
	// }()

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Panic("Error loading the environment variables")
	// }

	// var lock sync.Mutex
	// var wg sync.WaitGroup

	// userIDs, err := API.GetDriveDetails()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var detailedList []models.GitResponse

	// wg.Add(len(userIDs))
	// for _, user := range userIDs {
	// 	go func(user string) {

	// 		defer wg.Done()

	// 		if res, err := API.GetUserDetails(user); err != nil {
	// 			log.Println(err)
	// 		} else {
	// 			lock.Lock()
	// 			detailedList = append(detailedList, res)
	// 			lock.Unlock()
	// 		}
	// 	}(user)
	// }
	// wg.Wait()

	// utils.Printer(detailedList)
	// // GET USER DETAILS FROM GITHUB

}
