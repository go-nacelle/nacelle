package config

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

var replacer = strings.NewReplacer(
	"\n", `\n`,
	"\t", `\t`,
	"\r", `\r`,
)

func toJSON(data []byte, v interface{}) bool {
	if json.Unmarshal(data, v) == nil {
		return true
	}

	ptr := reflect.ValueOf(v)

	if ptr.Kind() == reflect.Ptr && reflect.Indirect(ptr).Kind() == reflect.String {
		if json.Unmarshal(quoteJSON(data), v) == nil {
			return true
		}
	}

	return false
}

func quoteJSON(data []byte) []byte {
	return []byte(fmt.Sprintf(`"%s"`, replacer.Replace(string(data))))
}
