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
		switch option.Type {
		case discordgo.ApplicationCommandOptionUser:
			userValue := getOption.UserValue(session)
			fmt.Printf("User value for %s: %v\n", option.Name, userValue)
			if userValue != nil {
				fieldName := attributeMap[option.Name]
				fmt.Printf("Setting field %s to user %s\n", fieldName, userValue.Username)
				reflectField := reflectValue.FieldByName(fieldName)
				if reflectField.IsValid() && reflectField.CanSet() {
					if reflectField.Kind() == reflect.String {
						reflectField.SetString(userValue.Username)
					} else if reflectField.Kind() == reflect.Struct && reflectField.Type().Name() == "User" {
						reflectField.Set(reflect.ValueOf(*userValue))
					} else {
						fmt.Printf("Field %s is not a string or User struct, cannot set value\n", fieldName)
					}
				} else {
					fmt.Printf("Field %s is not valid or cannot be set\n", fieldName)
				}
			}
		case discordgo.ApplicationCommandOptionString:
			stringValue := getOption.StringValue()
			fmt.Printf("String value for %s: %s\n", option.Name, stringValue)
			if stringValue != "" {
				fieldName := attributeMap[option.Name]
				fmt.Printf("Setting field %s to value %s\n", fieldName, stringValue)
				reflectField := reflectValue.FieldByName(fieldName)
				if reflectField.IsValid() && reflectField.CanSet() {
					if !option.Required {
						reflectField.Set(reflect.ValueOf(&stringValue))
					} else {
						reflectField.SetString(stringValue)
					}
				} else {
					fmt.Printf("Field %s is not valid or cannot be set\n", fieldName)
				}
			}
		}
	}

	fmt.Printf("Options after parsing: %v\n", options)
}

func SetField(options any, fieldName string, userValue any) error {
	v := reflect.ValueOf(options)

	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("options must be a non-nil pointer to a struct")
	}

	v = v.Elem()

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("options must point to a struct")
	}

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("no such field: %s", fieldName)
	}
	if !field.CanSet() {
		return fmt.Errorf("cannot set field: %s", fieldName)
	}

	val := reflect.ValueOf(userValue)

	if val.Type().AssignableTo(field.Type()) {
		field.Set(val)
	} else if val.Type().ConvertibleTo(field.Type()) {
		field.Set(val.Convert(field.Type()))
	} else {
		return fmt.Errorf("cannot assign %v to field %s (expected %v)", val.Type(), fieldName, field.Type())
	}

	return nil
}
