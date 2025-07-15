package bot

import (
	"dis-scrim/bot/commands"
	"dis-scrim/bot/handlers"
	"log"

	"github.com/bwmarrin/discordgo"
)

var Session *discordgo.Session

func StartBot(token string) {
	var err error
	Session, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Failed to create Discord session: %v", err)
	}

	Session.AddHandler(handlers.HandleInteractions)

	err = Session.Open()
	if err != nil {
		log.Fatalf("Failed to open session: %v", err)
	}

	app, err := Session.Application("@me")
	if err != nil {
		log.Fatalf("Failed to get app info: %v", err)
	}

	err = commands.RegisterCommands(Session, app.ID)
	if err != nil {
		log.Fatalf("Failed to register commands: %v", err)
	}

	log.Println("Bot is running...")
}
