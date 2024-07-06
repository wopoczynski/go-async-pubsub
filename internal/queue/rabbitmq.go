package queue

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler func(ctx context.Context, mesaage *Message) error

type Message struct {
	Body      string
	QueueName string
}

type AMQP struct {
	Connection *amqp.Connection
	Ch         *amqp.Channel
	Db         *gorm.DB
	Cache      *redis.Client
}
