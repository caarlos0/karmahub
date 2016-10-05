package karma_test

import (
	"errors"
	"testing"

	"github.com/caarlos0/karmahub/internal/karma"
	"github.com/stretchr/testify/assert"
)

func TestAuthors(t *testing.T) {
	assert := assert.New(t)
	var c int
	fn := func(query string) (total int, err error) {
		c++
		return 10 * c, nil
	}

	res, err := karma.Authors(fn, "caarlos0", "is:pr")
	assert.NoError(err)
	assert.Equal([]int{10, 10, 10}, res)
}

func TestReviews(t *testing.T) {
	assert := assert.New(t)
	var c int
	results := []int{10, 15, 18, 20, 31, 39}
	fn := func(query string) (total int, err error) {
		c++
		return results[c-1], nil
	}

	res, err := karma.Reviews(fn, "caarlos0", "is:pr")
	assert.NoError(err)
	assert.Equal([]int{10, 6, 5}, res)
}

func TestErrs(t *testing.T) {
	assert := assert.New(t)
	fn := func(query string) (total int, err error) {
		return 0, errors.New("BREAK")
	}

	res, err := karma.Authors(fn, "caarlos0", "is:pr")
	assert.Error(err)
	assert.Empty(res)
	res, err = karma.Reviews(fn, "caarlos0", "is:pr")
	assert.Error(err)
	assert.Empty(res)
}
