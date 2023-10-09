package main

import "context"

type Notifier interface {
	Type() string
	FormatText(string) string
	Send(ctx context.Context, msg Message) error
	Close() error
}
