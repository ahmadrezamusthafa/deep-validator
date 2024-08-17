package utils

import (
	"bytes"
	"fmt"
)

func ReplacePlaceholders(str string, attributes map[string]interface{}) string {
	var buf, bufTempKey bytes.Buffer
	inPlaceholder := false

	for i := 0; i < len(str); i++ {
		char := str[i]

		if char == '{' {
			if i+1 < len(str) && str[i+1] == '{' {
				inPlaceholder = true
				i++
				bufTempKey.Reset()
				continue
			}
		} else if char == '}' {
			if inPlaceholder && i+1 < len(str) && str[i+1] == '}' {
				if value, ok := attributes[bufTempKey.String()]; ok {
					buf.WriteString(fmt.Sprint(value))
				} else {
					buf.WriteString("{{" + bufTempKey.String() + "}}")
				}
				inPlaceholder = false
				i++
				continue
			}
		}

		if inPlaceholder {
			bufTempKey.WriteByte(char)
		} else {
			buf.WriteByte(char)
		}
	}

	return buf.String()
}
