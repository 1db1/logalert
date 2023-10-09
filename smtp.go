package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"sync"
)

type EmailConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
	To       string `yaml:"to"`
}

type SmtpNotifier struct {
	sync.Mutex
	hostname string
	auth     smtp.Auth
	from     string
	to       string
	client   *smtp.Client
}

func NewSmtpNotifier(cfg EmailConfig) (*SmtpNotifier, error) {
	hostname := cfg.Host + ":" + cfg.Port

	client, err := smtp.Dial(hostname)
	if err != nil {
		return nil, err
	}

	err = client.Hello("localhost")
	if err != nil {
		return nil, err
	}

	if ok, _ := client.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: cfg.Host}
		if err = client.StartTLS(config); err != nil {
			return nil, err
		}
	}

	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)

	if ok, _ := client.Extension("AUTH"); ok {
		if err = client.Auth(auth); err != nil {
			return nil, err
		}
	}

	return &SmtpNotifier{sync.Mutex{}, hostname, auth, cfg.From, cfg.To, client}, nil
}

func (sn *SmtpNotifier) Send(ctx context.Context, msg string) error {
	sn.Lock()
	defer sn.Unlock()

	err := sn.client.Mail(sn.from)
	if err != nil {
		return fmt.Errorf("SMTP client mail error: %v", err)
	}

	err = sn.client.Rcpt(sn.to)
	if err != nil {
		return fmt.Errorf("SMTP client rcpt error: %v", err)
	}

	w, err := sn.client.Data()
	if err != nil {
		return fmt.Errorf("SMTP client data error: %v", err)
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("SMTP client write error: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("SMTP client write close error: %v", err)
	}

	return nil
}

func (sn *SmtpNotifier) From() string {
	return sn.from
}

func (sn *SmtpNotifier) To() string {
	return sn.to
}

func (sn *SmtpNotifier) Type() string {
	return "email"
}

func (sn *SmtpNotifier) Close() error {
	return sn.client.Close()
}

func NewSmtpMessage(from, to, subject, body string) string {
	return fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\r\n", from, to, subject, body)
}
