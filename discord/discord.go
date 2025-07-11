package discord

import (
	"fmt"

	"github.com/SHA65536/TimezoneBot/database"
	"github.com/bwmarrin/discordgo"
)

// DiscordServer wraps the Discord session and database
// Implements Start and Stop methods
// Delegates slash command handling to timezone.go

type DiscordServer struct {
	session *discordgo.Session
	db      *database.Queries
}

// MakeDiscordServer creates a new DiscordServer
func MakeDiscordServer(token string, db *database.Queries) (*DiscordServer, error) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord session: %w", err)
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions

	return &DiscordServer{
		session: dg,
		db:      db,
	}, nil
}

// Start opens the Discord session
func (s *DiscordServer) Start() error {
	if err := s.session.Open(); err != nil {
		return fmt.Errorf("error opening Discord session: %w", err)
	}

	// Register slash commands and handler
	if err := RegisterTimezoneCommand(s.session, s.db); err != nil {
		return fmt.Errorf("failed to register slash command: %w", err)
	}

	if err := RegisterConvertHandler(s.session, s.db); err != nil {
		return fmt.Errorf("failed to register convert handler: %w", err)
	}

	fmt.Println("Bot is now running.")
	return nil
}

// Stop closes the Discord session
func (s *DiscordServer) Stop() error {
	return s.session.Close()
}
