package domain

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ID              uuid.UUID  `json:"id"`
	UserID          uuid.UUID  `json:"user_id"`
	RestaurantID    uuid.UUID  `json:"restaurant_id"`
	OrderID         uuid.UUID  `json:"order_id"`
	Rating          int8       `json:"rating"`
	Comment         *string    `json:"comment,omitempty"`
	PhotoURLs       []string   `json:"photo_urls"`
	Reply           *string    `json:"reply,omitempty"`
	RepliedAt       *time.Time `json:"replied_at,omitempty"`
	IsHidden        bool       `json:"is_hidden"`
	AuthorName      *string    `json:"author_name,omitempty"`
	AuthorAvatarURL *string    `json:"author_avatar_url,omitempty"`
	RestaurantName  *string    `json:"restaurant_name,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type ReviewPhoto struct {
	URL             string     `json:"url"`
	ReviewID        uuid.UUID  `json:"review_id"`
	AuthorName      *string    `json:"author_name,omitempty"`
	AuthorAvatarURL *string    `json:"author_avatar_url,omitempty"`
	Rating          int8       `json:"rating"`
	Comment         *string    `json:"comment,omitempty"`
	Reply           *string    `json:"reply,omitempty"`
	RepliedAt       *time.Time `json:"replied_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type CreateReviewInput struct {
	UserID       uuid.UUID
	RestaurantID uuid.UUID
	OrderID      uuid.UUID
	Rating       int8
	Comment      *string
	PhotoURLs    []string
}

func NewReview(input CreateReviewInput) *Review {
	now := time.Now().UTC()
	return &Review{
		ID:           uuid.New(),
		UserID:       input.UserID,
		RestaurantID: input.RestaurantID,
		OrderID:      input.OrderID,
		Rating:       input.Rating,
		Comment:      input.Comment,
		PhotoURLs:    input.PhotoURLs,
		IsHidden:     false, // reviews are not hidden by default
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (r *Review) AddReply(reply string) {
	now := time.Now().UTC()
	r.Reply = &reply
	r.RepliedAt = &now
	r.UpdatedAt = now
}
