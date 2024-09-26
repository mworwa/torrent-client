package bencode

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Bencode encodes a Go value into Bencode format.
func Encode(value interface{}) ([]byte, error) {
	var builder strings.Builder
	err := bencodeValue(value, &builder)
	if err != nil {
		return nil, err
	}
	return []byte(builder.String()), nil
}

func bencodeValue(value interface{}, builder *strings.Builder) error {
	switch v := value.(type) {
	case int, int8, int16, int32, int64:
		builder.WriteString("i")
		builder.WriteString(fmt.Sprintf("%d", v))
		builder.WriteString("e")
	case uint, uint8, uint16, uint32, uint64:
		builder.WriteString("i")
		builder.WriteString(fmt.Sprintf("%d", v))
		builder.WriteString("e")
	case string:
		builder.WriteString(strconv.Itoa(len(v)))
		builder.WriteString(":")
		builder.WriteString(v)
	case []byte:
		builder.WriteString(strconv.Itoa(len(v)))
		builder.WriteString(":")
		builder.Write(v)
	case []interface{}:
		builder.WriteString("l")
		for _, item := range v {
			if err := bencodeValue(item, builder); err != nil {
				return err
			}
		}
		builder.WriteString("e")
	case map[string]interface{}:
		builder.WriteString("d")
		keys := make([]string, 0, len(v))
		for key := range v {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			builder.WriteString(strconv.Itoa(len(key)))
			builder.WriteString(":")
			builder.WriteString(key)
			if err := bencodeValue(v[key], builder); err != nil {
				return err
			}
		}
		builder.WriteString("e")
	default:
		val := reflect.ValueOf(value)
		switch val.Kind() {
		case reflect.Slice, reflect.Array:
			builder.WriteString("l")
			for i := 0; i < val.Len(); i++ {
				if err := bencodeValue(val.Index(i).Interface(), builder); err != nil {
					return err
				}
			}
			builder.WriteString("e")
		case reflect.Map:
			if val.Type().Key().Kind() != reflect.String {
				return fmt.Errorf("unsupported map key type: %s", val.Type().Key().Kind())
			}
			builder.WriteString("d")
			keys := val.MapKeys()
			keyStrs := make([]string, 0, len(keys))
			keyMap := make(map[string]reflect.Value)
			for _, key := range keys {
				k := key.String()
				keyStrs = append(keyStrs, k)
				keyMap[k] = val.MapIndex(key)
			}
			sort.Strings(keyStrs)
			for _, k := range keyStrs {
				builder.WriteString(strconv.Itoa(len(k)))
				builder.WriteString(":")
				builder.WriteString(k)
				if err := bencodeValue(keyMap[k].Interface(), builder); err != nil {
					return err
				}
			}
			builder.WriteString("e")
		default:
			return fmt.Errorf("bencode encoding error: unsupported type: %T", v)
		}
	}
	return nil
}
