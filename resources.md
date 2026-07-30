# Consolidated Reference List

Every link below is also cited in its own topic's README (with more context on *why* it's the right
read). This file exists purely as a single scannable index — every link here was verified to actually
resolve (live HTTP check) as of the last update; if a site restructures its URLs later, search the
site/channel directly for the topic name.

Primary sources:
- **[hellointerview.com](https://www.hellointerview.com/learn/system-design/in-a-hurry/introduction)** — System Design in a Hurry guide (HLD fundamentals + per-topic deep dives) and a dedicated Low-Level Design track (framework + problem breakdowns). Subscription unlocks the full problem breakdowns and mock interview tooling; many community-submitted breakdowns are free.
- **YouTube** — hellointerview's own channel, plus freely available system design/LLD content from ByteByteGo, Gaurav Sen, Exponent, Geekific (design patterns), and individual FAANG engineers.
- **[refactoring.guru](https://refactoring.guru/design-patterns)** — the reference used for all 12 design patterns in `lld/week-02-03-patterns/`.

## HLD — Weeks

| Week | hellointerview | YouTube |
|---|---|---|
| [Week 1 — Foundations & Mental Models](hld/week-01-foundations/README.md) | *(framework-level content — see individual designs below for topic references)* | |
| [Week 2 — Core Infrastructure Building Blocks](hld/week-02-infra/README.md) | *(framework-level content — see individual designs below for topic references)* | |
| [Week 3 — Database & Messaging Systems](hld/week-03-db-messaging/README.md) | *(framework-level content — see individual designs below for topic references)* | |
| [Week 4 — High-Traffic & Search Systems](hld/week-04-high-traffic/README.md) | *(framework-level content — see individual designs below for topic references)* | |
| [Week 5 — Feed & Streaming Systems](hld/week-05-feed-streaming/README.md) | [System Design in a Hurry — Introduction](https://www.hellointerview.com/learn/system-design/in-a-hurry/introduction) | [Distributed Transactions Explained: 2PC vs Saga](https://www.youtube.com/watch?v=DOFflggE_0Q) |
| [Week 6 — Storage & Notification Systems](hld/week-06-storage-notifications/README.md) | [Real-time Updates Pattern](https://www.hellointerview.com/learn/system-design/patterns/realtime-updates) | [Notification Service for Billions of Users](https://www.youtube.com/watch?v=CUwt9_l0DOg) |
| [Week 7 — Location & Scheduling Systems](hld/week-07-location-scheduling/README.md) | [Proximity Search in System Design Interviews](https://www.hellointerview.com/learn/system-design/deep-dives/proximity-search) | |
| [Week 8 — Payments, Reliability & Mock Interviews](hld/week-08-payments-reliability/README.md) | *(see [Payment Gateway/UPI](hld/designs/payment-gateway-upi/README.md) below)* | [2PC vs Saga Pattern](https://www.youtube.com/watch?v=DOFflggE_0Q) |

## HLD — 16 Designs

| Design | hellointerview | YouTube |
|---|---|---|
| [Consistent Hashing](hld/designs/consistent-hashing/README.md) | [Concept guide](https://www.hellointerview.com/learn/system-design/core-concepts/consistent-hashing) · [Quick reference](https://www.hellointerview.com/learn/system-design/core-concepts/consistent-hashing/quick-reference) | [Consistent Hashing Explained](https://www.youtube.com/watch?v=CfwmTPzTdUc) |
| [URL Shortener](hld/designs/url-shortener/README.md) | [Design a URL Shortener Like Bitly](https://www.hellointerview.com/learn/system-design/problem-breakdowns/bitly) | [Design Bitly w/ Ex-Meta Staff Engineer](https://www.youtube.com/watch?v=iUU4O1sWtJA) |
| [Rate Limiter](hld/designs/rate-limiter/README.md) | [Distributed Rate Limiter breakdown](https://www.hellointerview.com/learn/system-design/problem-breakdowns/distributed-rate-limiter) · [Rate Limiter LLD](https://www.hellointerview.com/learn/low-level-design/problem-breakdowns/rate-limiter) | [Ex-Meta Staff Engineer walkthrough](https://www.youtube.com/watch?v=MIJFyUPG4Z4) |
| [Search Autocomplete](hld/designs/search-autocomplete/README.md) | [Design Typeahead Search](https://www.hellointerview.com/community/questions/typeahead-search-system/cm7l2wazy00t7105qdvnemtwy) | [Design Search Autocomplete System](https://www.youtube.com/watch?v=TZ_LSourdUc) |
| [Twitter / Instagram Feed](hld/designs/twitter-feed/README.md) | [Design Facebook's News Feed](https://www.hellointerview.com/learn/system-design/problem-breakdowns/fb-news-feed) | [Design FB News Feed w/ Ex-Meta Sr. Manager](https://www.youtube.com/watch?v=Qj4-GruzyDU) |
| [YouTube / Netflix](hld/designs/youtube-netflix/README.md) | [Design a Video Streaming Platform Like YouTube](https://www.hellointerview.com/learn/system-design/problem-breakdowns/youtube) | [Design YouTube w/ Ex-Meta Staff Engineer](https://www.youtube.com/watch?v=IUrQ5_g3XKs) |
| [WhatsApp / Messaging](hld/designs/whatsapp-messaging/README.md) | [Design a Messaging App Like WhatsApp](https://www.hellointerview.com/learn/system-design/problem-breakdowns/whatsapp) | [Design WhatsApp w/ Ex-Meta Sr. Manager](https://www.youtube.com/watch?v=cr6p0n0N-VA) |
| [Key-Value Store](hld/designs/key-value-store/README.md) | [Design a Key-Value Store](https://www.hellointerview.com/community/questions/key-value-store/cm8gcrkz800b7epmpcj06fkwk) · [DynamoDB deep dive](https://www.hellointerview.com/learn/system-design/deep-dives/dynamodb) · [Redis deep dive](https://www.hellointerview.com/learn/system-design/deep-dives/redis) | [Designing a Distributed Key-Value Store (Dynamo-style)](https://www.youtube.com/watch?v=j8iDY_RudJw) |
| [Distributed Cache](hld/designs/lru-distributed-cache/README.md) | [Design a Distributed Cache Like Redis](https://www.hellointerview.com/learn/system-design/problem-breakdowns/distributed-cache) · [Community breakdown](https://www.hellointerview.com/community/questions/distributed-cache-system/cm6d9gnep03c46hpqrwc062ir) | [Design a Distributed LRU Cache (full mock interview)](https://www.youtube.com/watch?v=lZ5QuFLCVn0) |
| [Notification System](hld/designs/notification-system/README.md) | [Design a Notification System](https://www.hellointerview.com/learn/system-design/problem-breakdowns/notification-system) · [Real-time Updates Pattern](https://www.hellointerview.com/learn/system-design/patterns/realtime-updates) | [Notification System: SMS, Email, Push](https://www.youtube.com/watch?v=1E3oeYkJ1P8) |
| [Uber / Ride-Sharing](hld/designs/uber-ride-sharing/README.md) | [Design a Ride-Sharing Service Like Uber](https://www.hellointerview.com/learn/system-design/problem-breakdowns/uber) | [Design Uber w/ Ex-Meta Staff Engineer](https://www.youtube.com/watch?v=lsKU38RKQSo) |
| [Google Maps](hld/designs/google-maps/README.md) | [Design Google Maps](https://www.hellointerview.com/community/questions/map-service-design/cm7wcazsa010t133rbusc7igc) · [Proximity Search deep dive](https://www.hellointerview.com/learn/system-design/deep-dives/proximity-search) | [Google Maps System Design Interview Question](https://www.youtube.com/watch?v=1pmcoh4hc_A) |
| [Distributed Task Scheduler](hld/designs/task-scheduler/README.md) | [Design a Distributed Job Scheduler Like Airflow](https://www.hellointerview.com/learn/system-design/problem-breakdowns/job-scheduler) | [Distributed Job Scheduler Deep Dive w/ Google SWE](https://www.youtube.com/watch?v=WTxG5880EH8) |
| [Payment Gateway / UPI](hld/designs/payment-gateway-upi/README.md) | [Design a Payment System Like Stripe](https://www.hellointerview.com/learn/system-design/problem-breakdowns/payment-system) | [Payment System: PSP, Idempotency, Saga, Ledger](https://www.youtube.com/watch?v=Zrd1wdkAZLY) |
| [Leaderboard / Distributed Counters](hld/designs/leaderboard-counters/README.md) | [Design an Online Game Leaderboard](https://www.hellointerview.com/community/questions/cm4t0qbr9004988ilmum8jm06) · [YouTube's Top K Videos](https://www.hellointerview.com/learn/system-design/problem-breakdowns/top-k) | [Real-Time Leaderboards: Redis Sorted Sets](https://www.youtube.com/watch?v=9yEPu8oSrhI) |
| [Logging & Monitoring](hld/designs/logging-monitoring/README.md) | [Design a Metrics Monitoring Platform Like Datadog](https://www.hellointerview.com/learn/system-design/problem-breakdowns/metrics-monitoring) · [Log Collection System](https://www.hellointerview.com/community/questions/log-ingestion-system/cmp3gf5f203st08adwrv5jev4) | [Real Time Metrics Database at Datadog](https://www.youtube.com/watch?v=uQrRbvLyJ4M) |

## LLD — Foundations & Concurrency

| Topic | hellointerview | YouTube |
|---|---|---|
| [Week 1 — LLD Foundations](lld/week-01-foundations/README.md) | [Design Principles](https://www.hellointerview.com/learn/low-level-design/in-a-hurry/design-principles) · [How to Prepare for an LLD Interview](https://www.hellointerview.com/blog/how-to-prepare-lld) | [LLD Interview Approach](https://www.youtube.com/watch?v=GixkAcu3eEw) |
| [Concurrency Primitives](lld/concurrency/README.md) | [Introduction to Concurrency](https://www.hellointerview.com/learn/low-level-design/concurrency/intro) | [Java ExecutorService — Introduction](https://www.youtube.com/watch?v=6Oo-9Can3H8) |
| General LLD delivery framework | [LLD Interview Delivery Framework](https://www.hellointerview.com/learn/low-level-design/in-a-hurry/delivery) — cited from nearly every problem README | |

## LLD — 12 Design Patterns

All 12 use [refactoring.guru](https://refactoring.guru/design-patterns) as the primary conceptual reference, plus [Hello Interview's Design Patterns overview](https://www.hellointerview.com/learn/low-level-design/in-a-hurry/patterns); table below has the pattern-specific refactoring.guru page and YouTube walkthrough (mostly Geekific's Java series).

| Pattern | refactoring.guru | YouTube |
|---|---|---|
| [Singleton](lld/week-02-03-patterns/singleton/README.md) | [Singleton](https://refactoring.guru/design-patterns/singleton) | [Singleton Explained (Geekific)](https://www.youtube.com/watch?v=hUE_j6q0LTQ) |
| [Factory](lld/week-02-03-patterns/factory/README.md) | [Factory Method](https://refactoring.guru/design-patterns/factory-method) | [Factory Method Explained (Geekific)](https://www.youtube.com/watch?v=EcFVTgRHJLM) |
| [Builder](lld/week-02-03-patterns/builder/README.md) | [Builder](https://refactoring.guru/design-patterns/builder) | [Builder Explained (with code)](https://www.youtube.com/watch?v=ALzvPK9_r0A) |
| [Prototype](lld/week-02-03-patterns/prototype/README.md) | [Prototype](https://refactoring.guru/design-patterns/prototype) | [Prototype Introduction](https://www.youtube.com/watch?v=f1BG1tkqZQU) |
| [Adapter](lld/week-02-03-patterns/adapter/README.md) | [Adapter](https://refactoring.guru/design-patterns/adapter) | [Adapter Explained (Geekific)](https://www.youtube.com/watch?v=e1i4CQCZeaQ) |
| [Decorator](lld/week-02-03-patterns/decorator/README.md) | [Decorator](https://refactoring.guru/design-patterns/decorator) | [Decorator Explained (Geekific)](https://www.youtube.com/watch?v=GtWvgTfxRDI) |
| [Facade](lld/week-02-03-patterns/facade/README.md) | [Facade](https://refactoring.guru/design-patterns/facade) | [Facade: Easy Guide for Beginners](https://www.youtube.com/watch?v=xv74RW5IAvo) |
| [Proxy](lld/week-02-03-patterns/proxy/README.md) | [Proxy](https://refactoring.guru/design-patterns/proxy) | [Proxy Explained (Geekific)](https://www.youtube.com/watch?v=TS5i-uPXLs8) |
| [Observer](lld/week-02-03-patterns/observer/README.md) | [Observer](https://refactoring.guru/design-patterns/observer) | [Observer Explained (Geekific)](https://www.youtube.com/watch?v=-oLDJ2dbadA) |
| [Strategy](lld/week-02-03-patterns/strategy/README.md) | [Strategy](https://refactoring.guru/design-patterns/strategy) | [Strategy Explained (Geekific)](https://www.youtube.com/watch?v=Nrwj3gZiuJU) |
| [Chain of Responsibility](lld/week-02-03-patterns/chain-of-responsibility/README.md) | [Chain of Responsibility](https://refactoring.guru/design-patterns/chain-of-responsibility) | [CoR Explained & Implemented](https://www.youtube.com/watch?v=FafNcoBvVQo) |
| [State](lld/week-02-03-patterns/state/README.md) | [State](https://refactoring.guru/design-patterns/state) | [State Explained (Geekific)](https://www.youtube.com/watch?v=abX4xzaAsoc) |

## LLD — 18 Problems

| Problem | hellointerview | YouTube |
|---|---|---|
| [Chess](lld/problems/chess/README.md) | [Design an Online Chess Platform](https://www.hellointerview.com/community/questions/online-chess-platform/cm4szywgj003ivszmblm3pmoa) | [Design Chess — OOD Interview](https://www.youtube.com/watch?v=C4tyv9k0i9M) |
| [Snake & Ladder](lld/problems/snake-and-ladder/README.md) | *(no dedicated page — see [Elevator LLD](https://www.hellointerview.com/learn/low-level-design/problem-breakdowns/elevator) for the same turn-driven-simulation skills)* | [Snake and Ladder — LLD Interview](https://www.youtube.com/watch?v=nkeXstII8vQ) |
| [Tic-Tac-Toe](lld/problems/tic-tac-toe/README.md) | [LeetCode 348: Design Tic-Tac-Toe](https://www.hellointerview.com/community/questions/design-tic-tac-toe/cm5eh7nri04wo838oitc9peu4) | [Tic Tac Toe — LLD/OOD Interview](https://www.youtube.com/watch?v=V7aFobyuLrU) |
| [Parking Lot](lld/problems/parking-lot/README.md) | [Parking Lot LLD breakdown](https://www.hellointerview.com/learn/low-level-design/problem-breakdowns/parking-lot) | [Parking Lot — LLD Interview](https://www.youtube.com/watch?v=FC-rVMlsbHk) |
| [Elevator System](lld/problems/elevator-system/README.md) | [Elevator LLD breakdown](https://www.hellointerview.com/learn/low-level-design/problem-breakdowns/elevator) | [Elevator System Design — LLD Interview](https://www.youtube.com/watch?v=x8CtiPRWq04) |
| [Trading System](lld/problems/trading-system/README.md) | [Design a Stock Trading Platform Like Robinhood](https://www.hellointerview.com/learn/system-design/problem-breakdowns/robinhood) *(system-design-level; this problem is the order-book/matching-engine class design underneath it)* | [Design a Stock Exchange — LLD Interview](https://www.youtube.com/watch?v=XY6pRVpB1Rw) |
| [Splitwise](lld/problems/splitwise/README.md) | [Design an LLD of Splitwise](https://www.hellointerview.com/community/questions/splitwise-lld/cm6jwwh6700bxui4bzs4jmddl) | [Design Splitwise — LLD Interview](https://www.youtube.com/watch?v=Yhu-1H8UWv4) |
| [Food Delivery](lld/problems/food-delivery/README.md) | [Design DoorDash](https://www.hellointerview.com/community/questions/food-delivery-platform/cm5omjbka00023b6q9bfevwjl) | [Food Delivery App like Swiggy/Zomato — LLD](https://www.youtube.com/watch?v=Zy2QQ6L0z3s) |
| [Ride-Sharing](lld/problems/ride-sharing/README.md) | [Design a Ride-Sharing Service Like Uber](https://www.hellointerview.com/learn/system-design/problem-breakdowns/uber) *(system-design-level; this problem is the Trip/Rider/Driver class design underneath it)* | [Uber System Design — LLD Interview](https://www.youtube.com/watch?v=oX6NRQtqvIY) |
| [Kafka LLD](lld/problems/kafka-lld/README.md) | [Kafka Deep Dive for System Design Interviews](https://www.hellointerview.com/learn/system-design/deep-dives/kafka) | [Apache Kafka Explained](https://www.youtube.com/watch?v=uvb00oaa3k8) |
| [Payment Gateway](lld/problems/payment-gateway/README.md) | [Design a Payment System Like Stripe](https://www.hellointerview.com/learn/system-design/problem-breakdowns/payment-system) | [Payment Processing: Idempotency, Security](https://www.youtube.com/watch?v=ai7RH1DuoMg) |
| [LRU Cache](lld/problems/lru-cache/README.md) | [Implement LRU Cache](https://www.hellointerview.com/community/questions/lru-cache-implementation/cmk5avhlo00x708adxjp1vuil) · [Distributed Cache HLD](https://www.hellointerview.com/learn/system-design/problem-breakdowns/distributed-cache) | [LRU Cache — Design & Coding Interview](https://www.youtube.com/watch?v=JEABxEdfV5Q) |
| [Rate Limiter](lld/problems/rate-limiter/README.md) | [Rate Limiter LLD breakdown](https://www.hellointerview.com/learn/low-level-design/problem-breakdowns/rate-limiter) | [Distributed Rate Limiter w/ Ex-Meta Staff Engineer](https://www.youtube.com/watch?v=MIJFyUPG4Z4) |
| [Task Scheduler](lld/problems/task-scheduler/README.md) | [Distributed Job Scheduler Like Airflow](https://www.hellointerview.com/learn/system-design/problem-breakdowns/job-scheduler) *(system-design-level)* · [Coordination in LLD](https://www.hellointerview.com/learn/low-level-design/concurrency/coordination) | [Job Scheduler — System Design Interview](https://www.youtube.com/watch?v=WTxG5880EH8) |
| [Notification System](lld/problems/notification-system/README.md) | [Design a Notification System](https://www.hellointerview.com/learn/system-design/problem-breakdowns/notification-system) | [Notification System: SMS, Email, Push](https://www.youtube.com/watch?v=1E3oeYkJ1P8) |
| [Logging Framework](lld/problems/logging-framework/README.md) | [Logging Service LLD breakdown](https://www.hellointerview.com/learn/low-level-design/problem-breakdowns/logging-service) | [Design a Logging Framework — LLD Interview](https://www.youtube.com/watch?v=xpDnVSmNFX0) |
| [In-Memory File System](lld/problems/in-memory-file-system/README.md) | [File System LLD breakdown](https://www.hellointerview.com/learn/low-level-design/problem-breakdowns/file-system) | [In-Memory File System — LLD Interview](https://www.youtube.com/watch?v=DQqfNwbeXvE) |
| [Vending Machine](lld/problems/vending-machine/README.md) | [Design a Vending Machine System](https://www.hellointerview.com/community/questions/beverage-vending-system/cmkr2w9l80cph08adwuu5vnfh) | [Vending Machine System Design — LLD Interview](https://www.youtube.com/watch?v=Qk2ze3gYQU8) |

Every URL above was checked with a live HTTP request and confirmed to resolve. If you find one that's gone
stale, that page's own README (linked in the first column) is the source of truth to re-derive it from.
