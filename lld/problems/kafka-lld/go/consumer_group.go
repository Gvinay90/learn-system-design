package kafkalld

import "sync"

type partitionKey struct {
	topic       string
	partitionID int
}

// ConsumerGroup tracks a committed offset per (topic, partition), independent
// of any other group reading the same topic.
type ConsumerGroup struct {
	ID string

	mu      sync.Mutex
	offsets map[partitionKey]int64
}

func NewConsumerGroup(id string) *ConsumerGroup {
	return &ConsumerGroup{ID: id, offsets: make(map[partitionKey]int64)}
}

func (g *ConsumerGroup) committedOffset(topic string, partitionID int) int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.offsets[partitionKey{topic, partitionID}]
}

func (g *ConsumerGroup) commit(topic string, partitionID int, offset int64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.offsets[partitionKey{topic, partitionID}] = offset
}
