package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SHA65536/TimezoneBot/database"
	"github.com/SHA65536/TimezoneBot/discord"

	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "TimezoneBot",
		Usage: "A custom made discord bot.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "db-host",
				Usage:    "Database host",
				EnvVars:  []string{"DB_HOST"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "db-port",
				Usage:    "Database port",
				EnvVars:  []string{"DB_PORT"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "db-user",
				Usage:    "Database user",
				EnvVars:  []string{"DB_USER"},
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
			&cli.StringFlag{
				Name:     "dc-token",
				Usage:    "Database name",
				EnvVars:  []string{"DC_TOKEN"},
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

			db, err := database.MakeDatabase(db_cfg)
			if err != nil {
				return fmt.Errorf("error creating db: %w", err)
			}

			srv, err := discord.MakeDiscordServer(c.String("dc-token"), db)
			if err != nil {
				return fmt.Errorf("error creating srv: %w", err)
			}

			if err := srv.Start(); err != nil {
				return fmt.Errorf("error starting srv: %w", err)
			}

			stop := make(chan os.Signal, 1)
			signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
			<-stop

			return srv.Stop()
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
