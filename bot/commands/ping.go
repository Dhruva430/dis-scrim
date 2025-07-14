package commands

import (
	"github.com/bwmarrin/discordgo"
)

func RegisterCommands(s *discordgo.Session, appID string) error {
	_, err := s.ApplicationCommandCreate(appID, "", &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Replies with Pong!",
	})
	return err
}

func PingCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
		},
	})
}
