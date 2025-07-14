package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Dhruva430/dis-scrim.git/configs"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var token = configs.GetToken()

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	token = os.Getenv("TOKEN")
}

func main() {
	loadEnv()

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			switch i.ApplicationCommandData().Name {
			case "ping":
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Pong!",
					},
				})
			}
		}
	})

	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening Discord session: %v", err)
	}
	defer dg.Close()

	app, err := dg.Application("@me")
	if err != nil {
		log.Fatalf("Unable to get application info: %v", err)
	}

	_, err = dg.ApplicationCommandCreate(app.ID, "", &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Replies with Pong!",
	})
	if err != nil {
		log.Fatalf("Cannot create slash command: %v", err)
	}

	log.Println("Bot is running...")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop
}
