package main

import (
	"context"
	"fmt"
	"log"
	"time"

	github "github.com/xhermitx/gitpulse-tracker/github-service/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server at localhost:8080: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	client := github.NewGithubClient(conn)

	res, err := client.FetchData(ctx, &github.Profile{UserID: 1, JobID: 2, Usernames: []string{"test"}})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Successful Response: ", res)

}
