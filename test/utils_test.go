package testutils_test

import (
	"testing"

	testutils "github.com/axelarnetwork/utils/test"
)

func TestTestCases(t *testing.T) {
	sum := 0
	testutils.AsTestCases([]int{1, 2, 3}...).
		ForEach(
			func(t *testing.T, tc int) {
				sum += tc
			}).
		Run(t)

	sum = 0
}
