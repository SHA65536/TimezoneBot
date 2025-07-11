package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

const CONN_STR = "postgres://%s:%s@%s:%s/%s"

type DatabaseConfig struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

func MakeDatabase(cfg DatabaseConfig) (*Queries, error) {
	connStr := fmt.Sprintf(CONN_STR, cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name)

	pool, err := pgxpool.New(context.TODO(), connStr)
	if err != nil {
		return nil, err
	}

	return New(pool), nil
}
