package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const NotifierTypeTelegram = "telegram"

type TelegramConfig struct {
	Token  string `yaml:"token"`
	ChatID int64  `yaml:"chatID"`
}

type TelegramNotifier struct {
	name   string
	bot    *bot.Bot
	chatID int64
}

func NewTelegramNotifier(cfg NotificationConfig) (*TelegramNotifier, error) {
	opts := []bot.Option{
		bot.WithCheckInitTimeout(time.Second * 10),
	}

	b, err := bot.New(cfg.TelegramConfig.Token, opts...)
	if err != nil {
		return nil, err
	}

	return &TelegramNotifier{cfg.Name, b, cfg.TelegramConfig.ChatID}, nil
}

func (tn *TelegramNotifier) Send(ctx context.Context, msg Message) error {
	msg.BuildText()
	msg.Text = tn.FormatText(msg.Text)

	_, err := tn.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    tn.chatID,
		Text:      msg.Text,
		ParseMode: models.ParseModeMarkdown,
	})

	if err != nil {
		return fmt.Errorf("send error: %v", err)
	}

	return nil
}

func (tn TelegramNotifier) Name() string {
	return tn.name
}

func (tn TelegramNotifier) Type() string {
	return NotifierTypeTelegram
}

func (tn TelegramNotifier) FormatText(text string) string {
	return bot.EscapeMarkdown(text)
}

func (tn *TelegramNotifier) Close() error {
	_, err := tn.bot.Close(context.TODO())
	return err
}
