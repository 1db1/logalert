package main

import "context"

type Notifier interface {
	Type() string
	Send(ctx context.Context, msg string) error
	Close() error
}
