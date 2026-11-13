# S2 algorithms core — common brief for every builder (read this first, then your own brief)

Repo: /Users/devinanderson/Desktop/calcbench (module github.com/CrimsonCoder42/calcbench, Go 1.27). The law:
/Users/devinanderson/Desktop/calcbench/.claude/CLAUDE.md and every file in .claude/rules/ (go-functions.md,
go-testing.md, go-complexity.md, numeric.md, sources.md). Read them. Model your code on S1:
algo/window/rolling_sum.go, algo/search/lower_bound.go, algo/heap/top_k.go, algo/selection/kth_smallest.go
and their _test.go files and doc.go files (doc comment shape, error style, test style, benchmark style).

## Role: BUILDER
- Edit ONLY your owned paths (your brief). Other subject agents and other builders are live in this tree.
  A file you do not own is never touched, even if it is red: report it.
- Allowed checks on your own files: gofmt -l, go build ./algo/..., go vet ./algo/<your packages>/... . Do NOT
  run go test or benchmarks (one single test by name only if you cannot continue without the answer; say so).
- Write every test and benchmark for what you build. Never run them.
- No git writes. Never rm: a delete is `mkdir -p` then `mv -n <path> /Users/devinanderson/Desktop/calcbench/_delete/<same path>`.
- Standard library only (container/heap, slices, maps, cmp, iter, math/bits, hash/fnv, hash/maphash, math/rand/v2
  sources passed in by the caller, time.Time passed in). No goroutines, no globals (except sentinel errors and
  constant tables), no clock, no randomness unless a rand.Source / *rand.Rand is a parameter.
- Absolute paths; print the path before each write.

## What every function needs (see the rules for detail)
1. A catalog entry FIRST, in your section file doc/algorithms_core_2026-11-12/catalog/<package>.md, in the S1
   shape: ID, signature, primary source cited to the line (author, title, edition/journal, year, section /
   algorithm / page / DOI), best known time and space AND the lower bound where known, the approach chosen and
   the one that lost (and why), the worked example that becomes the first test, errors, whether input is
   modified, consumer evidence (file:line from your brief) when given, LeetCode problem numbers ONLY under Tags.
2. Code: one concept per file, snake_case file named for the concept; doc comment exactly per
   go-functions.md (first line starts with the name; ID, Time, Space, Source (one or more), Tags in order;
   Time/Space define every variable; "Space" says extra vs result). Generic over cmp.Ordered or a small
   constraint where type-independent. Large outputs: dst slices (dst == nil allocates once) or iter.Seq
   iterators. Inputs never modified unless the name says so (…InPlace) and the doc says so.
   Invalid input returns a sentinel error (package-level var ErrX = errors.New("<pkg>: ..."), wrapped with the
   parameter name) or a documented zero value; never panic on input. Integer overflow checked where inputs
   are unbounded (math/bits).
3. Tests (<concept>_test.go): the source's worked example cited in a comment; table cases (empty, one, two,
   duplicates, sorted, reversed, negative, zero, extremes); every documented error; an oracle against a slow
   obvious version on a few hundred generated inputs with rand.New(rand.NewPCG(1, 2)); the LeetCode problem's
   examples as plain test cases written in OUR OWN WORDS (no copied statements); AllocsPerRun asserting 0
   where you claim no allocation; test names Test<Function>_<case>.
4. Benchmark Benchmark<Function>/<size> with b.ReportAllocs() at 1e3/1e4/1e5 (or the sizes your complexity
   needs; e.g. O(n²) DP at 1e2/1e3/1e4 lengths that finish in a second; O(2ⁿ) enumerations at small n).
5. doc.go: package comment, "ID short form: <SHORT> (IDs are ALG-<SHORT>-NNN)", conventions, and the list of
   every function with its ID (S1 doc.go shape). For an EXISTING package (search, selection, window, heap,
   interval) you own doc.go too: ADD your lines, keep the existing ones word for word.

## Sources (.claude/rules/sources.md)
Acceptable primary: CLRS 3rd ed., Knuth TAOCP vols 1–4A, Sedgewick & Wayne Algorithms 4th ed., Skiena (Algorithm
Design Manual), Bentley (Programming Pearls), Gusfield (strings), Crochemore & Rytter, de Berg et al.
(Computational Geometry), and the original papers (Tarjan 1975, Dijkstra 1959, Bellman 1958, Floyd 1962, Kahn
1962, Kruskal 1956, Prim 1957, Manacher 1975, Knuth–Morris–Pratt 1977, Karp–Rabin 1987, Fenwick 1994, Bloom 1970,
Kirsch–Mitzenmacher 2006, Vitter 1985, Stein 1967, Andrew 1979, Wagner–Fischer 1974, Hirschberg 1975,
Hunt–Szymanski, Welford 1962, …), the Go language specification (range-over-func), Go standard library docs.
NOT acceptable as primary: LeetCode, code-heim, blogs, Q&A sites, Wikipedia (pointer only), any library source.
Confirm each citation (WebSearch/WebFetch: DOI on Crossref, table of contents, or the paper itself). A function
whose method has no citable source is NOT built: list it as HELD (no source) in your report.
LeetCode/code-heim: they tell you WHICH problems recur; never copy statements, examples, code, structure or
identifiers. Our names are the practitioner's terms.

## Existing IDs (never reuse; continue numbering)
ALG-SEARCH-001…004 (LowerBound, UpperBound, BinarySearchFirst, BinarySearchLast) · ALG-SELECT-001 KthSmallest,
-002 Median · ALG-WINDOW-001…004 RollingSum/Mean/Min/Max · ALG-HEAP-001 TopK · ALG-MERGE-001 MergeKSorted ·
ALG-INTVL-001 MergeOverlapping. A new package numbers from 001. Reuse an existing function by CALLING it (algo
packages may import each other, acyclic); never copy it.

## Report (≤ 10 lines, to the lead)
Functions done with IDs and file:line · tests and benchmarks written · HELD functions and why · anything
needed outside your paths · anything unfinished and why. Do not run tests.

## Rule update (.claude/rules/go-testing.md, changed during this run)
Every oracle and property test caps its input size and iteration count; no test may run longer than 60 s alone,
and the whole package must run under 10 minutes; a Fuzz function is run only with -fuzztime. Size your oracles
(and exponential enumerations) accordingly.
