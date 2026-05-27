package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func LoadConfig(configPath string) map[string]string {
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	config := make(map[string]string)

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

func ExpandHomePath(path string) string {
	if strings.HasPrefix(path, "~") {
		dirName, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(dirName, path[1:])
		}
	}
	return path
}
