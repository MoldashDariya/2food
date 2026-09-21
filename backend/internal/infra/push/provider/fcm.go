package provider

import (
	"context"
	"encoding/base64"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type FCMConfig struct {
	CredentialsB64 string `envconfig:"FCM_CREDENTIALS_B64" required:"true"`
}

type FCM struct {
	client *messaging.Client
}

func NewFCM(cfg FCMConfig) (*FCM, error) {
	// decode credentials from base64
	decoded, err := base64.StdEncoding.DecodeString(cfg.CredentialsB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode credentials: %w", err)
	}

	// create firebase app with credentials
	app, err := firebase.NewApp(context.Background(), nil,
		option.WithAuthCredentialsJSON(option.ServiceAccount, decoded),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create firebase app: %w", err)
	}

	// create messaging client
	client, err := app.Messaging(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to create messaging client: %w", err)
	}

	return &FCM{client: client}, nil
}

func (f *FCM) Send(ctx context.Context, token string, title, body string, data map[string]any) error {
	// Adapting the map to FCM requirements
	fcmData := make(map[string]string, len(data))
	for k, v := range data {
		fcmData[k] = fmt.Sprint(v)
	}
	_, err := f.client.Send(ctx, &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: fcmData,
	})

	return err
}

func (f *FCM) SendMulti(ctx context.Context, tokens []string, title, body string, data map[string]any) ([]string, error) {
	// Adapting the map to FCM requirements
	fcmData := make(map[string]string, len(data))
	for k, v := range data {
		fcmData[k] = fmt.Sprint(v)
	}
	resp, err := f.client.SendEachForMulticast(ctx, &messaging.MulticastMessage{
		Tokens: tokens,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: fcmData,
	})
	if err != nil {
		return nil, err
	}

	// collect stale tokens
	var stale []string
	for i, r := range resp.Responses {
		if r.Error != nil {
			if messaging.IsUnregistered(r.Error) || messaging.IsInvalidArgument(r.Error) {
				stale = append(stale, tokens[i])
			}
		}
	}

	return stale, nil
}
