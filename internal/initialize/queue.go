package initialize

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/wopoczynski/playground/internal/queue"
)

type RabbitConfig struct {
	DSN           string `env:"DSN"`
	Queue         string `env:"QUEUE"`
	WorkersNumber int    `env:"WORKERS_NO,default=1"`
}

type RabbitInterface interface{}

func Start(cfg *RabbitConfig, db *gorm.DB, cache *redis.Client) (*queue.AMQP, error) {
	connection, err := amqp.Dial(cfg.DSN)
	if err != nil {
		return nil, err
	}
	ch, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		cfg.Queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to declare queue")
	}

	return &queue.AMQP{
		Connection: connection,
		Ch:         ch,
		Db:         db,
		Cache:      cache,
	}, nil
}
