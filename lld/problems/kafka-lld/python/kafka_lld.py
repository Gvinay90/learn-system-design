"""A simplified, in-process pub/sub message broker modeled on Kafka: topics,
partitions, offsets, and consumer groups. It is not a networked service —
everything runs in a single process.
"""

from __future__ import annotations

import threading
from dataclasses import dataclass
from itertools import count
from typing import Dict, List, Protocol, Tuple


class TopicExistsError(Exception):
    """Raised when creating a topic that already exists."""


class TopicNotFoundError(Exception):
    """Raised when referencing a topic that hasn't been created."""


class PartitionNotFoundError(Exception):
    """Raised when referencing a partition index that doesn't exist."""


class InvalidPartitionCountError(Exception):
    """Raised when a topic is created with numPartitions <= 0."""


@dataclass(frozen=True)
class Message:
    key: str
    value: str
    offset: int


class Partition:
    """An append-only log of messages. Offsets are assigned sequentially
    starting at 0, in append order."""

    def __init__(self) -> None:
        self._lock = threading.Lock()
        self._messages: List[Message] = []

    def append(self, key: str, value: str) -> int:
        """Assigns the next sequential offset to (key, value) while holding
        the partition's lock, so concurrent producers can never race on
        offset assignment (no lost updates, no duplicate offsets)."""
        with self._lock:
            offset = len(self._messages)
            self._messages.append(Message(key=key, value=value, offset=offset))
            return offset

    def read(self, from_offset: int, max_messages: int) -> List[Message]:
        with self._lock:
            size = len(self._messages)
            if from_offset < 0 or from_offset >= size:
                return []
            end = size
            if max_messages > 0 and from_offset + max_messages < end:
                end = from_offset + max_messages
            return list(self._messages[from_offset:end])

    def __len__(self) -> int:
        with self._lock:
            return len(self._messages)


class PartitionerStrategy(Protocol):
    """Selects a partition index for a produced message."""

    def partition(self, key: str, num_partitions: int) -> int: ...


class KeyHashPartitioner:
    """Routes by hash(key) % numPartitions so all messages for the same key
    land on the same partition (preserving per-key order). Keyless messages
    are spread round robin across partitions."""

    def __init__(self) -> None:
        self._round_robin = count()
        self._round_robin_lock = threading.Lock()

    def partition(self, key: str, num_partitions: int) -> int:
        if not key:
            with self._round_robin_lock:
                n = next(self._round_robin)
            return n % num_partitions
        return self._fnv1a32(key) % num_partitions

    @staticmethod
    def _fnv1a32(s: str) -> int:
        """FNV-1a 32-bit hash, matching Go's hash/fnv New32a used by the
        reference implementation, so key routing is consistent across ports."""
        hash_ = 0x811C9DC5
        prime = 0x01000193
        for b in s.encode("utf-8"):
            hash_ ^= b
            hash_ = (hash_ * prime) & 0xFFFFFFFF
        return hash_


class Topic:
    """Owns a fixed set of partitions created at topic-creation time."""

    def __init__(self, name: str, num_partitions: int, partitioner: PartitionerStrategy) -> None:
        if num_partitions <= 0:
            raise InvalidPartitionCountError("numPartitions must be > 0")
        self.name = name
        self.partitions: List[Partition] = [Partition() for _ in range(num_partitions)]
        self._partitioner = partitioner

    def select_partition(self, key: str) -> int:
        return self._partitioner.partition(key, len(self.partitions))


class ConsumerGroup:
    """Tracks a committed offset per (topic, partition), independent of any
    other group reading the same topic."""

    def __init__(self, group_id: str) -> None:
        self.id = group_id
        self._lock = threading.Lock()
        self._offsets: Dict[Tuple[str, int], int] = {}

    def committed_offset(self, topic: str, partition_id: int) -> int:
        with self._lock:
            return self._offsets.get((topic, partition_id), 0)

    def commit(self, topic: str, partition_id: int, offset: int) -> None:
        with self._lock:
            self._offsets[(topic, partition_id)] = offset


class Broker:
    """Top-level orchestrator: it owns topics and consumer groups."""

    def __init__(self) -> None:
        self._lock = threading.RLock()
        self._topics: Dict[str, Topic] = {}
        self._consumer_groups: Dict[str, ConsumerGroup] = {}

    def create_topic(self, name: str, num_partitions: int) -> None:
        with self._lock:
            if name in self._topics:
                raise TopicExistsError("topic already exists")
            self._topics[name] = Topic(name, num_partitions, KeyHashPartitioner())

    def produce(self, topic_name: str, key: str, value: str) -> Tuple[int, int]:
        """Returns (partition_id, offset)."""
        with self._lock:
            topic = self._topics.get(topic_name)
        if topic is None:
            raise TopicNotFoundError("topic not found")
        partition_id = topic.select_partition(key)
        offset = topic.partitions[partition_id].append(key, value)
        return partition_id, offset

    def _get_or_create_consumer_group(self, group_id: str) -> ConsumerGroup:
        with self._lock:
            group = self._consumer_groups.get(group_id)
            if group is None:
                group = ConsumerGroup(group_id)
                self._consumer_groups[group_id] = group
            return group

    def consume(
        self, group_id: str, topic_name: str, partition_id: int, max_messages: int
    ) -> List[Message]:
        """Returns up to max_messages messages after the group's last
        committed offset for (topic, partition_id), then auto-commits past
        the returned batch. Auto-commit-on-read keeps this exercise's API
        small and deterministic; the alternative is an explicit commit_offset
        call with at-least-once redelivery on crash."""
        with self._lock:
            topic = self._topics.get(topic_name)
        if topic is None:
            raise TopicNotFoundError("topic not found")
        if partition_id < 0 or partition_id >= len(topic.partitions):
            raise PartitionNotFoundError("partition not found")

        group = self._get_or_create_consumer_group(group_id)
        from_offset = group.committed_offset(topic_name, partition_id)
        messages = topic.partitions[partition_id].read(from_offset, max_messages)
        if messages:
            group.commit(topic_name, partition_id, messages[-1].offset + 1)
        return messages

    def committed_offset(self, group_id: str, topic_name: str, partition_id: int) -> int:
        group = self._get_or_create_consumer_group(group_id)
        return group.committed_offset(topic_name, partition_id)


if __name__ == "__main__":
    broker = Broker()
    broker.create_topic("orders", 3)

    for v in ["a", "b", "c"]:
        partition_id, offset = broker.produce("orders", "k1", v)
        print(f"produced {v} -> partition {partition_id} offset {offset}")

    messages = broker.consume("group-1", "orders", 0, 10)
    print("group-1 read:", [(m.offset, m.value) for m in messages])

    again = broker.consume("group-1", "orders", 0, 10)
    print(f"group-1 read again (after auto-commit): {len(again)} messages")
