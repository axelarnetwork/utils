package file

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"testing"
)

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
	err := ioutil.WriteFile(resolve(filepath), contents, perm)
	if err != nil {
		panic(err)
	}
}

func readFile(filepath string) []byte {
	contents, err := ioutil.ReadFile(resolve(filepath))
	if err != nil {
		panic(err)
	}
	return contents
}

func resolve(path string) string {
	path = resolveHome(path)
	path = resolveCurrentDirectory(path)
	return path
}

func resolveCurrentDirectory(path string) string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return replacePathPrefix(path, ".", wd)
}

func resolveHome(path string) string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return replacePathPrefix(path, "~", usr.HomeDir)
}

func replacePathPrefix(path string, prefix string, replacement string) string {
	if path == prefix {
		path = replacement
	} else if strings.HasPrefix(path, fmt.Sprintf("%v/", prefix)) {
		path = filepath.Join(replacement, path[2:])
	}
	return path
}
