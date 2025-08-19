package pkg

import (
	"os"
	"path/filepath"
)

func GetPathRelativeToExe(path string) (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	exeDir := filepath.Dir(exePath)
	return filepath.Join(exeDir, path), nil
}
