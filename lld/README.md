# 8-Week Low-Level Design (LLD) Learning Path — Roadmap

Goal: master OOP, SOLID, design patterns, and 18 real LLD interview problems — every problem coded in
**Go, Java, and Python**, runnable and tested locally.

## Weekly checklist

- [ ] **Week 1 — [Foundations & Thinking Framework](week-01-foundations/README.md)**
  - OOP mastery: abstraction, encapsulation, polymorphism — 🎯 Asked at Flipkart
  - SOLID Principles — applied, not theoretical — 🎯 Asked at Microsoft
  - 5-step LLD interview framework
- [ ] **Week 2-3 — [Design Patterns in Practice](week-02-03-patterns/README.md)**
  - Creational: Singleton, Factory, Builder, Prototype
  - Structural: Adapter, Decorator, Facade, Proxy, Composition — 🎯 Asked at Swiggy
  - Behavioral: Observer, Strategy, Chain of Responsibility, State — 🎯 Asked at Zomato
  - Pattern selection under interview conditions
- [ ] **Week 4-5 — [Real LLD Problems Asked in Interviews](week-04-05-problems/README.md)**
  - [Concurrency Control for LLD](concurrency/README.md) — threads, locks, ReadWriteLock, ThreadPool/ExecutorService, producer-consumer, task scheduler
  - [Chess](problems/chess/README.md), [Snake & Ladder](problems/snake-and-ladder/README.md), [Tic-Tac-Toe](problems/tic-tac-toe/README.md) — 🎯 Asked at Google
  - [Parking Lot](problems/parking-lot/README.md), [Elevator](problems/elevator-system/README.md), [Trading System](problems/trading-system/README.md) — 🎯 Asked at Amazon
  - [Splitwise](problems/splitwise/README.md), [Food Delivery](problems/food-delivery/README.md), [Ride Sharing](problems/ride-sharing/README.md)
  - [LLD of Kafka](problems/kafka-lld/README.md), [Payment Gateway](problems/payment-gateway/README.md)
- [ ] **Week 6-7 — [Advanced LLD Problems — Core Systems & Patterns](week-06-07-advanced/README.md)**
  - [LRU Cache](problems/lru-cache/README.md) (HashMap + Doubly Linked List, O(1) ops) — 🎯 Asked at Uber
  - [Rate Limiter](problems/rate-limiter/README.md) (Token Bucket, Leaky Bucket, Sliding Window) — 🎯 Asked at Razorpay
  - [Task Scheduler](problems/task-scheduler/README.md) (priority execution, delayed jobs, retries)
  - [Notification System](problems/notification-system/README.md) (email, SMS, push)
  - [Logging Framework](problems/logging-framework/README.md) (log levels, appenders, extensibility)
  - [In-Memory File System](problems/in-memory-file-system/README.md) (directories, files, operations)
  - [Vending Machine](problems/vending-machine/README.md) (state machine, inventory, payment)
- [ ] **Week 8 — [Mock Interviews](week-08-mock-interviews/README.md)**
  - Full 1-on-1 mock interview sessions
  - Personalized feedback & improvement plan

## All 18 LLD problems (quick index)

| Problem | Asked at | Go | Java | Python |
|---|---|---|---|---|
| [Parking Lot](problems/parking-lot/README.md) | Amazon | ✅ | ✅ | ✅ |
| [Chess](problems/chess/README.md) | Google | ✅ | ✅ | ✅ |
| [Snake & Ladder](problems/snake-and-ladder/README.md) | Flipkart | ✅ | ✅ | ✅ |
| [Tic-Tac-Toe](problems/tic-tac-toe/README.md) | Amazon | ✅ | ✅ | ✅ |
| [Elevator System](problems/elevator-system/README.md) | Zomato | ✅ | ✅ | ✅ |
| [Trading System](problems/trading-system/README.md) | Microsoft | ✅ | ✅ | ✅ |
| [Splitwise](problems/splitwise/README.md) | PhonePe | ✅ | ✅ | ✅ |
| [Food Delivery](problems/food-delivery/README.md) | Swiggy | ✅ | ✅ | ✅ |
| [Ride-Sharing](problems/ride-sharing/README.md) | Uber | ✅ | ✅ | ✅ |
| [Kafka LLD](problems/kafka-lld/README.md) | Uber | ✅ | ✅ | ✅ |
| [Payment Gateway](problems/payment-gateway/README.md) | Razorpay | ✅ | ✅ | ✅ |
| [LRU Cache](problems/lru-cache/README.md) | Uber | ✅ | ✅ | ✅ |
| [Rate Limiter](problems/rate-limiter/README.md) | Razorpay | ✅ | ✅ | ✅ |
| [Task Scheduler](problems/task-scheduler/README.md) | Spotify | ✅ | ✅ | ✅ |
| [Notification System](problems/notification-system/README.md) | Swiggy | ✅ | ✅ | ✅ |
| [Logging Framework](problems/logging-framework/README.md) | Microsoft | ✅ | ✅ | ✅ |
| [In-Memory File System](problems/in-memory-file-system/README.md) | Netflix | ✅ | ✅ | ✅ |
| [Vending Machine](problems/vending-machine/README.md) | Meesho | ✅ | ✅ | ✅ |

## 12 Design patterns (quick index)

Creational: [Singleton](week-02-03-patterns/singleton), [Factory](week-02-03-patterns/factory), [Builder](week-02-03-patterns/builder), [Prototype](week-02-03-patterns/prototype)
Structural: [Adapter](week-02-03-patterns/adapter), [Decorator](week-02-03-patterns/decorator), [Facade](week-02-03-patterns/facade), [Proxy](week-02-03-patterns/proxy)
Behavioral: [Observer](week-02-03-patterns/observer), [Strategy](week-02-03-patterns/strategy), [Chain of Responsibility](week-02-03-patterns/chain-of-responsibility), [State](week-02-03-patterns/state)

Each pattern folder has a Go/Java/Python runnable example — see [`week-02-03-patterns/README.md`](week-02-03-patterns/README.md).

## Running the code

```bash
# Go — from interview-prep/
go test ./...

# Java — from any problem's java/ dir (see that problem's README for exact class names)
javac src/*.java -d out && java -cp out Main

# Python — from interview-prep/lld/
pip install -r requirements.txt
pytest
```
