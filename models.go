package main

import (
	"net/http"
)

type TopicMessagesDto struct {
	Topic    string   `json:"topic" validate:"required"`
	Messages []string `json:"messages" validate:"required|min_len:1"`
}

type ProducerRequestDto struct {
	Data []TopicMessagesDto `json:"data" validate:"required|min_len:1"`
}

func (dto *ProducerRequestDto) Bind(r *http.Request) error {

	return nil
}
