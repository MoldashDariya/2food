package locale

import (
	"strings"

	"gitlab.com/2foood/backend/internal/domain"
)

func FromAcceptLanguage(acceptLanguage string) domain.Language {
	switch {
	case strings.Contains(acceptLanguage, "kk"):
		return domain.LanguageKazakh
	case strings.Contains(acceptLanguage, "en"):
		return domain.LanguageEnglish
	case strings.Contains(acceptLanguage, "ru"):
		return domain.LanguageRussian
	case strings.Contains(acceptLanguage, "uk"):
		return domain.LanguageUkrainian
	case strings.Contains(acceptLanguage, "zh"):
		return domain.LanguageChinese
	default:
		return domain.LanguageEnglish
	}
}
