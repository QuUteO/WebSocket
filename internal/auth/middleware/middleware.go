package authmiddleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/QuUteO/video-communication/internal/auth/jwt"
)

type ctxKey string

const UserIDKey ctxKey = "user_id"

func JWT(jwt *authjwt.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			if h == "" {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(h, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			userID, err := jwt.Parse(parts[1])
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
