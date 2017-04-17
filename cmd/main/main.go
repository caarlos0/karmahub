package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/karmahub"
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
			Name:  "user, u",
			Usage: "To collect data from. Defaults to the GitHub token owner",
		},
		cli.StringFlag{
			Name:  "filter",
			Usage: "Additional filters, github syntax. E.g.: user:Docker will filter for issues and pulls in Docker organization",
		},
	}
	app.Action = func(c *cli.Context) error {
		var token = c.String("token")
		var user = c.String("user")
		var filter = c.String("filter")
		if token == "" {
			return cli.NewExitError("missing github token", 1)
		}
		var spin = spin.New("%s Gathering data...")
		spin.Start()

		var ctx = context.Background()
		var client = github.NewClient(oauth2.NewClient(
			ctx,
			oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
		))
		var fn = karmahub.GitHubSearch(ctx, client)
		if user == "" {
			guser, _, err := client.Users.Get(ctx, "")
			if err != nil {
				spin.Stop()
				return cli.NewExitError(err.Error(), 1)
			}
			user = *guser.Login
		}
		prs, err := karmahub.Authors(fn, user, filter)
		if err != nil {
			spin.Stop()
			return cli.NewExitError(err.Error(), 1)
		}
		crs, err := karmahub.Reviews(fn, user, filter)
		if err != nil {
			spin.Stop()
			return cli.NewExitError(err.Error(), 1)
		}
		spin.Stop()

		fmt.Printf(
			"\033[1;36mAction    \t%v\t%v\t%v\033[0m\n",
			month(0),
			month(-1),
			month(-2),
		)
		fmt.Printf("Authored    \t%v\t%v\t%v\n", prs[0], prs[1], prs[2])
		fmt.Printf("Reviewed    \t%v\t%v\t%v\n", crs[0], crs[1], crs[2])
		fmt.Printf("Karma score \t%v\t%v\t%v\n", crs[0]/prs[0], crs[1]/prs[1], crs[2]/prs[2])
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func month(i int) string {
	return time.Now().AddDate(0, i, 0).UTC().Format("Jan")
}
