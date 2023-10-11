package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type App struct {
	config    Config
	notifiers []Notifier
	filters   []*Filter
	watchers  []*Watcher
}

func NewApp(cfg Config) *App {
	return &App{config: cfg}
}

func (app *App) BuildNotifiers() *App {
	for _, cfg := range app.config.Notifications {
		notifier, err := NewNotifier(cfg)
		if err != nil {
			log.Fatalf("[ERROR] New notifier '%s' error: %v", cfg.Name, err)
		}
		app.notifiers = append(app.notifiers, notifier)
	}
	return app
}

func (app *App) BuildFilters() *App {
	for _, filterCfg := range app.config.Filters {
		filter, err := NewFilter(filterCfg, app.config.Hostname, app.notifiers)
		if err != nil {
			log.Fatalf("[ERROR] NewFilter error: %v", err)
		}
		app.filters = append(app.filters, filter)
	}
	return app
}

func (app *App) BuildWatchers() *App {
	for _, fileCfg := range app.config.Files {
		watcher, err := NewWatcher(fileCfg, app.filters)
		if err != nil {
			log.Fatalf("[ERROR] NewWatcher error: %v", err)
		}
		app.watchers = append(app.watchers, watcher)
	}
	return app
}

func (app *App) Watch() {
	log.Printf("[INFO] LogAlert is running")

	// Handle ctrl+c/ctrl+x interrupt
	var stopChan = make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTSTP)

	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(len(app.watchers))

	for _, watcher := range app.watchers {
		go watcher.watch(ctx, &wg)
	}

	interrupt := <-stopChan
	cancel()
	wg.Wait()

	log.Printf("[INFO] LogAlert is shutting down due to %+v", interrupt)
}
