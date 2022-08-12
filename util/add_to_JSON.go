package util

import (
	"encoding/json"
	"strings"
)

// AddItemToJSONArray Decodes the original string and appends to it the given object.
// Expects the incoming item to be json-ifyable.
// and be the same type as the original array.
// Returns the updated string.
func AddItemToJSONArray[T any](original string, newObject T) (string, error) {

	if strings.TrimSpace(original) == "" {
		original = "[]"
	}

	var arrayItems []T
	err := json.Unmarshal([]byte(original), &arrayItems)
	if err != nil {
		return "", err
	}

	arrayItems = append(arrayItems, newObject)

	jsonBytes, err := json.Marshal(arrayItems)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil

}
