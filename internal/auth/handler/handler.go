package authhandler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/QuUteO/video-communication/internal/auth/service"
	"github.com/QuUteO/video-communication/internal/model"
	"github.com/go-chi/render"
)

type Handler struct {
	service authservice.Service
	logger  *slog.Logger
}

func NewHandler(service authservice.Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "internal/auth/handler/Register"
	log := h.logger.With("op: ", op)

	var req *model.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("error decoding request", slog.String("error", err.Error()))
		render.Status(r, http.StatusUnprocessableEntity)
		return
	}

	token, err := h.service.Register(r.Context(), req)
	if err != nil {
		log.Error("error registering user", slog.String("error", err.Error()))
		render.Status(r, http.StatusUnprocessableEntity)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, model.AuthResponse{Token: token})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "internal/auth/handler/Login"
	log := h.logger.With("op: ", op)

	var req *model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("error decoding request", slog.String("error", err.Error()))
		render.Status(r, http.StatusUnprocessableEntity)
		return
	}

	token, err := h.service.Login(r.Context(), req)
	if err != nil {
		log.Error("error login", slog.String("error", err.Error()))
		render.Status(r, http.StatusUnprocessableEntity)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, model.LoginResponse{Token: token})
}
