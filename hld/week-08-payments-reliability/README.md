# Week 8 — Payments, Reliability & Mock Interviews

Part of the [8-week HLD learning path](../README.md).

## Concept: Distributed transactions (2PC, Saga), idempotency, exactly-once semantics

- **Two-phase commit (2PC)**: a coordinator asks all participants to "prepare" (vote yes/no), then tells
  all to "commit" only if everyone voted yes. Gives strong (ACID-like) consistency across systems, but
  is blocking (participants hold locks while awaiting the coordinator's decision) and fragile (a
  coordinator crash mid-protocol can leave participants blocked indefinitely) — rarely viable when
  participants are external systems (other companies' banks/services) that won't join your coordinator.
- **Saga pattern**: replaces one distributed ACID transaction with a sequence of local transactions, each
  in its own system, plus an explicit compensating action for each step to undo it if a later step fails.
  Trades strong consistency for eventual consistency and non-blocking execution — the standard answer
  for cross-service/cross-company distributed transactions (e.g. payments spanning two banks).
- **Orchestration vs. choreography**: an orchestrated saga has a central coordinator explicitly calling
  each step and deciding on compensation (easier to trace/debug); a choreographed saga has each service
  react to events from the previous step with no central coordinator (more decoupled, harder to observe
  as a single flow).
- **Idempotency**: the property that processing the same request/command twice has the same effect as
  processing it once. Implemented via a client-supplied idempotency key, checked-and-stored atomically
  (unique constraint) before executing — essential for safe retries, since a network timeout doesn't tell
  the client whether the original request actually completed server-side.
- **Exactly-once semantics**: true exactly-once delivery across a network is not achievable in general
  (a message can always be lost or duplicated); what's actually built is at-least-once delivery combined
  with idempotent processing, which produces an exactly-once *effect* even though the message itself may
  arrive more than once.

**References**
- Watch: [Distributed Transactions Explained: 2 Phase Commit vs Saga Pattern (YouTube)](https://www.youtube.com/watch?v=DOFflggE_0Q)
- Read: [System Design: Distributed Transactions — Two-Phase Commit, Saga Pattern, and the Outbox Pattern](https://www.techinterview.org/post/3233465289/system-design-distributed-transactions/) (no dedicated hellointerview.com problem-breakdown page found specifically for 2PC/saga/idempotency as a standalone concept at time of writing — see the [Payment Gateway/UPI design](../designs/payment-gateway-upi/README.md) this week for the applied, hellointerview-referenced walkthrough)

## Designs this week

- [Payment System / UPI-like Gateway](../designs/payment-gateway-upi/README.md) — 🎯 Asked at PhonePe
- [Leaderboard / Distributed Counters](../designs/leaderboard-counters/README.md) *(best learned hands-on with a real Redis instance — see that design's write-up)*
- [Logging & Monitoring System](../designs/logging-monitoring/README.md) — 🎯 Asked at Microsoft

## Practice prompt
Design the payment flow from this week's UPI-like gateway end to end as a saga: a transfer debits the
payer's PSP and credits the payee's PSP, two systems that cannot share a transaction. Whiteboard the
state machine (states + transitions), where the idempotency key is checked, and exactly which failure
at which step should trigger a compensating transaction vs. a safe retry of the same step.

## Self-run mock HLD interview format (45-60 min)

Weeks 1-7 built the concepts and the 16 designs; this week is about running full designs under interview
conditions, timed, without stopping to look anything up.

1. **Pick a design blind** — use the rotation below, or have a peer/friend pick one you haven't reread
   in the last 2 weeks.
2. **Set a single 45-60 minute timer, no pausing**, following the time budget in that design's own
   "How to narrate this in the interview" section (every design doc in [`designs/`](../designs/) has
   one) — roughly: 5 min requirements/clarifying questions, 5 min estimation, 10 min API+data model,
   15-20 min high-level design + diagram, 10-15 min deep dives and trade-offs.
3. **Talk out loud the whole time**, even solo — narrate every decision as if an interviewer is
   listening; silent design time is the single biggest tell of an unpracticed candidate.
4. **Stop at the timer**, even mid-deep-dive — an interview doesn't pause for you either.

If running this with a peer, swap interviewer/candidate roles each session, and have the interviewer
ask the "what if this needs to scale 10x" or component-failure follow-up from the design's own
narration section rather than inventing one — that keeps the mock calibrated to a real bar.

### Design rotation

[URL Shortener](../designs/url-shortener/README.md) · [Rate Limiter](../designs/rate-limiter/README.md) ·
[Search Autocomplete](../designs/search-autocomplete/README.md) · [Twitter/Instagram Feed](../designs/twitter-feed/README.md) ·
[YouTube/Netflix](../designs/youtube-netflix/README.md) · [WhatsApp/Messaging](../designs/whatsapp-messaging/README.md) ·
[Key-Value Store](../designs/key-value-store/README.md) · [Distributed Cache](../designs/lru-distributed-cache/README.md) ·
[Notification System](../designs/notification-system/README.md) · [Uber/Ride-Sharing](../designs/uber-ride-sharing/README.md) ·
[Google Maps](../designs/google-maps/README.md) · [Distributed Task Scheduler](../designs/task-scheduler/README.md) ·
[Payment Gateway/UPI](../designs/payment-gateway-upi/README.md) · [Leaderboard/Counters](../designs/leaderboard-counters/README.md) ·
[Logging & Monitoring](../designs/logging-monitoring/README.md)

### After the mock

Diff your whiteboard against that design's own "Deep dives" and "Trade-offs" sections, not the
high-level diagram — the deep dives are where interviews are actually won or lost. Note which
follow-up questions you'd have fumbled live, and revisit that design again in a few days.

### Finding company-specific question lists

This repo doesn't maintain a per-company question bank (those go stale fast). For current,
crowd-sourced reports of what's actually being asked at specific companies, hellointerview.com's
company guides and r/ExperiencedDevs / Blind system-design threads are commonly used — cross-reference
a couple of sources rather than trusting any single list.
