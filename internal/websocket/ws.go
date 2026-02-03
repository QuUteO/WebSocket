package websocket

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/QuUteO/video-communication/internal/user/service"
	"github.com/gorilla/websocket"
)

type HandlerWS struct {
	upgrader websocket.Upgrader
	logger   *slog.Logger
	hub      *Hub
	service  service.Service
}

func NewHandlerWS(hub *Hub, service service.Service, logger *slog.Logger) *HandlerWS {
	return &HandlerWS{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// TODO: CORS правильно настроить
				return true
			},
		},
		hub:     hub,
		service: service,
		logger:  logger,
	}
}

func (h *HandlerWS) WebSocketHTTP(w http.ResponseWriter, r *http.Request) {
	const op = "WebSocketHTTP"
	h.logger.With("op: ", op)

	userID := r.URL.Query().Get("user_id")
	username := r.URL.Query().Get("username")

	if userID == "" {
		userID = "user_" + time.Now().Format("20060102150405")
	}
	if username == "" {
		username = "User_" + userID[len(userID)-6:]
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("Error upgrading websocket connection", slog.String("error", err.Error()))
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	client := NewClient(userID, username, conn, h.service, h.hub, h.logger)

	// запуск обработчиков
	go client.ReadPump()
	go client.WritePump()
}
