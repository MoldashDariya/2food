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
	"gitlab.com/2foood/backend/internal/pkg/respond"
)

type UploadService interface {
	PresignUpload(ctx context.Context, folder string, userID uuid.UUID, contentType string) (string, string, error)
	Delete(ctx context.Context, fileURL string) error
}

type Handler struct {
	log      *slog.Logger
	service  UploadService
	validate *validator.Validate
}

func NewHandler(log *slog.Logger, service UploadService, validate *validator.Validate) *Handler {
	return &Handler{
		log:      log,
		service:  service,
		validate: validate,
	}
}

// Presign godoc
// @Summary      Get presigned upload URL
// @Tags         upload
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        input body PresignRequest true "Upload data"
// @Success      200 {object} PresignResponse
// @Failure      400 {object} respond.ErrorResponse "Invalid JSON"
// @Failure      401 {object} respond.ErrorResponse "Invalid authorization header"
// @Failure      406 {object} respond.ErrorResponse "Unsupported content type"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /upload/presign [post]
func (h *Handler) Presign(w http.ResponseWriter, r *http.Request) {
	const op = "upload.Handler.Presign"
	log := h.log.With(slog.String("op", op))

	var req PresignRequest
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

	uploadURL, fileURL, err := h.service.PresignUpload(r.Context(), req.Target, userID, req.ContentType)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, PresignResponse{
		UploadURL: uploadURL,
		FileURL:   fileURL,
	})
}

// Delete godoc
// @Summary      Delete uploaded file
// @Tags         upload
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        input body DeleteRequest true "Delete data"
// @Success      204
// @Failure      400 {object} respond.ErrorResponse "Invalid JSON"
// @Failure      401 {object} respond.ErrorResponse "Invalid authorization header"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /upload [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "upload.Handler.Delete"
	log := h.log.With(slog.String("op", op))

	var req DeleteRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	if err := h.service.Delete(r.Context(), req.FileURL); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.NoContent(w, r)
}
