package handlers

import "testing"

func TestProcessFormatsJSON(t *testing.T) {
	raw := `{"name":"alice"}`
	formatted, matches, keys, err := Process(raw, "name", "alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if formatted == raw {
		t.Fatalf("expected pretty-printed output, got same string")
	}
	if len(matches) != 1 || matches[0] != "alice" {
		t.Fatalf("expected match alice, got %v", matches)
	}
	if len(keys) != 1 || keys[0] != "name" {
		t.Fatalf("expected key match name, got %v", keys)
	}
}

func TestProcessInvalidJSON(t *testing.T) {
	_, _, _, err := Process("{", "name", "alice")
	if err == nil {
		t.Fatalf("expected error for invalid JSON")
	}
}

func TestFindKeyValuesNested(t *testing.T) {
	raw := `{"a": {"b": [{"c": 1}, {"c": 2}]}}`
	formatted, matches, keys, err := Process(raw, "c", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = formatted
	if len(matches) != 2 {
		t.Fatalf("expected two matches, got %v", matches)
	}
	if len(keys) != 1 || keys[0] != "c" {
		t.Fatalf("expected key match c, got %v", keys)
	}
}
