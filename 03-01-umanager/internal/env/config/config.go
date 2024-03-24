package config

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type Config struct {
	UsersDB UsersDBConfig `env:",prefix=USERS_DB_"`
	LinksDB LinksDBConfig `env:",prefix=LINKS_DB_"`
}

type LinksDBConfig struct {
	MongoConfig
}

type MongoConfig struct {
	Name           string        `env:"NAME,default=links"`
	Host           string        `env:"HOST,default=127.0.0.1"`
	Port           int           `env:"PORT,default=27017"`
	User           string        `env:"USER,default=mongo"`
	Password       string        `env:"PASSWORD,default=mongo"`
	MinPoolSize    uint64        `env:"MIN_POOL_SIZE,default=5"`
	MaxPoolSize    uint64        `env:"MAX_POOL_SIZE,default=50"`
	ConnectTimeout time.Duration `env:"CONNECT_TIMEOUT,default=5s"`
}

func (m MongoConfig) ConnectionString() string {
	return fmt.Sprintf("mongodb://%s:%d", m.Host, m.Port)
}

type UsersDBConfig struct {
	PostgresConfig
}

type PostgresConfig struct {
	Name         string        `env:"NAME,default=users" json:",omitempty"`
	User         string        `env:"USER,default=postgres" json:",omitempty"`
	Host         string        `env:"HOST,default=localhost" json:",omitempty"`
	Port         int           `env:"PORT,default=4321" json:",omitempty"`
	SSLMode      string        `env:"SSLMODE,default=disable" json:",omitempty"`
	ConnTimeout  int           `env:"CONN_TIMEOUT,default=5" json:",omitempty"`
	Password     string        `env:"PASSWORD,default=postgres" json:"-"`
	PoolMinConns int           `env:"POOL_MIN_CONNS,default=10" json:",omitempty"`
	PoolMaxConns int           `env:"POOL_MAX_CONNS,default=50" json:",omitempty"`
	DBTimeout    time.Duration `env:"TIMEOUT,default=5s"`
}

func (c PostgresConfig) ConnectionURL() string {
	host := c.Host
	if v := c.Port; v != 0 {
		host = host + ":" + strconv.Itoa(c.Port)
	}

	u := &url.URL{
		Scheme: "postgres",
		Host:   host,
		Path:   c.Name,
	}

	if c.User != "" || c.Password != "" {
		u.User = url.UserPassword(c.User, c.Password)
	}

	q := u.Query()
	if v := c.ConnTimeout; v > 0 {
		q.Add("connect_timeout", strconv.Itoa(v))
	}
	if v := c.SSLMode; v != "" {
		q.Add("sslmode", v)
	}

	u.RawQuery = q.Encode()

	return u.String()
}
