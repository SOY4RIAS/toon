package decode

import (
	"testing"

	"github.com/toon-format/toon-go"
)

func TestDecodeSimpleObject(t *testing.T) {
	input := `id: 1
name: Alice
active: true`

	scanResult, err := ToParsedLines(input, 2, true)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	cursor := NewLineCursor(scanResult.Lines, scanResult.BlankLines)
	options := &toon.DecodeOptions{Indent: 2, Strict: true}

	result, err := DecodeValueFromLines(cursor, options)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	obj, ok := result.(toon.JsonObject)
	if !ok {
		t.Fatalf("expected JsonObject, got %T", result)
	}

	if obj["id"] != 1.0 {
		t.Errorf("expected id=1.0, got %v", obj["id"])
	}
	if obj["name"] != "Alice" {
		t.Errorf("expected name=Alice, got %v", obj["name"])
	}
	if obj["active"] != true {
		t.Errorf("expected active=true, got %v", obj["active"])
	}
}

func TestDecodeInlinePrimitiveArray(t *testing.T) {
	input := `tags[3]: admin,user,guest`

	scanResult, err := ToParsedLines(input, 2, true)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	cursor := NewLineCursor(scanResult.Lines, scanResult.BlankLines)
	options := &toon.DecodeOptions{Indent: 2, Strict: true}

	result, err := DecodeValueFromLines(cursor, options)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	obj, ok := result.(toon.JsonObject)
	if !ok {
		t.Fatalf("expected JsonObject, got %T", result)
	}

	tags, ok := obj["tags"].(toon.JsonArray)
	if !ok {
		t.Fatalf("expected tags to be JsonArray, got %T", obj["tags"])
	}

	if len(tags) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(tags))
	}

	if tags[0] != "admin" || tags[1] != "user" || tags[2] != "guest" {
		t.Errorf("unexpected tag values: %v", tags)
	}
}

func TestDecodeTabularArray(t *testing.T) {
	input := `items[2]{id,name,price}:
  1,Widget,9.99
  2,Gadget,14.50`

	scanResult, err := ToParsedLines(input, 2, true)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	cursor := NewLineCursor(scanResult.Lines, scanResult.BlankLines)
	options := &toon.DecodeOptions{Indent: 2, Strict: true}

	result, err := DecodeValueFromLines(cursor, options)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	obj, ok := result.(toon.JsonObject)
	if !ok {
		t.Fatalf("expected JsonObject, got %T", result)
	}

	items, ok := obj["items"].(toon.JsonArray)
	if !ok {
		t.Fatalf("expected items to be JsonArray, got %T", obj["items"])
	}

	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}

	item1, ok := items[0].(toon.JsonObject)
	if !ok {
		t.Fatalf("expected item1 to be JsonObject, got %T", items[0])
	}

	if item1["id"] != 1.0 || item1["name"] != "Widget" || item1["price"] != 9.99 {
		t.Errorf("unexpected item1 values: %v", item1)
	}
}

func TestDecodeListArray(t *testing.T) {
	input := `items[3]:
  - hello
  - 42
  - true`

	scanResult, err := ToParsedLines(input, 2, true)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	cursor := NewLineCursor(scanResult.Lines, scanResult.BlankLines)
	options := &toon.DecodeOptions{Indent: 2, Strict: true}

	result, err := DecodeValueFromLines(cursor, options)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	obj, ok := result.(toon.JsonObject)
	if !ok {
		t.Fatalf("expected JsonObject, got %T", result)
	}

	items, ok := obj["items"].(toon.JsonArray)
	if !ok {
		t.Fatalf("expected items to be JsonArray, got %T", obj["items"])
	}

	if len(items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(items))
	}

	if items[0] != "hello" {
		t.Errorf("expected items[0]='hello', got %v", items[0])
	}
	if items[1] != 42.0 {
		t.Errorf("expected items[1]=42.0, got %v", items[1])
	}
	if items[2] != true {
		t.Errorf("expected items[2]=true, got %v", items[2])
	}
}

func TestDecodeNestedObject(t *testing.T) {
	input := `user:
  id: 1
  name: Alice`

	scanResult, err := ToParsedLines(input, 2, true)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	cursor := NewLineCursor(scanResult.Lines, scanResult.BlankLines)
	options := &toon.DecodeOptions{Indent: 2, Strict: true}

	result, err := DecodeValueFromLines(cursor, options)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	obj, ok := result.(toon.JsonObject)
	if !ok {
		t.Fatalf("expected JsonObject, got %T", result)
	}

	user, ok := obj["user"].(toon.JsonObject)
	if !ok {
		t.Fatalf("expected user to be JsonObject, got %T", obj["user"])
	}

	if user["id"] != 1.0 || user["name"] != "Alice" {
		t.Errorf("unexpected user values: %v", user)
	}
}
