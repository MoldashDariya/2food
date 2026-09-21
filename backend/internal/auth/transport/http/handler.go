package http

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gitlab.com/2foood/backend/internal/domain"
	"gitlab.com/2foood/backend/internal/pkg/ctxkey"
	"gitlab.com/2foood/backend/internal/pkg/locale"
	"gitlab.com/2foood/backend/internal/pkg/respond"
)

const (
	HeaderAcceptLanguage = "Accept-Language"
)

type AuthService interface {
	SendPhoneOTP(ctx context.Context, phone string) error
	SendEmailOTP(ctx context.Context, email string) error
	LoginWithPhone(ctx context.Context, phone, code string, lang domain.Language) (*domain.User, *domain.TokenPair, error)
	LoginWithGoogle(ctx context.Context, idToken string) (*domain.OAuthResult, error)
	CompleteGoogleLogin(ctx context.Context, sessionID, phone string, code string, lang domain.Language) (*domain.User, *domain.TokenPair, error)
	LinkGoogle(ctx context.Context, userID uuid.UUID, idToken string) error
	RefreshTokens(ctx context.Context, refreshToken string) (*domain.TokenPair, error)
}

type Handler struct {
	log      *slog.Logger
	service  AuthService
	validate *validator.Validate
}

func NewHandler(log *slog.Logger, service AuthService, validate *validator.Validate) *Handler {
	return &Handler{
		log:      log,
		service:  service,
		validate: validate,
	}
}

// PhoneSendOTP godoc
// @Summary      Send OTP via SMS
// @Description  Requests a One-Time Password (OTP) to be sent to the specified phone number
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body PhoneSendOTPRequest true "Data for sending OTP"
// @Success      200 {object} respond.MessageResponse "Success message (code sent)"
// @Failure      400 {object} respond.ErrorResponse "Validation error or invalid JSON"
// @Failure      429 {object} respond.ErrorResponse "Too many attempts"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /auth/phone/otp [post]
func (h *Handler) PhoneSendOTP(w http.ResponseWriter, r *http.Request) {
	const op = "auth.Handler.PhoneSendOTP"
	log := h.log.With(slog.String("op", op))

	var req PhoneSendOTPRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	if err := h.service.SendPhoneOTP(r.Context(), req.Phone); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	respond.Message(w, r, "code sent")
}

// EmailSendOTP godoc
// @Summary      Send OTP via email
// @Description  Requests a One-Time Password (OTP) to be sent to the specified email address
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body EmailSendOTPRequest true "Data for sending OTP"
// @Success      200 {object} respond.MessageResponse "Success message (code sent)"
// @Failure      400 {object} respond.ErrorResponse "Validation error or invalid JSON"
// @Failure      429 {object} respond.ErrorResponse "Too many attempts"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /auth/email/otp [post]
func (h *Handler) EmailSendOTP(w http.ResponseWriter, r *http.Request) {
	const op = "auth.Handler.EmailSendOTP"
	log := h.log.With(slog.String("op", op))

	var req EmailSendOTPRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	if err := h.service.SendEmailOTP(r.Context(), req.Email); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	respond.Message(w, r, "code sent")
}

// PhoneLogin godoc
// @Summary      Login with phone number
// @Description  Verifies the OTP code and authorizes the user, returning access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body PhoneLoginRequest true "Login credentials"
// @Success      200 {object} AuthResponse "Successful authorization"
// @Failure      400 {object} respond.ErrorResponse "Invalid code, validation error, or invalid JSON"
// @Failure      401 {object} respond.ErrorResponse "Invalid or expired token/code"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /auth/phone/login [post]
func (h *Handler) PhoneLogin(w http.ResponseWriter, r *http.Request) {
	const op = "auth.Handler.PhoneLogin"
	log := h.log.With(slog.String("op", op))

	var req PhoneLoginRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	// get language from accept language header
	lang := locale.FromAcceptLanguage(r.Header.Get(HeaderAcceptLanguage))
	
	user, tokens, err := h.service.LoginWithPhone(r.Context(), req.Phone, req.Code, lang)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, AuthResponse{
		Status: StatusSuccess,
		User:   user,
		Tokens: tokens,
	})
}

// GoogleLogin godoc
// @Summary      Login with Google
// @Description  Authorizes the user via Google ID Token. If the account is new and requires phone binding, returns StatusNeedPhone and a pre-auth token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body GoogleLoginRequest true "Google ID Token"
// @Success      200 {object} AuthResponse "Successful authorization"
// @Success      202 {object} OAuthNeedPhoneResponse "User needs to bind phone number"
// @Failure      400 {object} respond.ErrorResponse "Invalid token or validation error"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /auth/google/login [post]
func (h *Handler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	const op = "auth.Handler.GoogleLogin"
	log := h.log.With(slog.String("op", op))

	var req GoogleLoginRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	result, err := h.service.LoginWithGoogle(r.Context(), req.IDToken)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	// if user needs to bind phone number, return session id
	if result.NeedPhone {
		render.Status(r, http.StatusAccepted)
		render.JSON(w, r, OAuthNeedPhoneResponse{
			Status:    StatusNeedPhone,
			SessionID: result.SessionID,
		})
		return
	}

	render.JSON(w, r, AuthResponse{
		Status: StatusSuccess,
		Tokens: result.Tokens,
		User:   result.User,
	})
}

// CompleteGoogleLogin godoc
// @Summary      Complete Google login
// @Description  Completes the Google login process by linking the phone number to the account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body CompleteGoogleLoginRequest true "Data for completing Google login"
// @Success      200 {object} AuthResponse "Google login completed successfully"
// @Failure      400 {object} respond.ErrorResponse "Validation error, invalid JSON, or invalid code"
// @Failure      401 {object} respond.ErrorResponse "Invalid or expired session ID"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /auth/google/complete [post]
func (h *Handler) CompleteGoogleLogin(w http.ResponseWriter, r *http.Request) {
	const op = "auth.Handler.CompleteGoogleLogin"
	log := h.log.With(slog.String("op", op))

	var req CompleteGoogleLoginRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	// get language from accept language header
	lang := locale.FromAcceptLanguage(r.Header.Get(HeaderAcceptLanguage))

	user, tokens, err := h.service.CompleteGoogleLogin(r.Context(), req.SessionID, req.Phone, req.Code, lang)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, AuthResponse{
		Status: StatusSuccess,
		Tokens: tokens,
		User:   user,
	})
}

// LinkGoogle godoc
// @Summary      Link Google account
// @Description  Links a Google ID to an already existing authenticated account
// @Tags         auth
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        input body LinkGoogleRequest true "Google ID Token"
// @Success      200 {object} respond.MessageResponse "Google account successfully linked"
// @Failure      400 {object} respond.ErrorResponse "Invalid token or validation error"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      409 {object} respond.ErrorResponse "Google account already linked"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /auth/link/google [post]
func (h *Handler) LinkGoogle(w http.ResponseWriter, r *http.Request) {
	const op = "auth.Handler.LinkGoogle"
	log := h.log.With(slog.String("op", op))

	var req LinkGoogleRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	userID, ok := ctxkey.GetUserID(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrInternalServerError)
		return
	}

	if err := h.service.LinkGoogle(r.Context(), userID, req.IDToken); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	respond.Message(w, r, "google linked")
}

// Refresh godoc
// @Summary      Refresh tokens
// @Description  Issues a new pair of Access and Refresh tokens using a valid Refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body RefreshRequest true "Refresh token"
// @Success      200 {object} RefreshResponse "New pair of tokens"
// @Failure      400 {object} respond.ErrorResponse "Validation error or invalid JSON"
// @Failure      401 {object} respond.ErrorResponse "Refresh token is expired or invalid"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /auth/refresh [post]
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	const op = "auth.Handler.Refresh"
	log := h.log.With(slog.String("op", op))

	var req RefreshRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	tokens, err := h.service.RefreshTokens(r.Context(), req.RefreshToken)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, RefreshResponse{
		Tokens: tokens,
	})
}
