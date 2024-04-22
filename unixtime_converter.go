package freak_conventer

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func ConvertToUnixTime(v interface{}) (map[string]interface{}, error) {
	val := reflect.ValueOf(v)

	if val.IsZero() {
		return nil, fmt.Errorf("input data is nil")
	}

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input data should be a struct")
	}

	data := make(map[string]interface{})

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)

		if !field.IsExported() {
			continue
		}

		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		valueField := val.Field(i)

		if valueField.Kind() == reflect.Ptr {
			if IsOmitEmpty(field) && val.Field(i).IsNil() {
				continue
			}
			valueField = valueField.Elem()
		} else {
			if IsOmitEmpty(field) && val.Field(i).IsZero() {
				continue
			}
		}

		if valueField.Kind() == reflect.Struct {
			if field.Anonymous {
				for k, v := range GetStructWithUnixTime(valueField) {
					data[k] = v
				}
			} else {
				structField := GetStructWithUnixTime(valueField)
				data[jsonTag] = structField
			}
		} else if valueField.Type() == reflect.TypeOf(time.Time{}) {
			data[jsonTag] = valueField.Interface().(time.Time).Unix()
		} else {
			data[jsonTag] = valueField.Interface()
		}
	}

	return data, nil
}

func GetStructWithUnixTime(value reflect.Value) map[string]interface{} {
	data := make(map[string]interface{})

	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)

		if !field.IsExported() {
			continue
		}

		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		valueField := value.Field(i)

		if valueField.Kind() == reflect.Ptr {
			if IsOmitEmpty(field) && value.Field(i).IsNil() {
				continue
			}
			valueField = valueField.Elem()
		} else {
			if IsOmitEmpty(field) && value.Field(i).IsZero() {
				continue
			}
		}

		if valueField.Kind() == reflect.Struct {
			if field.Type.Name() == "" {
				for k, v := range GetStructWithUnixTime(valueField) {
					data[k] = v
				}
			} else {
				structField := GetStructWithUnixTime(valueField)
				data[jsonTag] = structField
			}
		} else if valueField.Type() == reflect.TypeOf(time.Time{}) {
			data[jsonTag] = valueField.Interface().(time.Time).Unix()
		} else {
			data[jsonTag] = valueField.Interface()
		}
	}

	return data
}

func IsOmitEmpty(field reflect.StructField) bool {
	jsonTag := field.Tag.Get("json")
	return strings.Contains(jsonTag, "omitempty")
}
