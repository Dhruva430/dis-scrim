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

func ParseOptions(data discordgo.ApplicationCommandInteractionData, session *discordgo.Session, options any, cmd CommandBuilder) {
	attributeMap := make(map[string]string)

	fmt.Println("Parsing options for command:", cmd)
	reflectOptions := reflect.ValueOf(options)
	if reflectOptions.Kind() != reflect.Ptr || reflectOptions.IsNil() {
		panic("HandlerOptions must be a non-nil pointer")
	}

	reflectOptions = reflectOptions.Elem()
	reflectOptionsType := reflectOptions.Type()
	// if reflectOptionsType.Kind() != reflect.Struct {
	// 	panic("HandlerOptions must be a pointer to a struct")
	// }

	for i := 0; i < reflectOptions.NumField(); i++ {
		field := reflectOptionsType.Field(i)
		tag := field.Tag.Get("name")
		if tag == "" {
			tag = strings.ToLower(field.Name)
		}
		attributeMap[tag] = field.Name
	}
	fmt.Println("Attribute Map:", attributeMap)

	for _, option := range cmd.Options {
		fmt.Println("Processing option:", option.Name)
		getOption := data.GetOption(option.Name)
		if getOption == nil {

			fmt.Printf("Option %s not found in interaction data\n", option.Name)
			continue
		}
		switch option.Type {
		case discordgo.ApplicationCommandOptionUser:
			userValue := getOption.UserValue(session)

			if userValue != nil {
				fmt.Println("User Value:", userValue.Username)
				fmt.Println("Setting field:", attributeMap[option.Name])
				fmt.Println("Value:", reflect.ValueOf(userValue))
				fmt.Println("Options before setting:", options)
				fmt.Println("Options type:", reflect.TypeOf(options))
				fieldName := attributeMap[option.Name]
				fmt.Println("Setting field:", fieldName)

				reflect.ValueOf(options).Elem().FieldByName(fieldName).Set(reflect.ValueOf(userValue))

			}

		case discordgo.ApplicationCommandOptionChannel:
			channelValue := getOption.ChannelValue(session)
			if channelValue != nil {
				reflect.ValueOf(options).Elem().FieldByName(attributeMap[option.Name]).Set(reflect.ValueOf(channelValue))
			}

		}
	}
}
