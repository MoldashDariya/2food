package decode

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gitlab.com/2foood/backend/internal/domain"
)

// UUID decodes a UUID from the request URL parameter with the given key
func UUID(r *http.Request, key string) (uuid.UUID, error) {
	val := chi.URLParam(r, key)
	if val == "" {
		return uuid.Nil, domain.ErrMissingParam
	}

	id, err := uuid.Parse(val)
	if err != nil {
		return uuid.Nil, domain.ErrInvalidParam
	}

	return id, nil
}

// Int decodes an integer from the request URL parameter with the given key
func Int(r *http.Request, key string) (int, error) {
	val := chi.URLParam(r, key)
	if val == "" {
		return 0, domain.ErrMissingParam
	}

	num, err := strconv.Atoi(val)
	if err != nil {
		return 0, domain.ErrInvalidParam
	}

	return num, nil
}
