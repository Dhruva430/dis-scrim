package commands

import (
	"fmt"
	"reflect"

	"github.com/bwmarrin/discordgo"
)

type InferResult struct {
	optionType discordgo.ApplicationCommandOptionType
	optional   bool
}

func inferOptionType(t reflect.Type) (*InferResult, error) {
	optional := false
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		optional = true
	}
	switch t.Kind() {
	case reflect.String:
		return &InferResult{
			optionType: discordgo.ApplicationCommandOptionString,
			optional:   optional,
		}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &InferResult{
			optionType: discordgo.ApplicationCommandOptionInteger,
			optional:   optional,
		}, nil
	case reflect.Bool:
		return &InferResult{
			optionType: discordgo.ApplicationCommandOptionBoolean,
			optional:   optional,
		}, nil
	case reflect.Float32, reflect.Float64:
		return &InferResult{
			optionType: discordgo.ApplicationCommandOptionNumber,
			optional:   optional,
		}, nil
	}
	if t == reflect.TypeOf(discordgo.Member{}) {
		return &InferResult{
			optionType: discordgo.ApplicationCommandOptionUser,
			optional:   optional,
		}, nil
	}
	if t == reflect.TypeOf(discordgo.Channel{}) {
		return &InferResult{
			optionType: discordgo.ApplicationCommandOptionChannel,
			optional:   optional,
		}, nil
	}
	if t == reflect.TypeOf(discordgo.Role{}) {
		return &InferResult{
			optionType: discordgo.ApplicationCommandOptionRole,
			optional:   optional,
		}, nil
	}
	if t == reflect.TypeOf(discordgo.MessageAttachment{}) {
		return &InferResult{
			optionType: discordgo.ApplicationCommandOptionAttachment,
			optional:   optional,
		}, nil
	}

	return nil, fmt.Errorf("unsupported type: %s", t)
}

type CommandBuilder struct {
	Name        string
	Description string
	Options     []*discordgo.ApplicationCommandOption
	Handler     func(session *discordgo.Session, interaction *discordgo.InteractionCreate)
}

func Builder() *CommandBuilder {
	return &CommandBuilder{}
}

func (b *CommandBuilder) SetName(name string) *CommandBuilder {
	b.Name = name
	return b
}

func (b *CommandBuilder) SetDescription(description string) *CommandBuilder {
	b.Description = description
	return b
}

func (b *CommandBuilder) SetOptions(options any) *CommandBuilder {
	if options == nil {
		b.Options = nil
		return b
	}
	opts := InferOptions(options)
	b.Options = opts
	return b
}

func (b *CommandBuilder) Build() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        b.Name,
		Description: b.Description,
		Options:     b.Options,
	}
}

func (b *CommandBuilder) SetHandler(handler func(session *discordgo.Session, interaction *discordgo.InteractionCreate)) *CommandBuilder {
	if handler == nil {
		panic("Command handler cannot be nil")
	}

	b.Handler = handler
	return b
}
