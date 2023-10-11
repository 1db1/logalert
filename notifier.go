package main

import (
	"context"
	"fmt"
)

type Notifier interface {
	Name() string
	Type() string
	FormatText(string) string
	Send(ctx context.Context, msg Message) error
	Close() error
}

func NewNotifier(cfg NotificationConfig) (Notifier, error) {
	var (
		notifier Notifier
		err      error
	)

	switch cfg.Type {
	case NotifierTypeMail:
		notifier, err = NewSmtpNotifier(cfg)
		if err != nil {
			return nil, err
		}
	case NotifierTypeTelegram:
		notifier, err = NewTelegramNotifier(cfg)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Notifier type '%s' is unsupported", cfg.Type)
	}

	return notifier, nil
}
