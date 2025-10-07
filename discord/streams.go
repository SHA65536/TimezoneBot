package discord

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
	_ "time/tzdata"

	"github.com/SHA65536/TimezoneBot/database"
	"github.com/bwmarrin/discordgo"
)

// RegisterStreamsAlertCommand registers the /streamalert slash command and its handler
func RegisterStreamsAlertCommand(s *discordgo.Session, db *database.Queries) error {
	command := &discordgo.ApplicationCommand{
		Name:        "streamalert",
		Description: "Set the alert channel for stream alerts",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "channel",
				Description: "Pick the channel to alert when a stream starts",
				Required:    true,
			},
		},
	}

	_, err := s.ApplicationCommandCreate(s.State.User.ID, "", command)
	if err != nil {
		return fmt.Errorf("cannot create slash command: %w", err)
	}

	// Handle command execution
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}
		if i.ApplicationCommandData().Name != "streamalert" {
			return
		}

		member, err := s.GuildMember(i.GuildID, i.Member.User.ID)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Could not verify your permissions.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}
		perms, err := s.State.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Could not verify your permissions.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}
		if perms&discordgo.PermissionAdministrator != 0 {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You must be an administrator to use this command.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}

		var channelID string
		for _, opt := range i.ApplicationCommandData().Options {
			if opt.Name == "channel" {
				channelID = opt.StringValue()
				break
			}
		}

		err = db.SetAlertChannel(context.Background(), database.SetAlertChannelParams{
			GuildID:      i.GuildID,
			AlertChannel: channelID,
		})
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Failed to save stream alert channel.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Stream alerts will now be sent",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	})

	return nil
}

var streamCooldownTable = map[string]time.Time{}
var streamCooldownLock sync.RWMutex

func RegisterStreamsHandler(s *discordgo.Session, db *database.Queries) error {
	s.AddHandler(func(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {

		if rand.Intn(1000) == 0 {
			cleanstreamcooldown()
		}

		if !vs.VoiceState.SelfStream {
			return
		}

		if b, err := db.IsStreamAlertEnabled(context.Background(), vs.GuildID); err != nil || !b {
			return
		}

		streamCooldownLock.RLock()
		if time.Since(streamCooldownTable[vs.UserID]) < 10*time.Minute {
			streamCooldownLock.RUnlock()
			return
		}
		streamCooldownLock.RUnlock()

		alertChannel, err := db.GetAlertChannel(context.Background(), vs.GuildID)
		if err != nil {
			fmt.Println("Error getting alert channel:", err)
			return
		}

		_, err = s.ChannelMessageSendComplex(alertChannel, &discordgo.MessageSend{
			Content: "OMG A STREAM HAS STARTED",
			Embed: &discordgo.MessageEmbed{
				Title:       "Stream started",
				Description: "Stream started in <#" + vs.ChannelID + "> by " + vs.Member.Nick,
				Color:       0x0099ff,
			},
		})
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}

		streamCooldownLock.Lock()
		streamCooldownTable[vs.UserID] = time.Now()
		streamCooldownLock.Unlock()
	})
	return nil
}

func cleanstreamcooldown() {
	var newTable = map[string]time.Time{}
	streamCooldownLock.Lock()
	defer streamCooldownLock.Unlock()
	for id, t := range streamCooldownTable {
		if time.Since(t) < 10*time.Minute {
			newTable[id] = t
		}
	}
}
