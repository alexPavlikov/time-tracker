package router

import (
	"net/http"

	"github.com/alexPavlikov/time-tracker/internal/server/locations"
)

type Repository struct {
	handler locations.Handler
}

func New(handler locations.Handler) *Repository {
	return &Repository{
		handler: handler,
	}
}

func (r *Repository) Build() {
	http.HandleFunc("/v1/info", r.handler.InfoHandler)
	http.HandleFunc("/v1/add", r.handler.AddHandler)
}
