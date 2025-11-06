package toon_test

import (
	"fmt"
	"testing"

	"github.com/toon-format/toon-go"
)

func TestBasicDecode(t *testing.T) {
	input := `
name: John Doe
age: 30
active: true
`

	result, err := toon.Decode(input, nil)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	obj, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map[string]interface{}, got %T", result)
	}

	if obj["name"] != "John Doe" {
		t.Errorf("Expected name to be 'John Doe', got %v", obj["name"])
	}

	if obj["age"] != 30.0 {
		t.Errorf("Expected age to be 30, got %v", obj["age"])
	}

	if obj["active"] != true {
		t.Errorf("Expected active to be true, got %v", obj["active"])
	}
}

func TestInlineArray(t *testing.T) {
	input := `tags[3]: web, backend, api`

	result, err := toon.Decode(input, nil)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	obj, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map[string]interface{}, got %T", result)
	}

	tags, ok := obj["tags"].([]interface{})
	if !ok {
		t.Fatalf("Expected tags to be []interface{}, got %T", obj["tags"])
	}

	if len(tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(tags))
	}

	expected := []string{"web", "backend", "api"}
	for i, exp := range expected {
		if tags[i] != exp {
			t.Errorf("Expected tags[%d] to be '%s', got %v", i, exp, tags[i])
		}
	}
}

func ExampleDecode() {
	input := `
name: John Doe
age: 30
active: true
tags[3]: web, backend, api
`

	result, err := toon.Decode(input, nil)
	if err != nil {
		panic(err)
	}

	obj := result.(map[string]interface{})
	fmt.Printf("Name: %s\n", obj["name"])
	fmt.Printf("Age: %.0f\n", obj["age"])
	fmt.Printf("Active: %t\n", obj["active"])
	fmt.Printf("Tags: %v\n", obj["tags"])

	// Output:
	// Name: John Doe
	// Age: 30
	// Active: true
	// Tags: [web backend api]
}
