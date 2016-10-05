package search

import (
	"log"
	"time"

	"github.com/google/go-github/github"
)

// Fn is function that searches somewhere and returns the count of
// results only (or error)
type Fn func(query string) (total int, err error)

var filter = &github.SearchOptions{
	ListOptions: github.ListOptions{
		PerPage: 1,
	},
}

// Github is a Fn impl for github
func Github(client *github.Client) Fn {
	var fn Fn
	fn = func(query string) (total int, err error) {
		result, _, err := client.Search.Issues(query, filter)
		if _, ok := err.(*github.RateLimitError); ok {
			log.Println("Rate limit, waiting 10s...")
			time.Sleep(time.Second * 10)
			return fn(query)
		}
		if result.Total != nil {
			return *result.Total, nil
		}
		return 0, err
	}
	return fn
}
