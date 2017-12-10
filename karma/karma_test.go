package karma

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthors(t *testing.T) {
	var c int
	fn := func(query string) (total int, err error) {
		c++
		return 10 * c, nil
	}

	res, err := Authors(fn, "caarlos0", "is:pr", 3)
	assert.NoError(t, err)
	assert.Equal(t, []int{10, 10, 10}, res)
}

func TestReviews(t *testing.T) {
	var c int
	results := []int{10, 15, 18, 20, 31, 39}
	fn := func(query string) (total int, err error) {
		c++
		return results[c-1], nil
	}

	res, err := Reviews(fn, "caarlos0", "is:pr", 3)
	assert.NoError(t, err)
	assert.Equal(t, []int{10, 6, 5}, res)
}

func TestErrs(t *testing.T) {
	fn := func(query string) (total int, err error) {
		return 0, errors.New("BREAK")
	}

	res, err := Authors(fn, "caarlos0", "is:pr", 3)
	assert.Error(t, err)
	assert.Empty(t, res)
	res, err = Reviews(fn, "caarlos0", "is:pr", 3)
	assert.Error(t, err)
	assert.Empty(t, res)
}
