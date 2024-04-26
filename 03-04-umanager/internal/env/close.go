package env

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v4/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type Closer struct {
	usersDBConn *pgxpool.Pool
	linksDBConn *mongo.Client
	amqpConn    *amqp.Connection
	amqpChannel *amqp.Channel
}

func NewCloser(usersDBConn *pgxpool.Pool, linksDBConn *mongo.Client, amqpConn *amqp.Connection, amqpChannel *amqp.Channel) *Closer {
	return &Closer{
		usersDBConn: usersDBConn,
		linksDBConn: linksDBConn,
		amqpConn:    amqpConn,
		amqpChannel: amqpChannel,
	}
}

func (c *Closer) Close(ctx context.Context) {
	defer c.amqpChannel.Close()
	defer c.amqpConn.Close()
	defer func() {
		err := c.linksDBConn.Disconnect(ctx)
		if err != nil {
			slog.Error("closing", slog.Any("err", err))
		}
	}()
	defer c.usersDBConn.Close()
}
