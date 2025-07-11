package main

import (
	"log"
	"os"

	"github.com/SHA65536/TimezoneBot/database"

	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "Migrations",
		Usage: "Run migrations for the turtlemancer app",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "db-host",
				Usage:    "Database host",
				EnvVars:  []string{"DB_HOST"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "db-user",
				Usage:    "Database user",
				EnvVars:  []string{"DB_USER"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "db-port",
				Usage:    "Database port",
				EnvVars:  []string{"DB_PORT"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "db-pass",
				Usage:    "Database password",
				EnvVars:  []string{"DB_PASS"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "db-name",
				Usage:    "Database name",
				EnvVars:  []string{"DB_NAME"},
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			var db_cfg = database.DatabaseConfig{
				Host: c.String("db-host"),
				Port: c.String("db-port"),
				User: c.String("db-user"),
				Pass: c.String("db-pass"),
				Name: c.String("db-name"),
			}

			return database.RunMigrations(db_cfg)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
