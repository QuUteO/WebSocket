package websocket

import (
	"log/slog"
	"sync"
	"time"

	"github.com/QuUteO/video-communication/internal/model"
	"github.com/gofrs/uuid"
)

type Hub struct {
	channels map[string]map[*Client]bool // мапа для хранения пользователей в канале

	register   chan *ClientRegistration // канал для регистрации в канал
	unregister chan *ClientRegistration // канал для ухода из канала
	broadcast  chan model.Message       // канал для трансляции всем пользователем в канале

	mu     *sync.RWMutex
	logger *slog.Logger
}

type ClientRegistration struct {
	Client  *Client
	Channel string
}

func NewHub(logger *slog.Logger) *Hub {
	return &Hub{
		channels:   make(map[string]map[*Client]bool),
		register:   make(chan *ClientRegistration),
		unregister: make(chan *ClientRegistration),
		broadcast:  make(chan model.Message),
		mu:         &sync.RWMutex{},
		logger:     logger,
	}
}

func (h *Hub) Run() {
	for {

		select {

		case registration := <-h.register:
			h.registerClientToChannel(registration.Client, registration.Channel)

		case registration := <-h.unregister:

			h.unregisterClientToChannel(registration.Client, registration.Channel)

		case msg := <-h.broadcast:

			h.broadcastToChannel(msg)
		}

	}
}

func (h *Hub) registerClientToChannel(client *Client, channel string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Если канал не создан, то создаем его
	if _, ok := h.channels[channel]; !ok {
		h.channels[channel] = make(map[*Client]bool)
		h.logger.Info("channel created", slog.String("channel", channel))
	}

	// Добавление клиента в канал
	h.channels[channel][client] = true

	// Добавление название канала в информацию о клиенте
	client.CurrentChannel = channel

	// Создание системного сообщения
	systemMsg := model.Message{
		ID:      uuid.Must(uuid.NewV4()),
		User:    "System",
		Msg:     client.Username + " присоединился к каналу",
		Channel: channel,
		Time:    time.Now(),
	}

	for c := range h.channels[channel] {
		if c != client {
			select {
			case c.Send <- systemMsg:
			default:
				h.logger.Warn("client channel overflow, disconnecting",
					slog.String("client_id", c.ID))
				delete(h.channels[channel], c)
				close(c.Send)
			}
		}
	}
}

// Отписка клиента от канала
func (h *Hub) unregisterClientToChannel(client *Client, channel string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if ch, ok := h.channels[channel]; ok {
		delete(ch, client)

		client.CurrentChannel = ""

		systemMsg := model.Message{
			ID:      uuid.Must(uuid.NewV4()),
			User:    "System",
			Msg:     client.Username + " покинул канал",
			Channel: channel,
			Time:    time.Now(),
		}

		for c := range h.channels[channel] {
			if c != client {
				select {
				case c.Send <- systemMsg:
				default:
					h.logger.Warn("client channel overflow, disconnecting",
						slog.String("client_id", c.ID))
					delete(h.channels[channel], c)
					close(c.Send)
				}
			}
		}

		if len(ch) == 0 {
			delete(h.channels, channel)
			h.logger.Info("channel deleted (empty)", slog.String("channel", channel))
		}
	}

}

// Рассылка в определенный канал
func (h *Hub) broadcastToChannel(msg model.Message) {
	h.mu.Lock()
	defer h.mu.Unlock()

	channelName := msg.Channel

	if channel, ok := h.channels[channelName]; ok {
		for c := range channel {
			select {
			case c.Send <- msg:
			default:
				h.logger.Warn("client channel overflow, disconnecting",
					slog.String("client_id", c.ID))
				delete(channel, c)
				close(c.Send)
			}

		}
	} else {
		h.logger.Warn("channel does not exist", slog.String("channel", channelName))
	}
}

func (h *Hub) GetChannels() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	channels := make([]string, 0, len(h.channels))
	for name := range h.channels {
		channels = append(channels, name)
	}
	return channels
}

func (h *Hub) GetClientsInChannel(channelName string) []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients := make([]string, 0)
	if channel, exists := h.channels[channelName]; exists {
		for client := range channel {
			clients = append(clients, client.Username)
		}
	}
	return clients
}
