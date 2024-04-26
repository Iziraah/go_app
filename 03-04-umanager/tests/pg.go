package tests

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func CreateSchema(connStr string) error {
	ctx := context.TODO()

	usersDBConn, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		return err
	}

	_, err = usersDBConn.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS users
	(
		id         UUID NOT NULL,
		username   TEXT NOT NULL,
		password   TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

		CONSTRAINT pk_users_idx PRIMARY KEY (id),
		CONSTRAINT users_username_uniq_idx UNIQUE (username)
	)`)
	return err
}
