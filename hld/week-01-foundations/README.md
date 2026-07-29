# Week 1 — Foundations & Mental Models

Part of the [8-week HLD learning path](../README.md).

## Concept: How to approach any HLD interview (the RESHADED framework)

- **The framework**: **R**equirements → **E**stimation → **S**torage/data model → **H**igh-level design
  → **A**PI design → **D**eep dives → **E**dge cases → **D**iscussion of trade-offs (mnemonic groupings
  vary by source; the order that matters is: clarify scope before estimating, estimate before designing,
  design at a high level before diving into any one component). Interviewers are grading your *process*,
  not just the final diagram — a candidate who narrates a structured process while drawing a mediocre
  diagram consistently beats one who jumps straight to a brilliant diagram with no narration.
- **Why order matters**: designing before requirements are clear is the single most common failure mode
  — candidates who start drawing boxes in minute 2 end up re-architecting in minute 25 when a
  functional requirement they never asked about (e.g. "does this need to support edits?") invalidates
  their data model.
- **Time budgeting in a 45-minute interview**: roughly 5 min requirements + estimation, 5 min API/data
  model, 10 min high-level design, 20 min deep dives (the interviewer usually steers these), 5 min
  wrap-up/trade-offs. Deep dives are where the signal is — don't over-invest in the high-level diagram
  at the expense of deep-dive time.
- **The first 5 minutes, concretely** — a script:
  1. "Let me make sure I understand the problem — restate it in my own words, and list 2-3 things I'm
     assuming are in scope." *(shows you listen and scope things down)*
  2. Ask 3-5 clarifying questions covering: who are the users, what's the core action, read-heavy or
     write-heavy, any non-functional constraint that changes the design (strong consistency? global
     users? massive scale?). "Roughly what scale are we targeting — is this more like a startup's first
     million users, or Twitter-scale?"
  3. Explicitly separate functional requirements ("what the system does") from non-functional ones
     ("how well it does it" — latency, availability, consistency, scale) and say both out loud.
  4. Only then say "Let me do a quick back-of-envelope estimate before I sketch the design" and move to
     estimation.
- **Common failure modes to avoid**: designing for requirements nobody asked for (over-engineering),
  silently assuming scale instead of asking, going deep on one component in minute 10 before the
  high-level shape exists, and never circling back to trade-offs (every design decision should get a
  one-sentence "I chose X over Y because...").

## Concept: CAP theorem, consistency models (eventual, strong, linearizable)

- 🎯 Asked at Amazon
- **CAP theorem**: in the presence of a network **P**artition, a distributed system must choose between
  **C**onsistency (every read sees the latest write) and **A**vailability (every request gets a
  non-error response). You don't get to choose C vs A in the happy path — only during an actual
  partition; the rest of the time you can have both. This is why "CAP" is really about partition
  behavior, not a permanent system-wide label.
- **Strong consistency**: every read reflects the most recent completed write, as if there were only one
  copy of the data. Requires coordination (e.g. reads and writes both go through a leader, or a quorum),
  which costs latency and can sacrifice availability during a partition.
- **Linearizability**: a stricter, real-time form of strong consistency — operations appear to take
  effect atomically at some point between their start and end, and that global order is consistent with
  real time across all clients. It's the strongest practical consistency model; systems like etcd/ZooKeeper
  target it for coordination primitives (locks, leader election) where "everyone agrees on the exact
  current state" is non-negotiable.
- **Eventual consistency**: writes propagate asynchronously; reads may return stale data for a window,
  but all replicas converge to the same value once writes stop arriving. Much cheaper (low write latency,
  high availability, works fine during partitions) and is the right default for data where staleness for
  seconds is harmless (like counts, follower feeds) but wrong for data where it isn't (bank balances).
- **Picking one in an interview**: name the specific field/operation, not the whole system — "the
  message store needs strong consistency for read-your-own-writes, but the follower count on a profile
  can be eventually consistent" is a much stronger answer than declaring the whole system "AP" or "CP".

```mermaid
flowchart TB
    subgraph Normal["No partition — both C and A achievable"]
        C1[Client] --> N1[Node A] <-->|replicates| N2[Node B]
    end
    subgraph Partitioned["Network partition between A and B"]
        direction LR
        Client1[Client 1] --> PA[Node A]
        Client2[Client 2] --> PB[Node B]
        PA -.x.-|partition: no replication| PB
        PA -->|"CP choice: reject write\n(stay consistent, sacrifice availability)"| R1[Error / unavailable]
        PB -->|"AP choice: accept write\n(stay available, risk stale/divergent reads)"| R2[200 OK, may diverge from A]
    end
```

## Concept: SQL vs NoSQL — when and why

- **What SQL (RDBMS) buys you**: a fixed schema with enforced types/constraints, ACID transactions
  across multiple rows/tables, and rich query flexibility (joins, aggregations, ad-hoc queries) via SQL.
  Best when data is relational (orders ↔ users ↔ products), correctness constraints matter (foreign
  keys, uniqueness), and query patterns aren't fully known upfront.
- **What NoSQL buys you**: schema flexibility (documents can vary shape), horizontal scalability that's
  usually easier to reason about (most NoSQL stores are built partition-first), and data models suited to
  specific access patterns — key-value (Redis, DynamoDB), document (MongoDB), wide-column (Cassandra,
  HBase), graph (Neo4j). The cost is usually weaker consistency/transaction guarantees and query
  flexibility traded for write/read throughput at scale.
- **The real decision axis isn't "SQL vs NoSQL" as a binary** — it's: how relational is the data, how
  well-known are the access patterns, do you need multi-row ACID transactions, and what's the actual
  write/read scale. A social network's user profile can live happily in either; a payments ledger
  strongly wants SQL-style ACID; a chat message store with a single dominant access pattern
  ("get last N messages for a conversation") is a textbook wide-column/key-value fit.
- **Scaling angle**: SQL databases can scale (read replicas, sharding) but it requires more deliberate
  engineering (choosing a shard key, dealing with cross-shard joins/transactions); many NoSQL stores ship
  with horizontal partitioning as a first-class primitive, which is why they're the default pick once a
  single SQL primary clearly can't hold the write volume.
- **In an interview**: don't declare "I'll use NoSQL because it scales better" as a blanket statement —
  justify it against the specific access pattern and consistency requirement you identified in the
  requirements step; interviewers probe this exact justification.

## Concept: Back-of-envelope estimation (QPS, storage, bandwidth)

- 🎯 Asked at Google
- **Why it matters**: estimation isn't about getting an exact number — it's about sanity-checking design
  decisions ("do we actually need sharding at this scale, or is a single beefy Postgres instance fine?")
  and surfacing the dominant constraint (read-heavy? write-heavy? storage-heavy?) before you design.
- **QPS (queries per second)**: start from a daily/monthly active user count and an actions-per-user
  estimate, convert to average QPS via `total requests / 86,400 seconds`, then apply a peak multiplier
  (commonly 2-5x average) for traffic spikes.
- **Storage**: estimate bytes per record, multiply by records/day, multiply by retention period; round
  generously (order-of-magnitude accuracy is the goal, not precision).
- **Bandwidth**: `QPS × average request/response size`, separately for read and write paths since they
  often differ by orders of magnitude (a redirect request is tiny; a video chunk response is not).
- **Worked example — URL shortener**: 100M new URLs/day (writes), 1000:1 read:write ratio.
  - Write QPS: `100,000,000 / 86,400 ≈ 1,160 writes/sec` average.
  - Read QPS: `1,160 × 1000 ≈ 1,160,000 reads/sec` average — this single ratio is the whole ballgame:
    it immediately tells you to optimize aggressively for read latency (cache in front of the DB) and
    that write throughput is comparatively a non-issue.
  - Peak read QPS at a 3x multiplier: `~3.5M reads/sec` — this number alone justifies a CDN/edge cache
    for popular redirects rather than hitting an origin DB per read.
  - Storage: assume ~500 bytes per URL record (long URL + code + metadata). `100M/day × 500B ≈ 50GB/day`;
    over 5 years of retention, `~90TB` — large but very manageable for a sharded key-value store, and this
    number tells you indexing/compression matters more than exotic storage tiers.
- **Rule of thumb numbers worth memorizing**: 1 day ≈ 86,400s (often rounded to ~100K for quick mental
  math); a typical small record/row ≈ 100 bytes–1KB; a single modern DB node handles roughly
  low-thousands to tens-of-thousands of simple QPS depending on query complexity and hardware.

## Concept: API design and data model

- **API design principles**: design endpoints around resources and the client's actual use cases, not
  around your database schema. Prefer REST-ish resource verbs (`POST /orders`, `GET /orders/{id}`) or
  RPC-style (`POST /createOrder`) consistently — pick one style and justify it. Always specify request/
  response shapes explicitly in the interview; vague endpoints ("an API to get feed data") signal you
  haven't thought through the client contract.
- **Idempotency in API design**: any write endpoint that a client might retry (network timeout, mobile
  flakiness) should accept an idempotency key so retries don't double-create — worth mentioning even in
  week 1, it recurs constantly through week 8's payments material.
- **Pagination**: for any endpoint returning a list, decide offset-based (`?page=2&size=20`, simple but
  breaks under concurrent inserts) vs. cursor-based (`?after=<opaque_cursor>`, stable under inserts,
  standard for feeds/infinite-scroll) and say which and why.
- **Data model**: translate functional requirements into entities and relationships before worrying
  about SQL vs NoSQL — get the entities, their key fields, and their access patterns right first, *then*
  pick a storage engine that serves those access patterns well. State the primary access pattern
  explicitly ("we'll always fetch messages by conversation_id, ordered by timestamp") because that
  single sentence usually determines the partition/index key.
- **Common mistake**: designing the data model as a literal SQL schema even when the store will be
  NoSQL, or vice versa — the data model step should be storage-agnostic (entities + relationships +
  access patterns); the storage engine choice comes after, informed by SQL vs NoSQL trade-offs above.

## How to bring this up in the interview

- **When to mention it**: this week's material isn't a "concept to name-drop" — it *is* the shape of the
  entire interview. Use the RESHADED ordering as your internal checklist from second 1, and use CAP/
  consistency-model language whenever you justify a storage or replication choice in a deep dive later.
- **A good opening line**: "Before I design anything, let me clarify requirements and do a quick
  estimate — that'll tell me whether this is a single-server problem or a genuinely distributed-systems
  problem." This signals structure immediately and buys you permission to ask questions instead of
  rushing to whiteboard.
- **A question to ask the interviewer early**: "Should I optimize this design assuming huge scale (like
  the real company), or should I start simple and call out where I'd change the design as scale grows?"
  — this single question often changes the whole 45 minutes, and asking it explicitly is itself a signal
  of seniority.
- **Common follow-up 1**: *"You said eventual consistency is fine here — what would actually go wrong if
  it weren't, concretely?"* Answer with a specific bad outcome tied to the field in question (e.g. "a
  user who just posted might not see their own post for a few seconds on a different device — annoying
  but not harmful; contrast with a balance check that must never show stale funds").
- **Common follow-up 2**: *"Your storage estimate assumed X bytes per record — where did that number come
  from, and how does the design change if it's 10x larger?"* Answer by breaking the record down field by
  field, and name the specific design lever you'd pull (e.g. "media wouldn't live in the row at all —
  it'd move to blob storage with a URL reference in the row instead").

## Designs to revisit after this week

- [URL Shortener](../designs/url-shortener/README.md) — 🎯 Asked at Flipkart — a small, self-contained
  design that's ideal for practicing the RESHADED framework end to end (requirements → estimation → API
  → data model → high-level design → deep dives) before tackling week 4's larger systems.

## Practice prompt
Pick any system you use daily (e.g. a note-taking app's sync feature). Spend exactly 5 minutes running
through the first-5-minutes script above out loud: restate the problem, ask yourself 3-5 clarifying
questions as if an interviewer answered them, separate functional from non-functional requirements, then
do a back-of-envelope QPS and storage estimate. Notice which of your early assumptions the estimate
either confirms or invalidates — that's the muscle this week is building.
