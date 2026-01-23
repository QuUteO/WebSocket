package websocket

import (
	"log"

	"github.com/QuUteO/video-communication/internal/model"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn      *websocket.Conn
	broadcast chan *model.Message
	hub       *Hub
}

func NewClient(conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		conn:      conn,
		broadcast: make(chan *model.Message),
		hub:       hub,
	}
}

// ReadPump слушает, что передает браузер
func (c *Client) ReadPump() {
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("Error closing websocket connection %v: \n", err.Error())
		}
	}(c.conn)

	for {
		var msg *model.Message
		if err := c.conn.ReadJSON(&msg); err != nil {
			log.Printf("Error reading websocket message:")
			break
		}
		c.hub.broadcast <- msg
	}
}

// WritePump пишет браузеру
func (c *Client) WritePump() {
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("Error closing websocket connection %v: \n", err.Error())
		}
	}(c.conn)

	for msg := range c.broadcast {
		if err := c.conn.WriteJSON(msg); err != nil {
			log.Printf("Error writing websocket message: %v \n", err.Error())
			break
		}
	}
}
