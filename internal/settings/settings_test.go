package settings

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultAutoCopyCodesIsTrue(t *testing.T) {
	if !Default().AutoCopyCodes {
		t.Errorf("Default().AutoCopyCodes = false, want true")
	}
}

func TestLoadLegacyFilePreservesAutoCopyDefault(t *testing.T) {
	dir := t.TempDir()
	legacy := filepath.Join(dir, "settings.json")
	// A settings file written before this feature existed (no autoCopyCodes key).
	if err := os.WriteFile(legacy, []byte(`{"port":"2525","notifications":true,"theme":"dark"}`), 0644); err != nil {
		t.Fatalf("write legacy file: %v", err)
	}

	mu.Lock()
	orig := filePath
	filePath = legacy
	mu.Unlock()
	defer func() {
		mu.Lock()
		filePath = orig
		mu.Unlock()
	}()

	s := Load()
	if !s.AutoCopyCodes {
		t.Errorf("legacy load AutoCopyCodes = false, want true (missing key must keep default)")
	}
	if s.Theme != "dark" {
		t.Errorf("legacy load Theme = %q, want \"dark\" (sanity: file was read)", s.Theme)
	}
}
