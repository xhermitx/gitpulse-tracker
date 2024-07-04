package api

import (
	"log"
	"testing"

	"github.com/xhermitx/gitpulse-tracker/github-service/models"
)

func TestGetUserDetails(t *testing.T) {
	input := "xhermitx"

	expected := uint(11)

	profile, err := GetUserDetails(input)
	if err != nil {
		t.Fatal("failed to get details")
		return
	}

	candidate := models.Candidate{
		JobId:         0,
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

	if candidate.Followers != expected {
		t.Fatal("Failed the test")
	}

	log.Println(candidate)
}
