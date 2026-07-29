# Week 1 — LLD Foundations

🎯 Asked at: Flipkart (OOP fundamentals — abstraction/encapsulation/polymorphism are routinely probed
before any design problem), Microsoft (SOLID applied to real class designs, not recited from memory)

## References
- Read first: [Design Principles — Hello Interview](https://www.hellointerview.com/learn/low-level-design/in-a-hurry/design-principles)
- Prep overview: [How to Prepare for a Low-Level Design Interview — Hello Interview](https://www.hellointerview.com/blog/how-to-prepare-lld)
- Watch: [L01: LLD Interview Approach | Low Level Design (YouTube)](https://www.youtube.com/watch?v=GixkAcu3eEw)

## 1. OOP mastery

The four pillars, as they actually get probed in an interview — not the textbook definitions.

- **Abstraction** — expose *what* an object does (its interface) and hide *how* it does it (the
  implementation). In Go this is a small interface backed by a concrete type the caller never sees directly.
- **Encapsulation** — bundle state and the behavior that mutates it together, and restrict direct access
  to that state so invariants can't be violated from outside. In Go: unexported fields + exported methods.
- **Polymorphism** — the same call (`Area()`, `Notify()`, `CalculateFee()`) behaves differently depending
  on the concrete type behind an interface, so calling code doesn't need a type switch. 🎯 Asked at Flipkart
  — expect "what's the difference between polymorphism and a type switch, and why do we prefer the former?"
- **Inheritance** — a mechanism for one type to reuse another's structure/behavior. Go has no classical
  inheritance; struct embedding gives reuse without an "is-a" contract, which is exactly why most Go LLD
  answers reach for composition + interfaces instead. 🎯 Asked at Flipkart — be ready to justify *why*
  composition is preferred over inheritance for extensibility (avoids fragile base-class problems, avoids
  forcing an "is-a" relationship where "has-a"/"behaves-like" is more accurate).

## 2. SOLID — applied, not theoretical

🎯 Asked at Microsoft: interviewers here rarely ask "define SRP." They hand you a class, ask what's wrong
with it, and expect you to refactor it live.

- **S — Single Responsibility**: a class/type should have exactly one reason to change.
- **O — Open/Closed**: open for extension, closed for modification — new behavior should be addable
  without editing existing, tested code.
- **L — Liskov Substitution**: a subtype must be usable anywhere its base type is expected, without
  breaking correctness (the classic tell: a subtype that throws/no-ops on a method its base promises).
- **I — Interface Segregation**: don't force a type to implement methods it doesn't need — prefer several
  small interfaces over one fat one.
- **D — Dependency Inversion**: high-level modules should depend on abstractions, not on concrete
  low-level implementations.

### Worked example: SRP violation → fix

**Before** — `Invoice` computes the total *and* knows how to persist itself. Two reasons to change: pricing
logic changing, or the storage mechanism changing (file → DB → API).

```go
type Invoice struct {
    Items []LineItem
}

func (i *Invoice) Total() float64 {
    var total float64
    for _, it := range i.Items {
        total += it.Price * float64(it.Qty)
    }
    return total
}

// Violates SRP: Invoice now also owns persistence concerns.
func (i *Invoice) SaveToFile(path string) error {
    data := fmt.Sprintf("Total: %.2f\n", i.Total())
    return os.WriteFile(path, []byte(data), 0644)
}
```

**After** — persistence is pulled out into its own type behind an interface. `Invoice` only ever changes
if pricing rules change; storage can change (file → S3 → Postgres) without touching `Invoice` at all.

```go
type Invoice struct {
    Items []LineItem
}

func (i *Invoice) Total() float64 {
    var total float64
    for _, it := range i.Items {
        total += it.Price * float64(it.Qty)
    }
    return total
}

type InvoiceStore interface {
    Save(inv *Invoice) error
}

type FileInvoiceStore struct{ Path string }

func (s FileInvoiceStore) Save(inv *Invoice) error {
    data := fmt.Sprintf("Total: %.2f\n", inv.Total())
    return os.WriteFile(s.Path, []byte(data), 0644)
}
```

Why it matters in interviews: this is the single most common "smell" interviewers plant on purpose — a
domain object with a stray `Save`/`Print`/`Send` method — precisely so they can watch you notice and split
it out.

## 3. The 5-step LLD interview framework

1. **Clarify requirements & scope** — separate functional requirements (what the system must do) from
   non-functional ones (concurrency, extensibility, scale). Say the scope out loud and get the interviewer
   to confirm it before designing anything; this is the step most candidates rush and pay for later.
2. **Identify core objects/entities** — pull nouns out of the requirements (e.g. `Vehicle`, `ParkingSpot`,
   `Ticket`) and decide which are entities (have identity/lifecycle) vs. value objects vs. pure services.
3. **Define relationships & class diagram** — decide composition vs. association vs. inheritance between
   the entities from step 2, and sketch it (mermaid `classDiagram` on a whiteboard/doc) before writing code.
4. **Identify key interactions/APIs** — write the method signatures on the main coordinating class first
   (e.g. `parkVehicle(vehicle) Ticket`) — this forces you to commit to inputs/outputs before implementation
   details leak in.
5. **Code the core flows, discuss trade-offs** — implement only the 2-3 flows that matter (not every
   getter/setter), then proactively raise trade-offs (thread-safety, extensibility via Strategy/Factory,
   what you'd change at 10x scale) rather than waiting to be asked.
