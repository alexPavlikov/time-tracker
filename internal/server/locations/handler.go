package locations

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/alexPavlikov/time-tracker/internal/domain"
	"github.com/alexPavlikov/time-tracker/internal/server/service"
)

type Handler struct {
	service service.Services
}

func New(service service.Services) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InfoHandler(w http.ResponseWriter, r *http.Request) {
	// if r.Method == "GET" {

	// }
}

func (h *Handler) AddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var req domain.RequestAdd

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			slog.Error("failed to decode passport", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		slice := strings.Split(req.Passport, " ")

		if len(slice) != 2 {
			slog.Error("failed to convert all passport")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		var err error
		var resp domain.ResponseAdd

		resp.PassportSeries, err = strconv.Atoi(slice[0])
		if err != nil {
			slog.Error("failed to convert passport series", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		resp.PassportNumber, err = strconv.Atoi(slice[1])
		if err != nil {
			slog.Error("failed to convert passport number", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		var APIUrl string

		user, err := h.ApiInfoHandler(APIUrl, resp)
		if err != nil {
			if errors.Is(err, errors.New("bad request")) {
				slog.Error("failed send to API", "error", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			} else if errors.Is(err, errors.New("internal server error")) {
				slog.Error("failed to convert passport number", "error", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		user = domain.User{
			ID:             1,
			PassportSeries: 2,
			PassportNumber: 3,
			Surname:        "123",
			Name:           "456",
			Patronymic:     "789",
			Address:        "000",
		}

		user.ID, err = h.service.Add(user)
		if err != nil {
			slog.Error("failed on service add", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		slog.Info("user add successfull", "id", user.ID)
	}
}

func (h *Handler) ApiInfoHandler(APIUrl string, resp domain.ResponseAdd) (user domain.User, err error) {
	r, err := http.Get(APIUrl)
	if err != nil {
		return domain.User{}, err
	}

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&user); err != nil {
		return domain.User{}, err
	}

	user.PassportNumber = resp.PassportNumber
	user.PassportSeries = resp.PassportSeries

	return user, nil
}
