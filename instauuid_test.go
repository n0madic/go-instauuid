package instauuid

import (
	"encoding/base64"
	"encoding/hex"
	"testing"
	"time"
)

func BenchmarkGenerateID(b *testing.B) {
	generator := NewGenerator(1, 0)
	for i := 0; i < b.N; i++ {
		generator.GenerateID()
	}
}

func TestIDUniqueness(t *testing.T) {
	generator := NewGenerator(1, 0)

	idSet := make(map[uint64]bool)
	count := 1024

	for i := 0; i < count; i++ {
		id := generator.GenerateID()
		if _, exists := idSet[id]; exists {
			t.Fatalf("Duplicate ID found: %d", id)
		}
		idSet[id] = true
	}
}

func TestIDFormat(t *testing.T) {
	generator := NewGenerator(1, 0)

	for i := 0; i < 100; i++ {
		base64ID := generator.GenerateBase64()
		hexID := generator.GenerateHex()

		if _, err := base64.RawURLEncoding.DecodeString(base64ID); err != nil {
			t.Errorf("Invalid Base64 format: %s", base64ID)
		}

		if _, err := hex.DecodeString(hexID); err != nil {
			t.Errorf("Invalid Hex format: %s", hexID)
		}
	}
}

func TestPerformance(t *testing.T) {
	generator := NewGenerator(1, 0)

	start := time.Now()
	count := 1024

	for i := 0; i < count; i++ {
		generator.GenerateID()
	}

	elapsed := time.Since(start)
	t.Logf("Generated %d IDs in %s", count, elapsed)

	if elapsed.Seconds() > 1 {
		t.Errorf("Performance issue: took too long to generate IDs")
	}
}

func TestSequenceOverflow(t *testing.T) {
	generator := NewGenerator(1, 0)

	// Simulate the sequence number reaching its maximum
	generator.sequence = maxSeq

	// Generate one more ID to trigger overflow
	id1 := generator.GenerateID()
	id2 := generator.GenerateID()

	if id1 == id2 {
		t.Errorf("ID collision detected when sequence overflows")
	}
}

func TestShardIDConsistency(t *testing.T) {
	generator1 := NewGenerator(1, 0)
	generator2 := NewGenerator(2, 0)

	id1 := generator1.GenerateID()
	id2 := generator2.GenerateID()

	if id1>>shardShift == id2>>shardShift {
		t.Errorf("Shard IDs are not differentiating IDs correctly")
	}
}
