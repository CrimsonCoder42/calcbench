# algo/seq (SHORT = SEQ)

Builder W4-B15. New package (lead decision L3: `algo/seq`, not `algo/iter`). Not yet in `doc/taxonomy.md`
(lead to add: `ALG | algo/seq | SEQ | range-over-func helpers over iter.Seq / iter.Seq2`).

**Primary sources for every function** (confirmed 2026-11-12):
- The Go Programming Language Specification (Go 1.23+), "For statements with range clause", range over a
  function `func(yield func(K) bool)`: the loop body runs once per yield call; when the body terminates
  (break, return, panic) yield returns false and must not be called again. https://go.dev/ref/spec#For_range.
  (The fetched page was truncated before that section; the heading and rule are cited from the spec as
  published since Go 1.23 — the reviewer should open the anchor once.)
- Go standard library, package `iter` documentation, https://pkg.go.dev/iter: "Single-Use Iterators", "Pulling
  Values" — `Pull` returns (next, stop); stop "must be called when the caller is no longer interested in next
  values and next has not yet signaled that the sequence is over"; calling stop more than once, or after the
  end, is allowed (read 2026-11-12).

**Common shape.** Every helper is a pure function returning an iterator; nothing is read ahead; the helper's own
`for v := range upstream` returns as soon as the consumer's yield returns false, so the upstream's yield returns
false at once and it is never asked for another element (tested per helper with a probe upstream that records
produced / stopped / finished). A nil iterator is the empty sequence. A bare count or size returns
`(iter.Seq, error)` (run convention, L32). Lower bound for all: Ω(m) to read m elements; every helper is O(1)
(amortized) per element.

| ID | Name | Signature |
|---|---|---|
| ALG-SEQ-001 | Map | `func Map[T, U any](seq iter.Seq[T], fn func(T) U) iter.Seq[U]` |
| ALG-SEQ-002 | Filter | `func Filter[T any](seq iter.Seq[T], keep func(T) bool) iter.Seq[T]` |
| ALG-SEQ-003 | Take | `func Take[T any](seq iter.Seq[T], n int) (iter.Seq[T], error)` |
| ALG-SEQ-004 | Chunk | `func Chunk[T any](seq iter.Seq[T], size int) (iter.Seq[[]T], error)` |
| ALG-SEQ-005 | Zip | `func Zip[A, B any](first iter.Seq[A], second iter.Seq[B]) iter.Seq2[A, B]` |
| ALG-SEQ-006 | Flatten | `func Flatten[T any](seqs iter.Seq[iter.Seq[T]]) iter.Seq[T]` |
| ALG-SEQ-007 | Batch | `func Batch[T any](seq iter.Seq[T], maxSize int, endsBatch func(T) bool) (iter.Seq[[]T], error)` |
| ALG-SEQ-008 | Window | `func Window[T any](seq iter.Seq[T], width int) (iter.Seq[[]T], error)` |

### ALG-SEQ-001 Map / ALG-SEQ-002 Filter
- **Best known**: O(1) per element + the cost of fn/keep; O(1) extra. Met.
- **Approach**: range over seq, yield fn(v) / v when keep(v). **Lost**: a version taking iter.Seq2 too (one
  concept per function; a Seq2 variant can be added with its own ID if a consumer asks).
- **Errors**: none; nil seq or nil fn/keep is the empty sequence (documented zero value, never a panic).
- **Worked example**: Map(1 2 3, square) → 1 4 9; Filter(1…6, even) → 2 4 6. Hand-derived.
- **Allocation**: one closure per call, one per range loop (the loop body closure escapes into seq); none per
  element (tested: allocations over 10 and 10 000 elements are equal).

### ALG-SEQ-003 Take
- **Best known**: O(min(n, m)). Met: after the n-th element Take returns from its range loop, so the upstream
  is stopped at its n-th yield and never produces element n+1 (tested on an infinite probe).
- **Errors**: ErrCount (n < 0). n = 0 never starts the upstream.
- **Worked example**: Take(0 1 2 …, 3) → 0 1 2, upstream produced exactly 3. Hand-derived.

### ALG-SEQ-004 Chunk
- **Best known**: O(1) amortized per element, O(size) extra. Met (O(min(size, m)): the buffer grows by append,
  so a huge size never allocates up front or panics in make).
- **Decision (buffer)**: ONE reused slice per range loop, rewritten between groups (valid until the body
  returns; slices.Clone to keep). **Lost**: a fresh slice per group (one allocation per group; callers that
  only stream the groups pay for nothing they use). Reason: run convention L32 (enumeration iterators yield
  one reused buffer). Documented in doc.go and the doc comment; a test checks the backing array is reused.
- **Errors**: ErrSize (size < 1). Last group holds the remaining 1…size elements.
- **Worked example**: Chunk(1…7, 3) → [1 2 3] [4 5 6] [7]. Hand-derived. Oracle: re-slicing the input.
- **Not slices.Chunk**: the standard library's slices.Chunk splits a slice; this groups an iterator.

### ALG-SEQ-005 Zip
- **Best known**: O(1) per pair. Met.
- **Approach**: range over first; pull second with `iter.Pull`, created lazily at first's first element
  (second is not started when first is empty); `stop` is deferred and therefore ALWAYS called: on consumer
  break, when first ends, when second ends, and when the loop body panics (each tested). first is stopped by
  returning from its range loop. iter.Pull runs second as a coroutine (no goroutine of ours, no channel).
  **Lost**: pulling both with two iter.Pull (two coroutine switches per pair instead of one); collecting
  either side into a slice (O(m) memory, breaks infinite iterators).
- **Errors**: none; nil on either side is the empty sequence.
- **Worked example**: Zip(1 2 3, "x" "y") → (1,x) (2,y). Hand-derived.

### ALG-SEQ-006 Flatten
- **Best known**: O(1) per element and per inner iterator. Met.
- **Approach**: nested range loops; on stop, both the current inner and the outer iterator are stopped and
  no later inner is started (tested). nil inner iterators are skipped.
- **Errors**: none. **Worked example**: [1 2] [] [3] → 1 2 3. Hand-derived.

### ALG-SEQ-007 Batch
- **Decision**: KEPT, not dropped: Batch is NOT identical to Chunk. A batch closes at maxSize elements OR when
  `endsBatch(v)` is true for the element just added (that element is the batch's last). With a nil
  endsBatch it groups exactly as Chunk (tested). Use: records up to a terminator, messages up to a commit
  marker, with a size cap. **Lost**: dropping it (the predicate flush has no other home); a time-based flush
  (needs a clock or a timer goroutine: not a pure iterator).
- **Best known**: O(1) amortized per element; O(maxSize) extra. Met. Same reused-buffer rule as Chunk.
- **Errors**: ErrSize (maxSize < 1).
- **Worked example**: 5 0 1 2 3 4 0 6, terminator 0, maxSize 4 → [5 0] [1 2 3 4] [0] [6]. Hand-derived.

### ALG-SEQ-008 Window
- **Source (ring)**: Knuth, TAOCP Vol. 1, 3rd ed., 1997, §2.2.2 (circular array, indices modulo M).
- **Best known**: O(1) per element, O(width) extra, and no copy of the window per step. Met.
- **Approach**: a ring of 2·width slots; element i is written to slot i mod width AND slot i mod width +
  width, so the window starting at slot s is the contiguous slice [s, s+width): two writes per step, never a
  copy. The ring is allocated once the first window is full (before that, the elements are collected by
  append), so a huge width on a short seq costs nothing and never panics in make.
  **Lost**: a ring of width slots copied into a separate output slice per window (O(width) per step);
  `queue.CircularBuffer` (its elements are not contiguous, so each window would need a copy).
- **Errors**: ErrSize (width < 1). m < width gives no windows.
- **Worked example**: Window(1…5, 3) → [1 2 3] [2 3 4] [3 4 5]. Hand-derived. Oracle: in[i:i+width].

## Channel patterns: NOT library material (they own goroutines)
calcbench functions start no goroutine the caller did not ask for (CLAUDE.md, "What every function is"). The
code-heim channel exercises are listed here only to record why they are not built:
- go_61 errgroup (run tasks, cancel on first error): owns goroutines and a context; that is
  golang.org/x/sync/errgroup's job, not a pure function.
- go_62 or-done (wrap a channel so a done signal ends the read): a goroutine forwarding between channels.
- go_62 tee (copy one channel to two): a goroutine per tee with unbounded blocking coupling.
- go_56 / go_60 generators (a goroutine feeding a channel): replaced by iter.Seq, which needs no goroutine.
- fan-out / fan-in (several workers, merged results): goroutines and synchronisation by definition.

## Tags (problem sites, never in identifiers)
Take/Chunk/Window/Zip/Flatten: iterator-pattern exercises (code-heim go_56/go_60 generators reworked as
iter.Seq); Window: sliding-window pattern (LeetCode 239/643 use the same windows).
