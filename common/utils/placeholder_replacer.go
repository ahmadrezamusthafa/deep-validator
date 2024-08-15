package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func ReplacePlaceholders(str string, attributes map[string]interface{}) string {
	for key, value := range attributes {
		placeholder := fmt.Sprintf("{{%s}}", key)
		var valueStr string
		switch v := value.(type) {
		case string:
			valueStr = v
		case int, int32, int64:
			valueStr = strconv.FormatInt(reflect.ValueOf(value).Int(), 10)
		case float32, float64:
			valueStr = strconv.FormatFloat(reflect.ValueOf(value).Float(), 'f', -1, 64)
		case bool:
			valueStr = strconv.FormatBool(v)
		default:
			valueStr = fmt.Sprintf("%v", v)
		}
		str = strings.ReplaceAll(str, placeholder, valueStr)
	}
	return str
}
