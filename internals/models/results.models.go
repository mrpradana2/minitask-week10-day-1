package models

type Message struct {
	Status int    `json:"status"`
	Msg    string `json:"message"`
	Result any    `json:"result,omitempty"`
	Token  string `json:"token,omitempty"`
}

// swagger:model
type MessageInternalServerError struct {
	Status int    `json:"status" example:"500"`
	Msg    string `json:"message"`
}

type MessageBadRequest struct {
	Status int    `json:"status" example:"400"`
	Msg    string `json:"message"`
}

type MessageNotFound struct {
	Status int    `json:"status" example:"404"`
	Msg    string `json:"message"`
}

type MessageCreated struct {
	Status int    `json:"status" example:"201"`
	Msg    string `json:"message"`
}

type MessageOK struct {
	Status int    `json:"status" example:"200"`
	Msg    string `json:"message"`
}

type MessageUnauthorized struct {
	Status int    `json:"status" example:"401"`
	Msg    string `json:"message"`
}

type MessageConflict struct {
	Status int    `json:"status" example:"409"`
	Msg    string `json:"message"`
}

// swagger:model
type MessageLogin struct {
	Status int    `json:"status" example:"200"`
	Msg    string `json:"message"`
	Token  string `json:"token,omitempty"`
}

// swagger:model
type MessageResult struct {
	Status int    `json:"status" example:"200"`
	Msg    string `json:"message"`
	Result any    `json:"result,omitempty"`
}