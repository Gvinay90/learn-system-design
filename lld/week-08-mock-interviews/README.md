# Week 8 — Mock Interviews

Part of the [8-week LLD program roadmap](../README.md).

## Why this week is different

Weeks 1-7 build the skills; week 8 is about applying them under real interview conditions — timed,
narrated out loud, and reviewed against a rubric — rather than learning anything new. What follows is a
**self-run (or peer-run) mock format** you can use solo or with a study partner.

## Self-run mock format (45 minutes)

1. **Pick a problem blind** — use the rotation list below, or have a peer/friend pick one you haven't
   solved in the last 2 weeks. Don't reread its README first.
2. **Set a single 45-minute timer, no pausing**, split roughly as:
   - 5 min — requirements gathering (say functional vs. non-functional out loud, ask 2-3 clarifying
     questions even if you're solo — write down what you'd ask)
   - 15 min — class design: identify entities, sketch the class diagram, write the key method
     signatures on the main coordinating class
   - 15 min — code the 2-3 core flows (not every getter/setter)
   - 10 min — trade-offs: state at least one alternative design you rejected and why, and answer one
     self-posed "what if this needs to be thread-safe / scale 10x" follow-up
3. **Stop at 45 minutes even if incomplete** — an incomplete-but-well-reasoned design under time pressure
   is the realistic outcome you're training for, not a finished implementation.

If running this with a peer: swap roles (interviewer/candidate) each session, and have the interviewer
deliberately introduce one ambiguous requirement (e.g. "what happens if two users pay at the same time")
partway through, since real interviewers course-correct based on your answers rather than reading from a
script.

## Problem rotation

Roll a die, use a random picker, or have a peer choose — don't self-select the problem you feel most
confident about.

**LLD (18 problems)**: [Chess](../problems/chess/README.md) · [Snake & Ladder](../problems/snake-and-ladder/README.md) ·
[Tic-Tac-Toe](../problems/tic-tac-toe/README.md) · [Parking Lot](../problems/parking-lot/README.md) ·
[Elevator System](../problems/elevator-system/README.md) · [Trading System](../problems/trading-system/README.md) ·
[Splitwise](../problems/splitwise/README.md) · [Food Delivery](../problems/food-delivery/README.md) ·
[Ride-Sharing](../problems/ride-sharing/README.md) · [Kafka LLD](../problems/kafka-lld/README.md) ·
[Payment Gateway](../problems/payment-gateway/README.md) · [LRU Cache](../problems/lru-cache/README.md) ·
[Rate Limiter](../problems/rate-limiter/README.md) · [Task Scheduler](../problems/task-scheduler/README.md) ·
[Notification System](../problems/notification-system/README.md) · [Logging Framework](../problems/logging-framework/README.md) ·
[In-Memory File System](../problems/in-memory-file-system/README.md) · [Vending Machine](../problems/vending-machine/README.md)

**For a mixed HLD+LLD mock session**, pair one LLD problem above with an unrelated HLD design from
[`hld/README.md`](../../hld/README.md) — real onsite loops usually run one of each, not two of the same
flavor back to back.

## Self/peer feedback rubric

Score each 1-5 right after the timer stops, before looking anything up — accuracy of self-assessment
matters more than the score itself.

| Dimension | 1 (needs work) | 3 (solid) | 5 (strong) |
|---|---|---|---|
| Requirements gathering | Jumped straight to design | Asked 1-2 relevant questions | Scoped functional vs. non-functional explicitly, caught an edge case unprompted |
| Object modeling | Missing/wrong core entities | Correct entities, some awkward relationships | Clean entities, right composition/association choices, matches the reference class design's intent |
| API design | Method signatures bolted on during coding | Signatures written before implementation | Signatures anticipate the deep-dive/extension questions |
| Code quality & correctness | Doesn't compile / core flow broken | Core flow works, rough edges | Core flow works cleanly, handles the obvious edge cases |
| Trade-off articulation | Only mentions trade-offs if asked | States one trade-off unprompted | Proactively raises trade-offs *and* the pattern/concurrency angle |
| Communication | Long silent design phases | Narrates decisions as they're made | Narrates *and* checks in with the interviewer before committing to a direction |

## After the mock

Open the reference problem's README and specifically diff your design against its **"Design patterns
used"** and **"Key trade-offs / talking points"** sections — not the code line-by-line. The gap between
what you said out loud and what's in those two sections is exactly what to drill before the next mock.

## Finding a live 1-on-1 mock partner

The self-run format above gets you most of the way, but nothing fully replaces a live 1-on-1 with real
pushback. If you want that, commonly used options include:

- **A study partner or peer group** — swap interviewer/candidate roles on alternating problems from the
  rotation above; this is free and usually the highest-value option if you can find one.
- **Peer mock-interview platforms** (e.g. Pramp, interviewing.io) — pair you with another candidate or a
  vetted engineer for a timed live mock with feedback.
- **A senior engineer or mentor** you already know, given this repo's problem list and rubric as the
  format to run.

Whichever route you take, come with a specific weak spot from your self-run mocks (e.g. "I keep
skipping the trade-offs discussion") so the session targets it instead of repeating a generic pass.
