package aws

import (
	"os"
	"path/filepath"
	"strings"
)

func contextFilePath() string {
	return filepath.Join(os.Getenv("HOME"), ".awsctx")
}

func LoadCurrentContext() string {
	data, err := os.ReadFile(contextFilePath())
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func SaveCurrentContext(profileName string) error {
	return os.WriteFile(contextFilePath(), []byte(profileName), 0644)
}
