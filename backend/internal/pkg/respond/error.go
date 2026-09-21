package respond

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"gitlab.com/2foood/backend/internal/domain"
	"gitlab.com/2foood/backend/internal/pkg/validation"
)

type ErrorResponse struct {
	Code    string            `json:"code"`              // code for frontend
	Message string            `json:"message"`           // text for logs/debug
	Details map[string]string `json:"details,omitempty"` // for validation errors
}

func Error(w http.ResponseWriter, r *http.Request, log *slog.Logger, err error) {
	// handle validation errors
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		log.Warn("validation error",
			slog.String("path", r.URL.Path),
			slog.Int("status", http.StatusBadRequest),
			slog.String("err", err.Error()),
		)

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{
			Code:    "ERR_VALIDATION_FAILED",
			Message: "validation error",
			Details: validation.ToMap(validationErrors),
		})
		return
	}

	// handle domain errors
	statusCode := http.StatusInternalServerError
	msg := domain.ErrInternalServerError.Error()
	errCode := "ERR_INTERNAL_SERVER_ERROR"

	switch {
	// 400
	case errors.Is(err, domain.ErrInvalidJSON):
		statusCode = http.StatusBadRequest
		msg = domain.ErrInvalidJSON.Error()
		errCode = "ERR_INVALID_JSON"
	case errors.Is(err, domain.ErrCodeExpired):
		statusCode = http.StatusBadRequest
		msg = domain.ErrCodeExpired.Error()
		errCode = "ERR_CODE_EXPIRED"
	case errors.Is(err, domain.ErrCodeInvalid):
		statusCode = http.StatusBadRequest
		msg = domain.ErrCodeInvalid.Error()
		errCode = "ERR_CODE_INVALID"
	case errors.Is(err, domain.ErrInvalidOAuthToken):
		statusCode = http.StatusBadRequest
		msg = domain.ErrInvalidOAuthToken.Error()
		errCode = "ERR_INVALID_OAUTH_TOKEN"
	case errors.Is(err, domain.ErrInvalidParam):
		statusCode = http.StatusBadRequest
		msg = domain.ErrInvalidParam.Error()
		errCode = "ERR_INVALID_PARAM"
	case errors.Is(err, domain.ErrPickupWindowMismatch):
		statusCode = http.StatusBadRequest
		msg = domain.ErrPickupWindowMismatch.Error()
		errCode = "ERR_PICKUP_WINDOW_MISMATCH"
	case errors.Is(err, domain.ErrInsufficientStock):
		statusCode = http.StatusBadRequest
		msg = domain.ErrInsufficientStock.Error()
		errCode = "ERR_INSUFFICIENT_STOCK"
	case errors.Is(err, domain.ErrProductUnavailable):
		statusCode = http.StatusBadRequest
		msg = domain.ErrProductUnavailable.Error()
		errCode = "ERR_PRODUCT_UNAVAILABLE"
	case errors.Is(err, domain.ErrMissingParam):
		statusCode = http.StatusBadRequest
		msg = domain.ErrMissingParam.Error()
		errCode = "ERR_MISSING_PARAM"

		// 401
	case errors.Is(err, domain.ErrInvalidToken):
		statusCode = http.StatusUnauthorized
		msg = domain.ErrInvalidToken.Error()
		errCode = "ERR_INVALID_TOKEN"
	case errors.Is(err, domain.ErrExpiredToken):
		statusCode = http.StatusUnauthorized
		msg = domain.ErrExpiredToken.Error()
		errCode = "ERR_EXPIRED_TOKEN"
	case errors.Is(err, domain.ErrInvalidAuthorizationHeader):
		statusCode = http.StatusUnauthorized
		msg = domain.ErrInvalidAuthorizationHeader.Error()
		errCode = "ERR_INVALID_AUTHORIZATION_HEADER"
	case errors.Is(err, domain.ErrSessionNotFound):
		statusCode = http.StatusUnauthorized
		msg = domain.ErrSessionNotFound.Error()
		errCode = "ERR_SESSION_NOT_FOUND"

		// 403
	case errors.Is(err, domain.ErrPermissionDenied):
		statusCode = http.StatusForbidden
		msg = domain.ErrPermissionDenied.Error()
		errCode = "ERR_PERMISSION_DENIED"
	case errors.Is(err, domain.ErrInvalidPickupCode):
		statusCode = http.StatusForbidden
		msg = domain.ErrInvalidPickupCode.Error()
		errCode = "ERR_INVALID_PICKUP_CODE"

		// 404
	case errors.Is(err, domain.ErrNotFound):
		statusCode = http.StatusNotFound
		msg = domain.ErrNotFound.Error()
		errCode = "ERR_NOT_FOUND"

		// 406
	case errors.Is(err, domain.ErrUnsupportedContentType):
		statusCode = http.StatusNotAcceptable
		msg = domain.ErrUnsupportedContentType.Error()
		errCode = "ERR_UNSUPPORTED_CONTENT_TYPE"

		// 409
	case errors.Is(err, domain.ErrAlreadyExists):
		statusCode = http.StatusConflict
		msg = domain.ErrAlreadyExists.Error()
		errCode = "ERR_ALREADY_EXISTS"
	case errors.Is(err, domain.ErrInvalidStatusTransition):
		statusCode = http.StatusConflict
		msg = domain.ErrInvalidStatusTransition.Error()
		errCode = "ERR_INVALID_STATUS_TRANSITION"

		// 429
	case errors.Is(err, domain.ErrTooManyAttempts):
		statusCode = http.StatusTooManyRequests
		msg = domain.ErrTooManyAttempts.Error()
		errCode = "ERR_TOO_MANY_ATTEMPTS"
	}

	// log error or warning
	if statusCode >= 500 {
		log.Error("server error",
			slog.String("path", r.URL.Path),
			slog.Int("status", statusCode),
			slog.String("err", err.Error()),
		)
	} else {
		log.Warn("client error",
			slog.String("path", r.URL.Path),
			slog.Int("status", statusCode),
			slog.String("err", err.Error()),
		)
	}

	// respond with error
	render.Status(r, statusCode)
	render.JSON(w, r, ErrorResponse{
		Code:    errCode,
		Message: msg,
	})
}
