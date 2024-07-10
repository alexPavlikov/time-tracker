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
	http.HandleFunc("/v1/users", r.handler.UsersHandler)
	http.HandleFunc("/v1/user", r.handler.UserHandler)
	http.HandleFunc("/v1/add", r.handler.AddHandler)
	http.HandleFunc("/v1/update", r.handler.UpdateHandler)
	http.HandleFunc("/v1/metrics", r.handler.MetricsHandler)
}
