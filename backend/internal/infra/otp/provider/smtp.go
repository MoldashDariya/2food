package provider

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"net/mail"
	"net/smtp"
	"text/template"
	"time"

	"gitlab.com/2foood/backend/internal/pkg/ctxkey"
)

//go:embed templates/email.html
var emailTmpl string

type SMTPConfig struct {
	Host     string `envconfig:"SMTP_HOST" required:"true"`
	Port     string `envconfig:"SMTP_PORT" required:"true"`
	Username string `envconfig:"SMTP_USERNAME" required:"true"`
	Password string `envconfig:"SMTP_PASSWORD" required:"true"`
	From     string `envconfig:"SMTP_FROM" required:"true"`
}

// otp email data for template rendering
type otpEmailData struct {
	Code    string
	TTL     int
	Title   string
	Body    string
	Expires string
	Warning string
	Ignore  string
}

type SMTP struct {
	cfg  SMTPConfig
	auth smtp.Auth
	tmpl *template.Template
	tr   Translator
	ttl  time.Duration
}

func NewSMTP(cfg SMTPConfig, tr Translator, ttl time.Duration) (*SMTP, error) {
	tmpl, err := template.New("OTPEmail").Parse(emailTmpl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse otp email template: %w", err)
	}

	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	return &SMTP{cfg: cfg, auth: auth, tmpl: tmpl, tr: tr, ttl: ttl}, nil
}

func (s *SMTP) Send(ctx context.Context, to string, code string) error {
	from, err := mail.ParseAddress(s.cfg.From)
	if err != nil {
		return fmt.Errorf("failed to parse sender address: %w", err)
	}

	// extract language
	lang := ctxkey.GetLang(ctx)

	// get translated subject
	subject := s.tr.TFallback(lang, "Email_OTP_Subject", "Verification code", nil)

	// prepare data with fallbacks
	data := otpEmailData{
		Code:    code,
		TTL:     int(s.ttl.Minutes()),
		Title:   s.tr.TFallback(lang, "Email_OTP_Title", "Verify Your Action", nil),
		Body:    s.tr.TFallback(lang, "Email_OTP_Body", "To complete your request, please use the following one-time verification code.", nil),
		Expires: s.tr.TFallback(lang, "Email_OTP_Expires", "Code expires in", nil),
		Warning: s.tr.TFallback(lang, "Email_OTP_Warning", "Our employees will never ask for this code.", nil),
		Ignore:  s.tr.TFallback(lang, "Email_OTP_Ignore", "If you did not make this request, you can safely ignore this message.", nil),
	}

	var body bytes.Buffer
	if err := s.tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to render otp email: %w", err)
	}

	msg := s.buildMessage(from.String(), to, subject, body.String())
	addr := fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port)

	return smtp.SendMail(addr, s.auth, from.Address, []string{to}, msg)
}

func (s *SMTP) buildMessage(from, to, subject, htmlBody string) []byte {
	var b bytes.Buffer
	fmt.Fprintf(&b, "From: %s\r\n", from)
	fmt.Fprintf(&b, "To: %s\r\n", to)
	fmt.Fprintf(&b, "Subject: %s\r\n", subject)
	fmt.Fprintf(&b, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&b, "Content-Type: text/html; charset=\"UTF-8\"\r\n")
	fmt.Fprintf(&b, "\r\n")
	b.WriteString(htmlBody)
	return b.Bytes()
}
