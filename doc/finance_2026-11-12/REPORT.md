# Finance report (2026-11-12): catalog and build

**Counts** (190 vCalc rows): KEEP 105 · MERGE 41 · DROP 37 · MOVE 7 (6 to the proposed `fin/personal`, 1 to `life/automotive`),
plus 23 standard-set functions. Cataloged 137 IDs: 122 built, 13 held back (sources), 2 removed (duplicates).

| Package | Built | go test (final, out/go_test_final.txt) | Benchmarks (out/bench_final.txt) |
|---|---|---|---|
| fin/interest | 40 | ok 0.223s (219 test funcs) | 44; scalars 4–30 ns, 0 allocs; IRR 14.7µs→196µs→2.33ms (1e3/4/5), 0 allocs |
| fin/loan | 10 | ok 0.824s (52) | 18; schedule 15.7µs→182µs→2.12ms, 1 alloc (result) |
| fin/invest | 27 | ok 0.415s (124) | 41; 0 allocs everywhere; Macaulay 1.6µs→15µs→136µs |
| fin/business | 41 | ok 0.794s (206) | 45; scalars 0 allocs; depreciation schedules 1 alloc, linear |
| fin/economics | 4 | ok 0.951s (20) | 6; 0 allocs |

**Verify:**
- First run: 11 failures.
- Reviewer: found 7 defect classes, and fixes went out to the four package agents (summarised below).
- Re-run: 1 failure, the hidden-tie test, caused by FMA fusion.
- After that fix: all 5 packages pass.

**Reviewer findings → fixes:**
- **YieldToMaturity panicked** on a huge `periods` → it now uses the new FIN-INT-042 AnnuityRate (closed form, 0 allocations).
- **Loops could hang** on huge period counts → caps added (F-10).
- **Declining balance broke half-cent ties** at ordinary sizes → multiply before dividing.
- **Amortization:** the last-payment bound was understated and tiny loans went negative → the true bound is documented and the schedule ends early (F-11).
- **Float money results** were off by up to 6¢ near 2⁵³ → the F-5 contract is stated in every doc. EVA, FCFF and FCFE are now exact (F-9).
- **Black-Scholes was unstable at huge σ** → a stable d₁.
- **Four functions returned NaN/Inf** with a nil error → they now return ErrOverflow.
- **Three tests were wrong** (FixedChargeCoverageRatio expected value, CPI truncation, a payback tie) → the tests were fixed, not the code.
- **Two duplicates** → removed (F-8).
- **Found by the lead:** Go may fuse x·r − c into an FMA, which silently undid the exactness trick → explicit `float64(…)` conversions (F-9).

**Held back (D-015):** FIN-INT-028, 041; FIN-INVEST-008, 029, 031; FIN-BIZ-006, 009, 010, 011, 013, 039, 044, 045. Reasons are in catalog/finance.md, "Build status".

**Decisions:** F-1…F-13 and assumptions A-fin-1…9, in PENDING_SHARED.md for you to apply.

**FOUND, NOT FIXED:**
- The OpenStax Principles of Finance author is spelled four ways across Source lines. The book is the same; the strings differ.
- The reviewer had no network access, so its "15 URLs opened" check was not done. The package agents confirmed every citation, and I re-checked the doc/sources entries against the code text.
- The research files' "Covers rows" lines are empty. The Step 1 table in the catalog maps the rows instead.
- AnnualizedHoldingPeriodYield's `days` signature cannot express its source's example (m = 4), so that example is not tested.

**LEFT FOR DEVIN:**
- Apply PENDING_SHARED.md (INDEX, SCOPE, taxonomy, DECISIONS, ASSUMPTIONS, CHANGELOG, sources README).
- Run `kill 4323 29142`: two reviewer fuzz processes are stuck, and kill is denied to me.
- Sources for the 13 held-back functions.
- Rule on Harris read through a mirror, and on SOA and Damodaran (accepted by me as F-2/F-3).
