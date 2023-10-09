package main

import (
	"context"
	"testing"
)

func TestSmtp(t *testing.T) {
	t.Skip()

	cfgPath := "./config.yml"

	globalCfg, err := NewConfig(cfgPath)
	if err != nil {
		t.Error(err)
	}

	cfg := globalCfg.Notifications[0].EmailConfig

	smtpNotifier, err := NewSmtpNotifier(cfg)
	if err != nil {
		t.Fatal(err)
	}

	msg := Message{
		Text: "Logalert SMTP notifier test message",
	}

	err = smtpNotifier.Send(context.Background(), msg)
	if err != nil {
		t.Error(err)
	}
}
