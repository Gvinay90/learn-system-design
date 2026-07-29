# Week 7 — Location & Scheduling Systems

Part of the [8-week HLD learning path](../README.md).

## Concept: Geospatial indexing — Geohashing, QuadTrees, Uber's H3

- **Geohashing**: encodes a lat/lng pair into a single string by recursively subdividing the world into
  a grid, interleaving bits of latitude and longitude; nearby locations tend to share a string prefix, so
  a range/prefix query on the geohash approximates a spatial "nearby" query. Simple to store as a regular
  indexed string column, but grid cells are uneven in shape (rectangles that distort near the poles) and
  two physically close points can rarely land in very different geohash prefixes near a cell boundary.
- **QuadTrees**: an in-memory tree that recursively subdivides 2D space into 4 quadrants only where
  needed (dense areas get subdivided further, sparse areas stay coarse) — adapts resolution to point
  density, unlike geohashing's fixed grid. Good for in-memory nearest-neighbor/range queries; not
  natively a storage engine, so it's typically rebuilt/maintained in application memory over data
  persisted elsewhere.
- **Uber's H3**: a hierarchical hexagonal grid system. Hexagons give uniform adjacency (all 6 neighbors
  equidistant from center, unlike a square grid's edge- vs corner-neighbors being different distances),
  which makes distance/neighbor calculations more consistent — valuable for ride-matching and
  drive-time-radius queries where "how far is this cell from that cell" needs to be a stable answer.
- **Choosing one**: geohash is simplest to bolt onto an existing DB (index a geohash column, query by
  prefix range); quadtrees suit an in-memory service needing adaptive resolution; H3 is the pick when
  the domain genuinely benefits from uniform hex adjacency (ride-sharing, delivery zones) and you can
  afford Uber's H3 library/tooling.
- **Common pattern regardless of choice**: reduce a 2D nearest-neighbor problem to a 1D range/prefix
  query so it can run on a normal indexed data store, then do an exact distance filter on the (small)
  candidate set the index query returns.

**References**
- Background: [Proximity Search in System Design Interviews — Hello Interview](https://www.hellointerview.com/learn/system-design/deep-dives/proximity-search) — covers geohash, S2, and H3 directly, the exact concepts this week applies
- Watch: [Guide to Uber's H3 for Spatial Indexing (video companion article)](https://www.analyticsvidhya.com/blog/2025/03/ubers-h3-for-spatial-indexing/) — see also the Uber/Google Maps design write-ups this week for applied walkthroughs

## Designs this week

- [Uber / Ride-Sharing](../designs/uber-ride-sharing/README.md) — 🎯 Asked at Zomato *(built by sibling agent)*
- [Google Maps](../designs/google-maps/README.md) *(built by sibling agent)*
- [Distributed Task Scheduler](../designs/task-scheduler/README.md) — 🎯 Asked at Spotify

## Practice prompt
Design the "find nearby drivers" query for a ride-sharing app at Uber's scale: given a rider's lat/lng,
return the K nearest available drivers within 3km, refreshed as drivers move every few seconds.
Whiteboard which geospatial index you'd pick and why, how you'd keep the index updated as driver
locations change constantly (a write-heavy stream), and how you avoid a moving driver's index entry
being stale by the time a rider's request looks it up.
