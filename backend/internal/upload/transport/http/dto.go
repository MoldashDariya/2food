package http

type PresignRequest struct {
	Target      string `json:"target" validate:"required,oneof=avatar restaurant_logo restaurant_banner product category review"` // folder name
	ContentType string `json:"content_type" validate:"required,oneof=image/jpeg image/png image/webp"`
}

type PresignResponse struct {
	UploadURL string `json:"upload_url"`
	FileURL   string `json:"file_url"`
}

type DeleteRequest struct {
	FileURL string `json:"file_url" validate:"required,url"`
}
