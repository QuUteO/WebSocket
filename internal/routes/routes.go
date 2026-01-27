package routes

import (
	"github.com/QuUteO/video-communication/internal/static"
	"github.com/QuUteO/video-communication/internal/user/handler"
	"github.com/QuUteO/video-communication/internal/websocket"
	"github.com/go-chi/chi/v5"
)

type Route struct {
	UserHandler      *handler.UserHandler
	WebSocketHandler *websocket.HandlerWS
}

func NewRoute(userHandler *handler.UserHandler, WebSocketHandler *websocket.HandlerWS) *Route {
	return &Route{
		UserHandler:      userHandler,
		WebSocketHandler: WebSocketHandler,
	}
}

func (h *Route) RegisterRoutes(router chi.Router) {
	router.Get("/", static.ServeHtml("index.html"))

	router.Route("/ws", func(r chi.Router) {
		r.Get("/", h.WebSocketHandler.WebSocketHTTP)
	})

	router.Route("/users", func(r chi.Router) {
		r.Post("/", h.UserHandler.CreateUser) // POST /users
		r.Get("/", h.UserHandler.GetAllUsers) // GET /users

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.UserHandler.GetUserByID)   // GET /users/{id}
			r.Put("/", h.UserHandler.UpdateUser)    // PUT /users/{id}
			r.Delete("/", h.UserHandler.DeleteUser) // DELETE /users/{id}
		})
	})
}
