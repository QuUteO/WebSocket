package websocket

import (
	"log"

	"github.com/QuUteO/video-communication/internal/model"
	"github.com/QuUteO/video-communication/internal/user/service"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan model.Message
	hub  *Hub
	srv  *service.Service
}

func NewClient(conn *websocket.Conn, hub *Hub, srv *service.Service) *Client {
	return &Client{
		conn: conn,
		send: make(chan model.Message),
		hub:  hub,
		srv:  srv,
	}
}

// ReadPump слушает, что передает браузер
func (c *Client) ReadPump() {
	defer func(conn *websocket.Conn) {
		c.hub.unregister <- c
		err := c.conn.Close()
		if err != nil {
			log.Printf("Error closing websocket connection %v: \n", err.Error())
		}
	}(c.conn)

	c.hub.register <- c

	for {
		var msg model.Message

		if err := c.conn.ReadJSON(&msg); err != nil {
			log.Printf("Error reading websocket message:")
			break
		}

		c.hub.broadcast <- msg
	}
}

// WritePump пишет браузеру
func (c *Client) WritePump() {

	for msg := range c.send {
		if err := c.conn.WriteJSON(msg); err != nil {
			log.Printf("Error writing websocket message: %v \n", err.Error())
			break
		}
	}
}
