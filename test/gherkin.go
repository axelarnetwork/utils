package testutils

import (
	"strings"
	"testing"
)

const (
	_given = "GIVEN"
	_when  = "WHEN"
	_then  = "THEN"
	_and   = "AND"
)

// GivenStatement is used to set up unit test preconditions - do not implement yourself
type GivenStatement interface {
	Given(description string, setup func(t *testing.T)) WhenStatement
}

// WhenStatement is used to hit a trigger - do not implement yourself
type WhenStatement interface {
	And() GivenStatement
	Branch(...Runner) Runner
	When(description string, setup func(t *testing.T)) ThenStatement
}

// ThenStatement is used to check the test outcome - do not implement yourself
type ThenStatement interface {
	And() WhenStatement
	Branch(...Runner) Runner
	Then(description string, execution func(t *testing.T)) Runner
}

// ForEachStatement is used to create a set of tests for each of the test case items
type ForEachStatement[T any] interface {
	ForEach(createRunner func(testCase T) Runner) Runner
}

// Runner executes the test - do not implement yourself
type Runner interface {
	Run(t *testing.T, repeats ...int) bool
}

type given struct {
	label []string
	test  []func(t *testing.T)
}

type when struct {
	label []string
	test  []func(t *testing.T)
}

type then struct {
	label []string
	test  []func(t *testing.T)
}

type testCases[T any] struct {
	items []T
}

type multiRunner struct {
	runners []runner
}

type runner struct {
	label []string
	test  []func(t *testing.T)
}

// Given starts the test with the first precondition
func Given(description string, setup func(t *testing.T)) WhenStatement {
	return when{
		label: []string{_given, description},
		test:  []func(t *testing.T){setup},
	}
}

// When is an independent trigger that can be used to start a statement in a Branch
func When(description string, setup func(t *testing.T)) ThenStatement {
	return then{
		label: []string{_when, description},
		test:  []func(t *testing.T){setup},
	}
}

// Then is an independent outcome check that can be used to start a statement in a Branch
func Then(description string, setup func(t *testing.T)) Runner {
	return runner{
		label: []string{_then, description},
		test:  []func(t *testing.T){setup},
	}
}

// TestCases takes an array of arbitrary items to then generate a test for each of the test case items
func TestCases[T any](items []T) ForEachStatement[T] {
	return testCases[T]{
		items: items,
	}
}

// Given adds an additional precondition
func (g given) Given(description string, setup func(t *testing.T)) WhenStatement {
	return when{
		label: append(g.label, _given, description),
		test:  append(g.test, setup),
	}
}

// And allows to concatenate an additional precondition
func (w when) And() GivenStatement {
	return given{
		label: append(w.label, _and),
		test:  w.test,
	}
}

// When adds a trigger to the test path
func (w when) When(description string, setup func(t *testing.T)) ThenStatement {
	return then{
		label: append(w.label, _when, description),
		test:  append(w.test, setup),
	}
}

// Branch allows test branching by adding multiple sub-statements after a Given
func (w when) Branch(tests ...Runner) Runner {
	checkFirstWordOfLabels(tests, assertNotTHEN)
	return branch(w.label, w.test, tests)
}

// Branch allows test branching by adding multiple sub-statements after a When
func (t then) Branch(tests ...Runner) Runner {
	checkFirstWordOfLabels(tests, assertNotGIVEN)
	return branch(t.label, t.test, tests)
}

// And allows to concatenate an additional trigger
func (t then) And() WhenStatement {
	return when{
		label: append(t.label, _and),
		test:  t.test,
	}
}

// Then adds an outcome check to the test path
func (t then) Then(description string, execution func(t *testing.T)) Runner {
	return runner{
		label: append(t.label, _then, description),
		test:  append(t.test, execution),
	}
}

// ForEach generates a test for each test case
func (tc testCases[T]) ForEach(createRunner func(testCase T) Runner) Runner {

	runners := make([]Runner, len(tc.items))

	for i, v := range tc.items {
		runners[i] = createRunner(v)
	}

	return branch([]string{}, []func(t *testing.T){}, runners)
}

// Run executes all defined test paths. Optionally, each path is repeated a given number of times
func (r multiRunner) Run(t *testing.T, repeats ...int) bool {
	result := true
	for _, runner := range r.runners {
		// cannot inline this because if result is false the second part of && is not going to be evaluated
		newResult := runner.Run(t, repeats...)
		result = result && newResult
	}
	return result
}

// Run executes all defined test paths. Optionally, each path is repeated a given number of times
func (r runner) Run(t *testing.T, repeats ...int) bool {
	repeat := 1
	if len(repeats) == 1 && repeats[0] > 1 {
		repeat = repeats[0]
	}
	return t.Run(strings.Join(r.label, " "), Func(func(t *testing.T) {
		for _, f := range r.test {
			f(t)
		}
	}).Repeat(repeat))
}

func checkFirstWordOfLabels(tests []Runner, assertion func(string)) {
	for _, test := range tests {
		switch test := test.(type) {
		case runner:
			assertion(test.label[0])
		case multiRunner:
			for _, r := range test.runners {
				assertion(r.label[0])
			}
		default:
			panic("do not extend the test suite with custom types")
		}
	}
}

func assertNotGIVEN(label string) {
	if label == _given {
		panic("GIVEN is not allowed after WHEN")
	}
}

func assertNotTHEN(label string) {
	if label == _then {
		panic("THEN is not allowed after GIVEN")
	}
}

func branch(startLabel []string, startTest []func(t *testing.T), tests []Runner) Runner {
	var mr multiRunner
	for _, test := range tests {
		switch test := test.(type) {
		case runner:
			mr.runners = append(mr.runners, concatRunner(startLabel, startTest, test))
		case multiRunner:
			for _, r := range test.runners {
				mr.runners = append(mr.runners, concatRunner(startLabel, startTest, r))
			}
		default:
			panic("do not extend the test suite with custom types")
		}
	}
	return mr
}

func concatRunner(startLabel []string, startTest []func(t *testing.T), r runner) runner {
	return runner{
		label: mergeLabels(startLabel, r.label),
		test:  append(startTest, r.test...),
	}
}

func mergeLabels(startLabel, endLabel []string) []string {
	label := startLabel
	startLabelLength := len(startLabel)
	// second last word is always the last action word
	if startLabelLength >= 2 && startLabel[startLabelLength-2] == endLabel[0] {
		label = append(label, _and)
	}
	return append(label, endLabel...)
}
