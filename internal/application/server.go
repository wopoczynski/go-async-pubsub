package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	server "github.com/wopoczynski/playground/internal/http/echo"
	"github.com/wopoczynski/playground/internal/initialize"
	"github.com/wopoczynski/playground/internal/queue"

	_ "github.com/wopoczynski/playground/docs"
)

type Config struct {
	Http  *server.ServerConfig     `env:", prefix=HTTP_"`
	DB    *initialize.DBConfig     `env:", prefix=DB_"`
	Redis *initialize.RedisConfig  `env:", prefix=REDIS_"`
	AMQP  *initialize.RabbitConfig `env:", prefix=QUEUE_"`
}

type ApplicationContainer struct {
	cfg    *Config
	cache  *redis.Client
	db     *gorm.DB
	amqp   *queue.AMQP
	server server.Server
}

func New(cfg Config) (*ApplicationContainer, error) {
	db, err := initialize.DB(*cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("db initialize error: %w", err)
	}

	cache, err := initialize.Connect(cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("cache initialize error: %w", err)
	}
	amqp, err := initialize.Start(cfg.AMQP, db, cache)
	if err != nil {
		return nil, fmt.Errorf("amqp initialize error: %w", err)
	}
	return &ApplicationContainer{
		cfg:   &cfg,
		db:    db,
		cache: cache,
		amqp:  amqp,
	}, nil
}

func (s *ApplicationContainer) Init(ctx context.Context) {
	initialize.Automigrate(ctx, s.db)
}

func (s *ApplicationContainer) Start(ctx context.Context) {
	go func() {
		err := s.server.Start(":" + s.cfg.Http.HTTPServerPort)
		if errors.Is(err, http.ErrServerClosed) {
			log.Err(err).Msg("Server shutdown")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	const shutdownTimeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Msg("shutting down server...")
	}

	log.Info().Msg("server stopped gracefully")
}
