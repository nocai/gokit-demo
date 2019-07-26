package configs

import (
	"gopkg.in/yaml.v2"
	"strings"
	"time"
)

type mapType map[string]interface{}

// searchMap recursively searches for a value for path in source map.
// Returns nil if not found.
// Note: This assumes that the path entries and map keys are lower cased.
func (m mapType) searchMap(source map[string]interface{}, path []string) interface{} {
	if len(path) == 0 {
		return source
	}

	next, ok := source[path[0]]
	if ok {
		// Fast path
		if len(path) == 1 {
			return next
		}

		// Nested case
		switch next.(type) {
		case map[interface{}]interface{}:
			return m.searchMap(ToStringMap(next), path[1:])
		case map[string]interface{}:
			// Type assertion is safe here since it is only reached
			// if the type of `next` is the same as the type being asserted
			return m.searchMap(next.(map[string]interface{}), path[1:])
		default:
			// got a value but nested key expected, return "nil" for not found
			return nil
		}
	}
	return nil
}

func (m mapType) Get(key string) interface{} {
	key = strings.ToLower(key)
	if value, exist := m[key]; exist {
		return value
	}

	path := strings.Split(key, ".")
	val := m.searchMap(m, path)
	switch val.(type) {
	case bool:
		return ToBool(val)
	case string:
		return ToString(val)
	case int32, int16, int8, int:
		return ToInt(val)
	case uint:
		return ToUint(val)
	case uint32:
		return ToUint32(val)
	case uint64:
		return ToUint64(val)
	case int64:
		return ToInt64(val)
	case float64, float32:
		return ToFloat64(val)
	case time.Time:
		return ToTime(val)
	case time.Duration:
		return ToDuration(val)
	case []string:
		return ToStringSlice(val)
	}

	return val
}

func Unmarshal(in []byte) error {
	return yaml.Unmarshal(in, &defaultMapType)
}

var defaultMapType mapType

// GetString returns the value associated with the key as a string.
func GetString(key string) string { return defaultMapType.GetString(key) }
func (m *mapType) GetString(key string) string {
	return ToString(m.Get(key))
}

// GetBool returns the value associated with the key as a boolean.
func GetBool(key string) bool { return defaultMapType.GetBool(key) }
func (m *mapType) GetBool(key string) bool {
	return ToBool(m.Get(key))
}

// GetInt returns the value associated with the key as an integer.
func GetInt(key string) int { return defaultMapType.GetInt(key) }
func (m *mapType) GetInt(key string) int {
	return ToInt(m.Get(key))
}

// GetInt32 returns the value associated with the key as an integer.
func GetInt32(key string) int32 { return defaultMapType.GetInt32(key) }
func (m *mapType) GetInt32(key string) int32 {
	return ToInt32(m.Get(key))
}

// GetInt64 returns the value associated with the key as an integer.
func GetInt64(key string) int64 { return defaultMapType.GetInt64(key) }
func (m *mapType) GetInt64(key string) int64 {
	return ToInt64(m.Get(key))
}

// GetUint returns the value associated with the key as an unsigned integer.
func GetUint(key string) uint { return defaultMapType.GetUint(key) }
func (m *mapType) GetUint(key string) uint {
	return ToUint(m.Get(key))
}

// GetUint32 returns the value associated with the key as an unsigned integer.
func GetUint32(key string) uint32 { return defaultMapType.GetUint32(key) }
func (m *mapType) GetUint32(key string) uint32 {
	return ToUint32(m.Get(key))
}

// GetUint64 returns the value associated with the key as an unsigned integer.
func GetUint64(key string) uint64 { return defaultMapType.GetUint64(key) }
func (m *mapType) GetUint64(key string) uint64 {
	return ToUint64(m.Get(key))
}

// GetFloat64 returns the value associated with the key as a float64.
func GetFloat64(key string) float64 { return defaultMapType.GetFloat64(key) }
func (m *mapType) GetFloat64(key string) float64 {
	return ToFloat64(m.Get(key))
}

// GetTime returns the value associated with the key as time.
func GetTime(key string) time.Time { return defaultMapType.GetTime(key) }
func (m *mapType) GetTime(key string) time.Time {
	return ToTime(m.Get(key))
}

// GetDuration returns the value associated with the key as a duration.
func GetDuration(key string) time.Duration { return defaultMapType.GetDuration(key) }
func (m *mapType) GetDuration(key string) time.Duration {
	return ToDuration(m.Get(key))
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func GetStringSlice(key string) []string { return defaultMapType.GetStringSlice(key) }
func (m *mapType) GetStringSlice(key string) []string {
	return ToStringSlice(m.Get(key))
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func GetStringMap(key string) map[string]interface{} { return defaultMapType.GetStringMap(key) }
func (m *mapType) GetStringMap(key string) map[string]interface{} {
	return ToStringMap(m.Get(key))
}

// GetStringMapString returns the value associated with the key as a map of strings.
func GetStringMapString(key string) map[string]string { return defaultMapType.GetStringMapString(key) }
func (m *mapType) GetStringMapString(key string) map[string]string {
	return ToStringMapString(m.Get(key))
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func GetStringMapStringSlice(key string) map[string][]string {
	return defaultMapType.GetStringMapStringSlice(key)
}
func (m *mapType) GetStringMapStringSlice(key string) map[string][]string {
	return ToStringMapStringSlice(m.Get(key))
}
