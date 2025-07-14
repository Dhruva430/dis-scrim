package handlers

import (
	"github.com/Dhruva430/dis-scrim.git/bot/commands"
	"github.com/bwmarrin/discordgo"
)

func HandleInteractions(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		switch i.ApplicationCommandData().Name {
		case "ping":
			commands.PingCommand(s, i)
		}
	}
}
