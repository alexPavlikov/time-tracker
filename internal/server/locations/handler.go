package locations

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alexPavlikov/time-tracker/internal/domain"
	"github.com/alexPavlikov/time-tracker/internal/server/service"
	"github.com/gofrs/uuid/v5"
)

type Handler struct {
	ctx     context.Context
	service service.Services
}

func New(context context.Context, service service.Services) *Handler {
	return &Handler{
		ctx:     context,
		service: service,
	}
}

// url: localhost:9090/v1/users?offset=_&limit=_&address=_&name=_&patronymic=_&surname=_&pass_number=_&pass_series
// Обязательные параметры offset limit
func (h *Handler) UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		start := h.startDoTask()

		uuid, err := uuid.FromString(h.ctx.Value("UUID").(string))
		if err != nil {
			slog.Error("failed to get user UUID from context", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		var m = domain.Metrics{
			User_ID:  uuid,
			FuncName: "UsersHandler",
			Time:     0,
		}

		defer h.endDoTask(start, m)

		r.ParseForm()

		var opt domain.UserSortParameters

		offset, _ := strconv.Atoi(r.FormValue("offset"))
		limit, _ := strconv.Atoi(r.FormValue("limit"))

		opt.Address = r.FormValue("address")
		opt.Name = r.FormValue("name")
		opt.Patronymic = r.FormValue("patronymic")
		opt.Surname = r.FormValue("surname")
		opt.PassportNumber, _ = strconv.Atoi(r.FormValue("pass_number"))
		opt.PassportSeries, _ = strconv.Atoi(r.FormValue("pass_series"))

		users, err := h.service.GetPag(limit, offset, opt)
		if err != nil {
			slog.Error("failed to get pagination users", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(users)
		if err != nil {
			slog.Error("failed to encode pagination users", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) UserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		r.ParseForm()
		id, err := strconv.Atoi(r.FormValue("user_id"))
		if err != nil {
			slog.Error("failed to convert user_id", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if err := h.service.Delete(id); err != nil {
			slog.Error("failed to delete user", "id", id, "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) AddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		start := h.startDoTask()

		uuid, err := uuid.FromString(h.ctx.Value("UUID").(string))
		if err != nil {
			slog.Error("failed to get user UUID from context", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		var m = domain.Metrics{
			User_ID:  uuid,
			FuncName: "AddHandler",
			Time:     0,
		}

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

		var APIUrl string = ""

		user, err := h.ApiInfoHandler(APIUrl, resp) // how create url to api
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

		user = domain.User{ //пример
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

		err = h.endDoTask(start, m)
		if err != nil {
			slog.Error("failed on endDoTask", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
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

func (h *Handler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPatch {
		start := h.startDoTask()

		uuid, err := uuid.FromString(h.ctx.Value("UUID").(string))
		if err != nil {
			slog.Error("failed to get user UUID from context", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		var m = domain.Metrics{
			User_ID:  uuid,
			FuncName: "UpdateHandler",
			Time:     0,
		}

		defer h.endDoTask(start, m)

		var user domain.User

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&user); err != nil {
			slog.Error("failed to decode passport", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if err := h.service.Update(user); err != nil {
			slog.Error("failed to update user", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

// Функция начала отсчета работы функции
func (h *Handler) startDoTask() time.Time {
	return time.Now()
}

// Функция посчета времени работы функция и записи ее в базу данных
func (h *Handler) endDoTask(start time.Time, m domain.Metrics) error {
	m.Time = time.Since(start)

	err := h.service.AddMetrics(m)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		uuid, err := uuid.FromString(h.ctx.Value("UUID").(string))
		if err != nil {
			slog.Error("failed to get user UUID from context", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		metrics, err := h.service.GetMetrics(uuid)
		if err != nil {
			slog.Error("failed to get metrics users", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&metrics)
		if err != nil {
			slog.Error("failed to encode users metrics", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
