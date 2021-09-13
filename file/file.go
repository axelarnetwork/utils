package file

import (
	"os/user"
	"path/filepath"
)

// GetAxelarHome returns the absolute path to the current user's working directory for Axelar
func GetAxelarHome() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	home := u.HomeDir
	return filepath.Join(home, ".axelar"), nil
}
