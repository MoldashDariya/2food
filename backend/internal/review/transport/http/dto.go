package http

import "gitlab.com/2foood/backend/internal/domain"

type ReviewResponse struct {
	Review *domain.Review `json:"review"`
}

type ReviewsResponse struct {
	Reviews []*domain.Review `json:"reviews"`
}

type ReviewPhotosResponse struct {
	Photos []*domain.ReviewPhoto `json:"photos"`
}

type CreateReviewRequest struct {
	Rating    int8     `json:"rating" validate:"required,min=1,max=5"`
	Comment   *string  `json:"comment" validate:"omitempty,min=1,max=1000"`
	PhotoURLs []string `json:"photo_urls" validate:"omitempty,max=5,dive,url"`
}

type AddReplyRequest struct {
	Reply string `json:"reply" validate:"required,min=1,max=1000"`
}
