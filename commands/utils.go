package commands

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func InferOptions(optionStruct any) []*discordgo.ApplicationCommandOption {
	cmd := reflect.TypeOf(optionStruct)
	if cmd.Kind() != reflect.Struct {
		panic("BanCommand must be a struct")
	}
	options := make([]*discordgo.ApplicationCommandOption, 0, cmd.NumField())
	for i := 0; i < cmd.NumField(); i++ {
		field := cmd.Field(i)
		opt, err := inferOptionType(field.Type)
		if err != nil {
			panic(fmt.Sprintf("failed to infer option type for field %s: %v", field.Name, err))
		}
		name := field.Tag.Get("name")
		if name == "" {
			name = strings.ToLower(field.Name)
		}
		option := &discordgo.ApplicationCommandOption{
			Name:        name,
			Description: field.Tag.Get("desc"),
			Type:        opt.optionType,
			Required:    !opt.optional,
		}
		options = append(options, option)
	}
	return options
}
