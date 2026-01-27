package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Message struct {
	ID      uuid.UUID `json:"id"`
	User    string    `json:"user"`    // отправитель
	Msg     string    `json:"msg"`     // текст пользователя
	Channel string    `json:"channel"` // канал, в котором пользователь зарегистрировался
	Time    time.Time `json:"time"`    // время отправки сообщения отправителем
}
