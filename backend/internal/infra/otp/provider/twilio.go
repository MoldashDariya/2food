package provider

import (
	"context"
	"fmt"

	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
	"gitlab.com/2foood/backend/internal/pkg/ctxkey"
)

type TwilioConfig struct {
	AccountSID string `envconfig:"TWILIO_ACCOUNT_SID" required:"true"`
	AuthToken  string `envconfig:"TWILIO_AUTH_TOKEN" required:"true"`
	FromNumber string `envconfig:"TWILIO_FROM_NUMBER" required:"true"`
}

type Twilio struct {
	client     *twilio.RestClient
	fromNumber string
	tr         Translator
}

func NewTwilio(cfg TwilioConfig, tr Translator) *Twilio {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.AccountSID,
		Password: cfg.AuthToken,
	})

	return &Twilio{
		client:     client,
		fromNumber: cfg.FromNumber,
		tr:         tr,
	}
}

func (p *Twilio) Send(ctx context.Context, to string, code string) error {
	params := &api.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(p.fromNumber)

	// extract language
	lang := ctxkey.GetLang(ctx)

	// prepare message
	fallback := fmt.Sprintf("Your verification code is: %s", code)
	text := p.tr.TFallback(lang, "SMS_OTP_Text", fallback, map[string]any{"code": code})
	params.SetBody(text)

	_, err := p.client.Api.CreateMessage(params)
	return err
}
