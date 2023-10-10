package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type FilterConfig struct {
	Name          string   `yaml:"name"`
	Pattern       string   `yaml:"pattern"`
	Exceptions    []string `yaml:"exceptions"`
	Message       string   `yaml:"message"`
	Subject       string   `yaml:"subject"`
	Notifications []string `yaml:"notifications"`
}

type LogFileConfig struct {
	Path           string         `yaml:"path"`
	DateFormat     string         `yaml:"dateFormat"`
	ReadBufferSize string         `yaml:"readBufferSize"`
	IntervalSec    uint           `yaml:"interval"`
	Filters        []FilterConfig `yaml:"filters"`
}

type NotificationConfig struct {
	Name           string `yaml:"name"`
	Type           string `yaml:"type"`
	EmailConfig    `yaml:",inline"`
	TelegramConfig `yaml:",inline"`
}

type Config struct {
	Hostname      string               `yaml:"hostname"`
	Notifications []NotificationConfig `yaml:"notifications"`
	LogFiles      []LogFileConfig      `yaml:"logFiles"`
}

func NewConfig(configPath string) (Config, error) {
	config := Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return config, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return config, fmt.Errorf("YAML config decode error: %v", err)
	}

	if config.Hostname == "" {
		config.Hostname, _ = os.Hostname()
	}

	return config, nil
}

func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}

	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}

	return nil
}
