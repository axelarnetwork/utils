package testutils_test

import (
	"strings"
	"testing"

	test "github.com/axelarnetwork/utils/test"
	"github.com/stretchr/testify/assert"
)

func TestGherkin(t *testing.T) {
	var testLabel string
	var testPaths int
	testSetup := test.Given("some setup", func(t *testing.T) { testLabel = "GIVEN some setup" }).
		Or(
			test.Given("additional setup", func(t *testing.T) { testLabel += " AND GIVEN additional setup" }).
				And().Given("even more setup", func(t *testing.T) { testLabel += " AND GIVEN even more setup" }).
				When("a condition is added", func(t *testing.T) { testLabel += " WHEN a condition is added" }).
				Or(
					test.When("a second condition is added", func(t *testing.T) { testLabel += " AND WHEN a second condition is added" }).
						And().When("a third condition is added", func(t *testing.T) { testLabel += " AND WHEN a third condition is added" }).
						Then("we finally execute", func(t *testing.T) {
							testLabel += " THEN we finally execute"
							assertNameEquals(t, testLabel)
							testPaths++
						}),
					test.Then("we execute directly", func(t *testing.T) {
						testLabel += " THEN we execute directly"
						assertNameEquals(t, testLabel)
						testPaths++
					}),
				),
			test.When("we directly add a condition", func(t *testing.T) { testLabel += " WHEN we directly add a condition" }).
				Then("we execute even earlier", func(t *testing.T) {
					testLabel += " THEN we execute even earlier"
					assertNameEquals(t, testLabel)
					testPaths++
				}),
		)

	testSetup.Run(t)
	assert.Equal(t, 3, testPaths)

	// do the same execution again, so tests will end in "#01"
	testPaths = 0
	testSetup.Run(t, 15)
	assert.Equal(t, 3*15, testPaths)
}

func assertNameEquals(t *testing.T, testLabel string) bool {
	// testname has form "testfunc/test_run_name#repetition"
	name := t.Name()
	name = strings.Split(name, "/")[1]
	name = strings.Split(name, "#")[0]
	name = strings.ReplaceAll(name, "_", " ")

	return assert.Equal(t, testLabel, name)
}
