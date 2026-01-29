package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Manager struct {
	secret []byte
	ttl    time.Duration
}

func NewJWT(secret string, ttl time.Duration) *Manager {
	return &Manager{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

func (m *Manager) Generate(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(m.ttl).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *Manager) Parse(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return m.secret, nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	return claims["sub"].(string), nil
}
