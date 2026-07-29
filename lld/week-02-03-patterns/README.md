# Design Patterns — Week 2-3

Twelve classic GoF design patterns, each implemented in Go/Java/Python with its own README
(problem it solves, when to use it, class diagram, and a working test suite). This page is the
index: what each pattern is for, and how to pick the right one under interview time pressure.

## Creational — how objects get constructed

| Pattern | When to reach for it |
|---|---|
| [Singleton](singleton/README.md) | Exactly one instance of a coordinating object (config, connection pool, cache) must exist and be reachable globally, with safe concurrent first-access. |
| [Factory](factory/README.md) | Callers need an object of one of several related types (notification channel, shape, document format) chosen by a runtime flag, without ever naming the concrete type themselves. |
| [Builder](builder/README.md) | An object has many optional/staged fields and you want a fluent, chainable construction API that only hands back a fully-valid (ideally immutable) object at the end. |
| [Prototype](prototype/README.md) | You need many independent variations of an already-configured "template" object, cloned (deep-copied) rather than rebuilt from scratch each time. |

## Structural — how objects are composed

| Pattern | When to reach for it |
|---|---|
| [Adapter](adapter/README.md) | You must make an existing, unmodifiable interface (legacy code, a third-party library) work with the interface the rest of your system expects. |
| [Decorator](decorator/README.md) | You need to layer optional behavior/cost onto an object (add-ons, middleware-style wrapping) without a combinatorial explosion of subclasses. |
| [Facade](facade/README.md) | A subsystem has several interdependent classes with a fiddly init/call order, and callers just want one simple entry point that hides that complexity. |
| [Proxy](proxy/README.md) | You need to control or intercept access to an object — permission checks, lazy/deferred construction, caching, logging — without changing the object itself. |

## Behavioral — how objects communicate and vary behavior

| Pattern | When to reach for it |
|---|---|
| [Observer](observer/README.md) | Multiple decoupled components must react whenever one subject's state changes, and the subject shouldn't know the concrete types of its dependents. |
| [Strategy](strategy/README.md) | An object needs to support several interchangeable algorithms (pricing, sorting, routing) selected at runtime, without branching logic baked into the object. |
| [Chain of Responsibility](chain-of-responsibility/README.md) | A request should be handled by one of several possible handlers, determined at runtime, without the sender knowing which handler (or how many) exist. |
| [State](state/README.md) | An object's behavior and legal transitions depend on a finite set of internal states, and you want illegal transitions to be structurally impossible rather than an `if` you forgot. |

## Pattern selection under interview conditions

When an interviewer's requirements hint at a shape, these are the fastest associations to reach for:

1. **Multiple interchangeable algorithms picked at runtime** (pricing rules, sort orders, route
   planning) -> **Strategy**.
2. **Object needs staged or optional construction** with many optional parameters and you want the
   final object immutable -> **Builder**.
3. **Need to react to state changes across decoupled components** (multiple displays/services
   watching one data source) -> **Observer**.
4. **Behavior and legal transitions depend on a finite set of states** (order lifecycle, traffic
   light, game state) -> **State**.
5. **A request could be handled by one of several handlers, decided by runtime data, and the sender
   shouldn't hardcode which** (approval chains, middleware pipelines) -> **Chain of Responsibility**.
6. **You must integrate with code you can't change** (legacy API, third-party SDK) whose interface
   doesn't match what your system expects -> **Adapter**.
7. **You want to add optional behavior/cost in layers without subclassing every combination**
   (coffee add-ons, I/O stream wrapping) -> **Decorator**.
8. **You need to gate, defer, or instrument access to an object** (auth checks, lazy loading,
   caching) while keeping the same interface -> **Proxy**.
9. **A complex subsystem needs one simple front door** for the common case, without hiding the
   subsystem from callers who need it directly -> **Facade**.
10. **Exactly one shared instance must exist, constructed safely under concurrency** (global config,
    connection pool) -> **Singleton** — but reach for it last; it's frequently over-used and adds
    hidden global state that hurts testability.
11. **You need one of several related concrete types produced from a runtime flag**, and don't want
    call sites to import/know every concrete type -> **Factory**.
12. **You need many independent copies of an expensive-to-build "template" object**, and mutating
    one copy must never leak into another -> **Prototype** (and be ready to explain deep vs. shallow
    copy — this is the actual interview content, not just "call Clone()").

A good tell in interviews: if you find yourself about to write `switch` on a type or a growing pile
of `if/else` to pick behavior, that's almost always Strategy, State, or Factory in disguise — say so
out loud, then justify which one fits based on whether the "modes" are algorithms (Strategy), a
lifecycle (State), or object variants (Factory).
