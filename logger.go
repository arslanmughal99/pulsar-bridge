package main

import (
	"os"
	"strconv"

	"github.com/onsi/ginkgo/reporters/stenographer/support/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func initLogger() {
	level, _ := strconv.Atoi(os.Getenv("LOG_LEVEL"))
	log.Logger = log.Level(zerolog.Level(level)).Output(
		zerolog.ConsoleWriter{
			Out: colorable.NewColorableStdout(), TimeFormat: "01/02/06 3:04 PM",
		},
	)
}
