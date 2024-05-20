package utils

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/xhermitx/gitpulse-tracker/models"
)

func Printer(candidates []models.GitResponse) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"USERNAME", "FOLLOWERS", "CONTRIBUTIONS", "MOST POPULAR REPOSITORY -> STARS"})

	for _, candidate := range candidates {
		t.AppendRow([]interface{}{

			candidate.Data.User.Login, // USERNAME

			candidate.Data.User.Followers.TotalCount, // FOLLOWERS

			candidate.Data.User.ContributionsCollection.ContributionCalendar.TotalContributions, // CONTRIBUTIONS

			func() string { // MOST POPULAR REPOSITORIES
				if len(candidate.Data.User.Repositories.Nodes) > 0 {
					return fmt.Sprintf("%s : %d", candidate.Data.User.Repositories.Nodes[0].Name, candidate.Data.User.Repositories.Nodes[0].StargazerCount)
				}
				return ""
			}(),
		})
		t.AppendSeparator()
	}
	t.Render()
}
