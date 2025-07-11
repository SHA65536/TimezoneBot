package database

import (
	"embed"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed sql/schemas/*.sql
var embedMigrations embed.FS

func RunMigrations(cfg DatabaseConfig) error {
	goose.SetBaseFS(embedMigrations)
	gooseConnStr := fmt.Sprintf(CONN_STR, cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name)

	db, err := goose.OpenDBWithDriver("postgres", gooseConnStr)
	if err != nil {
		return fmt.Errorf("failed to open database for migrations: %w", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set Goose dialect: %w", err)
	}

	if err := goose.Up(db, "sql/schemas"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
