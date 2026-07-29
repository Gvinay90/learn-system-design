package kafkalld

import (
	"sync"
	"testing"
)

func TestProduceConsumeInOrder(t *testing.T) {
	b := NewBroker()
	if err := b.CreateTopic("orders", 1); err != nil {
		t.Fatalf("CreateTopic: %v", err)
	}

	for i, v := range []string{"a", "b", "c"} {
		partitionID, offset, err := b.Produce("orders", "k1", v)
		if err != nil {
			t.Fatalf("Produce: %v", err)
		}
		if partitionID != 0 {
			t.Fatalf("expected partition 0, got %d", partitionID)
		}
		if offset != int64(i) {
			t.Fatalf("expected offset %d, got %d", i, offset)
		}
	}

	messages, err := b.Consume("g1", "orders", 0, 10)
	if err != nil {
		t.Fatalf("Consume: %v", err)
	}
	if len(messages) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(messages))
	}
	for i, m := range messages {
		if m.Offset != int64(i) {
			t.Fatalf("expected offset %d, got %d", i, m.Offset)
		}
	}

	more, err := b.Consume("g1", "orders", 0, 10)
	if err != nil {
		t.Fatalf("Consume: %v", err)
	}
	if len(more) != 0 {
		t.Fatalf("expected no new messages after auto-commit, got %d", len(more))
	}
}

func TestConsumerGroupsTrackOffsetsIndependently(t *testing.T) {
	b := NewBroker()
	if err := b.CreateTopic("orders", 1); err != nil {
		t.Fatalf("CreateTopic: %v", err)
	}
	for _, v := range []string{"a", "b", "c"} {
		if _, _, err := b.Produce("orders", "", v); err != nil {
			t.Fatalf("Produce: %v", err)
		}
	}

	g1Messages, err := b.Consume("group-1", "orders", 0, 2)
	if err != nil {
		t.Fatalf("Consume group-1: %v", err)
	}
	if len(g1Messages) != 2 {
		t.Fatalf("expected 2 messages for group-1, got %d", len(g1Messages))
	}

	g2Messages, err := b.Consume("group-2", "orders", 0, 10)
	if err != nil {
		t.Fatalf("Consume group-2: %v", err)
	}
	if len(g2Messages) != 3 {
		t.Fatalf("expected group-2 to see all 3 messages independently, got %d", len(g2Messages))
	}

	g1Rest, err := b.Consume("group-1", "orders", 0, 10)
	if err != nil {
		t.Fatalf("Consume group-1 rest: %v", err)
	}
	if len(g1Rest) != 1 {
		t.Fatalf("expected group-1 to resume at its own offset and see 1 message, got %d", len(g1Rest))
	}
	if g1Rest[0].Offset != 2 {
		t.Fatalf("expected offset 2, got %d", g1Rest[0].Offset)
	}
}

func TestEdgeCases(t *testing.T) {
	b := NewBroker()
	if err := b.CreateTopic("orders", 1); err != nil {
		t.Fatalf("CreateTopic: %v", err)
	}
	if _, _, err := b.Produce("orders", "", "only-message"); err != nil {
		t.Fatalf("Produce: %v", err)
	}

	messages, err := b.Consume("g1", "orders", 0, 10)
	if err != nil {
		t.Fatalf("Consume: %v", err)
	}
	if len(messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(messages))
	}

	past, err := b.Consume("g1", "orders", 0, 10)
	if err != nil {
		t.Fatalf("unexpected error consuming past end: %v", err)
	}
	if len(past) != 0 {
		t.Fatalf("expected no messages past end, got %d", len(past))
	}

	if _, _, err := b.Produce("unknown-topic", "k", "v"); err != ErrTopicNotFound {
		t.Fatalf("expected ErrTopicNotFound, got %v", err)
	}
	if _, err := b.Consume("g1", "unknown-topic", 0, 10); err != ErrTopicNotFound {
		t.Fatalf("expected ErrTopicNotFound, got %v", err)
	}
	if _, err := b.Consume("g1", "orders", 5, 10); err != ErrPartitionNotFound {
		t.Fatalf("expected ErrPartitionNotFound, got %v", err)
	}
}

// TestConcurrentProduceIntoSamePartition asserts many goroutines racing to
// append to the same partition never lose a message or assign a duplicate
// offset — the mutex in Partition.Append must serialize them.
func TestConcurrentProduceIntoSamePartition(t *testing.T) {
	b := NewBroker()
	if err := b.CreateTopic("orders", 1); err != nil {
		t.Fatalf("CreateTopic: %v", err)
	}

	const n = 500
	var wg sync.WaitGroup
	offsets := make(chan int64, n)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, offset, err := b.Produce("orders", "", "v")
			if err != nil {
				t.Errorf("Produce: %v", err)
				return
			}
			offsets <- offset
		}()
	}
	wg.Wait()
	close(offsets)

	seen := make(map[int64]bool, n)
	for o := range offsets {
		if seen[o] {
			t.Fatalf("duplicate offset %d", o)
		}
		seen[o] = true
	}
	if len(seen) != n {
		t.Fatalf("expected %d unique offsets, got %d", n, len(seen))
	}
	for i := int64(0); i < n; i++ {
		if !seen[i] {
			t.Fatalf("gap in offsets: missing %d", i)
		}
	}
}
