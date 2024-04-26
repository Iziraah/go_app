package tests

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	postgresImg = "postgres"
	postgresTag = "14.2-alpine"
	dbName      = "final"
	dbUser      = "postgres"
	dbPassword  = "postgres"
	dbPort      = 5434

	mongoImg  = "mongo"
	mongoTag  = "6-jammy"
	mongoPort = 27018

	rabbitImg  = "rabbitmq"
	rabbitTag  = "3.13-management-alpine"
	rabbitPort = 5674
)

func StartPG() (*dockertest.Pool, *dockertest.Resource) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("cannot connect to docker: %s", err)
	}

	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: postgresImg,
			Tag:        postgresTag,
			Env: []string{
				fmt.Sprintf("POSTGRES_DB=%s", dbName),
				fmt.Sprintf("POSTGRES_USER=%s", dbUser),
				fmt.Sprintf("POSTGRES_PASSWORD=%s", dbPassword),
			},
			PortBindings: map[docker.Port][]docker.PortBinding{
				"5432/tcp": {
					{HostIP: "localhost", HostPort: fmt.Sprintf("%d/tcp", dbPort)},
				},
			},
		}, func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		},
	)
	if err != nil {
		log.Fatalf("cannot start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		url := fmt.Sprintf("postgres://%s:%s@localhost:%d/%s", dbUser, dbPassword, dbPort, dbName)
		conf, err := pgxpool.ParseConfig(url)
		if err != nil {
			log.Fatalf("cannot parse config: %s", err)
		}

		var db *pgxpool.Pool
		db, err = pgxpool.ConnectConfig(context.Background(), conf)
		if err != nil {
			return err
		}
		defer db.Close()

		err = db.Ping(context.Background())
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Fatalf("cannot not connect to container: %s", err)
	}

	return pool, resource
}

func StartMongo() (*dockertest.Pool, *dockertest.Resource) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("cannot connect to docker: %s", err)
	}

	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: mongoImg,
			Tag:        mongoTag,
			PortBindings: map[docker.Port][]docker.PortBinding{
				"27017/tcp": {
					{HostIP: "localhost", HostPort: fmt.Sprintf("%d/tcp", mongoPort)},
				},
			},
		}, func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		},
	)
	if err != nil {
		log.Fatalf("cannot start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		mongoURL := fmt.Sprintf("mongodb://localhost:%s", resource.GetPort("27017/tcp"))
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURL))
		if err != nil {
			return err
		}
		defer client.Disconnect(context.Background())

		err = client.Ping(context.Background(), nil)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Fatalf("cannot not connect to container: %s", err)
	}

	return pool, resource
}

func StartRabbit() (*dockertest.Pool, *dockertest.Resource) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: rabbitImg,
			Tag:        rabbitTag,
			PortBindings: map[docker.Port][]docker.PortBinding{
				"5672/tcp": {
					{HostIP: "localhost", HostPort: fmt.Sprintf("%d/tcp", rabbitPort)},
				},
			},
		}, func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		},
	)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		rabbitURL := fmt.Sprintf("amqp://guest:guest@localhost:%s/", resource.GetPort("5672/tcp"))
		conn, err := amqp.Dial(rabbitURL)
		if err != nil {
			return err
		}
		defer conn.Close()

		ch, err := conn.Channel()
		if err != nil {
			return err
		}
		defer ch.Close()

		return nil
	}); err != nil {
		log.Fatalf("cannot not connect to container: %s", err)
	}

	return pool, resource
}

func Stop(pool *dockertest.Pool, resource *dockertest.Resource) {
	if err := pool.Purge(resource); err != nil {
		fmt.Printf("Could not purge resource: %s\n", err)
	}

	fmt.Printf("Purge resource: %s\n", "OK")
}
