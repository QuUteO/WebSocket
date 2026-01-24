package websocket

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocket struct {
	logger   *slog.Logger
	upgrader *websocket.Upgrader
	hub      *Hub
}

func NewWebSocket(hub *Hub, logger *slog.Logger) *WebSocket {
	return &WebSocket{
		upgrader: &websocket.Upgrader{},
		hub:      hub,
		logger:   logger,
	}
}

func (ws *WebSocket) WebSocketHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading websocket connection: %v \n", err.Error())
	}
	ws.logger.Info("WebSocket connection: %s", conn.RemoteAddr().String())

	client := NewClient(conn, ws.hub)
	client.hub.register <- client

	go client.WritePump()
	go client.ReadPump()
}
