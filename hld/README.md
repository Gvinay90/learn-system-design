# 8-Week HLD Learning Path — Roadmap

Goal: build a rock-solid System Design foundation, then apply it to 15+ of the most commonly asked
HLD interview questions, ending with mock interviews.

Each week folder has a README with: concept summaries, curated hellointerview.com + YouTube links,
"asked at" tags, and a practice prompt. Designs live in [`designs/`](designs/) — each design gets a
full write-up (requirements → API → data model → high-level diagram → deep dives → trade-offs); designs
with a natural runnable core (rate limiter, consistent hashing, LRU cache, autocomplete) also get a
small Go implementation you can run and test.

## Weekly checklist

- [ ] **Week 1 — [Foundations & Mental Models](week-01-foundations/README.md)**
  - How to approach any HLD interview in top product-based companies
  - CAP theorem, consistency models (eventual, strong, linearizable) — 🎯 Asked at Amazon
  - SQL vs NoSQL — when and why
  - Back-of-envelope estimation (QPS, storage, bandwidth) — 🎯 Asked at Google
  - API Design, Data Model
- [ ] **Week 2 — [Core Infrastructure Building Blocks](week-02-infra/README.md)**
  - Networking concepts, load balancers (L4 vs L7, consistent hashing)
  - CDN, DNS, reverse proxies — 🎯 Asked at Netflix
  - Caching strategies (write-through, write-back, cache-aside, TTL)
  - Redis deep dive — data structures, eviction policies, pub/sub (best explored hands-on with a real Redis instance) — 🎯 Asked at Spotify
- [ ] **Week 3 — [Database & Messaging Systems](week-03-db-messaging/README.md)**
  - Sharding strategies (range, hash, directory-based) — 🎯 Asked at Microsoft
  - Replication (leader-follower, multi-leader, leaderless)
  - Database indexing internals (B-Tree, LSM tree)
  - Message queues & event-driven architecture (Kafka, SQS) — 🎯 Asked at Uber
  - Async processing, fan-out patterns
- [ ] **Week 4 — [High-Traffic & Search Systems](week-04-high-traffic/README.md)**
  - Concept: Rate limiting at scale (Token Bucket, Sliding Window at distributed level)
  - Design: [URL Shortener](designs/url-shortener/README.md) — 🎯 Asked at Flipkart
  - Design: [Rate Limiter](designs/rate-limiter/README.md) — 🎯 Asked at Razorpay
  - Design: [Search Autocomplete](designs/search-autocomplete/README.md)
- [ ] **Week 5 — [Feed & Streaming Systems](week-05-feed-streaming/README.md)**
  - Concept: Circuit breaker, bulkhead, retry with backoff
  - Design: [Twitter / Instagram Feed](designs/twitter-feed/README.md) — 🎯 Asked at Meesho
  - Design: [YouTube / Netflix](designs/youtube-netflix/README.md)
  - Design: [WhatsApp / Messaging System](designs/whatsapp-messaging/README.md)
- [ ] **Week 6 — [Storage & Notification Systems](week-06-storage-notifications/README.md)**
  - Concept: WebSockets, Server-Sent Events & Long Polling
  - Design: [Key-Value Store](designs/key-value-store/README.md) (best explored hands-on — see design doc)
  - Design: [Distributed Cache](designs/lru-distributed-cache/README.md)
  - Design: [Notification System](designs/notification-system/README.md) (best explored hands-on — see design doc) — 🎯 Asked at Swiggy
- [ ] **Week 7 — [Location & Scheduling Systems](week-07-location-scheduling/README.md)**
  - Concept: Geospatial indexing — Geohashing, QuadTrees, Uber's H3
  - Design: [Uber / Ride-Sharing](designs/uber-ride-sharing/README.md) — 🎯 Asked at Zomato
  - Design: [Google Maps](designs/google-maps/README.md)
  - Design: [Distributed Task Scheduler](designs/task-scheduler/README.md)
- [ ] **Week 8 — [Payments, Reliability & Mock Interviews](week-08-payments-reliability/README.md)**
  - Concept: Distributed transactions (2PC, Saga), idempotency, exactly-once semantics
  - Design: [Payment System / UPI-like Gateway](designs/payment-gateway-upi/README.md) — 🎯 Asked at PhonePe
  - Design: [Leaderboard / Distributed Counters](designs/leaderboard-counters/README.md) (best explored hands-on — see design doc)
  - Design: [Logging & Monitoring System](designs/logging-monitoring/README.md)
  - 2-3 full mock HLD interviews

## All 15 designs (quick index)

| Design | Asked at | Runnable demo |
|---|---|---|
| [URL Shortener](designs/url-shortener/README.md) | Flipkart | README only |
| [Twitter/Instagram Feed](designs/twitter-feed/README.md) | Meesho | README only |
| [YouTube/Netflix](designs/youtube-netflix/README.md) | Netflix | README only |
| [WhatsApp/Messaging](designs/whatsapp-messaging/README.md) | Microsoft | README only |
| [Uber/Ride-Sharing](designs/uber-ride-sharing/README.md) | Zomato | README only |
| [Google Maps](designs/google-maps/README.md) | Google | README only |
| [Notification System](designs/notification-system/README.md) | Swiggy | README only |
| [Rate Limiter](designs/rate-limiter/README.md) | Razorpay | ✅ Go |
| [Key-Value Store](designs/key-value-store/README.md) | Amazon | README only |
| [Distributed Cache](designs/lru-distributed-cache/README.md) | Uber | ✅ Go |
| [Search Autocomplete](designs/search-autocomplete/README.md) | Google | ✅ Go |
| [Payment Gateway/UPI](designs/payment-gateway-upi/README.md) | PhonePe | README only |
| [Distributed Task Scheduler](designs/task-scheduler/README.md) | Spotify | README only |
| [Leaderboard/Counters](designs/leaderboard-counters/README.md) | Netflix | README only |
| [Logging & Monitoring](designs/logging-monitoring/README.md) | Microsoft | README only |
| [Consistent Hashing](designs/consistent-hashing/README.md) *(cross-cutting concept, not in the 15)* | — | ✅ Go |
