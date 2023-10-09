package main

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type TelegramConfig struct {
	Token  string `yaml:"token"`
	ChatID int64  `yaml:"chatID"`
}

type TelegramNotifier struct {
	bot    *bot.Bot
	chatID int64
}

func NewTelegramNotifier(cfg TelegramConfig) (*TelegramNotifier, error) {
	opts := []bot.Option{}

	b, err := bot.New(cfg.Token, opts...)
	if err != nil {
		return nil, err
	}

	return &TelegramNotifier{b, cfg.ChatID}, nil
}

func (tn *TelegramNotifier) Send(ctx context.Context, msg string) error {
	_, err := tn.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    tn.chatID,
		Text:      msg,
		ParseMode: models.ParseModeMarkdown,
	})

	if err != nil {
		return fmt.Errorf("send error: %v", err)
	}

	return nil
}

func (tn TelegramNotifier) Type() string {
	return "telegram"
}

func (tn *TelegramNotifier) Close() error {
	_, err := tn.bot.Close(context.TODO())
	return err
}