package testutils_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/axelarnetwork/utils/test"
)

func TestGherkinSyntax(t *testing.T) {
	var testLabel string
	var testPaths int
	testSetup := Given("some setup", func() { testLabel = "GIVEN some setup" }).
		Branch(
			Given("additional setup", func() { testLabel += " AND GIVEN additional setup" }).
				Given("even more setup", func() { testLabel += " AND GIVEN even more setup" }).
				When("a trigger happens", func() { testLabel += " WHEN a trigger happens" }).
				Branch(
					When("a second trigger happens", func() { testLabel += " AND WHEN a second trigger happens" }).
						When("a third trigger happens", func() { testLabel += " AND WHEN a third trigger happens" }).
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
			When("we directly hit the trigger", func() { testLabel += " WHEN we directly hit the trigger" }).
				Then("we check the outcome even earlier", func(t *testing.T) {
					testLabel += " THEN we check the outcome even earlier"
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

func TestGherkinPanicsGIVENAfterWHEN(t *testing.T) {
	assert.Panics(t, func() {
		Given("precondition", func() {}).
			When("trigger", func() {}).
			Branch(
				Given("forbidden statement", func() {}).
					When("trigger", func() {}).
					Then("check", func(*testing.T) {}),
			)
	})
}

func TestGherkinPanicsTHENAfterGIVEN(t *testing.T) {
	assert.Panics(t, func() {
		Given("precondition", func() {}).
			Branch(
				Then("check", func(*testing.T) {}),
			)
	})
}

func TestGherkinSeparateStatements(t *testing.T) {
	var testLabel string
	givenSetup := Given("the setup", func() { testLabel = "GIVEN the setup" })
	moreSetup := Given("more setup", func() { testLabel += " AND GIVEN more setup" })
	somethingHappens := When("something happens", func() { testLabel += " WHEN something happens" })
	somethingElseHappens := When("something else happens", func() { testLabel += " AND WHEN something else happens" })
	assertSomething := Then("assert something", func(t *testing.T) {
		testLabel += " THEN assert something"
		assertNameEquals(t, testLabel)
	})

	assertOneThing := Then("assert something", func(t *testing.T) {
		testLabel += " THEN assert something"
	})

	assertSomethingElse := Then("assert something else", func(t *testing.T) {
		testLabel += " AND THEN assert something else"
		assertNameEquals(t, testLabel)
	})

	givenSetup.Given2(moreSetup).When2(somethingHappens).Then2(assertSomething).Run(t)
	givenSetup.Branch(
		somethingHappens.Then2(assertSomething),
		moreSetup.When2(somethingHappens).When2(somethingElseHappens).Then2(assertOneThing).Then2(assertSomethingElse),
	).Run(t)
}

func assertNameEquals(t *testing.T, testLabel string) bool {
	// testname has form "testfunc/test_run_name#repetition"
	name := t.Name()
	name = strings.Split(name, "/")[1]
	name = strings.Split(name, "#")[0]
	name = strings.ReplaceAll(name, "_", " ")

	return assert.Equal(t, testLabel, name)
}
