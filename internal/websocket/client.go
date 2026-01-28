package websocket

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/QuUteO/video-communication/internal/model"
	"github.com/QuUteO/video-communication/internal/user/service"
	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID             string             // Уникальный ID клиента
	Conn           *websocket.Conn    // WebSocket соединение
	Send           chan model.Message // Канал для отправки сообщений
	Hub            *Hub               // Хаб
	Srv            service.Service    // Слой сервиса для работы с БД
	CurrentChannel string             // Текущий канал
	Username       string             // Имя пользователя
	Logger         *slog.Logger
}

func NewClient(clientID, username string, conn *websocket.Conn, srv service.Service, hub *Hub, logger *slog.Logger) *Client {
	return &Client{
		ID:       clientID,
		Conn:     conn,
		Send:     make(chan model.Message, 256),
		Hub:      hub,
		Srv:      srv,
		Username: username,
		Logger:   logger,
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.mu.Lock()
		if c.CurrentChannel != "" {
			c.Hub.unregister <- &ClientRegistration{
				Client:  c,
				Channel: c.CurrentChannel,
			}
		}
		c.Hub.mu.Unlock()

		// Закрываем соединение
		if c.Conn != nil {
			c.Conn.Close()
		}

		c.Logger.Info("read pump stopped", slog.String("client_id", c.ID))
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
				websocket.CloseNormalClosure) {
				c.Logger.Debug("websocket closed", slog.String("error", err.Error()))
			} else {
				c.Logger.Error("error reading from websocket", slog.String("error", err.Error()))
			}
			break
		}

		c.handleReadMessage(message)
	}
}

// чтение сообщения
func (c *Client) handleReadMessage(data []byte) {
	var rawMsg map[string]interface{}
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		c.Logger.Error("Error unmarshalling message:", slog.String("error", err.Error()))
		return
	}

	msgType, ok := rawMsg["type"].(string)
	if !ok {
		c.Logger.Error("message type is required")
		return
	}

	switch msgType {
	case "join":
		c.handleJoinMessage(rawMsg)
	case "message":
		c.handleMessage(rawMsg)
	case "leave":
		c.handleLeave()
	}
}

// обработка присоединения к каналу
func (c *Client) handleJoinMessage(rawMsg map[string]interface{}) {
	channel, ok := rawMsg["channel"].(string)
	if !ok {
		c.Logger.Error("channel is required")
		return
	}

	if c.CurrentChannel != "" {
		c.Hub.unregister <- &ClientRegistration{
			Client:  c,
			Channel: c.CurrentChannel,
		}
	}

	c.CurrentChannel = channel

	c.Hub.register <- &ClientRegistration{
		Client:  c,
		Channel: c.CurrentChannel,
	}

	// загрузка истории сообщения
	go c.loadChannelHistory(channel)

	response := map[string]interface{}{
		"type":    "joined",
		"channel": channel,
		"user":    c.Username,
	}
	c.Conn.WriteJSON(response)
}

// загрузка сообщений из БД
func (c *Client) loadChannelHistory(channel string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	messages, err := c.Srv.GetMessageByChannel(ctx, channel)
	if err != nil {
		c.Logger.Error("Error loading history:", slog.String("error", err.Error()))
		return
	}

	for _, message := range messages {
		select {
		case c.Send <- message:
		case <-ctx.Done():
			c.Logger.Error("Error loading history:", slog.String("error", ctx.Err().Error()))
			return
		}
	}
}

func (c *Client) handleMessage(rawMsg map[string]interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if c.CurrentChannel == "" {
		c.Logger.Warn("client not in any channel", slog.String("client_id", c.ID))
		return
	}

	msgText, ok := rawMsg["msg"].(string)
	if !ok || msgText == "" {
		c.Logger.Error("message type is required")
		return
	}

	msg := model.Message{
		ID:      uuid.Must(uuid.NewV4()),
		User:    c.Username,
		Msg:     msgText,
		Channel: c.CurrentChannel,
		Time:    time.Now(),
	}

	if err := c.Srv.SaveMsg(ctx, msg); err != nil {
		c.Logger.Error("Error saving message:", slog.String("error", err.Error()))
		return
	}

	c.Hub.broadcast <- msg
}

func (c *Client) handleLeave() {
	if c.CurrentChannel == "" {
		c.Logger.Warn("client not in any channel", slog.String("client_id", c.ID))
		return
	}

	c.Hub.unregister <- &ClientRegistration{
		Client:  c,
		Channel: c.CurrentChannel,
	}

	response := map[string]interface{}{
		"type":    "leave",
		"channel": c.CurrentChannel,
		"user":    c.Username,
	}
	c.Conn.WriteJSON(response)

	c.CurrentChannel = ""
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()

		if c.Conn != nil {
			c.Conn.Close()
		}

		close(c.Send)

		c.Logger.Info("write pump stopped", slog.String("client_id", c.ID))
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Logger.Warn("client send channel closed")
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			defer func() {
				if r := recover(); r != nil {
					c.Logger.Error("panic in WritePump", slog.Any("recover", r))
				}
			}()

			if err := c.Conn.WriteJSON(message); err != nil {
				c.Logger.Error("Error sending message:", slog.String("error", err.Error()))
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.Logger.Error("Error sending ping:", slog.String("error", err.Error()))
				return
			}
		}
	}
}
