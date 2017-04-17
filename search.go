package karmahub

import (
	"context"
	"time"

	"github.com/google/go-github/github"
)

// SearchFn is function that searches somewhere and returns the count of
// results only (or error)
type SearchFn func(query string) (total int, err error)

var filter = &github.SearchOptions{
	ListOptions: github.ListOptions{
		PerPage: 1,
	},
}

// GitHubSearch is a SearchFn impl for github
func GitHubSearch(ctx context.Context, client *github.Client) SearchFn {
	var fn SearchFn
	fn = func(query string) (total int, err error) {
		result, _, err := client.Search.Issues(ctx, query, filter)
		if _, ok := err.(*github.RateLimitError); ok {
			time.Sleep(time.Second * 10)
			return fn(query)
		}
		if _, ok := err.(*github.AcceptedError); ok {
			time.Sleep(time.Second * 2)
			return fn(query)
		}
		if result.Total != nil {
			return *result.Total, nil
		}
		return 0, err
	}
	return fn
}
