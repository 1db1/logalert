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
	for _, notif := range msg.Match.Notifications {
		notifier, ok := s.notifiers[notif]
		if !ok {
			return fmt.Errorf("Unknown notifier %s", notif)
		}

		msg.Text = notifier.FormatText(msg.Text)
		msg.BuildSubject()
		msg.BuildText()

		if err := notifier.Send(ctx, msg); err != nil {
			return fmt.Errorf("%s message send error: %v msg: %s", notifier.Type(), err, msg.Text)
		}
	}

	return nil
}

func (s *Sender) Close() {
	for _, n := range s.notifiers {
		n.Close()
	}
}
