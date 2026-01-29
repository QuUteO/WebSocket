package model

import (
	"time"

	"github.com/gofrs/uuid"
)

// User Структура для работы с БД
type User struct {
	Id        uuid.UUID `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

// DTOResponse Структура server для ответа
type DTOResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

// DTORequest Структура server для запроса
type DTORequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Response Структура для ответа с сервера
type Response struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Error      string      `json:"error"`
}
