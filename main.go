package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

	run(cfg)
}

func run(cfg Config) {
	log.Printf("[INFO] LogAlert is running")

	sender := NewSender()

	for _, notifCfg := range cfg.Notifications {
		sender.RegisterNotifier(notifCfg)
	}

	// Handle ctrl+c/ctrl+x interrupt
	var stopChan = make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTSTP)

	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(len(cfg.LogFiles))

	for _, logFileCfg := range cfg.LogFiles {
		watcher := NewLogWatcher(cfg.Hostname, logFileCfg, sender)
		go watcher.watch(ctx, &wg)
	}

	interrupt := <-stopChan
	cancel()
	wg.Wait()

	sender.Close()

	log.Printf("LogAlert is shutting down due to %+v", interrupt)
}
