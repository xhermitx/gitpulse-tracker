package API

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	models "github.com/xhermitx/gitpulse-tracker/models"
)

func GetUserDetails(username string) (models.GitResponse, error) {

	// LOAD ENV VARIABLES
	token := os.Getenv("GITHUB_TOKEN")

	query := `query GetUserDetails($username: String!) {
		user(login: $username) {
			login
			__typename
			contributionsCollection {
				contributionCalendar {
					totalContributions
				}
			}
			followers {
				totalCount
			}
			repositories(orderBy: {field: STARGAZERS, direction: DESC}, first: 1) {
				nodes {
					name
					stargazerCount
				}
			}
		}
	}`

	gqlQuery := models.GitQuery{
		Query: query,
		Variables: map[string]string{
			"username": username,
		},
	}

	body, err := json.Marshal(gqlQuery)
	if err != nil {
		return models.GitResponse{}, err
	}

	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(body))
	if err != nil {
		return models.GitResponse{}, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.GitResponse{}, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.GitResponse{}, err
	}

	var gitResponse models.GitResponse
	err = json.Unmarshal(responseBody, &gitResponse)
	if err != nil {
		return models.GitResponse{}, err
	}

	if len(gitResponse.Errors) > 0 {
		fmt.Println("Errors occurred:", gitResponse.Errors)
	} else if gitResponse.Data.User.Typename != "User" {
		return models.GitResponse{}, fmt.Errorf("invalid username")
	}

	return gitResponse, nil
	// } else {
	// 	fmt.Println("User:", gitResponse.Data.User.Login)
	// 	fmt.Println("Total Contributions:", gitResponse.Data.User.ContributionsCollection.ContributionCalendar.TotalContributions)
	// 	fmt.Println("Total Followers:", gitResponse.Data.User.Followers.TotalCount)
	// 	if len(gitResponse.Data.User.Repositories.Nodes) > 0 {
	// 		fmt.Println("Most Starred Repo:", gitResponse.Data.User.Repositories.Nodes[0].Name)
	// 		fmt.Println("Stars Count:", gitResponse.Data.User.Repositories.Nodes[0].StargazerCount)
	// 	}
	// }
	// return nil
}
