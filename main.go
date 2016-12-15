package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/caarlos0/karmahub/internal/karma"
	"github.com/caarlos0/karmahub/internal/search"
	"github.com/caarlos0/spin"
	"github.com/google/go-github/github"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
)

var version = "master"

func main() {
	app := cli.NewApp()
	app.Name = "karmahub"
	app.Version = version
	app.Author = "Carlos Alexandro Becker (caarlos0@gmail.com)"
	app.Usage = "get your history of reviews/comments and pull requests/issues opened"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			EnvVar: "GITHUB_TOKEN",
			Name:   "token",
			Usage:  "Your GitHub token",
		},
		cli.StringFlag{
			Name:  "user, org, u, o",
			Usage: "User/Organization to get data from.",
		},
		cli.StringFlag{
			Name:  "filter",
			Usage: "Additional filters, github syntax. E.g.: user:Docker will filter for issues and pulls in Docker organization",
		},
	}
	app.Action = func(c *cli.Context) error {
		token := c.String("token")
		if token == "" {
			return cli.NewExitError("Missing GitHub Token", 1)
		}
		user := c.String("user")
		if user == "" {
			return cli.NewExitError("Missing GitHub User/Org", 1)
		}
		filter := c.String("filter")
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(oauth2.NoContext, ts)
		client := github.NewClient(tc)
		fn := search.Github(client)

		spin := spin.New("%s Gathering data...")
		spin.Start()
		prs, err := karma.Authors(fn, user, filter)
		if err != nil {
			spin.Stop()
			return cli.NewExitError(err.Error(), 1)
		}
		crs, err := karma.Reviews(fn, user, filter)
		if err != nil {
			spin.Stop()
			return cli.NewExitError(err.Error(), 1)
		}
		spin.Stop()

		fmt.Println("\033[1mAction    \t1m\t2m\t3m\033[0m")
		fmt.Println("Authored\t" + intsToStr(prs))
		fmt.Println("Reviewed\t" + intsToStr(crs))
		return nil
	}
	app.Run(os.Args)
}

func intsToStr(arr []int) string {
	var strs []string
	for _, n := range arr {
		strs = append(strs, strconv.Itoa(n))
	}
	return strings.Join(strs, "\t")
}
