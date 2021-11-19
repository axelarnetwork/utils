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
	Or(...Runner) Runner
	When(description string, setup func(t *testing.T)) ThenStatement
}

// ThenStatement is used to check the test outcome - do not implement yourself
type ThenStatement interface {
	And() WhenStatement
	Or(...Runner) Runner
	Then(description string, execution func(t *testing.T)) Runner
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

// When is an independent trigger that can be used to start a statement in an Or condition
func When(description string, setup func(t *testing.T)) ThenStatement {
	return then{
		label: []string{_when, description},
		test:  []func(t *testing.T){setup},
	}
}

// Then is an independent outcome check that can be used to start a statement in an Or condition
func Then(description string, setup func(t *testing.T)) Runner {
	return runner{
		label: []string{_then, description},
		test:  []func(t *testing.T){setup},
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

// Or allows test branching by adding multiple sub-statements after a GIVEN
func (w when) Or(tests ...Runner) Runner {
	checkFirstWordOfLabels(tests, assertNotTHEN)
	return or(w.label, w.test, tests)
}

// Or allows test branching by adding multiple sub-statements after a WHEN
func (t then) Or(tests ...Runner) Runner {
	checkFirstWordOfLabels(tests, assertNotGIVEN)
	return or(t.label, t.test, tests)
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

func or(startLabel []string, startTest []func(t *testing.T), tests []Runner) Runner {
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
	// second last word is always the last action word
	if startLabel[len(startLabel)-2] == endLabel[0] {
		label = append(label, _and)
	}
	return append(label, endLabel...)
}
