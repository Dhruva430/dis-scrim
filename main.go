package main

import (
	"dis-scrim/commands"
	"dis-scrim/configs"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type BanCommand struct {
	Member discordgo.Member `name:"member" desc:"Member to ban"`
	Reason *string          `name:"reason" desc:"Reason for the ban"`
}

func main() {
	commandHandler := commands.NewCommandHandler()
	pingCommand := commands.Builder().
		SetName("ping").
		SetDescription("Ping the bot").
		SetHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
			response := &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Pong!",
				},
			}
			if err := session.InteractionRespond(interaction.Interaction, response); err != nil {
				log.Printf("Failed to respond to interaction: %v", err)
			}
		})
	banCommand := commands.Builder().
		SetName("ban").
		SetDescription("Ban a member from the server").
		SetOptions(BanCommand{}).
		SetHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
			member := interaction.ApplicationCommandData().Options[0].UserValue(session)
			reason := interaction.ApplicationCommandData().Options[1].StringValue()
			if member == nil {
				log.Println("No member specified for ban.")
				return
			}
			response := &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Banned %s for reason: %s", member.Username, reason),
				},
			}
			if err := session.InteractionRespond(interaction.Interaction, response); err != nil {
				log.Printf("Failed to respond to interaction: %v", err)
				return
			}
		})

	commandHandler.AddCommand(pingCommand)
	commandHandler.AddCommand(banCommand)
	session, err := discordgo.New("Bot " + configs.GetToken())
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	err = session.Open()
	if err != nil {
		log.Fatalf("Failed to open session: %v", err)
	}
	defer session.Close()
	app, err := session.Application("@me")
	if err != nil {
		log.Fatalf("Failed to get app info: %v", err)
	}
	if err := commandHandler.RegisterCommands(session, app.ID); err != nil {
		log.Fatalf("Failed to register commands: %v", err)
	}
	fmt.Println("Command created successfully!")
	fmt.Println("Logged in as:", app.Name)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop
	log.Println("Shutting down bot...")
}
