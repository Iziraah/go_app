package link_updater

import (
	"context"
)

func New(repository repository, consumer amqpConsumer) *Story {
	return &Story{repository: repository, consumer: consumer}
}

type Story struct {
	repository repository
	consumer   amqpConsumer
}

func (s *Story) Run(ctx context.Context) error {
	ch, err := s.consumer.Consume(s.queueName, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case m, ok := <-ch:
			if !ok {
				return errors.New("rabbitmq queue is closed")
			}
			err := s.processMsg(ctx, m)
			if err != nil {
				slog.Error("process message error", slog.Any("err", err))
			}
		}
	}
}