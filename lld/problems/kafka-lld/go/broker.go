package kafkalld

import (
	"errors"
	"sync"
)

var (
	ErrTopicExists       = errors.New("topic already exists")
	ErrTopicNotFound     = errors.New("topic not found")
	ErrPartitionNotFound = errors.New("partition not found")
)

// Broker is the top-level orchestrator: it owns topics and consumer groups.
type Broker struct {
	mu             sync.RWMutex
	topics         map[string]*Topic
	consumerGroups map[string]*ConsumerGroup
}

func NewBroker() *Broker {
	return &Broker{
		topics:         make(map[string]*Topic),
		consumerGroups: make(map[string]*ConsumerGroup),
	}
}

func (b *Broker) CreateTopic(name string, numPartitions int) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, exists := b.topics[name]; exists {
		return ErrTopicExists
	}
	topic, err := NewTopic(name, numPartitions, &KeyHashPartitioner{})
	if err != nil {
		return err
	}
	b.topics[name] = topic
	return nil
}

func (b *Broker) Produce(topicName, key, value string) (partitionID int, offset int64, err error) {
	b.mu.RLock()
	topic, ok := b.topics[topicName]
	b.mu.RUnlock()
	if !ok {
		return 0, 0, ErrTopicNotFound
	}
	partitionID = topic.SelectPartition(key)
	offset = topic.Partitions[partitionID].Append(key, value)
	return partitionID, offset, nil
}

func (b *Broker) getOrCreateConsumerGroup(groupID string) *ConsumerGroup {
	b.mu.Lock()
	defer b.mu.Unlock()
	g, ok := b.consumerGroups[groupID]
	if !ok {
		g = NewConsumerGroup(groupID)
		b.consumerGroups[groupID] = g
	}
	return g
}

// Consume returns up to maxMessages messages after the group's last
// committed offset for (topic, partitionID), then auto-commits past the
// returned batch. Auto-commit-on-read keeps this exercise's API small and
// deterministic; see the README trade-offs section for the alternative
// (explicit CommitOffset, at-least-once redelivery on crash).
func (b *Broker) Consume(groupID, topicName string, partitionID, maxMessages int) ([]Message, error) {
	b.mu.RLock()
	topic, ok := b.topics[topicName]
	b.mu.RUnlock()
	if !ok {
		return nil, ErrTopicNotFound
	}
	if partitionID < 0 || partitionID >= len(topic.Partitions) {
		return nil, ErrPartitionNotFound
	}

	group := b.getOrCreateConsumerGroup(groupID)
	fromOffset := group.committedOffset(topicName, partitionID)
	messages := topic.Partitions[partitionID].Read(fromOffset, maxMessages)
	if len(messages) > 0 {
		group.commit(topicName, partitionID, messages[len(messages)-1].Offset+1)
	}
	return messages, nil
}

func (b *Broker) CommittedOffset(groupID, topicName string, partitionID int) int64 {
	group := b.getOrCreateConsumerGroup(groupID)
	return group.committedOffset(topicName, partitionID)
}
