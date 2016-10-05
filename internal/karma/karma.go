package karma

import (
	"log"
	"time"

	"github.com/google/go-github/github"
)

// MONTHS of data gathered
const MONTHS = 3

// Pulls in each of the last MONTHS
func Pulls(client *github.Client, login, filter string) (result []int, err error) {
	return totals(client, filter+" author:"+login)
}

// Reviews in each of the last MONTHS
func Reviews(client *github.Client, login, filter string) (result []int, err error) {
	all, err := totals(client, filter+" commenter:"+login)
	if err != nil {
		return result, err
	}
	mine, err := totals(client, filter+" commenter:"+login+" author:"+login)
	if err != nil {
		return result, err
	}
	for i := 0; i < len(all); i++ {
		result = append(result, all[i]-mine[i])
	}
	return result, err
}

func totals(client *github.Client, search string) (result []int, err error) {
	var counts []int
	for i := 1; i <= MONTHS; i++ {
		d := time.Now().AddDate(0, i*-1, 0).Format("2006-01-02")
		count, err := total(client, search+" created:>"+d)
		if err != nil {
			return result, err
		}
		counts = append(counts, count)
	}
	result = append(result, counts[0])
	for i := 1; i <= MONTHS-1; i++ {
		result = append(result, counts[i]-counts[i-1])
	}
	return result, err
}

func total(client *github.Client, search string) (count int, err error) {
	result, _, err := client.Search.Issues(
		search,
		&github.SearchOptions{
			ListOptions: github.ListOptions{
				PerPage: 1,
			},
		},
	)
	if _, ok := err.(*github.RateLimitError); ok {
		log.Println("Rate limit, waiting 10s...")
		time.Sleep(time.Second * 10)
		return total(client, search)
	}
	if result.Total != nil {
		count = *result.Total
	}
	return count, err
}
