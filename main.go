package main

import (
	"flag"
	"log"
)

func parseFlags() (string, error) {
	var configPath string

	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")
	flag.Parse()

	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}

	return configPath, nil
}

func main() {
	cfgPath, err := parseFlags()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := NewConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	NewApp(cfg).
		BuildNotifiers().
		BuildFilters().
		BuildWatchers().
		Watch()
}
