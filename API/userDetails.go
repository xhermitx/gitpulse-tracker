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
		return models.GitResponse{}, fmt.Errorf("error occured while fetching username: %v", gitResponse.Errors)
	} else if gitResponse.Data.User.Typename != "User" {
		return gitResponse, fmt.Errorf("username %s is not of type user", gitResponse.Data.User.Login)
	}

	return gitResponse, nil
}
