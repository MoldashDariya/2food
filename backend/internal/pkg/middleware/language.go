package middleware

import (
	"net/http"

	"gitlab.com/2foood/backend/internal/pkg/ctxkey"
	"golang.org/x/text/language"
)

func (m *Manager) Language(next http.Handler) http.Handler {
	matcher := language.NewMatcher([]language.Tag{
		language.English,
		language.Russian,
		language.Kazakh,
		language.Ukrainian,
		language.Chinese,
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		langHeader := r.Header.Get("Accept-Language")
		langStr := ctxkey.DefaultLanguage // default language

		if langHeader != "" {
			tags, _, err := language.ParseAcceptLanguage(langHeader)
			if err == nil {
				tag, _, _ := matcher.Match(tags...) // trying to match languages

				base, _ := tag.Base() // returns the base language without regional suffixes
				langStr = base.String()
			}
		}

		ctx := ctxkey.WithLang(r.Context(), langStr)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
