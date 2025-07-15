package main

import (
	"fmt"
	"reflect"
)

type MyStruct struct {
	Name string `json:"name"`
}

type OwnerStruct struct {
	internal any
}

func NewOwnerStruct() *OwnerStruct {
	return &OwnerStruct{
		internal: MyStruct{},
	}
}

func main() {
	owner := NewOwnerStruct()
	original := reflect.ValueOf(owner.internal)
	ptrToCopy := reflect.New(original.Type())
	randomFunction(ptrToCopy.Interface())
	final := ptrToCopy.Elem()
	fmt.Println("Edited value:", final)
}

func randomFunction(t any) {
	EditValue(t)
}

func EditValue(s any) {
	val := reflect.ValueOf(s).Elem()
	typ := val.Type()

	if typ.Kind() != reflect.Struct {
		panic("s must be a pointer to a struct")
	}

	for i := 0; i < val.NumField(); i++ {
		// fieldVal := val.Field(i)
		fieldType := typ.Field(i)

		fmt.Printf("Field %d: %s, Type: %s, Tag: %s\n",
			i,
			fieldType.Name,
			fieldType.Type,
			fieldType.Tag.Get("json"))
		if fieldType.Name == "Name" {
			val.Field(i).Set(reflect.ValueOf("New Name"))
		}
	}
}
