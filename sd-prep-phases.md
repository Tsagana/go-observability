# System Design Interview Prep — Phases & Method

Companion to the prompt bank. The bank is *what* to practice; this is *how*, and in what order.

**Your situation in one line:** content-rich (DDIA + SD Vol II), output-poor (zero designs driven end to end), and new to the *format* — never sat in a real SD round. So the plan is not "learn more SD." It's "convert what you know into delivered designs, and learn the format by doing reps with feedback." Reading more is the trap; producing is the work.

---

## The core loop

Everything below is built on one repeatable cycle. Run it once per design. It *is* the method — the phases just change which prompts you feed it and how hard the grading gets.

1. **Pick one prompt** from the bank. Set a **45-minute timer**.
2. **Design blind.** Drive the full skeleton — requirements → estimates → API/data model → high-level → deep dive → failure modes — without looking at anything. Peeking mid-design turns the rep back into reading and wastes it. A flawed design you built beats a perfect one you read.
3. **Get it graded.** Paste it to me with an explicit instruction to grade *hard* and name the single weakest point. (Why "hard": left unprompted I drift generous; you have to pull the strict frame. Don't ask "is this good" — ask "where does this fall apart.")
4. **Compare against a key** when one exists (SD Vol II / HelloInterview). Catches whole missing dimensions — "you never discussed consistency," "you skipped estimates."
5. **Log the weak spot.** One line in a running list: the recurring failure pattern, not the per-problem detail. ("I keep skipping estimates." "I go shallow on the deep dive when the domain is unfamiliar.") This log is the highest-leverage artifact in the whole process — it's how you learn to self-grade.

The goal across the summer is to need step 3 less and trust step 5 more.

---

## Phase 1 — Bootstrap + calibrate the grader (now → ~Jul 13)

Roughly 4 weekends. Two jobs: make the skeleton reflexive, and verify the grader before you depend on it.

- **Reps 1–2: warm-ups** (URL shortener, rate limiter). Low domain load so you drill the *6-step structure*, not the content. The point is making the flow automatic.
- **Rep 3: your job scheduler.** You can't get lost in the domain — you built it. This is the confidence rep and your first real deep dive.
- **Calibration runs:** for at least two of these, pick prompts SD Vol II actually covers (rate limiter, news feed). Design blind, I grade, then you hold my grade against Xu's solution side by side. **Where my reasoning lines up with his on the prompts you *can* check, you earn trust for the prompts you *can't*.** That's how a non-expert bootstraps a grader: anchor on the covered cases, then spend the earned trust on the uncovered ones.
- **Milestone:** first timed SD mock **~Jul 15** (human or AI-as-interviewer, one full prompt under live questioning). Low setup, and the feedback is most useful with runway left to act on it.

By end of Phase 1 you should: run the skeleton without thinking about it, have one deep dive you're proud of, and *know whether you trust the grading* — because you checked it against Xu.

---

## Phase 2 — Output, weighted to your edge (~Jul 13 → ~Aug 16)

Roughly 5 weekends. One full design per weekend, timed, graded hard. Now bias the prompts toward where you win:

- **Prioritize Tier 2 (fintech) + the two 🏠 infra prompts.** Payment idempotency, ledger-under-concurrency, KYC pipeline, the job scheduler, the monitoring system. These are your differentiators and — critically — the prompts Xu *doesn't* cover cleanly, so the answer-key anchor is weakest exactly where you lean on me hardest.
- **Because the anchor is weak there, raise the bar on the grader instead of the key:** demand the *mechanism* behind every grade. A real grade says "weak *because* the processor-timeout case leaves you not knowing if the charge landed, and you didn't handle it" — checkable, you can reason about whether it holds. "Feels shallow" with no mechanism is the tell to distrust the score. This safeguard works even with no answer key.
- **Convert the 🏠 prompts twice:** once as a fresh design, once as "here's how I actually shipped it." Those become behavioral stories too.
- **Milestone:** by mid-Phase 2 you should be reaching a genuine deep dive *every time*, unprompted. If you're still staying broad-and-shallow, that's the weak-spot log telling you where to drill.

---

## Phase 3 — Delivery polish under pressure (~Aug 16 → end Aug)

Roughly 2–3 weekends. The shift here is from *content* to *delivery* — the one dimension reading-and-grading can't build.

- **Human mocks become the priority.** Telling a design out loud under time pressure, managing the clock, holding composure when the interviewer looks unconvinced — that's a different skill from producing a correct design, and only live reps build it. This is where a human (or a live AI-interviewer session) earns its place over written grading.
- **Back-to-back loops** (coding then SD in one sitting) to rehearse under fatigue, not just isolated rounds.
- **Revisit any prompt where you stalled** on the deep dive in Phase 2.
- **Milestone:** you can drive an unfamiliar prompt to a real deep dive in 45 minutes, out loud, and name your own design's weakest point before the interviewer does.

---

## The self-grading rubric

Learn these dimensions cold — they're what a senior SD round is actually scored on, and the muscle that makes you independent of any grader:

1. **Scope** — did you gather requirements and *cut* ruthlessly, or try to build everything?
2. **Estimates** — did you size it (QPS, storage, bandwidth) enough to justify later choices?
3. **Storage justification** — did you *defend* the DB/store choice against alternatives, not just name one?
4. **Deep dive** — did you reach genuine depth on 1–2 hard parts? **This is where senior is won or lost.** Most weak designs are detectably weak because they stayed broad and never went deep on anything.
5. **Failure modes** — bottlenecks, what breaks first at 10x, what you consciously punted.

**Tells of a bad design (self-detectable once you know them):** everything stayed a black box and never got concrete; you never reached a deep dive; you couldn't answer "why not X instead"; you quantified nothing; or scope creep ate the clock before the interesting part. Finish a design and check: did any of these fire? If none did, you probably went a reasonable direction.

---

## Grader safeguards (how to keep me honest)

- **Set the strict frame** every time — "grade hard, name the single weakest point, what would an interviewer dock." Don't ask "is this good."
- **Judge the reasoning, not the verdict.** You don't need SD expertise to catch a bad grade if I'm forced to justify it in terms you can follow. Demand the mechanism; if I can't name *why* in checkable terms, don't trust the score.
- **Anchor on covered prompts first** (Phase 1), spend the earned trust on uncovered ones (Phase 2).

---

## Division of labor

- **Content & reasoning — across all phases.** Whether the design is *right*: depth, tradeoffs, the deep dive. The larger gap, and the part written grading handles well.
- **Delivery & composure — Phase 3, human mocks.** Whether you can *perform* it live. Reading a design isn't watching you think out loud under pressure; that needs real reps.

Neither alone makes you ready. Content-correct but can't deliver = freezes in the room. Smooth but shallow = gets found out in the deep dive. The phases sequence them: get designs right first, make delivery smooth last.
