package api

import (
	"context"
	"log"

	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

type UserContributions struct {
	User struct {
		ContributionsCollection struct {
			ContributionCalendar struct {
				TotalContributions graphql.Int
			}
		}
	} `graphql:"user(login: $login)"`
}

func GetContributions(username string, token string) (int, error) {

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := graphql.NewClient("https://api.github.com/graphql", httpClient)

	var query UserContributions
	variables := map[string]interface{}{
		"login": graphql.String(username), // replace with the username you're querying
	}

	if err := client.Query(context.Background(), &query, variables); err != nil {
		log.Fatalf("Error querying GitHub: %v", err)
	}

	return int(query.User.ContributionsCollection.ContributionCalendar.TotalContributions), nil
}
