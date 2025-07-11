package discord

import (
	"context"
	"fmt"
	"time"

	"github.com/SHA65536/TimezoneBot/database"
	"github.com/SHA65536/TimezoneBot/parser"
	"github.com/bwmarrin/discordgo"
)

func RegisterConvertHandler(s *discordgo.Session, db *database.Queries) error {
	var tp = parser.NewTimeParser()

	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		fmt.Println(m.Content)

		// Try to parse time from the message
		if _, err := tp.ParseTimeFromMessage(m.Content); err != nil {
			return
		}

		// Check if owner has a timezone set
		if _, err := db.GetTimezone(context.Background(), m.Author.ID); err != nil {
			return
		}

		fmt.Println(s.MessageReactionAdd(m.ChannelID, m.ID, "⏰"))
	})

	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
		if m.Emoji.Name != "⏰" || m.UserID == s.State.User.ID {
			return
		}

		// Fetch the original message
		msg, err := s.ChannelMessage(m.ChannelID, m.MessageID)
		if err != nil {
			return
		}

		// Check if the original message author has a timezone set
		userTimezone, err := db.GetTimezone(context.Background(), msg.Author.ID)
		if err != nil {
			return
		}

		// Try to parse time from the original message content
		seconds, err := tp.ParseTimeFromMessage(msg.Content)
		if err != nil {
			return
		}

		// Load the user's timezone
		userLoc, err := time.LoadLocation(userTimezone)
		if err != nil {
			return
		}

		// Create the time in the user's timezone context
		now := time.Now().In(userLoc)
		parsedTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, userLoc).Add(time.Duration(seconds) * time.Second)

		// Convert to UTC for Discord timestamp
		utcTime := parsedTime.UTC()

		// Format the time message using Discord timestamp format
		unixTimestamp := utcTime.Unix()
		timeMessage := fmt.Sprintf("<t:%d:t>", unixTimestamp)

		_, _ = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: timeMessage,
			Reference: &discordgo.MessageReference{
				MessageID: msg.ID,
				ChannelID: msg.ChannelID,
				GuildID:   msg.GuildID,
			},
			AllowedMentions: &discordgo.MessageAllowedMentions{
				Parse: []discordgo.AllowedMentionType{},
			},
		})
	})

	return nil
}
