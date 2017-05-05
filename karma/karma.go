package karma

import "time"

// Authors in each of the last MONTHS
func Authors(search SearchFn, login, filter string, months int) (result []int, err error) {
	return totals(search, filter+" author:"+login, months)
}

// Reviews in each of the last MONTHS
func Reviews(search SearchFn, login, filter string, months int) (result []int, err error) {
	mine, err := totals(search, filter+" commenter:"+login+" author:"+login, months)
	if err != nil {
		return result, err
	}
	all, err := totals(search, filter+" commenter:"+login, months)
	if err != nil {
		return result, err
	}
	for i := 0; i < len(all); i++ {
		result = append(result, all[i]-mine[i])
	}
	return result, err
}

func totals(search SearchFn, query string, months int) (result []int, err error) {
	var counts []int
	for i := 1; i <= months; i++ {
		d := time.Now().AddDate(0, i*-1, 0).Format("2006-01-02")
		count, serr := search(query + " created:>" + d)
		if serr != nil {
			return result, serr
		}
		counts = append(counts, count)
	}
	result = append(result, counts[0])
	for i := 1; i <= months-1; i++ {
		result = append(result, counts[i]-counts[i-1])
	}
	return
}
