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

type GivenStatement interface {
	Given(description string, setup func(t *testing.T)) WhenStatement
}

type WhenStatement interface {
	And() GivenStatement
	Or(...Runner) Runner
	When(description string, setup func(t *testing.T)) ThenStatement
}

type ThenStatement interface {
	And() WhenStatement
	Or(...Runner) Runner
	Then(description string, execution func(t *testing.T)) Runner
}

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

func Given(description string, setup func(t *testing.T)) WhenStatement {
	return when{
		label: []string{_given, description},
		test:  []func(t *testing.T){setup},
	}
}

func When(description string, setup func(t *testing.T)) ThenStatement {
	return then{
		label: []string{_when, description},
		test:  []func(t *testing.T){setup},
	}
}

func Then(description string, setup func(t *testing.T)) Runner {
	return runner{
		label: []string{_then, description},
		test:  []func(t *testing.T){setup},
	}
}

func (g given) Given(description string, setup func(t *testing.T)) WhenStatement {
	return when{
		label: append(g.label, _given, description),
		test:  append(g.test, setup),
	}
}

func (w when) And() GivenStatement {
	return given{
		label: append(w.label, _and),
		test:  w.test,
	}
}

func (w when) When(description string, setup func(t *testing.T)) ThenStatement {
	return then{
		label: append(w.label, _when, description),
		test:  append(w.test, setup),
	}
}

func (w when) Or(tests ...Runner) Runner {
	return or(w.label, w.test, tests)
}

func (t then) Or(tests ...Runner) Runner {
	return or(t.label, t.test, tests)
}

func (t then) And() WhenStatement {
	return when{
		label: append(t.label, _and),
		test:  t.test,
	}
}

func (t then) Then(description string, execution func(t *testing.T)) Runner {
	return runner{
		label: append(t.label, _then, description),
		test:  append(t.test, execution),
	}
}

func (r multiRunner) Run(t *testing.T, repeats ...int) bool {
	result := true
	for _, runner := range r.runners {
		result = result && runner.Run(t, repeats...)
	}
	return result
}

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
