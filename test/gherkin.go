package testutils

import (
	"strings"
	"testing"

	"github.com/axelarnetwork/utils/slices"
)

const (
	_given = "GIVEN"
	_when  = "WHEN"
	_then  = "THEN"
	_and   = "AND"
)

// GivenStatement is used to set up unit test preconditions
type GivenStatement struct {
	label []string
	test  func()
}

// WhenStatement is used to define conditions under test
type WhenStatement struct {
	label []string
	test  func()
}

type thenStatement struct {
	label []string
	test  func(t *testing.T)
}

// ThenStatements are used to define test outcomes
type ThenStatements []thenStatement

// Given starts the test with the first precondition
func Given(description string, setup func()) GivenStatement {
	return GivenStatement{
		label: []string{_given, description},
		test:  setup,
	}
}

// When is an independent trigger that can be used to start a statement in a Branch
func When(description string, setup func()) WhenStatement {
	return WhenStatement{
		label: []string{_when, description},
		test:  setup,
	}
}

// Then is an independent outcome check that can be used to start a statement in a Branch
func Then(description string, setup func(t *testing.T)) ThenStatements {
	return []thenStatement{{
		label: []string{_then, description},
		test:  setup,
	}}
}

// Given adds an additional precondition
func (g GivenStatement) Given(description string, setup func()) GivenStatement {
	return GivenStatement{
		label: append(g.label, _and, _given, description),
		test:  func() { g.test(); setup() },
	}
}

func (g GivenStatement) Given2(g2 GivenStatement) GivenStatement {
	return GivenStatement{
		label: mergeLabels(g.label, g2.label),
		test:  func() { g.test(); g2.test() },
	}
}

// When adds a trigger to the test path
func (g GivenStatement) When(description string, setup func()) WhenStatement {
	return WhenStatement{
		label: append(g.label, _when, description),
		test:  func() { g.test(); setup() },
	}
}

func (g GivenStatement) When2(w WhenStatement) WhenStatement {
	return WhenStatement{
		label: mergeLabels(g.label, w.label),
		test:  func() { g.test(); w.test() },
	}
}

// Branch allows test branching by adding multiple sub-statements after a Given
func (g GivenStatement) Branch(thens ...ThenStatements) ThenStatements {
	slices.ForEach(thens, func(thens ThenStatements) { checkFirstWordOfLabels(thens, assertNotTHEN) })

	out := ThenStatements{}
	for _, then := range thens {
		for _, statement := range then {
			statement := statement
			out = append(out, thenStatement{
				label: mergeLabels(g.label, statement.label),
				test:  func(t *testing.T) { g.test(); statement.test(t) },
			})
		}
	}

	return out
}

// Branch allows test branching by adding multiple sub-statements after a When
func (w WhenStatement) Branch(thens ...ThenStatements) ThenStatements {
	slices.ForEach(thens, func(thens ThenStatements) { checkFirstWordOfLabels(thens, assertNotGIVEN) })

	out := ThenStatements{}
	for _, then := range thens {
		for _, statement := range then {
			statement := statement
			out = append(out, thenStatement{
				label: mergeLabels(w.label, statement.label),
				test:  func(t *testing.T) { w.test(); statement.test(t) },
			})
		}
	}

	return out
}

// When adds a trigger to the test path
func (w WhenStatement) When(description string, setup func()) WhenStatement {
	return WhenStatement{
		label: append(w.label, _and, _when, description),
		test:  func() { w.test(); setup() },
	}
}

func (w WhenStatement) When2(w2 WhenStatement) WhenStatement {
	return WhenStatement{
		label: mergeLabels(w.label, w2.label),
		test:  func() { w.test(); w2.test() },
	}
}

// Then adds an outcome check to the test path
func (w WhenStatement) Then(description string, execution func(t *testing.T)) ThenStatements {
	return []thenStatement{{
		label: append(w.label, _then, description),
		test:  func(t *testing.T) { w.test(); execution(t) },
	}}
}

// Then2 allows the use of a previously defined "then" statement
func (w WhenStatement) Then2(then ThenStatements) ThenStatements {
	return slices.Map(then, func(then thenStatement) thenStatement {
		return thenStatement{
			label: mergeLabels(w.label, then.label),
			test:  func(t *testing.T) { w.test(); then.test(t) },
		}
	})
}

// Run executes all defined test paths. Optionally, each path is repeated a given number of times
func (then thenStatement) Run(t *testing.T, repeats ...int) bool {
	repeat := 1
	if len(repeats) == 1 && repeats[0] > 1 {
		repeat = repeats[0]
	}
	return t.Run(strings.Join(then.label, " "), Func(then.test).Repeat(repeat))
}

// Run executes all defined test paths. Optionally, each path is repeated a given number of times
func (thens ThenStatements) Run(t *testing.T, repeats ...int) bool {
	return slices.Reduce(thens, true, func(result bool, then thenStatement) bool { return then.Run(t, repeats...) })
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

func checkFirstWordOfLabels(tests []thenStatement, assertion func(string)) {
	for _, test := range tests {
		assertion(test.label[0])
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
