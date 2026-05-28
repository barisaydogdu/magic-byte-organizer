package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func LoadConfig(configPath string) (map[string]string, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Println("cannot read config", err)
		return nil, err
	}

	config := make(map[string]string)

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Println("cannot parse config", err)
		return nil, err
	}

	return config, nil
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

func GetAbsoluteConfigPath() string {
	configFileName := "config.json"

	var err error
	//development envoirment
	if _, err = os.Stat(configFileName); err == nil {
		absPath, _ := filepath.Abs(configFileName)
		return absPath
	}

	exePath, err := os.Executable()
	if err == nil {
		binaryConfigPath := filepath.Join(filepath.Dir(exePath), configFileName)

		if _, err = os.Stat(binaryConfigPath); err == nil {
			return binaryConfigPath
		}
	}

	return configFileName
}
