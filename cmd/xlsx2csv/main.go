package main

import (
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
)

func main() {
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	log := zerolog.New(output).With().Timestamp().Logger()

	ctx := kong.Parse(&cli)
	ctx.Stdout = log
	err := ctx.Run(&Context{Logger: log})
	ctx.FatalIfErrorf(err)
}
