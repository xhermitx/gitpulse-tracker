package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/xhermitx/gitpulse-tracker/API"
	"github.com/xhermitx/gitpulse-tracker/gitfetch"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	// PERFORMANCE CHECKS
	t := time.Now()
	defer func() {

		f, err := os.OpenFile("Performance.md", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Print(err)
		}
		defer f.Close()
		_, _ = fmt.Fprintln(f, "# PERFORMANCE WITH CONCURRENTLY READING DRIVE DATA, PDFs AND FETCHING FROM GITHUB")
		_, _ = fmt.Fprintln(f, "TOTAL TIME TAKEN: ", time.Since(t).Seconds())
	}()

	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading the environment variables")
	}

	// var lock sync.Mutex
	// var wg sync.WaitGroup

	userIDs, err := API.GetDriveDetails()
	if err != nil {
		log.Fatal(err)
	}

	func(userIDs []string) {
		conn, err := grpc.NewClient("localhost:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("failed to connect to gRPC server at localhost:8080: %v", err)
		}
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		client := gitfetch.NewGithubClient(conn)

		res, err := client.FetchData(ctx, &gitfetch.Profile{UserID: 1, JobID: 2, Usernames: userIDs})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print("Successful Response: ", res)
	}(userIDs)

	// var detailedList []models.GitResponse

	// // GET USER DETAILS FROM GITHUB
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
}
