package testutils_test

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/axelarnetwork/utils/test"
	"github.com/stretchr/testify/assert"
)

func TestGherkinSyntax(t *testing.T) {
	var testLabel string
	var testPaths int
	testSetup := Given("some setup", func(t *testing.T) { testLabel = "GIVEN some setup" }).
		Branch(
			Given("additional setup", func(t *testing.T) { testLabel += " AND GIVEN additional setup" }).
				And().Given("even more setup", func(t *testing.T) { testLabel += " AND GIVEN even more setup" }).
				When("a trigger happens", func(t *testing.T) { testLabel += " WHEN a trigger happens" }).
				Branch(
					When("a second trigger happens", func(t *testing.T) { testLabel += " AND WHEN a second trigger happens" }).
						And().When("a third trigger happens", func(t *testing.T) { testLabel += " AND WHEN a third trigger happens" }).
						Then("we finally check the outcome", func(t *testing.T) {
							testLabel += " THEN we finally check the outcome"
							assertNameEquals(t, testLabel)
							testPaths++
						}),
					Then("we check the outcome directly", func(t *testing.T) {
						testLabel += " THEN we check the outcome directly"
						assertNameEquals(t, testLabel)
						testPaths++
					}),
				),
			When("we directly hit the trigger", func(t *testing.T) { testLabel += " WHEN we directly hit the trigger" }).
				Then("we check the outcome even earlier", func(t *testing.T) {
					testLabel += " THEN we check the outcome even earlier"
					assertNameEquals(t, testLabel)
					testPaths++
				}),
			TestCases([]int16{1, 2, 3}).ForEach(func(tc int16) Runner {
				description := fmt.Sprintf("test case %v", tc)
				return When(description, func(t *testing.T) { testLabel += fmt.Sprintf(" WHEN %s", description) }).
					Then("we check outcome for test case", func(t *testing.T) {
						testLabel += " THEN we check outcome for test case"
						assertNameEquals(t, testLabel)
						testPaths++
					})
			}),
		)

	testSetup.Run(t)
	assert.Equal(t, 6, testPaths)

	// do the same execution again, so tests will end in "#01"
	testPaths = 0
	testSetup.Run(t, 15)
	assert.Equal(t, 6*15, testPaths)
}

func TestGherkinPanicsGIVENAfterWHEN(t *testing.T) {
	assert.Panics(t, func() {
		Given("precondition", func(*testing.T) {}).
			When("trigger", func(*testing.T) {}).
			Branch(
				Given("forbidden statement", func(*testing.T) {}).
					When("trigger", func(*testing.T) {}).
					Then("check", func(*testing.T) {}),
			)
	})
}

func TestGherkinPanicsTHENAfterGIVEN(t *testing.T) {
	assert.Panics(t, func() {
		Given("precondition", func(*testing.T) {}).
			Branch(
				Then("check", func(*testing.T) {}),
			)
	})
}

func assertNameEquals(t *testing.T, testLabel string) bool {
	// testname has form "testfunc/test_run_name#repetition"
	name := t.Name()
	name = strings.Split(name, "/")[1]
	name = strings.Split(name, "#")[0]
	name = strings.ReplaceAll(name, "_", " ")

	return assert.Equal(t, testLabel, name)
}
