package freak_conventer

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func ConvertToUnixTime(v interface{}) (map[string]interface{}, error) {
	if v == nil {
		return nil, fmt.Errorf("input data is nil")
	}

	val := reflect.ValueOf(v)

	if val.IsZero() {
		return nil, fmt.Errorf("input data is zero")
	}

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input data should be a struct")
	}

	return GetStructWithUnixTime(val), nil
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

		if !valueField.IsValid() {
			continue
		}

		jsonTag = strings.Split(jsonTag, ",")[0]

		switch {
		case valueField.Type() == reflect.TypeOf(time.Time{}):
			data[jsonTag] = valueField.Interface().(time.Time).Unix()

		case valueField.Kind() == reflect.Slice:
			data[jsonTag] = GetArrayWithUnixTime(valueField)

		case valueField.Kind() == reflect.Struct:
			if field.Anonymous {
				for k, v := range GetStructWithUnixTime(valueField) {
					data[k] = v
				}
			} else {
				structField := GetStructWithUnixTime(valueField)
				data[jsonTag] = structField
			}

		default:
			data[jsonTag] = valueField.Interface()
		}
	}

	return data
}

func GetArrayWithUnixTime(value reflect.Value) []interface{} {
	sliceData := make([]interface{}, value.Len())

	for j := 0; j < value.Len(); j++ {
		sliceElem := value.Index(j)

		switch {
		case sliceElem.Type() == reflect.TypeOf(time.Time{}):
			sliceData[j] = sliceElem.Interface().(time.Time).Unix()

		case sliceElem.Kind() == reflect.Struct:
			sliceData[j] = GetStructWithUnixTime(sliceElem)

		default:
			sliceData[j] = sliceElem.Interface()
		}
	}

	return sliceData
}

func IsOmitEmpty(field reflect.StructField) bool {
	jsonTag := field.Tag.Get("json")
	return strings.Contains(jsonTag, "omitempty")
}
