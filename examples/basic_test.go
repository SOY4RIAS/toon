package examples

import (
	"fmt"
	"testing"

	"github.com/SOY4RIAS/toon"
)

func TestBasicEncode(t *testing.T) {
	data := map[string]interface{}{
		"users": []map[string]interface{}{
			{"id": 1, "name": "Alice", "role": "admin"},
			{"id": 2, "name": "Bob", "role": "user"},
		},
	}

	encoded, err := toon.Encode(data, nil)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	fmt.Println("Encoded TOON:")
	fmt.Println(encoded)

	// Expected output format:
	// users[2]{id,name,role}:
	//   1,Alice,admin
	//   2,Bob,user
}

func TestPrimitiveEncode(t *testing.T) {
	data := map[string]interface{}{
		"id":     123,
		"name":   "Alice",
		"active": true,
	}

	encoded, err := toon.Encode(data, nil)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	fmt.Println("Encoded TOON:")
	fmt.Println(encoded)
}

func TestArrayEncode(t *testing.T) {
	data := map[string]interface{}{
		"tags": []string{"admin", "ops", "dev"},
	}

	encoded, err := toon.Encode(data, nil)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	fmt.Println("Encoded TOON:")
	fmt.Println(encoded)

	// Expected: tags[3]: admin,ops,dev
}
