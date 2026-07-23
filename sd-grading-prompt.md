# SD Design — Grading Prompt

Paste this above your design each time you want it graded. Companion to the prompt bank and the prep-phases doc.

---

```text
Grade the system design below at the senior bar for my targets (Wise / Stripe / 
Datadog / Adyen tier). Be a strict grader by default, not a generous one — I'm 
here to find weaknesses, not to be reassured. Don't praise to soften the blow; if 
something's genuinely strong, one line, then spend the rest on where it falls apart.

Score it against these five, and say which dimension was weakest:
1. Scope — did I gather requirements and cut ruthlessly?
2. Estimates — did I size it (QPS/storage/bandwidth) enough to justify my choices?
3. Storage — did I defend the choice against alternatives, not just name a DB?
4. Deep dive — did I reach real depth on 1–2 hard parts? (weight this most)
5. Failure modes — bottlenecks, what breaks first at 10x, what I consciously punted.

Rules for how you grade:
- Every judgment must name a mechanism. Not "this is shallow" — say WHY as a 
  checkable claim I can reason about (e.g. "the processor-timeout case leaves you 
  not knowing if the charge landed, and you didn't handle it"). If you can't name 
  a concrete mechanism, drop the point.
- Give a verdict: pass / borderline / fail at the senior bar.
- Name the single weakest point in the whole design.
- Then play interviewer: ask the 3–4 questions you'd push on next — "why not X 
  instead," the load escalation that strains my design, what breaks first.
- If this prompt isn't covered by a standard answer key (most of my fintech/infra 
  ones), say so, and flag anywhere you're less certain so I don't take you as 
  ground truth.

Design:
[paste yours here]
```

---

## Notes on use

- You don't need the whole block every time once it's routine. The two load-bearing lines are **"strict by default"** (counters the drift to generous grading) and **"every judgment must name a mechanism"** (the safeguard that lets you catch a bad grade without SD expertise). Keep those two no matter what.
- When you want it to bite harder, add: **"before grading, argue the case that this design would fail the round."** Forcing the prosecution first stops a rationalized pass.
- After each grade, log the *recurring* weak spot (not the per-problem detail) in your running list — that log is how you grow out of needing the grader at all.
