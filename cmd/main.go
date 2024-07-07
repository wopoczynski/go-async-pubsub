package main

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"

	server "github.com/wopoczynski/playground/internal/application"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_ = godotenv.Load()

	var cfg server.Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		log.Fatal().Err(err).Msg("app config boot failed")
	}

	app, err := server.New(cfg)
	if err != nil {
		panic(fmt.Errorf("server error %v", err))
	}

	app.Init(ctx)
	app.Start(ctx)
}
