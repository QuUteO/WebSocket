package websocket

import (
	"log"
	"net/http"

	"github.com/QuUteO/video-communication/internal/model"
	"github.com/gorilla/websocket"
)

type WebSocket struct {
	upgrader *websocket.Upgrader
	hub      *Hub
}

func NewWebSocket(hub *Hub) *WebSocket {
	return &WebSocket{
		upgrader: &websocket.Upgrader{},
		hub:      hub,
	}
}

func (ws *WebSocket) WebSocketHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading websocket connection: %v \n", err.Error())
	}

	client := &Client{
		conn:      conn,
		broadcast: make(chan *model.Message),
		hub:       *ws.hub,
	}

	client.hub.register <- client

	go client.WritePump()
	client.ReadPump()
}
