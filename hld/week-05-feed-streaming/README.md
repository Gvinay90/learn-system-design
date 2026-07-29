# Week 5 — Feed & Streaming Systems

Part of the [8-week HLD learning path](../README.md).

## Concept: Circuit breaker, bulkhead, retry with backoff

- **Circuit breaker**: wraps calls to a downstream/external service; after a failure-rate threshold trips,
  it "opens" and fails fast (no network call) for a cooldown window, letting the downstream recover
  instead of getting hammered by retries. Half-open state periodically lets a trial request through to
  probe recovery before fully closing again.
- **Retry with exponential backoff + jitter**: transient failures (timeouts, brief unavailability) often
  succeed on a second attempt. Backoff (doubling delay each attempt) avoids hammering a struggling
  service; jitter (randomizing the delay) prevents synchronized "thundering herd" retries across many
  clients.
- **Bulkhead**: isolate resources (thread pools, connection pools) per downstream dependency so one
  overloaded/slow dependency can't exhaust resources needed to serve requests to other, healthy
  dependencies — named after ship compartments that contain flooding to one section.
- **Composing them**: production resilience typically layers all three — bulkhead limits concurrent
  calls to a dependency, retry-with-backoff handles transient blips within that limit, and circuit
  breaker stops retrying altogether once the dependency is clearly down.
- **Why feed/streaming systems care**: a feed-generation service fans out to many downstream services
  (ranking, media, ads); without these patterns, one slow dependency degrades the entire feed's latency
  or availability instead of just that one widget.

**References**
- Background: [System Design in a Hurry — Introduction](https://www.hellointerview.com/learn/system-design/in-a-hurry/introduction) (no dedicated hellointerview.com page found specifically for circuit-breaker/bulkhead/retry at time of writing)
- Watch: [Distributed Transactions Explained: 2 Phase Commit vs Saga Pattern (YouTube)](https://www.youtube.com/watch?v=DOFflggE_0Q) *(covers resilience/retry concepts in a distributed-transactions context; see week 8 for the dedicated saga/2PC treatment)*

## Designs this week

- [Twitter / Instagram Feed](../designs/twitter-feed/README.md) — 🎯 Asked at Meesho *(built by sibling agent)*
- [YouTube / Netflix](../designs/youtube-netflix/README.md) *(built by sibling agent)*
- [WhatsApp / Messaging System](../designs/whatsapp-messaging/README.md) *(built by sibling agent)*

## Practice prompt
Design a feed-generation service (like the Twitter/Instagram feed this week covers) that fans out to a
ranking service, a media-metadata service, and an ads service before returning a response. The ranking
service occasionally gets slow (p99 spikes to 5s) under load. Whiteboard exactly where you'd put a
circuit breaker, what its trip threshold and cooldown should be, and what the feed shows a user while
the ranking service's circuit is open (a graceful degradation — e.g. reverse-chronological fallback —
beats a failed page load).
