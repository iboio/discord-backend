package methods

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"strconv"
	"strings"
)

func StringToStruct[T any](data string) (T, error) {
	var result T
	err := json.Unmarshal([]byte(data), &result)
	return result, err
}
func StringToInt64(s string) (int64, error) {
	result, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error converting string to int64: %v", err)
	}
	return result, nil
}

func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}
func MapToStruct(data map[string]string, result interface{}) error {
	if err := mapstructure.Decode(data, result); err != nil {
		return fmt.Errorf("error decoding map to struct: %w", err)
	}
	return nil
}
func StringToSlice(str string) ([]string, error) {
	var result []string
	if str != "" {
		result = strings.Split(str, ",")
	}
	return result, nil
}
