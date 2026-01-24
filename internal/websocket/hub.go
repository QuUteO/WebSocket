package websocket

import (
	"github.com/QuUteO/video-communication/internal/model"
)

type Hub struct {
	channels   map[*Client]struct{}
	register   chan *Client
	unregister chan *Client
	broadcast  chan model.Message
}

func NewHub() *Hub {
	return &Hub{
		channels:   make(map[*Client]struct{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan model.Message),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.register:
			h.channels[client] = struct{}{}

		case client := <-h.unregister:
			if _, ok := h.channels[client]; ok {
				delete(h.channels, client)
				close(client.send)
			}

		case message := <-h.broadcast:
			for client := range h.channels {
				select {

				case client.send <- message:

				default:
					delete(h.channels, client)
					close(client.send)
				}
			}
		}
	}
}
