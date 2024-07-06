package initialize

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/wopoczynski/playground/internal/queue"
)

type RabbitConfig struct {
	DSN   string `env:"DSN"`
	Queue string `env:"QUEUE"`
}

type RabbitInterface interface{}

func Start(cfg *RabbitConfig, db *gorm.DB, cache *redis.Client) (*queue.AMQP, error) {
	connection, err := amqp.Dial(cfg.DSN)
	if err != nil {
		return nil, err
	}
	defer connection.Close()
	ch, err := connection.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	return &queue.AMQP{
		Connection: connection,
		Ch:         ch,
		Db:         db,
		Cache:      cache,
	}, nil
}
