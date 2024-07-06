package queue

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/wopoczynski/playground/internal/database"
)

func (a *AMQP) Consume(queueName string) {
	q, err := a.Ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to declare queue")
	}
	msgs, err := a.Ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to consume messages")
	}
	k := make(chan bool)

	go func() {
		ctx := context.Background()
		for d := range msgs {
			id := uuid.New()
			a.Db.Create(&database.PersistingStruct{
				ID:   id,
				Data: string(d.Body),
			})
			a.Cache.Set(ctx, id.String(), string(d.Body), 0)
			log.Info().Msgf("Received a message: %s", d.Body)
		}
	}()

	log.Info().Msgf("consuming from queue: %s", q.Name)
	<-k
}
