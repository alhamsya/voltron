package util

import (
	"reflect"
	"strings"
)

func GetFieldName(req any, structField, key string) string {
	typeOf := reflect.TypeOf(req)

	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}

	if field, ok := typeOf.FieldByName(structField); ok {
		tag := field.Tag.Get(key)
		if tag == "" {
			return structField
		}
		return strings.Split(tag, ",")[0]
	}

	return structField
}
