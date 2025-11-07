package utils

import (
	"reflect"
	"strings"
)

func GetUpdateMask(dto any) []string {
	v := reflect.ValueOf(dto)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	var updateMask []string

	for i := 0; i < v.NumField(); i++ {
		fieldVal := v.Field(i)
		fieldType := t.Field(i)

		if fieldVal.Kind() == reflect.Ptr {
			if !fieldVal.IsNil() {
				jsonTag := fieldType.Tag.Get("json")
				if jsonTag != "" {
					fieldName := strings.Split(jsonTag, ",")[0]
					updateMask = append(updateMask, fieldName)
				} else {
					updateMask = append(updateMask, fieldType.Type.Name())
				}
			} else {
				zeroValue := reflect.Zero(fieldVal.Type()).Interface()
				if fieldVal.Interface() != zeroValue {
					jsonTag := fieldType.Tag.Get("json")
					if jsonTag == "" {
						fieldName := strings.Split(jsonTag, ",")[0]
						updateMask = append(updateMask, fieldName)
					} else {
						updateMask = append(updateMask, fieldType.Type.Name())
					}
				}
			}
		}
	}
	return updateMask
}
