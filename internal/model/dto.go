package model

import "time"

type Message struct {
	User    string    `json:"user"`
	Msg     string    `json:"msg"`
	Channel string    `json:"channel"`
	Time    time.Time `json:"time"`
}
