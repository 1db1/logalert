package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"sync"
)

const NotifierTypeMail = "mail"

type MailConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
	To       string `yaml:"to"`
}

type SmtpNotifier struct {
	sync.Mutex
	name     string
	hostname string
	auth     smtp.Auth
	from     string
	to       string
	client   *smtp.Client
}

func NewSmtpNotifier(cfg NotificationConfig) (*SmtpNotifier, error) {
	hostname := cfg.MailConfig.Host + ":" + cfg.MailConfig.Port

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

	auth := smtp.PlainAuth("", cfg.MailConfig.Username, cfg.MailConfig.Password, cfg.MailConfig.Host)

	if ok, _ := client.Extension("AUTH"); ok {
		if err = client.Auth(auth); err != nil {
			return nil, err
		}
	}

	return &SmtpNotifier{
		sync.Mutex{},
		cfg.Name,
		hostname,
		auth,
		cfg.MailConfig.From,
		cfg.MailConfig.To,
		client,
	}, nil
}

func (sn *SmtpNotifier) Send(ctx context.Context, msg Message) error {
	sn.Lock()
	defer sn.Unlock()

	msg.BuildSubject()
	msg.BuildText()

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

	_, err = w.Write(newSmtpMessage(sn.From(), sn.To(), msg.Subject, msg.Text))
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

func (sn *SmtpNotifier) Name() string {
	return sn.name
}

func (sn *SmtpNotifier) Type() string {
	return NotifierTypeMail
}

func (sn *SmtpNotifier) FormatText(text string) string {
	return text
}

func (sn *SmtpNotifier) Close() error {
	return sn.client.Close()
}

func newSmtpMessage(from, to, subject, body string) []byte {
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\r\n", from, to, subject, body)
	return []byte(msg)
}
