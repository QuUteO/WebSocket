package routes

import (
	authhandler "github.com/QuUteO/video-communication/internal/auth/handler"
	authjwt "github.com/QuUteO/video-communication/internal/auth/jwt"
	authmiddleware "github.com/QuUteO/video-communication/internal/auth/middleware"
	"github.com/QuUteO/video-communication/internal/static"
	"github.com/QuUteO/video-communication/internal/user/handler"
	"github.com/QuUteO/video-communication/internal/websocket"
	"github.com/go-chi/chi/v5"
)

type Route struct {
	UserHandler      *handler.UserHandler
	WebSocketHandler *websocket.HandlerWS
	AuthHandler      *authhandler.Handler
	jwt              *authjwt.Manager
}

func NewRoute(
	userHandler *handler.UserHandler,
	WebSocketHandler *websocket.HandlerWS,
	AuthHandler *authhandler.Handler,
	jwt *authjwt.Manager) *Route {
	return &Route{
		UserHandler:      userHandler,
		WebSocketHandler: WebSocketHandler,
		AuthHandler:      AuthHandler,
		jwt:              jwt,
	}
}

func (h *Route) RegisterRoutes(router chi.Router) {
	router.Get("/", static.ServeHtml("index.html"))

	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.AuthHandler.Register)
		r.Post("/login", h.AuthHandler.Login)
	})

	router.Group(func(r chi.Router) {
		r.Use(authmiddleware.JWT(h.jwt))

		// websocket
		r.Route("/ws", func(r chi.Router) {
			r.Get("/", h.WebSocketHandler.WebSocketHTTP)
		})

		// users
		r.Route("/users", func(r chi.Router) {
			r.Get("/", h.UserHandler.GetAllUsers)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.UserHandler.GetUserByID)
				r.Put("/", h.UserHandler.UpdateUser)
				r.Delete("/", h.UserHandler.DeleteUser)
			})
		})
	})
}
