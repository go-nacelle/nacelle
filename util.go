package nacelle

import (
	"fmt"
	"reflect"
)

func serializeKey(v interface{}) string {
	if str, ok := v.(string); ok {
		return str
	}

	if stringer, ok := v.(fmt.Stringer); ok {
		return stringer.String()
	}

	t := reflect.TypeOf(v)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Name()
}
