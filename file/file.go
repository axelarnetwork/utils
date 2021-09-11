package file

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"
)

// Resolve resolves relative paths to absolute paths, including ~/ shortcuts
func Resolve(path string) (string, error) {
	path, err := resolveHome(path)
	if err != nil {
		return "", err
	}
	path, err = resolveCurrentDirectory(path)
	if err != nil {
		return "", err
	}
	return path, nil
}

func resolveCurrentDirectory(path string) (string, error) {
	return filepath.Abs(path)
}

func resolveHome(path string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return replacePathPrefix(path, "~", usr.HomeDir), nil
}

func replacePathPrefix(path string, prefix string, replacement string) string {
	if path == prefix {
		path = replacement
	} else if strings.HasPrefix(path, fmt.Sprintf("%v/", prefix)) {
		path = filepath.Join(replacement, path[2:])
	}
	return path
}
