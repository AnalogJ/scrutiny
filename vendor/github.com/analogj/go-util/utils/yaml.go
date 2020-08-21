package utils

import "fmt"

// stringifyKeysMapValue recurses into in and changes all instances of
// map[interface{}]interface{} to map[string]interface{}. This is useful to
// work around the impedence mismatch between JSON and YAML unmarshaling that's
// described here: https://github.com/go-yaml/yaml/issues/139
//
// Inspired by https://github.com/stripe/stripe-mock, MIT licensed
func StringifyYAMLMapKeys(in interface{}) interface{} {
	switch in := in.(type) {
	case []interface{}:
		res := make([]interface{}, len(in))
		for i, v := range in {
			res[i] = StringifyYAMLMapKeys(v)
		}
		return res
	case map[string]interface{}:
		res := make(map[string]interface{})
		for k, v := range in {
			res[k] = StringifyYAMLMapKeys(v)
		}
		return res
	case map[interface{}]interface{}:
		res := make(map[string]interface{})
		for k, v := range in {
			res[fmt.Sprintf("%v", k)] = StringifyYAMLMapKeys(v)
		}
		return res
	default:
		return in
	}
}
