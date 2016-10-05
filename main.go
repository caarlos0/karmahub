package main

import (
	"fmt"
	"os"

	"github.com/caarlos0/karmahub/internal/karma"
	"github.com/caarlos0/karmahub/internal/search"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	fn := search.Github(client)

	crs, _ := karma.Reviews(fn, "caarlos0", "user:ContaAzul")
	fmt.Println("Reviews:", crs)
	prs, _ := karma.Authors(fn, "caarlos0", "user:ContaAzul")
	fmt.Println("Authors:", prs)
}
