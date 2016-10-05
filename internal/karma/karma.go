package karma

import (
	"time"

	"github.com/caarlos0/karmahub/internal/search"
)

// MONTHS of data gathered
const MONTHS = 3

// Authors in each of the last MONTHS
func Authors(search search.Fn, login, filter string) (result []int, err error) {
	return totals(search, filter+" author:"+login)
}

// Reviews in each of the last MONTHS
func Reviews(search search.Fn, login, filter string) (result []int, err error) {
	mine, err := totals(search, filter+" commenter:"+login+" author:"+login)
	if err != nil {
		return result, err
	}
	all, err := totals(search, filter+" commenter:"+login)
	if err != nil {
		return result, err
	}
	for i := 0; i < len(all); i++ {
		result = append(result, all[i]-mine[i])
	}
	return result, err
}

func totals(search search.Fn, query string) (result []int, err error) {
	var counts []int
	for i := 1; i <= MONTHS; i++ {
		d := time.Now().AddDate(0, i*-1, 0).Format("2006-01-02")
		count, err := search(query + " created:>" + d)
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
