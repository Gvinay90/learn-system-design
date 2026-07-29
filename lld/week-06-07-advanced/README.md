# Week 6-7 — Advanced LLD Problems: Core Systems & Patterns

Part of the [8-week LLD program roadmap](../README.md).

## How this week differs from weeks 4-5

Weeks 4-5 were end-to-end business systems (a ride-sharing app, a food-delivery app) with several actors
and a broad surface area. This week is the opposite shape: each problem is a single, reusable
**infrastructure component** — the kind of thing you'd otherwise import as a library — and the interview
is almost entirely about getting *that one component's* API, complexity, and thread-safety exactly
right. Expect tighter scope and a higher bar on Big-O and concurrency correctness than weeks 4-5.

Several of these components have a direct HLD counterpart in this repo — the LLD problem here is "design
the class," the HLD design is "now distribute it across machines." Worth reading both back-to-back.

## Problems this week

- [LRU Cache](../problems/lru-cache/README.md) — 🎯 Asked at Uber — HashMap + doubly linked list, O(1)
  get/put. HLD counterpart: [Distributed Cache](../../hld/designs/lru-distributed-cache/README.md).
- [Rate Limiter](../problems/rate-limiter/README.md) — 🎯 Asked at Razorpay — Token Bucket vs. Sliding
  Window behind a common interface. HLD counterpart: [Rate Limiter (distributed)](../../hld/designs/rate-limiter/README.md).
- [Task Scheduler](../problems/task-scheduler/README.md) — 🎯 Asked at Spotify — priority-based
  execution, delayed jobs, retries. HLD counterpart: [Distributed Task Scheduler](../../hld/designs/task-scheduler/README.md).
- [Notification System](../problems/notification-system/README.md) — 🎯 Asked at Swiggy — multi-channel
  dispatch (email/SMS/push), retry-on-failure, template rendering. HLD counterpart: [Notification System (multi-channel at scale)](../../hld/designs/notification-system/README.md).
- [Logging Framework](../problems/logging-framework/README.md) — 🎯 Asked at Microsoft — log levels,
  pluggable appenders, extensibility.
- [In-Memory File System](../problems/in-memory-file-system/README.md) — 🎯 Asked at Netflix —
  directory-tree composite structure, path resolution.
- [Vending Machine](../problems/vending-machine/README.md) — 🎯 Asked at Meesho — State pattern,
  inventory + payment handling.

## How to approach these in an interview

- **Lead with complexity, not just correctness.** For LRU Cache and Rate Limiter especially, a
  correct-but-O(n) answer is treated as an incomplete answer — say the target complexity out loud
  before you design toward it (see [LRU Cache's "Key trade-offs"](../problems/lru-cache/README.md#key-trade-offs--talking-points)
  for why HashMap+DLL beats a tree-based structure here).
- **Assume concurrent access by default.** These are the components most likely to actually be shared
  across threads in production (a rate limiter guarding an API, a logger called from every request) —
  mention thread-safety in your requirements pass instead of waiting to be asked, and reuse the
  primitives from [this repo's concurrency doc](../concurrency/README.md).
- **Name the pattern and justify it, don't just apply it.** Vending Machine (State), Rate Limiter/
  Notification channels (Strategy), Logging Framework appenders (Strategy, arguably Chain of
  Responsibility) — interviewers at this level expect you to say *why* that pattern over a simpler
  if/else, not just produce code that happens to match it.
- **Have the "make it distributed" follow-up ready.** Since four of these seven have an HLD counterpart
  in this repo, expect "how would this change with multiple instances behind a load balancer" as a
  closing question — the honest answer usually involves moving in-memory state to Redis/a shared store
  and re-deriving which operations need to become atomic/distributed locks.

## Practice prompt

Pick two problems whose HLD counterpart exists in this repo (e.g. Rate Limiter and Task Scheduler).
Design the LLD class first, cold, with a 25-minute timebox. Then read that problem's HLD counterpart
doc and write down, in 5 bullet points, exactly which pieces of your in-memory design stop working with
multiple server instances and what you'd replace them with — that transition is the actual follow-up
question you'll get live.
