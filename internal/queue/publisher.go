package queue

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

func (a *AMQP) Publish(m *Message) {
	q, err := a.Ch.QueueDeclare(
		m.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Error().Err(err).Msg("queue not declared")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = a.Ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/json",
			Body:        []byte(m.Body),
		})
	if err != nil {
		log.Error().Err(err).Msg("message not published")
	}
}
