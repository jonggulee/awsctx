package aws

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadProfiles(t *testing.T) {
	t.Run("returns error when config file does not exist", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("HOME", tmpDir)

		profiles, err := LoadProfiles()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if profiles != nil {
			t.Errorf("expected nil profiles, got %v", profiles)
		}
	})
}

func TestLoadProfilesParsing(t *testing.T) {
	tests := []struct {
		name       string
		configFile string
		wantCount  int
		wantFirst  Profile
	}{
		{
			name: "expected only default profile",
			configFile: `
[default]
region = us-east-1
`,
			wantCount: 1,
			wantFirst: Profile{
				Name:   "default",
				Region: "us-east-1",
			},
		},
		{
			name: "expected profile keyword",
			configFile: `
[profile dev]
region = us-west-2
`,
			wantCount: 1,
			wantFirst: Profile{
				Name:   "dev",
				Region: "us-west-2",
			},
		},
		{
			name: "expected ignore 주석",
			configFile: `
# this is a comment
[profile prod]
region = us-east-1
`,
			wantCount: 1,
			wantFirst: Profile{
				Name:   "prod",
				Region: "us-east-1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			t.Setenv("HOME", tmpDir)

			awsDir := filepath.Join(tmpDir, ".aws")
			os.MkdirAll(awsDir, 0755)
			os.WriteFile(filepath.Join(awsDir, "config"), []byte(tt.configFile), 0644)

			profiles, err := LoadProfiles()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(profiles) != tt.wantCount {
				t.Fatalf("expected %d profiles, got %d", tt.wantCount, len(profiles))
			}
			if profiles[0] != tt.wantFirst {
				t.Errorf("expected %+v, got %+v", tt.wantFirst, profiles[0])
			}
		})
	}
}
