package metrics

import (
	"math/rand"
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumberStringSlice(t *testing.T) {
	for range [100]struct{}{} {
		a := make([]int, 100)
		s := make([]string, 100)
		for i := range a {
			a[i] = rand.Intn(100)
			s[i] = strconv.Itoa(a[i])
		}
		sort.Sort(sort.IntSlice(a))
		sort.Sort(StringNumberSlice(s))
		for i := range a {
			assert.Equal(t, s[i], strconv.Itoa(a[i]))
		}
	}
}
