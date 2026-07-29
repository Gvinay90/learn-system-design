// Package kafkalld implements a simplified, in-process pub/sub message broker
// modeled on Kafka: topics, partitions, offsets, and consumer groups. It is
// not a networked service — everything runs in a single process.
package kafkalld

import "sync"

type Message struct {
	Key    string
	Value  string
	Offset int64
}

type Partition struct {
	mu       sync.Mutex
	messages []Message
}

func NewPartition() *Partition {
	return &Partition{}
}

// Append assigns the next sequential offset to (key, value) while holding
// the partition's lock, so concurrent producers can never race on offset
// assignment (no lost updates, no duplicate offsets).
func (p *Partition) Append(key, value string) int64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	offset := int64(len(p.messages))
	p.messages = append(p.messages, Message{Key: key, Value: value, Offset: offset})
	return offset
}

func (p *Partition) Read(fromOffset int64, maxMessages int) []Message {
	p.mu.Lock()
	defer p.mu.Unlock()
	if fromOffset < 0 || fromOffset >= int64(len(p.messages)) {
		return nil
	}
	end := int64(len(p.messages))
	if maxMessages > 0 && fromOffset+int64(maxMessages) < end {
		end = fromOffset + int64(maxMessages)
	}
	out := make([]Message, end-fromOffset)
	copy(out, p.messages[fromOffset:end])
	return out
}

func (p *Partition) Len() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return int64(len(p.messages))
}
