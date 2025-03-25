package main

import (
	"crypto/tls"
	"fmt"

	"github.com/wizzldev/mailer/types"
	"gopkg.in/gomail.v2"
)

// MessageInit represents the initial data for sending an email.
type MessageInit struct {
	ToAddress string
	Subject   string
	Config    map[string][]string
}

// Service is an interface that defines the methods for sending emails.
type Service interface {
	// Send sends an email.
	Send(init MessageInit, content string) error
	// SendTemplate sends an email using a template.
	SendTemplate(init MessageInit, templateComponent string, data types.TemplateProps) error
	// SendText sends an email using plain text.
	SendText(init MessageInit, content string) error
	// SendHTML sends an email using HTML content.
	SendHTML(init MessageInit, content string) error
}

// Config is the MailService configuration.
type Config struct {
	// From is the sender's name, not the email address. Sender<email>
	From string
	// FromAddress is the sender's email address.
	FromAddress string

	// SMTP configuration
	SMTPHost string
	SMTPPort int
	SMTPUser string
	SMTPPass string
}

// MailService is a service for sending emails.
type MailService struct {
	// cfg is the MailService configuration.
	cfg *Config

	// tmpl is the MailService template.
	tmpl Template
}

// NewMailService creates a new MailService instance.
func NewMailService(cfg *Config, tmpl Template) Service {
	return &MailService{
		cfg:  cfg,
		tmpl: tmpl,
	}
}

// SendTemplate sends an email using a template.
func (ms *MailService) SendTemplate(init MessageInit, templateComponent string, data types.TemplateProps) error {
	if init.Config == nil {
		init.Config = make(map[string][]string)
	}

	init.Config["Content-Type"] = []string{"text/html"}

	content, err := ms.tmpl.ParseComponent(templateComponent, data)
	if err != nil {
		return err
	}

	return ms.Send(init, content)
}

func (ms *MailService) SendText(init MessageInit, content string) error {
	return ms.Send(init, content)
}

func (ms *MailService) SendHTML(init MessageInit, content string) error {
	init.Config["Content-Type"] = []string{"text/html"}
	return ms.Send(init, content)
}

func (ms *MailService) Send(init MessageInit, content string) error {
	var contentType = "text/plain"
	var from = ms.cfg.FromAddress

	// combine sender's name and email address
	if ms.cfg.From != "" {
		from = fmt.Sprintf("%s<%s>", ms.cfg.From, ms.cfg.FromAddress)
	}

	// set content type
	if typ, ok := init.Config["Content-Type"]; ok && len(typ) == 1 && typ[0] != "" {
		contentType = typ[0]
	}

	// configure email message
	message := gomail.NewMessage()
	message.SetHeader("From", from)
	message.SetHeader("To", init.ToAddress)
	message.SetHeader("Subject", init.Subject)
	message.SetBody(contentType, content)

	// creating a dialer
	d := gomail.NewDialer(
		ms.cfg.SMTPHost,
		ms.cfg.SMTPPort,
		ms.cfg.SMTPUser,
		ms.cfg.SMTPPass,
	)
	// using tls config
	d.TLSConfig = &tls.Config{
		ServerName: ms.cfg.SMTPHost,
	}

	return d.DialAndSend(message)
}
