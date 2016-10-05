package main

import (
	"fmt"
	"os"

	"github.com/caarlos0/karmahub/internal/karma"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	crs, _ := karma.Reviews(client, "caarlos0", "user:ContaAzul")
	fmt.Println("Reviews:", crs)
	prs, _ := karma.Pulls(client, "caarlos0", "user:ContaAzul")
	fmt.Println("Pulls:", prs)
}
