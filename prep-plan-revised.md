# Integrated Prep Plan — April to November 2026

**Target:** Offer-ready by mid-October, peak readiness through mid-November.
**Workstreams:** LeetCode (rusty baseline), System Design (foundations present, interview-shaped skill unbuilt), Project (V2 wrapping, V3+ as optionality).
**Language:** Go, raw stdlib only — no helper packages.

---

## Guiding principles

- **LC is the longest-ramp item** (4 months from rusty) so it starts now and stays primary through July.
- **SD marinates before it intensifies.** 2 hrs/week starting mid-May, ramping to 10 hrs/week in August.
- **Project work is backgrounded** after V2. V3+ only to the extent it generates SD/behavioral material.
- **Benchmark constantly.** Don't trust "I feel ready" — test against timed, realistic tasks.
- **Maintenance mode is the load-bearing assumption.** If you can't do 3 hrs/week of LC review in August–October, the whole plan shifts. Be honest with yourself about this early.

---

## April 20 – May 31 — LC Foundation (Weeks 1–6)

**LC:** 7–9 hrs/week. NeetCode 150 ordering: arrays/hashing → two pointers/sliding window → stack/binary search → linked lists → trees pt 1. Raw Go stdlib. Target: mediums in 35 min by end of May.

**SD:** None for first 3 weeks. Week 4 (mid-May): acquire Alex Xu Vol 1, read 1 chapter/week at ~2 hrs. No practice designs yet — just absorb the interview framework (requirements → estimation → high-level → deep-dive → tradeoffs).

**Project:** Finish V2 next week. Start V3 as weekend/evening work *only if you want to* — it's no longer on the critical path. See "V3+ strategy" below.

**Benchmark (end of May):** Attempt URL Shortener SD problem. 45 min, paper, alone, no lookups. Then compare to Hello Interview's reference. The gap = your June/July SD curriculum.

---

## June 1 – July 31 — LC Depth + SD Ramp (Weeks 7–14)

This is the hardest stretch. Graphs and DP are slow. SD ramps up. Expect to feel overloaded in late June.

**LC:** 7–9 hrs/week. Trees pt 2, tries, heap, intervals, backtracking, graphs, advanced graphs, 1-D DP, 2-D DP, greedy. By end of July: mediums in 25–30 min, hards attemptable with effort.

**SD:** Ramps from 2 hrs/week to 5 hrs/week over June, then 6–8 hrs/week in July.
- June: finish Alex Xu Vol 1, start Vol 2, watch Hello Interview's core flow videos.
- July: begin doing designs yourself. 1 per week, 45 min, out loud (alone is fine — record yourself). Canonical set: URL shortener, rate limiter, TinyURL, Twitter timeline, Instagram, messaging system, news feed, distributed cache.
- After each self-design, watch Hello Interview's or ByteByteGo's version and list the gaps.

**Project:** V3 as optional background work. See strategy section.

**Benchmark (end of July):**
- LC: solve a random NeetCode medium in <30 min unaided.
- SD: re-attempt URL shortener, now timed 45 min. Should be dramatically better than May attempt. If not, SD needs more time.

---

## August 1 – August 31 — SD Intensive + LC Maintenance (Weeks 15–18)

The flip month. LC drops to maintenance. SD becomes primary.

**LC:** ~3 hrs/week, maintenance only. Mixed medium deck (not topic-sorted) — 2–3 problems per session, timed. The goal is preserving pattern recognition, not learning new material. If you find maintenance isn't working — if LC feels rusty again by mid-August — add an hour or two.

**SD:** 10–12 hrs/week.
- 2–3 full self-designs per week, timed, recorded.
- 1 paid mock on Hello Interview with a real FAANG interviewer. This is expensive (~$200) but the single highest-leverage thing in the whole plan. Schedule mid-August.
- Deep-dive on 3 topics you're weakest on after the mock.
- Write up your own job processing system as an SD case study — 2 pages, as if presenting to an interviewer. This doubles as behavioral material.

**Behavioral:** Start STAR-format story prep, 2 hrs/week. Aim for 8–10 stories covering: hardest technical problem, disagreement with teammate, production failure, ambiguity/decision-making, learning something new, leading without authority, scope/tradeoff call. Your job processing system should anchor 2–3 of these.

**Benchmark (end of August):** Second paid mock with a different interviewer. Compare feedback to mid-August mock. Gap should be narrowing measurably.

---

## September 1 – September 30 — Mocks, Polish, Start Applying (Weeks 19–22)

You're in interview shape. Now calibrate against real loops.

**LC:** ~3 hrs/week maintenance continues. Add 1 hard per week for range-extending.

**SD:** 6–8 hrs/week. Mostly mocks now — aim for 2 mocks/week, mix paid (Hello Interview, interviewing.io) and free (Pramp, peer). Stop doing self-designs unless a mock reveals a specific weak topic.

**Behavioral:** Finalize stories. Practice out loud. 1 mock behavioral interview.

**Applications:** Start applying to **tier 2** companies — places you'd take an offer from but that aren't your top targets. This is calibration, not a fallback. First 2–3 real loops will go poorly regardless of prep; you want that to happen somewhere with acceptable consequences. Mid-tier tech, well-funded startups, strong infra companies outside the top 5. Save FAANG/Stripe/Databricks for October–November.

**Benchmark (end of September):** Real interview feedback from ≥2 loops. This is your actual signal — mocks only approximate.

---

## October 1 – November 15 — Peak Window (Weeks 23–30)

Peak hiring window. Apply to target companies.

**LC/SD:** Maintenance only. 2 hrs LC, 2 hrs SD per week — whatever keeps you warm between loops.

**Applications:** Dream-tier companies. Cluster loops close together — getting 3–4 offers in the same 2-week window maximizes negotiation leverage. Don't accept the first offer; let multiple complete.

**Negotiation:** Budget 1 week for this. Read Patrick McKenzie's "Salary Negotiation" essay. Do not give a number first. Do not accept on the phone. Get competing offers in writing.

---

## November 15+ — Fallback

If no acceptable offer by Nov 15, loops that started Nov 10+ will likely pause for the holidays and resume in January. Not a failure — January is the strongest hiring window of the year, and you'll go into it peaked. Maintain at 4–5 hrs/week through December to avoid decay.

---

## V3+ Strategy — what actually earns its keep

V3/V4 are **not on the critical path**. But if you're going to build them anyway (you will), make them interview-generative. Pick features that give you SD and behavioral material:

- **Distributed coordination** — multi-node workers, leader election, heartbeat/gossip. Directly answers "how would you scale horizontally?" in SD interviews.
- **Observability** — Prometheus metrics, OpenTelemetry traces, Grafana dashboard. Every senior loop probes this; most candidates wave hands.
- **Backpressure / rate limiting** on the Claude API worker — gold-tier SD material, pairs with real tradeoffs (queue vs. reject, per-tenant limits).
- **A real postmortem** — break your own system deliberately (e.g., kill Postgres mid-job, corrupt a retry counter), document the incident, recover, write it up. This is the single best behavioral-interview material you can manufacture. Staff interviews *love* well-written postmortems.

Do not build V3+ features that don't map to one of those buckets. Scope creep on a side project is a trap that eats prep hours and doesn't move the interview date.

---

## Checkpoints — when to adjust the plan

**End of May:** If LC mediums still take >45 min average, you need more LC hours per week, not more time. Push to 10 hrs/week in June.

**End of July:** If your URL shortener re-attempt is still shaky, SD ramp was too slow. Extend SD intensive into September, push interview start to October.

**Mid-August mock result:** If mock feedback is "not ready for senior SD," you need 4–6 more weeks of SD, not more LC. Shift applications to late October.

**End of September:** If tier-2 loops are going poorly in ways that reveal a specific gap (LC speed, SD communication, behavioral), pause applications for 2 weeks and drill that gap. Don't keep applying through weakness — each bad loop costs you a company for 6–12 months.

---

## Things to track, weekly

- **LC:** # problems solved, # solved unaided under 30 min, running list of "slow" or "hinted" problems for spaced-repetition review.
- **SD:** topics covered, self-designs attempted, gap list from Hello Interview comparisons.
- **Project:** what shipped, what you learned, any incident/failure notes (future behavioral material).
- **Energy:** honest 1–10 on how you're holding up. If it's under 6 for two weeks straight, cut hours. Burnout kills more prep timelines than bad plans do.

---

## What I'm less sure about

A few things in this plan are genuine guesses:

- **Whether SD ramps in 3 months or needs 4.** Depends how the mid-May and end-of-July benchmarks go. Plan has some buffer but not infinite.
- **Whether LC maintenance at 3 hrs/week holds.** Works for most people; may not work for you. First sign is mediums feeling slow again in mid-August.
- **Whether you'll actually do mocks.** The #1 skipped item in every prep plan. Schedule the first one by early August and pay for it so you can't flake.

If any of these go sideways, the timeline slips by 3–6 weeks. That's still inside the good autumn hiring window, so don't panic — just adjust.
