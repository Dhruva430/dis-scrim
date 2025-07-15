package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type CommandHandler struct {
	commands []*CommandBuilder
}

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		commands: make([]*CommandBuilder, 0),
	}
}

func (h *CommandHandler) ExecuteCommands(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Type != discordgo.InteractionApplicationCommand {
		return
	}
	name := interaction.ApplicationCommandData().Name
	for _, cmd := range h.commands {
		if cmd.Name == name {
			if cmd.Handler == nil {
				return
			}
			cmd.Handler(session, interaction)
		}
	}
}

func (h *CommandHandler) AddCommand(cmd *CommandBuilder) {
	h.commands = append(h.commands, cmd)
}

func (h *CommandHandler) RegisterCommands(session *discordgo.Session, appID string) error {
	for _, cmd := range h.commands {
		if _, err := session.ApplicationCommandCreate(appID, "", cmd.Build()); err != nil {
			return err
		}
		fmt.Printf("Command '%s' registered successfully.\n", cmd.Name)
	}
	session.AddHandler(h.ExecuteCommands)
	return nil
}
