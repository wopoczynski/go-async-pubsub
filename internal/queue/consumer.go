package queue

import (
	"context"
	"encoding/json"

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
			body := string(d.Body)
			persisted := &database.PersistingStruct{
				ID:   uuid.New(),
				Data: body,
			}
			a.Db.WithContext(ctx).Create(persisted)
			err := a.Cache.Set(ctx, persisted.ID.String(), body, 0).Err()
			if err != nil {
				log.Error().Err(err).Msg("failed ot persist to cache")
			}
			data, err := a.Cache.Get(ctx, persisted.ID.String()).Result()
			if err != nil {
				log.Error().Err(err).Msg("failed to retrieve from cache")
			}

			encoded, _ := json.Marshal(persisted)
			log.Info().Msgf("persisted in db: %s", encoded)
			log.Info().Msgf("persisted in cache: %s", data)
		}
	}()

	log.Info().Msgf("started consumer from queue: %s", q.Name)
	<-k
}
