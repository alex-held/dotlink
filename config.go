package main

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// FileLink represents a single file link configuration
type FileLink struct {
	File string `yaml:"file"`
	To   string `yaml:"to"`
}

// LinkConfig represents the link section of the config
type LinkConfig struct {
	To string `yaml:"to"`
}

// DotlinkConfig represents the entire dotlink.yaml configuration
type DotlinkConfig struct {
	Link  LinkConfig `yaml:"link"`
	Files []FileLink `yaml:"files"`
}

// LoadConfig reads and parses a dotlink.yaml file
func LoadConfig(filepath string) (*DotlinkConfig, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	config := &DotlinkConfig{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// GetConfigDir returns the directory containing the config file
func GetConfigDir(configPath string) string {
	return filepath.Dir(configPath)
}

// ExpandPath expands environment variables in a path
func ExpandPath(path string) (string, error) {
	expanded := os.ExpandEnv(path)
	abs, err := filepath.Abs(expanded)
	if err != nil {
		return "", err
	}
	return abs, nil
}

