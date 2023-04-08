package utils

import "reflect"

func GetStructTag(obj any, tagName string) string {
	field, ok := reflect.TypeOf(obj).Elem().FieldByName(tagName)
	if !ok {
		return ""
	}
	return string(field.Tag)
}
