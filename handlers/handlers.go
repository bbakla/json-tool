package handlers

import (
	"encoding/json"
	"errors"
	"sort"
)

// Process parses JSON, pretty-prints it, and finds values for the provided key as well as keys for the provided value.
func Process(raw, key, value string) (string, []string, []string, error) {
	if raw == "" {
		return "", nil, nil, errors.New("no JSON provided")
	}

	var payload any
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return "", nil, nil, err
	}

	pretty, err := marshalPretty(payload)
	if err != nil {
		return "", nil, nil, err
	}

	valueMatches := findKeyValues(payload, key)
	sort.Strings(valueMatches)

	keyMatches := findKeysForValue(payload, value)
	sort.Strings(keyMatches)

	return pretty, valueMatches, keyMatches, nil
}

func marshalPretty(v any) (string, error) {
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func findKeyValues(v any, key string) []string {
	if key == "" {
		return nil
	}
	var results []string
	walkForKey(v, key, &results)
	return results
}

func walkForKey(node any, key string, results *[]string) {
	switch val := node.(type) {
	case map[string]any:
		for k, v := range val {
			if k == key {
				switch m := v.(type) {
				case string:
					*results = append(*results, m)
				default:
					b, err := json.Marshal(m)
					if err == nil {
						*results = append(*results, string(b))
					}
				}
			}
			walkForKey(v, key, results)
		}
	case []any:
		for _, item := range val {
			walkForKey(item, key, results)
		}
	}
}

func findKeysForValue(v any, target string) []string {
	if target == "" {
		return nil
	}
	var results []string
	walkForValue(v, target, &results)
	return results
}

func walkForValue(node any, target string, results *[]string) {
	switch val := node.(type) {
	case map[string]any:
		for k, v := range val {
			if valueMatches(v, target) {
				*results = append(*results, k)
			}
			walkForValue(v, target, results)
		}
	case []any:
		for _, item := range val {
			walkForValue(item, target, results)
		}
	}
}

func valueMatches(v any, target string) bool {
	switch m := v.(type) {
	case string:
		return m == target
	default:
		b, err := json.Marshal(m)
		if err != nil {
			return false
		}
		return string(b) == target
	}
}
