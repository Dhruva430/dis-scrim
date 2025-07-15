package main

import (
	"fmt"
	"reflect"
)

type MyStruct struct {
	Name string
}

func main() {
	s := MyStruct{}
	EditValue(&s)
	fmt.Println(s)
}

func EditValue(s any) {
	val := reflect.ValueOf(s).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		// fieldVal := val.Field(i)
		fieldType := typ.Field(i)

		fmt.Printf("Field %d: %s, Type: %s, Tag: %s\n",
			i,
			fieldType.Name,
			fieldType.Type,
			fieldType.Tag.Get("json"))
	}
}
