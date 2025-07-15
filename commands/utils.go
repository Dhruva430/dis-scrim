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

func GetAttributeNameMap(options any) map[string]string {
	attributeMap := make(map[string]string)
	reflectOptions := reflect.ValueOf(options)
	if reflectOptions.Kind() != reflect.Ptr || reflectOptions.IsNil() {
		panic("HandlerOptions must be a non-nil pointer")
	}
	for reflectOptions.Kind() == reflect.Interface || reflectOptions.Kind() == reflect.Ptr {
		reflectOptions = reflectOptions.Elem()
	}
	reflectOptionsType := reflectOptions.Type()

	fmt.Println("Reflecting options:", reflectOptions)
	fmt.Println("Number of fields:", reflectOptions.NumField())

	for i := 0; i < reflectOptions.NumField(); i++ {
		field := reflectOptionsType.Field(i)
		tag := field.Tag.Get("name")
		if tag == "" {
			tag = strings.ToLower(field.Name)
		}
		attributeMap[tag] = field.Name
	}
	return attributeMap
}

func ParseOptions(data discordgo.ApplicationCommandInteractionData, session *discordgo.Session, options any, cmd *CommandBuilder) {
	attributeMap := GetAttributeNameMap(options)
	fmt.Println("Attribute map:", attributeMap)
	fmt.Println("Parsing options for command:", cmd)

	if options == nil {
		fmt.Println("Options are nil, skipping parsing")
		return
	}

	reflectValue := reflect.ValueOf(options)
	if reflectValue.Kind() != reflect.Ptr || reflectValue.IsNil() {
		panic("HandlerOptions must be a non-nil pointer")
	}
	for reflectValue.Kind() == reflect.Interface || reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}

	if reflectValue.Kind() != reflect.Struct {
		panic("HandlerOptions must be a struct")
	}

	for _, option := range cmd.Options {
		fmt.Println("Processing option:", option.Name)
		getOption := data.GetOption(option.Name)
		if getOption == nil {
			fmt.Printf("Option %s not found in interaction data\n", option.Name)
			continue
		}
		field := reflectValue.FieldByName(attributeMap[option.Name])
		switch option.Type {
		case discordgo.ApplicationCommandOptionUser:
			userValue := getOption.UserValue(session)
			setFieldValue(field, *userValue, option.Required)

		case discordgo.ApplicationCommandOptionString:
			stringValue := getOption.StringValue()
			setFieldValue(field, stringValue, option.Required)
		}
	}

	fmt.Printf("Options after parsing: %v\n", options)
}

func setFieldValue(field reflect.Value, value any, required bool) {
	if !field.IsValid() || !field.CanSet() {
		return
	}

	val := reflect.ValueOf(value)

	// Handle pointer vs non-pointer destination
	switch field.Kind() {
	case reflect.Ptr:
		// If field is a pointer and not nil, ensure correct type
		if val.Type().AssignableTo(field.Type()) {
			field.Set(val)
		} else if val.Type().AssignableTo(field.Type().Elem()) {
			ptr := reflect.New(field.Type().Elem())
			ptr.Elem().Set(val)
			field.Set(ptr)
		}
	case reflect.Struct, reflect.String, reflect.Int, reflect.Bool, reflect.Float64:
		if required {
			if val.Type().AssignableTo(field.Type()) {
				field.Set(val)
			}
		} else {
			// optional struct/string: store pointer to value
			if field.Type().Kind() == reflect.Ptr && val.Type().AssignableTo(field.Type().Elem()) {
				ptr := reflect.New(field.Type().Elem())
				ptr.Elem().Set(val)
				field.Set(ptr)
			}
		}
	}
}
