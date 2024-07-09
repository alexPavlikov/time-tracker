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
	http.HandleFunc("/info", r.handler.InfoHandler)
	http.HandleFunc("/add", r.handler.AddHandler)
}
