package kafkalld

import (
	"hash/fnv"
	"sync/atomic"
)

// PartitionerStrategy selects a partition index for a produced message.
type PartitionerStrategy interface {
	Partition(key string, numPartitions int) int
}

// KeyHashPartitioner routes by hash(key) % numPartitions so all messages for
// the same key land on the same partition (preserving per-key order).
// Keyless messages are spread round robin across partitions.
type KeyHashPartitioner struct {
	roundRobin int64
}

func (k *KeyHashPartitioner) Partition(key string, numPartitions int) int {
	if key == "" {
		n := atomic.AddInt64(&k.roundRobin, 1) - 1
		return int(n % int64(numPartitions))
	}
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() % uint32(numPartitions))
}
