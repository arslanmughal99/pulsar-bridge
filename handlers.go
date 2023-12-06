package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

type Handlers struct {
	service *Service
}

func NewHandlers(s *Service) *Handlers {
	h := new(Handlers)
	h.service = s
	return h
}

func (h *Handlers) HandleProduceRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	dto := new(ProducerRequestDto)
	if err := render.Bind(r, dto); err != nil {
		_, _ = w.Write([]byte("something went wrong"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dtoErrs := ValidateDto(dto)

	if dtoErrs != nil {
		e, err := json.Marshal(dtoErrs)
		if err != nil {
			_, _ = w.Write([]byte("Invalid request"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, _ = w.Write(e)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	total, err := h.service.ProduceMessagesTnx(dto)
	if err != nil {
		_, _ = w.Write([]byte("something went wrong"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte(fmt.Sprintf("{\"total\": %d}", *total)))
	w.WriteHeader(http.StatusOK)
	return
}
