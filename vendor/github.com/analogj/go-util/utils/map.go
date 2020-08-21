package utils

import (
	"encoding/json"
	"sort"
)

func MapKeys(m map[string]interface{}) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func MapDeepCopy(m map[string]interface{}) (map[string]interface{}, error) {
	jsonString, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	nm := map[string]interface{}{}
	err = json.Unmarshal([]byte(jsonString), &nm)
	if err != nil {
		return nil, err
	}

	return nm, nil
}
