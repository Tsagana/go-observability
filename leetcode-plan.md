# LeetCode Plan — 4-Month Ramp in Go

**Profile:** Backend engineer, rusty on LC, solid on easies / shaky on mediums. Targeting general strong interview prep (FAANG-capable) in 3–6 months, solving in Go.

**Core philosophy:** Depth over volume. ~150–200 well-understood problems beats 500 half-remembered ones. You already have good systems intuition — the job is translating that into fast pattern recognition under interview pressure.

---

## Why Go is fine (with one caveat)

Go is interview-acceptable at every company that isn't explicitly "Python/Java/C++ only" (rare). Pros: no context-switching from your day job, forces you to be explicit about types and allocations. Caveat: Go's stdlib for LC is thinner — no built-in heap with a clean API, no ordered set, verbose slice manipulation. Build a personal `lc-utils` package early with:

- `MinHeap[T]` / `MaxHeap[T]` wrappers around `container/heap`
- Stack/queue helpers on generic slices
- `Deque[T]` for sliding window / monotonic queue problems
- `TreeNode`, `ListNode` definitions + a `BuildTree(vals []any)` helper
- A `UnionFind` struct

Spending 2–3 hours on this in Week 1 saves you from re-implementing heap boilerplate 40 times.

---

## Time tiers — pick one, adjust monthly

| Tier | Hrs/week | Problems/week | Pace | Fits if... |
|------|----------|---------------|------|-----------|
| **Light** | 4–5 | 4–5 | Finish core 150 in ~4 months | You want to keep the job-system build as primary |
| **Standard** | 7–9 | 7–8 | Finish core 150 + 50 mixed in ~4 months | Balanced — recommended default |
| **Intense** | 12–15 | 12–15 | Core 150 + 100 mixed + mocks in 3 months | You want to interview in ~3 months |

**Recommendation for you:** Start **Standard**. Drop to Light in weeks where the job system is heating up (V2a/V2b reaper work will be mentally demanding). Ramp to Intense in the final 4–6 weeks before applying.

---

## The 4-month structure

### Month 1 — Rebuild the foundation (patterns, not problems)

The goal is pattern fluency, not problem count. For each topic: read the pattern, do 2 easies to lock mechanics, then 4–6 mediums. Review solutions even when you solve it — there's almost always a cleaner approach.

**Week 1 — Arrays & Hashing**
Two Sum, Contains Duplicate, Valid Anagram, Group Anagrams, Top K Frequent Elements, Product of Array Except Self, Valid Sudoku, Encode and Decode Strings, Longest Consecutive Sequence.  - DONE

**Week 2 — Two Pointers & Sliding Window**
Valid Palindrome, Two Sum II, 3Sum, Container With Most Water, Trapping Rain Water (hard, attempt), Best Time to Buy/Sell Stock, Longest Substring Without Repeating, Longest Repeating Character Replacement, Minimum Window Substring.

**Week 3 — Stack & Binary Search**
Valid Parentheses, Min Stack, Evaluate RPN, Generate Parentheses, Daily Temperatures, Car Fleet, Binary Search, Search 2D Matrix, Koko Eating Bananas, Find Minimum in Rotated Sorted Array, Search in Rotated Sorted Array, Time Based Key-Value Store.

**Week 4 — Linked Lists**
Reverse Linked List, Merge Two Sorted Lists, Reorder List, Remove Nth Node From End, Copy List with Random Pointer, Add Two Numbers, Linked List Cycle, Find the Duplicate Number, LRU Cache (important — intersects with your systems work).

End of month 1: you should be solving straightforward mediums in 25–35 min without looking at solutions.

### Month 2 — Trees, Tries, Heaps, Backtracking

**Week 5 — Trees (part 1)**
Invert Binary Tree, Max Depth, Diameter, Balanced Binary Tree, Same Tree, Subtree of Another Tree, Lowest Common Ancestor of BST, Level Order Traversal, Right Side View.

**Week 6 — Trees (part 2) & Tries**
Count Good Nodes, Validate BST, Kth Smallest in BST, Construct Tree from Preorder+Inorder, Binary Tree Max Path Sum (hard), Serialize/Deserialize Binary Tree (hard), Implement Trie, Design Add/Search Word, Word Search II (hard).

**Week 7 — Heap / Priority Queue + Intervals**
Kth Largest in Stream, Last Stone Weight, K Closest Points to Origin, Kth Largest in Array, Task Scheduler, Design Twitter, Find Median from Data Stream (hard), Insert Interval, Merge Intervals, Non-overlapping Intervals, Meeting Rooms I/II, Minimum Interval to Include Each Query (hard).

**Week 8 — Backtracking**
Subsets, Combination Sum, Permutations, Subsets II, Combination Sum II, Word Search, Palindrome Partitioning, Letter Combinations of Phone Number, N-Queens (hard).

### Month 3 — Graphs, DP, Greedy

This is where it gets hard. Budget more time per problem; these patterns take longer to internalize.

**Week 9 — Graphs (BFS/DFS)**
Number of Islands, Clone Graph, Max Area of Island, Pacific Atlantic Water Flow, Surrounded Regions, Rotting Oranges, Walls and Gates, Course Schedule, Course Schedule II, Redundant Connection, Number of Connected Components, Graph Valid Tree, Word Ladder (hard).

**Week 10 — Advanced Graphs**
Reconstruct Itinerary, Min Cost to Connect All Points (MST), Network Delay Time (Dijkstra), Swim in Rising Water (Dijkstra variant), Alien Dictionary (hard), Cheapest Flights Within K Stops (Bellman-Ford).

**Week 11 — 1-D DP**
Climbing Stairs, Min Cost Climbing Stairs, House Robber, House Robber II, Longest Palindromic Substring, Palindromic Substrings, Decode Ways, Coin Change, Maximum Product Subarray, Word Break, Longest Increasing Subsequence, Partition Equal Subset Sum.

**Week 12 — 2-D DP & Greedy**
Unique Paths, Longest Common Subsequence, Best Time to Buy/Sell Stock with Cooldown, Coin Change II, Target Sum, Edit Distance, Maximum Subarray, Jump Game, Jump Game II, Gas Station, Hand of Straights, Merge Triplets to Form Target, Partition Labels.

### Month 4 — Consolidation + mock mode

**Week 13 — Weak spots**
Look at your tracking sheet (see below). Re-do the 15–20 problems you solved slowest or needed hints on. Without looking at your old solution, re-solve from scratch. If you can't, study the pattern again.

**Week 14 — Bit manipulation, math, and odd patterns**
Single Number, Number of 1 Bits, Counting Bits, Reverse Bits, Missing Number, Sum of Two Integers, Rotate Image, Spiral Matrix, Set Matrix Zeroes, Happy Number, Plus One, Pow(x, n).

**Week 15 — Mock interviews**
4–5 mock sessions, 45 min each, on random unseen mediums. Use Pramp (free, peer-to-peer) or interviewing.io. Explain out loud while coding. This is the single biggest gap between "can solve LC" and "passes interviews" — and it's the skill that's cheapest to neglect and costliest to lack.

**Week 16 — Company-targeting + hards**
Once you pick target companies, pull their tagged question list (LC Premium or Grind 75). Do 15–20 of their tagged mediums. Attempt 5 hards in your weak areas.

---

## How to actually practice (the part most people skip)

**The 30-minute rule.** Set a timer. If you haven't made real progress in 30 min on a medium, stop and read the editorial. You're not learning from staring — you're pattern-matching against patterns you don't yet have. After reading, close it, wait 5 min, and re-solve from scratch.

**Write, don't just read.** After each problem, write 2–3 sentences in a tracking doc: what pattern, what the key insight was, where you got stuck. This is where retention comes from — not the solving itself.

**Spaced repetition.** Any problem that took >45 min or needed a hint goes on a review list. Re-solve it 3 days later, 1 week later, 3 weeks later. Most people skip this. It's probably a 3x multiplier on retention.

**Explain out loud.** Starting Week 3, narrate your thought process while coding — even alone. Interview performance is 50% communication. Silent solving builds a habit you'll have to break later.

**Don't chase the daily.** LC's daily problem is random and often niche. Stick to the structured list unless the daily happens to fit your current week's topic.

---

## Tracking — keep it dead simple

A single markdown file or spreadsheet with one row per problem:

```
| # | Name | Pattern | Date | Time | Solved unaided? | Notes |
```

Review the "no" and "slow" rows weekly. That's your real weak-spot signal — more reliable than generic topic tags.

---

## Realistic expectations

- **Month 1:** mediums feel hard, ~40 min average, you'll google often. Normal.
- **Month 2:** mediums in 25–30 min, patterns start clicking. You'll feel a real jump around week 6.
- **Month 3:** comfortable with most mediums, hards are attemptable. DP will be the hardest stretch — budget extra patience.
- **Month 4:** most mediums in <25 min, you can explain your approach cleanly, hards are 50/50.

At month 4 with standard-tier effort, you're ready for real interviews at strong companies. Not FAANG-guaranteed — nothing is — but in fighting shape.

---

## Resources

- **NeetCode 150** (free, neetcode.io) — the list above is approximately this. Excellent video explanations.
- **LeetCode Premium** — worth it for the last 4–6 weeks when company-targeting. Not before.
- **Grind 75** (grind75.com) — alternative ordering, good for mock-interview selection.
- **Pramp / interviewing.io** — free mock interviews with real humans.
- **Book (optional):** *Elements of Programming Interviews* has better problems than CTCI, if you want depth beyond LC.

---

## Pitfalls to avoid

- Solving and moving on without reviewing the editorial.
- Doing only easies because they feel good.
- Skipping DP because it's hard. DP is ~15% of real interviews — don't let it be 0% of your prep.
- Memorizing solutions. You want the pattern, not the code.
- Neglecting mock interviews until "you're ready." You're never ready. Start in week 15 no matter what.
