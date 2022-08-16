package slices_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/axelarnetwork/utils/funcs"
	"github.com/axelarnetwork/utils/slices"

	"github.com/axelarnetwork/utils/test/rand"
)

func TestMap(t *testing.T) {
	source := make([]string, 0, 20)

	for i := 0; i < 20; i++ {
		source = append(source, rand.StrBetween(5, 20))
	}

	out := slices.Map(source, func(s string) int { return len(s) })

	for i := range out {
		assert.Equal(t, len(source[i]), out[i])
	}
}

func TestAll(t *testing.T) {
	even := make([]int, 0, 20)

	for i := 0; i < 20; i += 2 {
		even = append(even, i)
	}

	assert.True(t, slices.All(even, func(i int) bool { return i%2 == 0 }))
	assert.False(t, slices.All(append(even, 5), func(i int) bool { return i%2 == 0 }))
}

func TestAny(t *testing.T) {
	even := make([]int, 0, 20)

	for i := 0; i < 20; i += 2 {
		even = append(even, i)
	}

	assert.False(t, slices.Any(even, func(i int) bool { return i%2 != 0 }))
	assert.True(t, slices.Any(append(even, 5), func(i int) bool { return i%2 != 0 }))
}

func TestFilter(t *testing.T) {
	source := make([]int, 0, 100)

	for i := 0; i < 100; i++ {
		source = append(source, i)
	}

	even := slices.Filter(source, func(i int) bool { return i%2 == 0 })

	for _, x := range even {
		assert.Equal(t, 0, x%2)
	}
}

func TestForEach(t *testing.T) {
	source := make([]int, 0, 100)
	total := 0

	for i := 0; i < 100; i++ {
		source = append(source, 1)
	}

	slices.ForEach(source, func(n int) {
		total += n
	})

	assert.Equal(t, 100, total)
}

func TestReduce(t *testing.T) {
	source := make([]string, 0, 100)
	for i := 0; i < 100; i++ {
		source = append(source, rand.Str(i))
	}

	assert.Equal(t, 4950, slices.Reduce(source, 0, func(v int, i string) int { return v + len(i) }))
}

func TestFlatten(t *testing.T) {
	source := make([][]int, 0, 10)

	n := 0
	for i := 0; i < 10; i++ {
		source = append(source, []int{})
		for j := 0; j < 10; j++ {
			source[i] = append(source[i], n)
			n++
		}
	}

	f := slices.Flatten(source)

	assert.Len(t, f, 100)

	for _, i := range f {
		assert.Equal(t, i, f[i])
	}
}

func TestExpand(t *testing.T) {
	out := slices.Expand(strconv.Itoa, 5)

	assert.Equal(t, []string{"0", "1", "2", "3", "4"}, out)
}

func TestWhile(t *testing.T) {
	source := []int{1, 2, 3, 4, 5, 6, 7}
	sum := 0
	slices.While(source, func(i int) bool {
		sum += i
		return sum < 10
	})

	assert.Equal(t, 10, sum)
}

func TestDistinct(t *testing.T) {
	out := slices.Distinct([]int{0, 3, 2, 7, 2, 1, 3, 0})
	assert.Equal(t, []int{0, 3, 2, 7, 1}, out)
}

func TestToMap(t *testing.T) {
	source := []int{0, 3, 2, 7, 3}
	m := slices.ToMap(source, strconv.Itoa)
	assert.Len(t, m, 4)
	assert.Equal(t, 0, m["0"])
	assert.Equal(t, 2, m["2"])
	assert.Equal(t, 3, m["3"])
	assert.Equal(t, 7, m["7"])

	assert.Panics(t, func() { slices.ToMap(source, strconv.Itoa, true) })
}

func TestCast(t *testing.T) {
	source := slices.Expand(funcs.Identity[int], 10)
	assert.Empty(t, slices.TryCast[int, string](source))

	assert.Len(t, slices.TryCast[int, interface{}](source), 10)
}

func TestReverse(t *testing.T) {
	source := []int{1, 2, 3, 4, 5}
	assert.Equal(t, []int{5, 4, 3, 2, 1}, slices.Reverse(source))
}

func TestGroupBy(t *testing.T) {
	source := []int{1, 2, 3, 4, 5}
	assert.Equal(t, map[bool][]int{true: {2, 4}, false: {1, 3, 5}}, slices.GroupBy(source, func(i int) bool { return i%2 == 0 }))
}
