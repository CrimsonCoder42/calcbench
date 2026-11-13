# algo/stack (SHORT = STACK)

Builder W1-B3. Not yet in `doc/taxonomy.md` (lead to add: `ALG | algo/stack | STACK | monotonic stacks,
nearest greater/smaller, stack evaluators`).

Common to the package: ordering is `cmp.Compare` / `cmp.Less` (D-002: NaN is the least value, -0 equals +0).
Functions that return one index per element append to `dst` and return the extended slice (D-003): nil
allocates once, `buf[:0]` reuses a buffer. No function modifies its input.

| ID | Name | Signature |
|---|---|---|
| ALG-STACK-001 | NextGreater | `func NextGreater[T cmp.Ordered](dst []int, values []T) []int` |
| ALG-STACK-002 | NextSmaller | `func NextSmaller[T cmp.Ordered](dst []int, values []T) []int` |
| ALG-STACK-003 | PreviousGreater | `func PreviousGreater[T cmp.Ordered](dst []int, values []T) []int` |
| ALG-STACK-004 | PreviousSmaller | `func PreviousSmaller[T cmp.Ordered](dst []int, values []T) []int` |
| ALG-STACK-005 | LargestRectangleHistogram | `func LargestRectangleHistogram(heights []int64) (int64, error)` |
| ALG-STACK-006 | MaximalRectangle | `func MaximalRectangle(grid [][]bool) (int, error)` |
| ALG-STACK-007 | DailyTemperatures | `func DailyTemperatures[T cmp.Ordered](dst []int, readings []T) []int` |
| ALG-STACK-008 | MinStack | `type MinStack[T cmp.Ordered]` with `Push(T)`, `Pop() (T, bool)`, `Top() (T, bool)`, `Min() (T, bool)`, `Len() int` |
| ALG-STACK-009 | EvaluateRPN | `func EvaluateRPN(tokens []string) (int64, error)` |
| ALG-STACK-010 | BalancedBrackets | `func BalancedBrackets(text, pairs string) (bool, error)` |
| — | SimplifyPath | DROPPED: duplicates `path.Clean` (below) |

### ALG-STACK-001 NextGreater · -002 NextSmaller · -003 PreviousGreater · -004 PreviousSmaller
- **Source**: the all-nearest-smaller-values (ANSV) problem and its sequential stack algorithm: Berkman,
  Schieber, Vishkin, "Optimal doubly logarithmic parallel algorithms based on finding all nearest smaller
  values", Journal of Algorithms 14(3):344–370, 1993, doi:10.1006/jagm.1993.1018 (verified on Crossref
  2026-11-12). Earlier: Gabow, Bentley, Tarjan, "Scaling and related techniques for geometry problems",
  STOC '84, pp. 135–143, 1984, doi:10.1145/800057.808675 (verified on Crossref). The stack held as a linked
  list: Knuth, TAOCP Vol. 1, 3rd ed., 1997, §2.2.3 (linked allocation of a stack).
- **Definition**: NextGreater[i] = the least j > i with values[j] > values[i] (strict), else −1. NextSmaller
  uses <. PreviousGreater/PreviousSmaller: the greatest j < i, strict, else −1.
- **Best known**: O(n) time, Ω(n) to read the input; the textbook stack is O(n) extra.
- **Approach**: the monotonic stack kept as a linked list through the result itself: scanning from the far
  end, the candidates for i are i±1, then result[i±1], then result[result[i±1]], …; every skipped index is
  skipped for good, so the total work is O(n) amortized and no stack is allocated: O(1) extra.
  **Lost**: an explicit index stack (same O(n) time, O(n) extra and a second allocation); the O(n²) scan.
- **Input modified**: no. **Errors**: none (no input is invalid).
- **Worked example**: 2 7 3 5 4 6 8 (hand-derived). NextGreater 1 6 3 5 5 6 −1; NextSmaller −1 2 −1 4 −1 −1 −1;
  PreviousGreater −1 −1 1 1 3 1 −1; PreviousSmaller −1 0 0 2 2 4 5.
- **Tags**: LeetCode 496, 503, 739, 901.

### ALG-STACK-005 LargestRectangleHistogram
- **Source**: every maximal rectangle under a histogram has the height of some bar i and spans exactly from
  its previous smaller bar to its next smaller bar (exclusive) — the ANSV application, Berkman, Schieber,
  Vishkin 1993 (above), with Gabow, Bentley, Tarjan 1984 (above).
- **Best known**: O(n) time, Ω(n); O(n) extra (the stack). Lower bound on extra space not known to us.
- **Approach**: one left-to-right pass with an index stack of non-decreasing heights; a bar's rectangle is
  closed when it is popped (its right limit is the popping index, its left limit the new top). Popping on
  ≥ makes the last of a run of equal bars compute the full width. **Lost**: PreviousSmaller + NextSmaller
  as two calls (O(n) extra twice, two result arrays); divide and conquer on the minimum (O(n log n), or O(n)
  only with a Cartesian tree, more memory).
- **Units**: heights in any unit u; the area is in u × bar widths.
- **Input modified**: no. **Errors**: `ErrNegativeHeight` (a height < 0); `ErrOverflow` (height × width
  exceeds int64). Empty → 0.
- **Worked example**: our own heights 3 1 4 4 3 5 1 → 12 (bars 2–5 at height 3, width 4). Hand-derived.
- **Tags**: LeetCode 84.

### ALG-STACK-006 MaximalRectangle
- **Source**: reduction to the histogram, row by row: the largest all-true rectangle has a bottom row r;
  above row r, column c holds a run of h_r[c] trues ending at r, and a rectangle with bottom row r fits
  exactly when its height ≤ min h_r over its columns, which is the histogram problem on h_r. This proof is
  written here; the histogram step is ALG-STACK-005 (ANSV, Berkman, Schieber, Vishkin 1993). The reduction
  is folklore in the largest-empty-rectangle literature (Naamad, Lee, Hsu, "On the maximum empty rectangle
  problem", Discrete Applied Mathematics 8(3):267–277, 1984, doi:10.1016/0166-218X(84)90124-0, verified on
  Crossref 2026-11-12); a primary source stating this exact reduction
  was not found and is FLAGGED for Devin.
- **Best known**: O(R·C) time, Ω(R·C) to read the grid; O(C) extra.
- **Approach**: update h_r in place per row; run the histogram pass with one reused stack. **Lost**: the
  dynamic-programming left/right/height arrays (also O(R·C), three O(C) arrays); enumerating rectangles
  (O(R²C²)).
- **Input modified**: no. **Errors**: `ErrRagged` (rows of different lengths). Empty grid → 0.
- **Worked example**: our own 4×5 grid (rows 01101, 11101, 11111, 01110) → 8 (columns 1–2, all four rows).
  Hand-derived.
- **Tags**: LeetCode 85.

### ALG-STACK-007 DailyTemperatures
- **Source**: NextGreater (ALG-STACK-001), by a call; the wait is next − i, or 0 when none.
- **Best known**: O(n), O(1) extra (as NextGreater).
- **Units**: readings in any unit, compared only; the result counts positions (days when there is one
  reading per day).
- **Worked example**: our own readings 15 18 16 14 17 20 19 19 → 1 4 2 1 1 0 0 0. Hand-derived.
- **Tags**: LeetCode 739.

### ALG-STACK-008 MinStack
- **Source**: Gajewska, Tarjan, "Deques with heap order", Information Processing Letters 22(4):197–200, 1986,
  doi:10.1016/0020-0190(86)90028-1 (find-min in O(1) worst case on a deque; the stack case keeps, with
  each element, the minimum of the elements below it).
- **Best known**: O(1) per operation; Push O(1) amortized here (slice growth), O(n) space.
- **Approach**: each slot stores (value, minimum so far). **Lost**: an auxiliary stack of minima pushed only
  on ≤ (saves space when minima are rare, but two slices and a second comparison per pop).
- **Errors**: none; Pop, Top and Min on an empty stack return (zero, false). Zero value is an empty stack.
- **Worked example**: push 5 3 7 3 → Min 3; pop, pop → Min 3; pop → Min 5. Hand-derived.
- **Tags**: LeetCode 155.

### ALG-STACK-009 EvaluateRPN
- **Source**: Hamblin, "Translation to and from Polish notation", The Computer Journal 5(3):210–213, 1962,
  doi:10.1093/comjnl/5.3.210 (verified on Crossref); Knuth, TAOCP Vol. 1, 3rd ed., 1997, §2.2.1 (stacks);
  Sedgewick & Wayne, *Algorithms*, 4th ed., 2011, §1.3, Exercise 1.3.11 (EvaluatePostfix) beside Dijkstra's
  two-stack evaluator in the same section.
- **Best known**: O(n) time for n tokens (each token once), O(n) extra for the operand stack.
- **Semantics**: tokens are base-10 int64 literals (an optional sign) or one of + − * /. Division truncates
  toward zero (the Go specification, "Arithmetic operators"). Results are exact or an error.
- **Errors**: `ErrMalformed` (no tokens, an operator with fewer than two operands, or more than one value
  left); `ErrBadToken` (neither a literal nor an operator); `ErrOverflow` (a literal or a result outside
  int64, including MinInt64 / −1); `ErrDivideByZero`.
- **Worked example**: Hamblin's reverse Polish order on our own expression (5 + 1) × (12 − 4) / 3 → tokens
  5 1 + 12 4 − * 3 / → 16. Hand-derived.
- **Tags**: LeetCode 150.

### ALG-STACK-010 BalancedBrackets
- **Source**: Sedgewick & Wayne, *Algorithms*, 4th ed., 2011, §1.3, Exercise 1.3.4 (Parentheses: a stack
  client that checks bracket balance); Knuth, TAOCP Vol. 1, 3rd ed., 1997, §2.2.1.
- **Semantics**: pairs is a string of opener/closer rune pairs, e.g. "()[]{}"; other runes are ignored.
  Balanced: every closer matches the most recent unmatched opener, and none is left open.
- **Best known**: O(n + p) time (n runes of text, p runes of pairs), O(n) extra worst case (the stack).
- **Errors**: `ErrPairs` (odd rune count, an opener equal to its closer, a rune in two roles, U+FFFD, or
  invalid UTF-8).
- **Worked example**: Sedgewick & Wayne's own: "[()]{}{[()()]()}" → true, "[(])" → false.
- **Tags**: LeetCode 20.

### SimplifyPath — DROPPED
Unix path canonicalisation (collapse "//", drop ".", resolve ".." lexically, root stays root) is exactly the
standard library's `path.Clean` on a rooted path: Go standard library, package path, func Clean (rules 1–4,
read 2026-11-12), which cites Pike, "Lexical File Names in Plan 9 or Getting Dot-Dot Right". POSIX pathname
resolution (IEEE Std 1003.1, Base Definitions, Pathname Resolution) resolves ".." through the file system
(symbolic links), which no pure function can do. Callers use `path.Clean`. No ID given.
