package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/sirupsen/logrus"
)

var (
	_                 = godotenv.Load()
	port              = os.Getenv("PORT")
	disablePulsarLogs = os.Getenv("DISABLE_PULSAR_LOGS")
)

func main() {
	initLogger()
	if strings.ToLower(disablePulsarLogs) == "yes" {
		logrus.SetOutput(io.Discard)
	}

	conn := NewConnection()
	conn.Connect()
	defer conn.Close()

	service := NewService(conn)
	handlers := NewHandlers(service)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Post("/produce", handlers.HandleProduceRequest)

	log.Info().Str("port", port).Msg("server started")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), r); err != nil {
		log.Fatal().Err(err).Msg("acute server failure")
	}
}
