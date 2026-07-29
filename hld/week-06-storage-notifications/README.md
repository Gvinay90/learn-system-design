# Week 6 — Storage & Notification Systems

Part of the [8-week HLD learning path](../README.md).

## Concept: WebSockets, Server-Sent Events & Long Polling

- **Long polling**: client sends a request, server holds it open until new data is available (or a
  timeout), then responds; client immediately re-requests. Works over plain HTTP/1.1, needs no special
  infra, but pays a new-connection cost per update cycle and adds latency.
- **Server-Sent Events (SSE)**: client opens one HTTP connection, server streams events down it
  indefinitely; unidirectional (server-to-client only). Good fit for one-way real-time feeds — live
  scores, notification streams, activity feeds — where the client never needs to push data back over
  the same channel.
- **WebSockets**: a single persistent, full-duplex TCP connection; either side sends messages at any
  time with minimal per-message overhead. Best fit when the client also needs to send frequent messages
  back (chat, collaborative editing, multiplayer) or update frequency/latency needs are the strictest.
- **Choosing between them**: pick the simplest option that meets the latency/direction requirement —
  long polling for infrequent updates and simplicity, SSE for one-way real-time streams, WebSockets only
  when true bidirectional low-latency messaging is actually needed (it costs more in connection-state
  management at scale, since every server holding open sockets must track that state).
- **Scaling stateful connections**: WebSocket/SSE servers hold per-connection state in memory, so a
  message meant for a user must be routed to whichever server holds their connection — typically via a
  pub/sub layer (Redis pub/sub, Kafka) that all connection-holding servers subscribe to and filter from.

**References**
- Read first: [Real-time Updates Pattern — Hello Interview](https://www.hellointerview.com/learn/system-design/patterns/realtime-updates)
- Watch: [Notification Service System Design Interview Question to handle Billions of users & Notifications (YouTube)](https://www.youtube.com/watch?v=CUwt9_l0DOg)

## Designs this week

- [Key-Value Store](../designs/key-value-store/README.md) (best explored hands-on — see design doc) — 🎯 Asked at Amazon
- [Distributed Cache](../designs/lru-distributed-cache/README.md) *(built by sibling agent)*
- [Notification System](../designs/notification-system/README.md) (best explored hands-on — see design doc) — 🎯 Asked at Swiggy

## Practice prompt
Design the delivery mechanism for the notification system above: a user has the app open and should see
a new notification appear instantly without refreshing. Whiteboard the choice between long polling, SSE,
and WebSockets for this specific case, then design how a notification generated on any backend server
gets routed to the exact connection-holding server the user's client is attached to, across a fleet of
many such servers.
