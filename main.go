package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"dis-scrim/commands"
	"dis-scrim/configs"

	"github.com/bwmarrin/discordgo"
)

type BanCommand struct {
	Member     discordgo.User               `name:"member" desc:"Member to ban"`
	Reason     *string                      `name:"reason" desc:"Reason for the ban"`
	Role       *discordgo.Role              `name:"role" desc:"Role to assign after ban"`
	Channel    *discordgo.Channel           `name:"channel" desc:"Channel to notify after ban"`
	Boolean    *bool                        `name:"boolean" desc:"A boolean option"`
	Attachment *discordgo.MessageAttachment `name:"attachment" desc:"Attachment to include in the ban message"`
}

type SuggestionCommand struct {
	Title string `name:"title" desc:"Title of the suggestion"`
}

func Dereference(v any) reflect.Value {
	val := reflect.ValueOf(v)

	for (val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface) && !val.IsNil() {
		val = val.Elem()
	}
	return val
}

func main() {
	commandHandler := commands.NewCommandHandler()
	pingCommand := commands.Builder().
		SetName("ping").
		SetDescription("Ping the bot").
		SetHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate, options any) {
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
		SetHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate, optionsData any) {
			options, ok := optionsData.(BanCommand)
			if !ok {
				log.Printf("Failed to cast optionsData to BanCommand: %v", optionsData)
				return
			}
			// Contruct the response message
			if options.Reason == nil {
				reason := "No reason provided"
				options.Reason = &reason
			}
			response := &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Banned %s for reason: %v", options.Member.Username, *options.Reason),
				},
			}
			if err := session.InteractionRespond(interaction.Interaction, response); err != nil {
				log.Printf("Failed to respond to interaction: %v", err)
				return
			}
		})

	commandHandler.AddCommand(pingCommand)
	commandHandler.AddCommand(banCommand)
	session, err := CreateBotSession()
	if err != nil {
		log.Fatalf("Failed to create bot session: %v", err)
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

func CreateBotSession() (*discordgo.Session, error) {
	session, err := discordgo.New("Bot " + configs.GetToken())
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %w", err)
	}
	if err := session.Open(); err != nil {
		return nil, fmt.Errorf("error opening Discord session: %w", err)
	}
	return session, nil
}
