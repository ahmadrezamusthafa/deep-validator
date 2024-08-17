package utils

import "reflect"

func StructsToMap(attributeNames map[string]interface{}, data interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	switch val := data.(type) {
	case []interface{}:
		for _, item := range val {
			nestedMap := StructsToMap(attributeNames, item)
			for k, v := range nestedMap {
				result[k] = v
			}
		}
	case interface{}:
		rValue := reflect.ValueOf(data)
		if rValue.Kind() == reflect.Ptr {
			if rValue.IsNil() {
				return result
			}
			rValue = rValue.Elem()
		}
		for i := 0; i < rValue.NumField(); i++ {
			field := rValue.Field(i)
			typeField := rValue.Type().Field(i)
			key := typeField.Name

			if field.Kind() == reflect.Struct {
				nestedMap := StructsToMap(attributeNames, field.Interface())
				for k, v := range nestedMap {
					result[k] = v
				}
			}
			if _, ok := attributeNames[key]; !ok {
				continue
			}
			if field.CanInterface() {
				result[key] = field.Interface()
			}
		}
	}
	return result
}
