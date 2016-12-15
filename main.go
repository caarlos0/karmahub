package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/caarlos0/karmahub/internal/karma"
	"github.com/caarlos0/karmahub/internal/search"
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
		user := c.String("user")
		if token == "" || user == "" {
			return cli.ShowAppHelp(c)
		}
		filter := c.String("filter")
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(oauth2.NoContext, ts)
		client := github.NewClient(tc)
		fn := search.Github(client)

		fmt.Println("Action    \t1m\t2m\t3m")
		prs, err := karma.Authors(fn, user, filter)
		if err != nil {
			return err
		}
		fmt.Println("Authored\t" + intsToStr(prs))
		crs, err := karma.Reviews(fn, user, filter)
		if err != nil {
			return err
		}
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
