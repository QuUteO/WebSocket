package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/QuUteO/video-communication/internal/model"
	"github.com/QuUteO/video-communication/internal/user/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type UserHandler struct {
	service service.Service
	logger  *slog.Logger
}

func NewUserHandler(srv service.Service, log *slog.Logger) *UserHandler {
	return &UserHandler{
		service: srv,
		logger:  log,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	const op = "./internal/server/handler/CreateUser"
	h.logger.With("op", op)

	var req model.DTORequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode JSON", slog.Any("error", err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.Response{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	ctx := r.Context()
	_, err := h.service.CreateUser(ctx, req.Email, req.Password)
	if err != nil {
		h.logger.Error("Failed to create server", slog.Any("error", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, model.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, model.Response{
		StatusCode: http.StatusCreated,
		Message:    "User created successfully",
		Error:      "nil",
	})
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	const op = "./internal/server/handler/FindUserById"
	h.logger.With("op", op)

	ctx := r.Context()
	users, err := h.service.FindAllUser(ctx)
	if err != nil {
		h.logger.Error("Failed to find users", slog.Any("error", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, model.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Users retrieved successfully",
		Data:       users,
		Error:      "nil",
	})
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	const op = "./internal/server/handler/FindUserById"
	h.logger.With("op", op)

	id := chi.URLParam(r, "id")

	ctx := r.Context()
	u, err := h.service.FindUserById(ctx, id)
	if err != nil {
		h.logger.Error("Failed to find server", slog.Any("error", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, model.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, model.Response{
		StatusCode: http.StatusOK,
		Message:    "User retrieved successfully",
		Data:       u,
		Error:      "nil",
	})
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	const op = "./internal/server/handler/UpdateUserByID"
	h.logger.With("op", op)

	id := chi.URLParam(r, "id")

	var req model.DTORequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode JSON", slog.Any("error", err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.Response{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	ctx := r.Context()
	if err := h.service.UpdateUser(ctx, id, req.Email, req.Password); err != nil {
		h.logger.Error("Failed to update server", slog.Any("error", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, model.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, model.Response{
		StatusCode: http.StatusOK,
		Message:    "User updated successfully",
		Error:      "nil",
	})
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	const op = "./internal/server/handler/DeleteUserByID"
	h.logger.With("op", op)

	id := chi.URLParam(r, "id")

	ctx := r.Context()
	if err := h.service.DeleteUser(ctx, id); err != nil {
		h.logger.Error("Failed to delete server", slog.Any("error", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, model.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, model.Response{
		StatusCode: http.StatusOK,
		Message:    "User deleted successfully",
		Error:      "nil",
	})
}
