package util

import (
	"encoding/json"
)

func ToJsonString[T any](v T) string {
	j, _ := json.Marshal(v)
	return string(j)
}
