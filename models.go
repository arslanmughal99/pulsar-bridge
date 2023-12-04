package main

import (
	"net/http"
)

type MessagePayload struct {
	Key         string            `json:"key" validate:"min_len:1"`
	Payload     string            `json:"payload" validate:"required"`
	OrderingKey string            `json:"orderingKey" validate:"min_len:1"`
	Properties  map[string]string `json:"properties" validate:"-"`
}

type TopicMessagesDto struct {
	Topic    string           `json:"topic" validate:"required"`
	Messages []MessagePayload `json:"messages" validate:"required|min_len:1"`
}

type ProducerRequestDto struct {
	Data []TopicMessagesDto `json:"data" validate:"required|min_len:1"`
}

func (dto *ProducerRequestDto) Bind(r *http.Request) error {

	return nil
}
