package websocket

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/QuUteO/video-communication/internal/user/service"
	"github.com/gorilla/websocket"
)

type WebSocket struct {
	logger   *slog.Logger
	upgrader *websocket.Upgrader
	hub      *Hub
	service  *service.Service
}

func NewWebSocket(hub *Hub, logger *slog.Logger, srv *service.Service) *WebSocket {
	return &WebSocket{
		logger:   logger,
		upgrader: &websocket.Upgrader{},
		hub:      hub,
		service:  srv,
	}
}

func (ws *WebSocket) WebSocketHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading websocket connection: %v \n", err.Error())
	}
	ws.logger.Info("WebSocket connection: %s", conn.RemoteAddr().String())

	client := NewClient(conn, ws.hub, *ws.service)
	client.hub.register <- client

	go client.WritePump()
	go client.ReadPump()
}
