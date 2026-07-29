# Week 4-5 — Real LLD Problems Asked in Interviews

Part of the [8-week LLD program roadmap](../README.md).

## How this week is different from weeks 1-3

Weeks 1-3 taught the building blocks (OOP/SOLID, the 5-step framework, the 12 design patterns). Weeks
4-5 apply them end-to-end against full interview problems, and introduce the twist that shows up in
almost every senior-level LLD round once your base design is accepted: **"now make this thread-safe" /
"now handle this with N concurrent workers."**

## Concept: Concurrency Control for LLD

🎯 Asked as a follow-up (not a standalone problem) on parking-lot, elevator, and rate-limiter designs in
this repo, and broadly across senior LLD rounds.

- Threads/goroutines, mutexes vs. `ReadWriteLock`, `ThreadPool`/`ExecutorService`, the producer-consumer
  pattern, and thread-lifecycle "scheduling" (as opposed to the functional task-scheduler problem in
  week 6-7) are all covered in depth, with a runnable Go/Java/Python pipeline and two mermaid diagrams
  (data-flow + shutdown sequence), in **[`concurrency/README.md`](../concurrency/README.md)**.
- The interview signal here isn't reciting definitions — it's correctly identifying *which* shared state
  in your own design needs protection, and reaching for the cheapest primitive that's actually correct
  (a single mutex before a lock-free structure, `RWMutex` only once reads meaningfully outnumber writes).

## Problems this week

Grouped by what's actually being tested — most of these look different on the surface but reduce to the
same handful of design skills.

**Board/turn-based games** — state machine + move validation + win detection:
- [Chess](../problems/chess/README.md) — 🎯 Asked at Google
- [Snake & Ladder](../problems/snake-and-ladder/README.md) — 🎯 Asked at Flipkart
- [Tic-Tac-Toe](../problems/tic-tac-toe/README.md) — 🎯 Asked at Amazon

**Physical-world state machines** — an entity moving through explicit states, usually with a dispatch/
matching strategy on top:
- [Parking Lot](../problems/parking-lot/README.md) — 🎯 Asked at Amazon
- [Elevator System](../problems/elevator-system/README.md) — 🎯 Asked at Zomato
- [Trading System](../problems/trading-system/README.md) — 🎯 Asked at Microsoft

**Multi-actor coordination systems** — several independent actors (riders, restaurants, friends
splitting a bill) whose interactions need a clean object model, usually with a sequence diagram doing
most of the explaining:
- [Splitwise](../problems/splitwise/README.md) — 🎯 Asked at PhonePe
- [Food Delivery](../problems/food-delivery/README.md) — 🎯 Asked at Swiggy
- [Ride-Sharing](../problems/ride-sharing/README.md) — 🎯 Asked at Uber

**Systems-flavored LLD** — the class design *behind* a piece of infrastructure you'd otherwise treat as
a black box:
- [Kafka LLD](../problems/kafka-lld/README.md) — 🎯 Asked at Uber
- [Payment Gateway](../problems/payment-gateway/README.md) — 🎯 Asked at Razorpay

## How to approach this week's problems in an interview

- **Board/turn-based games**: nail down the win/end condition and whose turn it is *before* touching
  piece movement — most candidates lose time re-deriving turn order mid-design. State the board
  representation (2D array vs. sparse map) out loud and justify it against the game's actual density.
- **Physical-world state machines**: draw the state diagram first (states + legal transitions) — code
  follows directly once the states are right, and an interviewer can spot a missing transition on the
  diagram in seconds instead of reading code.
- **Multi-actor coordination**: default to a sequence diagram over a class diagram as your first
  whiteboard artifact — these problems are mostly about *interaction order* (who calls whom, what's
  synchronous vs. async), and a class diagram alone hides that.
- **Systems-flavored LLD**: explicitly say which real system you're modeling and which parts you're
  deliberately simplifying (e.g. "I'm modeling Kafka's partition/offset/consumer-group model, not
  replication or the network protocol") — this shows scoping judgment instead of the interviewer having
  to ask you to narrow it.
- Across all four groups: apply the [5-step framework from week 1](../week-01-foundations/README.md#3-the-5-step-lld-interview-framework)
  and expect at least one interviewer follow-up that's really the [concurrency](../concurrency/README.md)
  concept above in disguise.

## Practice prompt

Pick one problem from each of the four groups above, cold — don't reread the reference solution first.
Timebox each to 30 minutes: 5 min requirements, 10 min class design (whiteboard the diagram before
writing any code), 10 min coding the 1-2 core flows, 5 min trade-offs. Then, for just one of the four,
add the concurrency twist: "two threads call this at once — walk me through what breaks and how you'd
fix it," and answer it using the primitives from this week's concurrency doc.
