package translator

import (
	"embed"
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.json
var localesFS embed.FS

type Translator struct {
	bundle *i18n.Bundle
}

func New() *Translator {
	bundle := i18n.NewBundle(language.English) // english is default (fallback) language
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	_, err := bundle.LoadMessageFileFS(localesFS, "locales/en.json")
	if err != nil {
		panic("failed to load en.json locale: " + err.Error())
	}

	_, err = bundle.LoadMessageFileFS(localesFS, "locales/ru.json")
	if err != nil {
		panic("failed to load ru.json locale: " + err.Error())
	}

	_, err = bundle.LoadMessageFileFS(localesFS, "locales/kk.json")
	if err != nil {
		panic("failed to load kk.json locale: " + err.Error())
	}

	_, err = bundle.LoadMessageFileFS(localesFS, "locales/zh.json")
	if err != nil {
		panic("failed to load zh.json locale: " + err.Error())
	}

	return &Translator{
		bundle: bundle,
	}
}

func (t *Translator) T(lang, messageID string, data map[string]any) string {
	localizer := i18n.NewLocalizer(t.bundle, lang)

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})

	if err != nil {
		return messageID
	}

	return msg
}

func (t *Translator) TFallback(lang, messageID, fallback string, data map[string]any) string {
	localizer := i18n.NewLocalizer(t.bundle, lang)

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})

	if err != nil {
		return fallback // if there is no translation, return fallback
	}

	return msg
}
