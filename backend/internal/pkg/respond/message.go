package respond

import (
	"net/http"

	"github.com/go-chi/render"
)

type MessageResponse struct {
	Message string `json:"message"`
}

func Message(w http.ResponseWriter, r *http.Request, message string) {
	render.JSON(w, r, MessageResponse{Message: message})
}
