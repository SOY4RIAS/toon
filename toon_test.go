package toon

import (
	"testing"
)

func TestDecodeSimpleObject(t *testing.T) {
	input := `id: 123
name: Alice
active: true`

	result, err := Decode(input, nil)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	obj, ok := result.(JsonObject)
	if !ok {
		t.Fatalf("Expected JsonObject, got %T", result)
	}

	if obj["id"] != 123.0 {
		t.Errorf("Expected id=123.0, got %v", obj["id"])
	}
	if obj["name"] != "Alice" {
		t.Errorf("Expected name=Alice, got %v", obj["name"])
	}
	if obj["active"] != true {
		t.Errorf("Expected active=true, got %v", obj["active"])
	}
}

func TestDecodeTabularArray(t *testing.T) {
	input := `items[2]{id,name,qty}:
  1,Alice,5
  2,Bob,3`

	result, err := Decode(input, nil)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	obj, ok := result.(JsonObject)
	if !ok {
		t.Fatalf("Expected JsonObject, got %T", result)
	}

	items, ok := obj["items"].(JsonArray)
	if !ok {
		t.Fatalf("Expected items to be JsonArray, got %T", obj["items"])
	}

	if len(items) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(items))
	}

	// Check first item
	item0, ok := items[0].(JsonObject)
	if !ok {
		t.Fatalf("Expected item 0 to be JsonObject, got %T", items[0])
	}
	if item0["id"] != 1.0 {
		t.Errorf("Expected item 0 id=1.0, got %v", item0["id"])
	}
	if item0["name"] != "Alice" {
		t.Errorf("Expected item 0 name=Alice, got %v", item0["name"])
	}
	if item0["qty"] != 5.0 {
		t.Errorf("Expected item 0 qty=5.0, got %v", item0["qty"])
	}
}

func TestDecodeInlinePrimitiveArray(t *testing.T) {
	input := `tags[3]: admin,ops,dev`

	result, err := Decode(input, nil)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	obj, ok := result.(JsonObject)
	if !ok {
		t.Fatalf("Expected JsonObject, got %T", result)
	}

	tags, ok := obj["tags"].(JsonArray)
	if !ok {
		t.Fatalf("Expected tags to be JsonArray, got %T", obj["tags"])
	}

	if len(tags) != 3 {
		t.Fatalf("Expected 3 tags, got %d", len(tags))
	}

	expected := []string{"admin", "ops", "dev"}
	for i, exp := range expected {
		if tags[i] != exp {
			t.Errorf("Expected tags[%d]=%s, got %v", i, exp, tags[i])
		}
	}
}

func TestDecodeNestedObject(t *testing.T) {
	input := `user:
  id: 123
  name: Ada`

	result, err := Decode(input, nil)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	obj, ok := result.(JsonObject)
	if !ok {
		t.Fatalf("Expected JsonObject, got %T", result)
	}

	user, ok := obj["user"].(JsonObject)
	if !ok {
		t.Fatalf("Expected user to be JsonObject, got %T", obj["user"])
	}

	if user["id"] != 123.0 {
		t.Errorf("Expected user.id=123.0, got %v", user["id"])
	}
	if user["name"] != "Ada" {
		t.Errorf("Expected user.name=Ada, got %v", user["name"])
	}
}

func TestDecodeListArray(t *testing.T) {
	input := `items[3]:
  - 1
  - a: 1
  - text`

	result, err := Decode(input, nil)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	obj, ok := result.(JsonObject)
	if !ok {
		t.Fatalf("Expected JsonObject, got %T", result)
	}

	items, ok := obj["items"].(JsonArray)
	if !ok {
		t.Fatalf("Expected items to be JsonArray, got %T", obj["items"])
	}

	if len(items) != 3 {
		t.Fatalf("Expected 3 items, got %d", len(items))
	}

	if items[0] != 1.0 {
		t.Errorf("Expected items[0]=1.0, got %v", items[0])
	}

	item1, ok := items[1].(JsonObject)
	if !ok {
		t.Fatalf("Expected items[1] to be JsonObject, got %T", items[1])
	}
	if item1["a"] != 1.0 {
		t.Errorf("Expected items[1].a=1.0, got %v", item1["a"])
	}

	if items[2] != "text" {
		t.Errorf("Expected items[2]=text, got %v", items[2])
	}
}
