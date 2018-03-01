package log

import "time"

type Fields map[string]interface{}

func (f Fields) clone() Fields {
	clone := Fields{}
	for k, v := range f {
		clone[k] = v
	}

	return clone
}

func (f Fields) normalizeTimeValues() Fields {
	for key, val := range f {
		switch v := val.(type) {
		case time.Time:
			f[key] = v.Format(JSONTimeFormat)
		default:
			f[key] = v
		}
	}

	return f
}
