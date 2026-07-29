# Interview Prep — HLD + LLD Learning Path

This is a self-contained, self-paced study repo covering two tracks:

- [`hld/`](hld/README.md) — 8-Week High-Level Design (system design) learning path
- [`lld/`](lld/README.md) — 8-Week Low-Level Design (LLD) learning path

Every topic gets: a concept summary, curated reading/watch links (hellointerview.com + YouTube +
other open-source references), a real "asked at" company tag where one is commonly reported for that
topic, and — where the topic is code-able — a runnable, tested implementation.

## How this repo is organized

```
interview-prep/
├── go.mod                # one Go module for the whole repo — `go test ./...` runs everything
├── resources.md          # consolidated link list (hellointerview pages, YouTube channels/playlists)
├── hld/                  # 8-week HLD roadmap + 16 design write-ups
└── lld/                  # 8-week LLD roadmap + 12 patterns + 18 problems, each in Go/Java/Python
```

## Tooling prerequisites

| Language | Version used | Run tests |
|---|---|---|
| Go | 1.25+ | `cd interview-prep && go test ./...` |
| Java | 17+ (plain `javac`/`java`, no build tool required) | see each problem's README for exact `javac`/`java` commands |
| Python | 3.11+ | `pip install -r lld/requirements.txt && pytest` |

No Maven/Gradle/pip global install is required beyond `pytest` for Python — Java examples are plain
`.java` files compiled directly so you can run them without setting up a build tool.

## How to use this repo day-to-day

1. Read the week's README for the concept overview + links. Watch/read before you look at any code.
2. For HLD: try to whiteboard the design yourself first (requirements → API → data model → high-level
   diagram → deep dives) using the "practice prompt" in the design's README, *then* compare against the
   write-up.
3. For LLD: try to write the class design and code yourself (pen/paper or a scratch file) before opening
   the reference implementation. Then run the reference implementation's tests to compare behavior.
4. Track progress with the checkboxes in [`hld/README.md`](hld/README.md) and [`lld/README.md`](lld/README.md).

## Progress checklist

- [ ] HLD Week 1–8 (see [`hld/README.md`](hld/README.md))
- [ ] LLD Week 1–8 (see [`lld/README.md`](lld/README.md))
- [ ] All 16 HLD designs whiteboarded from scratch at least once
- [ ] All 18 LLD problems coded from scratch in at least one language
- [ ] 2–3 full self-run mock HLD interviews completed (see [`hld/week-08-payments-reliability/README.md`](hld/week-08-payments-reliability/README.md))
- [ ] LLD self-run mock interview completed (see [`lld/week-08-mock-interviews/README.md`](lld/week-08-mock-interviews/README.md))

See [`resources.md`](resources.md) for the full reference list.
