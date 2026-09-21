package provider

type Translator interface {
	TFallback(lang, messageID, fallback string, data map[string]any) string
}
