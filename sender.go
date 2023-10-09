package main

import (
	"context"
	"fmt"
	"log"
)

const (
	NotifierTypeEmail    = "mail"
	NotifierTypeTelegram = "telegram"
)

type Sender struct {
	notifiers map[string]Notifier
}

func NewSender() *Sender {
	return &Sender{
		notifiers: make(map[string]Notifier),
	}
}

func (s *Sender) RegisterNotifier(notifCfg NotificationConfig) {
	var (
		notifier Notifier
		err      error
	)

	switch notifCfg.Type {
	case NotifierTypeEmail:
		notifier, err = NewSmtpNotifier(notifCfg.EmailConfig)
		if err != nil {
			log.Fatalf("[ERROR] New register email notifier '%s' error: %v", notifCfg.Name, err)
		}
	case NotifierTypeTelegram:
		notifier, err = NewTelegramNotifier(notifCfg.TelegramConfig)
		if err != nil {
			log.Fatalf("[ERROR] New register telegram notifier '%s' error: %v", notifCfg.Name, err)
		}
	default:
		log.Fatalf("[ERROR] Notifier type '%s' is unsupported", notifCfg.Type)
	}

	s.notifiers[notifCfg.Name] = notifier
}

func (s *Sender) SendMessage(ctx context.Context, msg Message) error {
	for _, notif := range msg.Notifications {
		notifier, ok := s.notifiers[notif]
		if !ok {
			return fmt.Errorf("Unknown notifier %s", notif)
		}

		switch notifier.Type() {
		case NotifierTypeEmail:
			notifier := notifier.(*SmtpNotifier)
			smtpMsg := NewSmtpMessage(notifier.From(), notifier.To(), msg.Subject, msg.Text)
			if err := notifier.Send(ctx, smtpMsg); err != nil {
				return fmt.Errorf("SMTP message send error: %v msg: %s", err, msg.Text)
			}
		case NotifierTypeTelegram:
			if err := notifier.Send(ctx, msg.Text); err != nil {
				return fmt.Errorf("Telegram message send error: %v msg: %s", err, msg.Text)
			}
		default:
			return fmt.Errorf("Unknown notifier %s", notifier.Type())
		}
	}

	return nil
}

func (s *Sender) Close() {
	for _, n := range s.notifiers {
		n.Close()
	}
}
