package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/xhermitx/gitpulse-tracker/gitfetch"
	"google.golang.org/grpc"
)

type myServer struct {
	gitfetch.UnimplementedGithubServer
}

func (s *myServer) FetchData(ctx context.Context, in *gitfetch.Profile) (*gitfetch.Response, error) {
	if len(in.Usernames) == 0 {
		return nil, fmt.Errorf("error processing the requests")
	}

	for _, u := range in.Usernames {
		fmt.Println(u)
	}
	return &gitfetch.Response{User: []string{"Hello", "World"}, Status: true}, nil
}

// SERVICE TO FETCH AND HANDLE GITHUB DATA
func main() {

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("cannot create listener: %v", err)
	}

	server := grpc.NewServer()

	gitfetch.RegisterGithubServer(server, &myServer{})
	log.Printf("gRPC server is listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}

}
