package aws

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadCurrentContext(t *testing.T) {
	t.Run("returns empty string when file does not exist", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("HOME", tmpDir)

		got := LoadCurrentContext()

		if got != "" {
			t.Errorf("expected empty string, got %q", got)
		}
	})

	t.Run("returns profile name from file", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("HOME", tmpDir)

		os.WriteFile(filepath.Join(tmpDir, ".awsctx"), []byte("test-profile"), 0644)

		got := LoadCurrentContext()

		if got != "test-profile" {
			t.Errorf("expected 'test-profile', got %q", got)
		}
	})

	t.Run("trims whitespace from file content", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("HOME", tmpDir)

		os.WriteFile(filepath.Join(tmpDir, ".awsctx"), []byte("  test-profile\n"), 0644)

		got := LoadCurrentContext()

		if got != "test-profile" {
			t.Errorf("expected 'test-profile', got %q", got)
		}
	})
}

func TestSaveCurrentContext(t *testing.T) {
	t.Run("saves profile name to file", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("HOME", tmpDir)

		err := SaveCurrentContext("test-profile")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := os.ReadFile(filepath.Join(tmpDir, ".awsctx"))
		if err != nil {
			t.Fatalf("failed to read context file: %v", err)
		}

		if string(data) != "test-profile" {
			t.Errorf("expected 'test-profile', got %q", string(data))
		}
	})
	t.Run("overwrites existing context", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("HOME", tmpDir)

		err := SaveCurrentContext("old-profile")
		if err != nil {
			t.Fatalf("unexpected error saving old-profile: %v", err)
		}

		err = SaveCurrentContext("new-profile")
		if err != nil {
			t.Fatalf("unexpected error saving new-profile: %v", err)
		}

		data, err := os.ReadFile(filepath.Join(tmpDir, ".awsctx"))
		if err != nil {
			t.Fatalf("failed to read context file: %v", err)
		}

		if string(data) != "new-profile" {
			t.Errorf("expected 'new-profile', got %q", string(data))
		}
	})
}
