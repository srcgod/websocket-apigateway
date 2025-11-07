package utils

import (
	"fmt"
	"reflect"
	"strconv"

	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func CopyToProto(dto any, proto any) {
	vDto := reflect.ValueOf(dto)
	if vDto.Kind() == reflect.Ptr {
		vDto = vDto.Elem()
	}
	vProto := reflect.ValueOf(proto)
	if vProto.Kind() == reflect.Ptr {
		vProto = vProto.Elem()
	}
	tDto := vDto.Type()

	for i := 0; i < vDto.NumField(); i++ {
		fieldDto := vDto.Field(i)
		if fieldDto.Kind() == reflect.Ptr && !fieldDto.IsNil() {
			protoField := vProto.FieldByName(tDto.Field(i).Name)
			if protoField.IsValid() && protoField.CanSet() {
				val := fieldDto.Elem().String()

				if protoField.Kind() == reflect.Ptr {
					protoField.Set(reflect.ValueOf(&val))
				} else {
					protoField.Set(reflect.ValueOf(val))
				}
			}
		}
	}
}

func StringToFieldMask(maskStr []string) *fieldmaskpb.FieldMask {
	return &fieldmaskpb.FieldMask{
		Paths: maskStr,
	}
}

func ConvertToInt64(n any) int64 {
	Int64, ok := n.(int64)
	if !ok {
		return 0
	}
	return Int64
}

func ConvertToInt32(value string) (int32, error) {
	if value == "" {
		return 0, nil
	}
	i64, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("could not parse int32: %w", err)
	}
	if i64 > (1<<31-1) || i64 < -(1<<31) {
		return 0, fmt.Errorf("value out of range for int32: %d", i64)
	}
	return int32(i64), nil
}
