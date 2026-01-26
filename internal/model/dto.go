package model

type Message struct {
	User string `json:"user"`
	Msg  string `json:"msg"`
	Time string `json:"time"`
}
