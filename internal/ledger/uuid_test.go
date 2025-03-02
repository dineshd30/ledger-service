package ledger_test

import (
	"testing"

	"github.com/dineshd30/ledger-service/internal/ledger"
	"github.com/google/uuid"
)

func TestGenerateNotEmpty(t *testing.T) {
	generator := ledger.NewUUIDGenerator()
	result := generator.Generate()

	if result == "" {
		t.Errorf("Expected non-empty UUID, got empty string")
	}
}

func TestGenerateValidUUID(t *testing.T) {
	generator := ledger.NewUUIDGenerator()
	result := generator.Generate()

	if _, err := uuid.Parse(result); err != nil {
		t.Errorf("Expected valid UUID, got error: %v", err)
	}
}

func TestGenerateUnique(t *testing.T) {
	generator := ledger.NewUUIDGenerator()
	uuid1 := generator.Generate()
	uuid2 := generator.Generate()

	if uuid1 == uuid2 {
		t.Errorf("Expected unique UUIDs, but got the same value: %s", uuid1)
	}
}
