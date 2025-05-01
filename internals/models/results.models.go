package models

type Message struct {
	Status string `json:"status"`
	Msg    string `json:"message"`
	Result any    `json:"result,omitempty"`
	Token  string `json:"token,omitempty"`
}