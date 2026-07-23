# System Design Prompt Bank

Calibrated to senior loops at Wise / Stripe / Datadog / Adyen / Booking / Revolut. Weighted toward fintech and infra because that's where your targets actually probe — and where you have lived depth most candidates don't.

**Legend:** 🏠 = you've built/own this in production (use as *both* a fresh design rep and a behavioral story) · 🎯 = high hit-probability at your specific targets.

---

## How to use this

Pick one per weekend in Phase 2. Time-box **45 minutes**, speak it out loud or write it end to end — the point is *output reps*, not a perfect design. The failure mode is reading the answer instead of producing one; don't look anything up until you've driven the whole thing and stalled.

Run every prompt through the same skeleton so the structure becomes automatic and you spend your scarce minutes on the deep dive, not on remembering what comes next:

1. **Requirements** — functional + non-functional. Ask clarifying questions out loud; scope ruthlessly. (~5 min)
2. **Estimates** — QPS, storage, bandwidth. Enough to justify later choices, not a math exercise. (~5 min)
3. **API + data model** — the contract and the schema. (~5 min)
4. **High-level architecture** — boxes and arrows, happy path. (~8 min)
5. **Deep dive** — the 1–2 hard parts below. *This is where senior is won or lost.* (~17 min)
6. **Failure modes & tradeoffs** — what breaks at 10x, what you'd do differently, what you consciously punted. (~5 min)

For each prompt below, the bullets are **the distinctive hard part** — the thing the interviewer is actually digging for. If your design doesn't reach it, you haven't passed the round, however clean the boxes-and-arrows were.

---

## Tier 0 — Warm-ups (do these first)

Low domain load, so you can drill the *structure* without drowning. Burn through these in week 1–2 to make the skeleton reflexive.

**URL shortener.** ID generation strategy (counter vs hash vs base62), read-heavy caching, collision handling, custom aliases. The whole value is rehearsing the 6-step flow on something you can't get lost in.

**Distributed unique ID generator.** Snowflake structure, clock skew and how you survive it, coordination-free generation vs a central allocator. Quick, and the reasoning recurs everywhere.

**🎯 Rate limiter.** Token bucket vs sliding-window-log vs sliding-window-counter, and *why* you'd pick each. The real deep dive: distributed state (the read-then-write race in Redis), where it lives (gateway vs sidecar vs library), and what happens when the rate-limiter store itself is down — fail open or closed? Asked at basically every infra/payments company.

**Typeahead / autocomplete.** Trie vs precomputed top-k, ranking, how you update suggestions, the latency budget. Trains the data-structure-choice-under-constraint muscle.

---

## Tier 1 — Core senior canon

The bread and butter. Each trains a distinct pattern; don't grind all of them, cover the patterns.

**🎯 Notification fanout / push service.** Fanout-on-write vs on-read, the celebrity/high-fanout problem, multi-channel (push/email/SMS) with per-channel reliability, dedup, user preferences, and delivery guarantees (at-least-once + idempotent consumers). The dispatcher shape here rhymes with your job processor.

**News feed.** Push vs pull vs hybrid, ranking, the hot-key problem for high-follower accounts. Mostly a vehicle for the fanout tradeoff discussion.

**Chat system.** WebSocket connection management at scale, presence, message ordering and delivery receipts, group-chat fanout, online/offline sync. Deep dive is usually ordering + delivery semantics.

**🎯 Booking / reservation system (Ticketmaster-style).** The reservation-hold pattern, preventing double-booking under concurrency, inventory consistency, checkout TTL and releasing abandoned holds. The concurrency core is *exactly* your `SELECT FOR UPDATE` territory — lean into it. Directly relevant to Booking.

**Web crawler.** Politeness/rate-per-domain, URL dedup (bloom filter / seen-set), the BFS frontier, freshness/recrawl. Deep dive: dedup at scale and the frontier design.

**File storage / sync (Drive-style).** Chunking, dedup, metadata vs blob separation, conflict resolution on concurrent edits, sync protocol. Deep dive is usually conflict handling.

---

## Tier 2 — Fintech / home turf 🎯

Highest hit-probability at your targets, and where your KYC/payments background is a genuine moat. Treat these as priority. For the 🏠 ones, prep them twice: once as a fresh design, once as "here's how I actually did it."

**🎯 Payment system / gateway.** The core of Stripe and Adyen. Distinctive parts, all of which they *will* push on: idempotency keys and exactly-once semantics over an unreliable network; the dual-write problem (your DB *and* the external processor — what if the processor call times out and you don't know if the charge landed?); sync authorization vs async settlement; a double-entry ledger as source of truth; reconciliation against processor records. This is the single most valuable prompt in the bank for you.

**🎯 Money transfer / remittance.** Wise's entire product. The insight they want: cross-border transfers don't actually move money across borders — you hold local liquidity pools on both sides and net/settle. FX rate locking and the window of risk, multi-currency ledgers, settlement timing, regulatory holds. If you walk into Wise without the local-pool insight, that's a miss.

**🎯 Digital wallet / ledger.** Double-entry bookkeeping, balance consistency under concurrent debits/credits, atomicity of the debit+credit pair, preventing negative balances under a race. This is your atomic-locking experience applied to money — your `SELECT FOR UPDATE SKIP LOCKED` reasoning transfers almost verbatim. Make that connection explicit in the room.

**🏠🎯 KYC / identity verification pipeline.** *You own this in production.* Async document processing, unreliable third-party verification vendors (the same shape as the AI-summarization endpoint you reviewed — blocking external call, needs status model + retries), a per-applicant state machine, the manual-review queue, audit/compliance and immutability requirements. You can out-design almost anyone here — and it doubles as your strongest behavioral narrative.

**🎯 Real-time fraud / risk.** The in-line latency budget (you must decide *during* the transaction), streaming feature computation vs batch, rules engine vs ML, and the precision/recall business tradeoff (a false positive blocks a real customer). Deep dive is the latency-vs-accuracy tension.

**Reconciliation system.** Matching internal ledger against external processor/bank statements, handling timing differences and partial matches, surfacing breaks for investigation. Less glamorous, very fintech-real, and a strong signal you've worked in the domain rather than read about it.

---

## Tier 3 — Infra / observability

Datadog-flavored, plus the two systems you've literally built. Turn your in-progress monitoring reading into output here first.

**🎯 Metrics / monitoring system.** You're on this chapter right now — design it *before* you finish reading it. Time-series ingestion at scale, the cardinality-explosion problem, downsampling/rollups, push vs pull collection, hot vs cold storage tiers, the query layer. Datadog core; converting input → output on this one is high leverage.

**🏠🎯 Distributed job scheduler / processor.** *You built this — V4 done.* At-least-once vs exactly-once execution, preventing duplicate execution on worker crash (your reaper), the dispatcher and leader election, atomic claim of work (`SELECT FOR UPDATE SKIP LOCKED`), retry/backoff, backpressure, fairness across tenants. "Design a distributed job scheduler" is a real, common prompt and you have *running code*. Almost no candidate can say "here's the design" then "here's how I actually shipped it." This is your single biggest differentiator — design rep and behavioral story both.

**Log aggregation / collection.** Ingestion pipeline, indexing for search, retention tiers, the write-vs-query-load tension. Pairs naturally with the metrics system.

**Alerting system.** Evaluating rules against streaming metrics, dedup and grouping, escalation, avoiding alert storms, the on-call notification path. Good follow-on once monitoring is solid.

---

## Suggested ordering for your summer

- **Phase 1 (now → mid-July):** Tier 0 to make the skeleton reflexive, then your *first full design* — make it the job scheduler, since you can't get lost in the domain and it builds confidence.
- **Phase 2 (mid-July → mid-Aug):** one per weekend, prioritize Tier 2 + the two 🏠 infra ones. First SD mock ~July 15 on a Tier 2 prompt.
- **Phase 3 (mid-Aug → end Aug):** mixed mock loops; revisit any prompt where you stalled on the deep dive.

The 🏠 prompts (KYC pipeline, job scheduler) are worth 2–3 passes each — they're where you convert "competent senior candidate" into "obviously been doing this for real."
