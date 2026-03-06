package aws

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type Profile struct {
	Name      string
	Region    string
	AccountID string
}

func LoadProfiles() ([]Profile, error) {
	configPath := filepath.Join(os.Getenv("HOME"), ".aws", "config")

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var profiles []Profile
	var current *Profile

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			if current != nil {
				profiles = append(profiles, *current)
			}

			section := line[1 : len(line)-1]
			if strings.HasPrefix(section, "profile ") {
				name := strings.TrimPrefix(section, "profile ")
				current = &Profile{Name: name}
			} else {
				current = nil
			}
			continue
		}

		if current != nil && strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			if key == "region" {
				current.Region = value
			}
		}
	}

	if current != nil {
		profiles = append(profiles, *current)
	}

	return profiles, scanner.Err()
}

func FetchAccountID(profileName string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profileName),
	)
	if err != nil {
		return "", err
	}

	client := sts.NewFromConfig(cfg)
	result, err := client.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}

	return *result.Account, nil
}
