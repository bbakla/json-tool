package handlers

import (
	"encoding/json"
	"errors"
	"sort"

	"gopkg.in/yaml.v3"
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

// Minify compacts JSON without whitespace.
func Minify(raw string) (string, error) {
	if raw == "" {
		return "", errors.New("No JSON providedddd")
	}
	var payload any
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return "", err
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ToYAML converts JSON into a YAML string.
func ToYAML(raw string) (string, error) {
	if raw == "" {
		return "", errors.New("no JSON provided")
	}
	var payload any
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return "", err
	}
	b, err := yaml.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ExtractKeyJSON finds all values for a key and returns them as JSON.
// If multiple values are found, it returns an array of objects {key: value}; if one, it returns a single object.
func ExtractKeyJSON(raw, key string) (string, error) {
	if raw == "" {
		return "", errors.New("no JSON provided")
	}
	if key == "" {
		return "", errors.New("no key provided")
	}
	var payload any
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return "", err
	}

	var values []any
	collectValues(payload, key, &values)
	if len(values) == 0 {
		return "", errors.New("key not found")
	}

	var out any
	if len(values) == 1 {
		out = map[string]any{key: values[0]}
	} else {
		arr := make([]map[string]any, 0, len(values))
		for _, v := range values {
			arr = append(arr, map[string]any{key: v})
		}
		out = arr
	}

	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func collectValues(node any, key string, results *[]any) {
	switch val := node.(type) {
	case map[string]any:
		for k, v := range val {
			if k == key {
				*results = append(*results, v)
			}
			collectValues(v, key, results)
		}
	case []any:
		for _, item := range val {
			collectValues(item, key, results)
		}
	}
}
