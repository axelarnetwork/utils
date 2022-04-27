// Package testutils provides general purpose utility functions for unit/integration testing.
package testutils

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/axelarnetwork/utils/slices"
)

// Func wraps a regular testing function so it can be used as a pointer function receiver
type Func func(t *testing.T)

// Repeat executes the testing function n times
func (f Func) Repeat(n int) Func {
	return func(t *testing.T) {
		for i := 0; i < n; i++ {
			f(t)
		}
	}
}

// Run is equivalent to calling f(t), it just provides an interface that reads better
func (f Func) Run(t *testing.T) {
	f(t)
}

// TestCases define alternative test setups that should all be tested
type TestCases[T any] []T

// AsTestCases is defined for convenience of casting a slice to TestCases
func AsTestCases[T any](cases ...T) TestCases[T] {
	return cases
}

// ForEach defines the test that should be run on each test case
func (tc TestCases[T]) ForEach(f func(t *testing.T, testCase T)) Func {
	return func(t *testing.T) {
		for _, testCase := range tc {
			f(t, testCase)
		}
	}
}

func (tc TestCases[T]) Map(f func(testCase T) Runner) Runner {
	runners := slices.Map(tc, f)
	var out ThenStatements

	for _, runner := range runners {
		switch runner := runner.(type) {
		case ThenStatement:
			out = append(out, runner)
		case ThenStatements:
			for _, then := range runner {
				out = append(out, then)
			}
		}
	}
	return out
}

// Events wraps sdk.Events
type Events []abci.Event

// Filter returns a collection of events filtered by the predicate
func (fe Events) Filter(predicate func(events abci.Event) bool) Events {
	var filtered Events
	for _, event := range fe {
		if predicate(event) {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

// ErrorCache is a struct that can be used to get at the error that is emitted by test assertions when passing it instead ot *testing.T
type ErrorCache struct {
	Error error
}

// Errorf records the given formatted string as an erro
func (ec *ErrorCache) Errorf(format string, args ...interface{}) {
	ec.Error = fmt.Errorf(format, args...)
}

// SetEnv safely sets an OS env var to the specified value and resets it to the original value upon test closure
func SetEnv(t *testing.T, key string, val string) {
	// TODO : enable with Go 1.17 >> it will automatically handle Cleanup
	// t.Setenv(key, val)
	orig := os.Getenv(key)
	os.Setenv(key, val)
	t.Cleanup(func() { os.Setenv(key, orig) })
}

// SetFileContents safely sets the contents of the specified filepath to the specified value, then resets it to the original value upon test closure
func SetFileContents(t *testing.T, filepath string, contents string) {
	SetFileContentsAndPermissions(t, filepath, contents, 0644)
}

// SetFileContentsAndPermissions safely sets the contents of the specified filepath to the specified value and FileMode, then resets it to the original value upon test closure
func SetFileContentsAndPermissions(t *testing.T, path string, contents string, perm fs.FileMode) {
	os.MkdirAll(filepath.Dir(path), os.ModePerm)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			panic(err)
		}
		defer file.Close()
	}

	original := readFile(path)
	writeFile(path, []byte(contents), perm)
	// this may be overkill to reset the file to its original value
	t.Cleanup(
		func() {
			writeFile(path, original, perm)
		},
	)
}

func writeFile(filepath string, contents []byte, perm fs.FileMode) {
	err := ioutil.WriteFile(filepath, contents, perm)
	if err != nil {
		panic(err)
	}
}

func readFile(filepath string) []byte {
	contents, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	return contents
}
