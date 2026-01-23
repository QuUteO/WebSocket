package websocket

import (
	"github.com/QuUteO/video-communication/internal/model"
)

type Hub struct {
	clients    map[*Client]struct{}
	register   chan *Client
	unregister chan *Client
	broadcast  chan *model.Message
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]struct{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *model.Message),
	}
}

func (h *Hub) Run() {
	for {

		select {

		case client := <-h.register:
			h.clients[client] = struct{}{}

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.broadcast)
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.broadcast <- message:
				default:
					delete(h.clients, client)
					close(client.broadcast)
				}
			}
		}
	}
}
