package discord

import (
	"context"
	"fmt"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/SHA65536/TimezoneBot/database"
	"github.com/SHA65536/TimezoneBot/timezones"
	"github.com/bwmarrin/discordgo"
)

// RegisterTimezoneCommand registers the /timezone slash command and its handler
func RegisterTimezoneCommand(s *discordgo.Session, db *database.Queries) error {
	command := &discordgo.ApplicationCommand{
		Name:        "timezone",
		Description: "Set your timezone",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "location",
				Description:  "Pick your timezone (IANA)",
				Required:     true,
				Autocomplete: true,
			},
		},
	}

	_, err := s.ApplicationCommandCreate(s.State.User.ID, "", command)
	if err != nil {
		return fmt.Errorf("cannot create slash command: %w", err)
	}

	// Handle autocomplete
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommandAutocomplete {
			return
		}
		if i.ApplicationCommandData().Name != "timezone" {
			return
		}

		var userInput string
		for _, opt := range i.ApplicationCommandData().Options {
			if opt.Name == "location" {
				userInput = opt.StringValue()
				break
			}
		}

		choices := getAutocompleteChoices(userInput)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: choices,
			},
		})
	})

	// Handle command execution
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}
		if i.ApplicationCommandData().Name != "timezone" {
			return
		}

		location := ""
		for _, opt := range i.ApplicationCommandData().Options {
			if opt.Name == "location" {
				location = opt.StringValue()
				break
			}
		}

		if location == "" {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "No timezone selected.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}

		// Validate that the timezone is valid
		if _, err := time.LoadLocation(location); err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Invalid timezone selected.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}

		err := db.SetTimezone(context.Background(), database.SetTimezoneParams{
			UserID:   i.Member.User.ID,
			Timezone: location,
		})
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Failed to save timezone.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Timezone set to %s", location),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	})

	return nil
}

// getAutocompleteChoices returns autocomplete choices based on user input
func getAutocompleteChoices(userInput string) []*discordgo.ApplicationCommandOptionChoice {
	timezoneNames := timezones.TimezoneLocations

	var choices []*discordgo.ApplicationCommandOptionChoice
	userInputLower := strings.ToLower(userInput)

	// Filter timezones that match the user input
	for _, tz := range timezoneNames {
		if strings.Contains(strings.ToLower(tz), userInputLower) {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  tz,
				Value: tz,
			})
			// Discord limits autocomplete results to 25 choices
			if len(choices) >= 25 {
				break
			}
		}
	}

	return choices
}
