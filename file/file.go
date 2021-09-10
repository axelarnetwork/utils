package file

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func Resolve(path string) (*string, error) {
	resolved, err := resolveHome(path)
	if err != nil {
		return nil, err
	}
	resolved, err = resolveCurrentDirectory(*resolved)
	if err != nil {
		return nil, err
	}
	return resolved, nil
}

func resolveCurrentDirectory(path string) (*string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return replacePathPrefix(path, ".", wd), nil
}

func resolveHome(path string) (*string, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	return replacePathPrefix(path, "~", usr.HomeDir), nil
}

func replacePathPrefix(path string, prefix string, replacement string) *string {
	if path == prefix {
		path = replacement
	} else if strings.HasPrefix(path, fmt.Sprintf("%v/", prefix)) {
		path = filepath.Join(replacement, path[2:])
	}
	return &path
}
