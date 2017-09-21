package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/karmahub/karma"
	"github.com/caarlos0/spin"
	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var token string
var user string
var filter string
var months int

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "karmahub",
	Short: "get your history of reviews/comments and pull requests/issues opened",
	Long: `Compares the amount of issues and pull requests you created with the
amount of comments  and code reviews you did.

The idea is to use it at your daily job organization, so you can get an idea
of how much are you actually contributing to the code review practices.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if token == "" {
			return fmt.Errorf("missing github token")
		}
		var spin = spin.New("%s Gathering data...").Start()

		var ctx = context.Background()
		var client = github.NewClient(oauth2.NewClient(
			ctx,
			oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
		))
		var fn = karma.GitHubSearch(ctx, client)
		if user == "" {
			me, _, err := client.Users.Get(ctx, "")
			if err != nil {
				spin.Stop()
				return err
			}
			user = *me.Login
		}
		prs, err := karma.Authors(fn, user, filter, months)
		if err != nil {
			spin.Stop()
			return err
		}
		crs, err := karma.Reviews(fn, user, filter, months)
		if err != nil {
			spin.Stop()
			return err
		}
		spin.Stop()

		// header
		fmt.Printf("\033[1;36mAction    ")
		for i := 0; i < months; i++ {
			fmt.Printf(
				"\t%v",
				time.Now().AddDate(0, i*-1, 0).UTC().Format("Jan"),
			)
		}
		fmt.Printf("\033[0m\n")

		// authored
		fmt.Printf("Authored    ")
		for i := 0; i < months; i++ {
			fmt.Printf("\t%v", prs[i])
		}
		fmt.Printf("\n")

		// reviewed
		fmt.Printf("Reviewed    ")
		for i := 0; i < months; i++ {
			fmt.Printf("\t%v", crs[i])
		}
		fmt.Printf("\n")

		// karma
		fmt.Printf("Karma    ")
		for i := 0; i < months; i++ {
			var karma = float32(crs[i]) / float32(prs[i])
			if prs[i] == 0 {
				karma = float32(crs[i])
			}
			fmt.Printf("\t%.1f", karma)
		}
		fmt.Printf("\n")
		return nil
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	Version = version
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(
		&token,
		"token",
		"t",
		os.Getenv("GITHUB_TOKEN"),
		"Your GitHub token",
	)
	RootCmd.PersistentFlags().StringVarP(
		&user,
		"user",
		"u",
		"",
		"User to collect data from",
	)
	RootCmd.PersistentFlags().IntVarP(
		&months,
		"months",
		"m",
		3,
		"Number of months to search",
	)
	RootCmd.PersistentFlags().StringVarP(
		&filter,
		"filter",
		"f",
		"",
		"Additional filters, github syntax. E.g.: is:pr will gather data for pull requests only",
	)
}
