package kafkalld

import "errors"

var ErrInvalidPartitionCount = errors.New("numPartitions must be > 0")

// Topic owns a fixed set of partitions created at topic-creation time.
type Topic struct {
	Name        string
	Partitions  []*Partition
	partitioner PartitionerStrategy
}

func NewTopic(name string, numPartitions int, partitioner PartitionerStrategy) (*Topic, error) {
	if numPartitions <= 0 {
		return nil, ErrInvalidPartitionCount
	}
	partitions := make([]*Partition, numPartitions)
	for i := range partitions {
		partitions[i] = NewPartition()
	}
	return &Topic{Name: name, Partitions: partitions, partitioner: partitioner}, nil
}

func (t *Topic) SelectPartition(key string) int {
	return t.partitioner.Partition(key, len(t.Partitions))
}
