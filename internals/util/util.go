package util

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

func PrettyPrint(data any, message string) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal data: %v", err)
	}
	fmt.Printf("%s: %s\n", message, string(jsonData))
}

func CollectFieldPtrs(v reflect.Value) []any {
	var result []any

	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			// If it's a struct, recurse
			if field.Kind() == reflect.Struct {
				result = append(result, CollectFieldPtrs(field)...)
				continue
			}
			// If it's a slice of struct, handle each element
			if field.Kind() == reflect.Slice && field.Type().Elem().Kind() == reflect.Struct {
				for j := 0; j < field.Len(); j++ {
					result = append(result, CollectFieldPtrs(field.Index(j))...)
				}
				continue
			}
			// Take addressable value
			if field.CanAddr() {
				result = append(result, field.Addr().Interface())
			}
		}
	case reflect.Slice:
		// If the value itself is a slice of struct
		if v.Type().Elem().Kind() == reflect.Struct {
			for i := 0; i < v.Len(); i++ {
				result = append(result, CollectFieldPtrs(v.Index(i))...)
			}
		}
	}

	return result
}

// func CollectFieldPtrs(v reflect.Value) []any {
// 	var result []any

// 	for i := 0; i < v.NumField(); i++ {
// 		field := v.Field(i)

// 		// If it's a struct, recurse
// 		if field.Kind() == reflect.Struct {
// 			result = append(result, CollectFieldPtrs(field)...)
// 			continue
// 		}

// 		// Take addressable value
// 		if field.CanAddr() {
// 			result = append(result, field.Addr().Interface())
// 		}
// 	}

// 	return result
// }
