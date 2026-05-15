package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type Settings struct {
	Port          string `json:"port"`
	Notifications bool   `json:"notifications"`
	Theme         string `json:"theme"`
	AutoCopyCodes bool   `json:"autoCopyCodes"`
}

var (
	mu       sync.Mutex
	filePath string
)

func init() {
	home, _ := os.UserHomeDir()
	filePath = filepath.Join(home, ".mailtap", "settings.json")
}

func Default() Settings {
	return Settings{
		Port:          "2525",
		Notifications: true,
		Theme:         "system",
		AutoCopyCodes: true,
	}
}

func Load() Settings {
	mu.Lock()
	defer mu.Unlock()

	s := Default()
	data, err := os.ReadFile(filePath)
	if err != nil {
		return s
	}
	json.Unmarshal(data, &s)
	return s
}

func Save(s Settings) error {
	mu.Lock()
	defer mu.Unlock()

	dir := filepath.Dir(filePath)
	os.MkdirAll(dir, 0755)

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}
