package ws

import (
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/QuUteO/video-communication/internal/model"
	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	logger    *slog.Logger
	upgrader  *websocket.Upgrader
	wsClients map[*websocket.Conn]struct{}
	broadcast chan *model.WsMessage
	mu        *sync.Mutex
}

func NewWsSocket(logger *slog.Logger) *WebSocketHandler {
	return &WebSocketHandler{
		logger:    logger,
		upgrader:  &websocket.Upgrader{},
		wsClients: make(map[*websocket.Conn]struct{}),
		broadcast: make(chan *model.WsMessage),
		mu:        &sync.Mutex{},
	}
}

// WsHandler подключение websocket
func (ws *WebSocketHandler) WsHandler(w http.ResponseWriter, r *http.Request) {
	const op = "./internal/server/handler/WsHandler"
	ws.logger.With("op: ", op)

	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		ws.logger.Error("Failed to upgrade connection", slog.Any("error", err))
		render.Status(r, http.StatusBadRequest)
		return
	}

	ws.logger.Info("Connected to server from: ", conn.RemoteAddr().String())

	ws.mu.Lock()
	ws.wsClients[conn] = struct{}{}
	ws.mu.Unlock()

	go ws.readFromClient(conn)
}

func (ws *WebSocketHandler) readFromClient(conn *websocket.Conn) {
	for {
		msg := new(model.WsMessage)
		err := conn.ReadJSON(msg)
		if err != nil {
			ws.logger.Error("Failed to read message", slog.Any("error", err))
			break
		}

		host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			ws.logger.Error("Failed to split host and port", slog.Any("error", err))
		}

		msg.IPAddress = host
		msg.Time = time.Now().Format("15:04")
	}
	ws.mu.Lock()
	delete(ws.wsClients, conn)
	ws.mu.Unlock()
}

func (ws *WebSocketHandler) WriteToClientsBroadcast() {
	for msg := range ws.broadcast {
		ws.mu.Lock()
		for conn := range ws.wsClients {
			func() {
				if err := conn.WriteJSON(msg); err != nil {
					ws.logger.Error("Failed to write message", slog.Any("error", err))
				}
			}()
		}
		ws.mu.Unlock()
	}
}
