# algo/geometry2d (SHORT = GEO2D)

Builder W4-B13. New package, assigned by the lead's brief; not yet in `doc/taxonomy.md` (lead to add:
`ALG | algo/geometry2d | GEO2D | computational geometry over planar point sets: orientation, hulls, closest
pair, point location, segment intersection`). `num/geometry` keeps the measurement formulas (PolygonArea,
PointDistance, RectangleArea = w·h); this package does not import `num/*`.

## Package decision: exact int64 coordinates only (builder choice, for the lead to log)

`Point{X, Y int64}`. Every predicate is exact on the whole int64 range: a coordinate difference is carried as
a 64-bit magnitude and a sign (65 bits), a product of two differences as an unsigned 128-bit magnitude
(`math/bits.Mul64`) and a sign, and a sum of two squares in 192 bits. Nothing overflows, so no input is
rejected for size. **Lost**: (a) `float64` points with Shewchuk's adaptive predicates (Shewchuk, "Adaptive
precision floating-point arithmetic and fast robust geometric predicates", Discrete & Computational Geometry
18:305–363, 1997, doi:10.1007/PL00009321): several hundred lines of expansion arithmetic for a predicate the
integer type gives exactly; callers with real-valued data scale to a fixed grid (for example micrometres in
int64) and lose nothing a float would keep. (b) `float64` with an epsilon: not robust (hull and
point-location answers become inconsistent near degeneracies). (c) plain int64 products with a coordinate cap (|coordinate| < 2³⁰ for the cross product to fit):
an invented cap that rejects valid input, for a saving of a few nanoseconds. The
consequence: no float segment-intersection POINT is offered (it would need the float type); the predicate
SegmentsIntersect is exact.

Units: coordinates are in any one length unit u the caller chooses; areas are in u².

| ID | Name | Signature |
|---|---|---|
| ALG-GEO2D-001 | Orientation | `func Orientation(a, b, c Point) int` |
| ALG-GEO2D-002 | ConvexHull | `func ConvexHull(dst, points []Point) []Point` |
| ALG-GEO2D-003 | ClosestPair | `func ClosestPair(points []Point) (i, j int, err error)` |
| ALG-GEO2D-004 | PointInPolygon | `func PointInPolygon(polygon []Point, q Point) (Location, error)` |
| ALG-GEO2D-005 | SegmentsIntersect | `func SegmentsIntersect(s, t Segment) bool` |
| ALG-GEO2D-006 | MaxPointsOnLine | `func MaxPointsOnLine(points []Point) int` |
| ALG-GEO2D-007 | RectanglesUnionArea | `func RectanglesUnionArea(r1, r2 Rect) (uint64, error)` |
| ALG-GEO2D-008 | HasVerticalSymmetryAxis | `func HasVerticalSymmetryAxis(points []Point) bool` |
| — | PerfectRectangle | HELD (no source) |

### ALG-GEO2D-001 Orientation
- **Source**: CLRS (Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009),
  §33.1 "Line-segment properties", the cross product (p1 − p0) × (p2 − p0) and DIRECTION. O'Rourke,
  Computational Geometry in C, 2nd ed., Cambridge University Press, 1998, §1.3 (Area2 / Left, the sign of
  twice the signed triangle area).
- **Definition**: the sign of (b − a) × (c − a): +1 when a, b, c turn counterclockwise (c left of a→b),
  −1 clockwise, 0 collinear (including any repeated point).
- **Best known**: O(1), exact. **Approach**: 65-bit differences, 128-bit products via `bits.Mul64`, signed
  compare. **Lost**: int64 products (overflow once |coordinate| ≥ 2³⁰); `math/big` (allocates).
- **Worked example**: hand: (0,0),(4,0),(0,3) → +1; swapped → −1; (0,0),(1,1),(5,5) → 0. Extremes:
  (MinInt64,MinInt64),(MaxInt64,MinInt64),(MinInt64,MaxInt64) → +1 (cross ≈ 2¹²⁸, exact).
- **Errors**: none. No allocation. **Tags**: orientation, ccw, cross-product, predicate.

### ALG-GEO2D-002 ConvexHull
- **Source**: Andrew, "Another efficient algorithm for convex hulls in two dimensions", Information
  Processing Letters 9(5):216–219, 1979, doi:10.1016/0020-0190(79)90072-3 (verified on Crossref
  2026-11-12) — the monotone chain. CLRS 3rd ed. §33.3 (Graham scan; the Ω(n log n) lower bound is stated
  there, with reference to Yao 1981).
- **Best known**: O(n log n), Ω(n log n) in the algebraic decision-tree model; output-sensitive O(n log h)
  (Kirkpatrick–Seidel 1986, Chan 1996) is better only when h ≪ n. Space O(n).
- **Approach**: clone, sort by (X, Y), drop duplicates, build lower then upper chain popping on
  Orientation ≤ 0. **Collinear policy**: only strict convex vertices; points in the interior of a hull edge
  and duplicates are dropped. Output counterclockwise starting at the least (X, Y) point; all collinear → the
  two extremes; one distinct point → that point; empty → nothing. **Lost**: Graham scan (angular sort is
  awkward exactly; same bound); Jarvis march O(nh); Chan O(n log h) (much more code, wins only for tiny h).
- **Output**: appended to dst (D-003). **Input modified**: no. **Errors**: none.
- **Worked example**: hand: square (0,0),(4,0),(4,4),(0,4) with (2,2) inside and (2,0) on an edge → (0,0),
  (4,0),(4,4),(0,4).
- **Tags**: convex-hull, monotone-chain. LeetCode: 587 (that problem keeps edge points; ours does not).

### ALG-GEO2D-003 ClosestPair
- **Source**: Shamos, Hoey, "Closest-point problems", 16th Annual Symposium on Foundations of Computer Science
  (SFCS 1975), IEEE, pp. 151–162, doi:10.1109/SFCS.1975.8 (verified on Crossref 2026-11-12). CLRS 3rd ed.
  §33.4 "Finding the closest pair of points" (divide and conquer with presorting; at most 7 strip
  successors).
- **Best known**: Θ(n log n) in the algebraic decision-tree model (element uniqueness reduces to it,
  Shamos–Hoey 1975); O(n) expected with randomisation and floor (Rabin 1976) is outside the model and needs
  hashing of real cells. Space O(n).
- **Approach**: sort indices by X, recurse, merge by Y on the way back (merge sort, so O(n log n) not
  O(n log² n)); strip test with exact squared distances in 192 bits. Recursion depth ⌈log₂ n⌉.
  **Lost**: O(n²) all pairs; sweep line with a balanced tree (same bound, needs an ordered set).
- **Result**: indices i < j of one closest pair (Euclidean). Duplicates are a pair at distance 0.
- **Errors**: `ErrTooFewPoints` (n < 2). **Input modified**: no.
- **Worked example**: hand: (0,0),(10,10),(3,4),(11,12),(20,0) → the pair (10,10),(11,12), indices 1, 3.
- **Tags**: closest-pair, divide-and-conquer, nearest. LeetCode: none exact.

### ALG-GEO2D-004 PointInPolygon
- **Source**: Haines, "Point in Polygon Strategies", in Heckbert (ed.), Graphics Gems IV, Academic Press, 1994,
  pp. 24–46 (crossings test, +X ray, half-open vertex rule). O'Rourke, Computational Geometry in C, 2nd ed.,
  1998, §7.4 "Point in Polygon" (§7.4.2 ray crossings, p. 239 ff.; confirmed from the book's table of
  contents 2026-11-12).
- **Best known**: Θ(n) for one query on an unpreprocessed polygon (every edge may matter). O(1) extra.
- **Approach**: one pass; each edge first tested for q on it (exact, → OnBoundary); otherwise an edge whose
  endpoints straddle y = q.Y (one strictly above, the other not) is a crossing when q is strictly on the side
  of it facing −X, decided by Orientation (no division). **Boundary policy**: OnBoundary is reported
  exactly, before parity. Even–odd rule for self-intersecting input. Either vertex order; closed implicitly.
  **Lost**: winding number (§7.4.1; same cost, differs only for self-intersecting polygons; a WindingNumber
  variant may follow with its own ID); float ray test (inexact at vertices).
- **Errors**: `ErrTooFewPoints` (< 3 vertices). No allocation. **Input modified**: no.
- **Worked example**: hand: U shape (0,0),(6,0),(6,6),(4,6),(4,2),(2,2),(2,6),(0,6): (1,4) Inside, (3,4)
  Outside (in the notch), (3,2) OnBoundary, (5,6) OnBoundary, (1,2) Inside (its ray runs along the notch's
  floor through two vertices: the half-open rule counts 3 crossings).
- **Tags**: point-in-polygon, crossing-number, ray-casting, even-odd.

### ALG-GEO2D-005 SegmentsIntersect
- **Source**: CLRS 3rd ed. §33.1, SEGMENTS-INTERSECT and ON-SEGMENT (closed segments; collinear overlap and
  touching endpoints count).
- **Best known**: O(1), exact with Orientation. No allocation. **Lost**: solving the 2×2 system in floats.
- **Degenerate**: a zero-length segment is a point; it intersects when it lies on the other.
- **Worked example**: hand: (0,0)–(4,4) × (0,4)–(4,0) true; (0,0)–(2,0) × (3,0)–(5,0) false; (0,0)–(3,0) ×
  (2,0)–(5,0) true (overlap); (0,0)–(2,2) × (2,2)–(4,0) true (touch).
- **Errors**: none. **Tags**: segment-intersection, predicate.

### ALG-GEO2D-006 MaxPointsOnLine
- **Source**: the problem (maximum number of collinear points) and its O(n²) bound: Gajentaan, Overmars, "On a
  class of O(n²) problems in computational geometry", Computational Geometry 5(3):165–185, 1995,
  doi:10.1016/0925-7721(95)00022-2 (verified on Crossref 2026-11-12): deciding whether any three are collinear
  is 3SUM-hard, so o(n²) is not expected. The per-anchor step is proved here (L15 pattern): a line through
  ≥ 2 of the points has a first point p (least index); every other point q on it gives the same direction
  q − p up to sign, and the direction reduced by gcd of |dx|, |dy| with a sign flag is a canonical key for
  that line through p; so the largest bucket for anchor p, plus copies of p, plus p, is the largest line
  whose first point is p. Gcd by Stein's binary algorithm (Stein, "Computational problems associated with Racah algebra", J. Comput. Phys. 1(3):397–405, 1967,
  doi:10.1016/0021-9991(67)90047-2), private (L4: the exported gcd belongs to num/integer).
- **Best known**: O(n²) (3SUM-hard). This: O(n²) expected (hash map), O(n) extra.
  **Lost**: sort directions per anchor (O(n² log n), deterministic); slope as float (wrong for steep/near
  slopes).
- **Duplicates**: counted separately (a multiset): every copy is on every line through that point.
- **Worked example**: hand: (0,0),(1,1),(2,2),(3,3),(0,1),(1,0),(2,0) → 4. Empty → 0, one → 1.
- **Errors**: none. **Tags**: collinear, lines, gcd. LeetCode: 149.

### ALG-GEO2D-007 RectanglesUnionArea
- Named per L10 (not RectangleArea: NUM-GEOM owns w·h).
- **Source**: inclusion–exclusion |A ∪ B| = |A| + |B| − |A ∩ B| (Stanley, Enumerative Combinatorics Vol. 1,
  2nd ed., Cambridge, 2012, §2.1, for finite sets; for areas it is additivity of Lebesgue measure). The
  intersection of two axis-aligned rectangles is the product of the intersections of their x and y
  intervals (proof: a point is in both iff its x is in both x ranges and its y in both y ranges). Klee's
  measure problem (Klee, Amer. Math. Monthly 84(4):284–285, 1977, doi:10.1080/00029890.1977.11994336) is
  the n-rectangle generalisation (sweep, O(n log n)); not built here.
- **Best known**: O(1). Exact: widths as 64-bit magnitudes, areas in 128 bits, sum in 192.
- **Result**: uint64 area in u²; `ErrOverflow` when the union exceeds 2⁶⁴ − 1. Degenerate (zero-width)
  rectangles are allowed (area 0). **Errors**: `ErrInvertedRect` (Min > Max on an axis), `ErrOverflow`.
- **Worked example**: hand: [−3,3]×[0,4] ∪ [0,9]×[−1,2] = 24 + 27 − 6 = 45.
- **Tags**: rectangle, union-area, inclusion-exclusion. LeetCode: 223.

### ALG-GEO2D-008 HasVerticalSymmetryAxis
- **Source**: Atallah, "On symmetry detection", IEEE Transactions on Computers C-34(7):663–666, 1985,
  doi:10.1109/TC.1985.1676605 (all symmetry axes of a planar point set in O(n log n)). The vertical-axis
  specialisation is proved here (L15 pattern): a reflection in x = c maps the set onto itself and so maps
  its least x to its greatest, hence 2c = minX + maxX; the set is symmetric iff every point's mirror
  (minX + maxX − x, y) is in the set.
- **Best known**: O(n) expected with hashing; Ω(n log n) deterministic in the algebraic decision-tree model
  (it decides set equality). This: O(n) expected time, O(n) extra (one map).
  **Lost**: sort + binary search (O(n log n), deterministic); Atallah's general method (all axes; more work).
- **Points as a set**: duplicates are ignored. Empty → true. Mirror computed without overflow.
- **Worked example**: hand: (1,1),(−1,1) → true (axis x = 0); (1,1),(−1,−1) → false; (0,0),(3,0),(1,5),(2,5)
  → true (axis x = 1.5).
- **Tags**: symmetry, reflection, line-reflection. LeetCode: 356 (the problem's name: LineReflection).

### PerfectRectangle — HELD (no source)
The area-plus-corner-parity test (an exact tiling iff total area equals the bounding box's and every corner
point except the four outer ones occurs an even number of times) was found only on problem-site and
tutorial pages (searched 2026-11-12); no textbook or paper stating and proving it. L22/L31: not built.
Candidate for Devin: a sweep over x with an interval set (a Klee-style measure sweep that also detects
overlap) is sourceable from Klee 1977 / Bentley's sweep if wanted.
