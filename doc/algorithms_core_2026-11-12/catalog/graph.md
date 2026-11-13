# algo/graph (SHORT = GRAPH)

Builder W2-B8. New package (in `doc/taxonomy.md`: `ALG | algo/graph | GRAPH | traversal, shortest path,
topological order`). Imports `algo/unionfind` (ConnectedComponents, Kruskal, and the Cell type) and
`algo/heap` (Dijkstra, Prim).

| ID | Name | Signature |
|---|---|---|
| ALG-GRAPH-001 | NewDirected | `func NewDirected(n int, edges []Edge) (*Graph, error)`; `type Edge struct{ From, To int; Weight int64 }` |
| ALG-GRAPH-002 | NewUndirected | `func NewUndirected(n int, edges []Edge) (*Graph, error)` |
| ALG-GRAPH-003 | Graph.AddEdge | `func (g *Graph) AddEdge(from, to int, weight int64) error` |
| ALG-GRAPH-004 | Graph.VertexCount | `func (g *Graph) VertexCount() int` |
| ALG-GRAPH-005 | Graph.EdgeCount | `func (g *Graph) EdgeCount() int` |
| ALG-GRAPH-006 | Graph.Directed | `func (g *Graph) Directed() bool` |
| ALG-GRAPH-007 | Graph.Neighbors | `func (g *Graph) Neighbors(v int) iter.Seq2[int, int64]` |
| ALG-GRAPH-008 | Graph.Edges | `func (g *Graph) Edges() []Edge` |
| ALG-GRAPH-009 | BFS | `func BFS(g *Graph, source int) (iter.Seq[int], error)` |
| ALG-GRAPH-010 | DFS | `func DFS(g *Graph, source int) (iter.Seq[int], error)` |
| ALG-GRAPH-011 | ShortestPathUnweighted | `func ShortestPathUnweighted(g *Graph, source, target int) ([]int, error)` |
| ALG-GRAPH-012 | Dijkstra | `func Dijkstra(g *Graph, source int) (dist []int64, parent []int, err error)` |
| ALG-GRAPH-013 | BellmanFord | `func BellmanFord(g *Graph, source int) (dist []int64, parent []int, err error)` |
| ALG-GRAPH-014 | FloydWarshall | `func FloydWarshall(g *Graph) ([][]int64, error)` |
| ALG-GRAPH-015 | TopologicalSort | `func TopologicalSort(g *Graph) ([]int, error)` |
| ALG-GRAPH-016 | ConnectedComponents | `func ConnectedComponents(g *Graph) (component []int, count int, err error)` |
| ALG-GRAPH-017 | IsBipartite | `func IsBipartite(g *Graph) (side []bool, ok bool, err error)` |
| ALG-GRAPH-018 | Kruskal | `func Kruskal(g *Graph) (tree []Edge, totalWeight int64, err error)` |
| ALG-GRAPH-019 | Prim | `func Prim(g *Graph) (tree []Edge, totalWeight int64, err error)` |
| ALG-GRAPH-020 | CloneGraph | `func CloneGraph(g *Graph) *Graph` |
| ALG-GRAPH-021 | WordLadder | `func WordLadder(start, target string, dictionary []string) []string` |
| ALG-GRAPH-022 | GridBFS | `func GridBFS[T any](grid [][]T, start Cell, passable func(T) bool, nb Neighborhood) (iter.Seq[Cell], error)` |
| ALG-GRAPH-023 | GridDFS | `func GridDFS[T any](grid [][]T, start Cell, passable func(T) bool, nb Neighborhood) (iter.Seq[Cell], error)` |
| ALG-GRAPH-024 | NumberOfIslands | `func NumberOfIslands[T any](grid [][]T, isLand func(T) bool, nb Neighborhood) (int, error)` |
| ALG-GRAPH-025 | FloodFillInPlace | `func FloodFillInPlace[T comparable](grid [][]T, start Cell, replacement T, nb Neighborhood) (int, error)` |
| ALG-GRAPH-026 | SurroundedRegionsInPlace | `func SurroundedRegionsInPlace[T comparable](grid [][]T, region, fill T, nb Neighborhood) (int, error)` |

`type Arc` is not exported: Neighbors yields (to, weight). `type Cell = unionfind.Cell` (alias); `type
Neighborhood int` with `FourNeighbors` (up, right, down, left) and `EightNeighbors` (clockwise from up).
`const Unreachable int64 = math.MaxInt64` (the distance of a vertex with no path).

## Sources (DOIs confirmed on Crossref, api.crossref.org, 2026-11-12)
- **CLRS** = Cormen, Leiserson, Rivest, Stein, *Introduction to Algorithms*, 3rd ed., MIT Press, 2009.
- Dijkstra, "A note on two problems in connexion with graphs", Numerische Mathematik 1(1):269–271, 1959,
  doi:10.1007/BF01386390 (problem 2 is single-source shortest paths; problem 1 is the spanning tree).
- Bellman, "On a routing problem", Quarterly of Applied Mathematics 16(1):87–90, 1958, doi:10.1090/qam/102435.
  Ford, "Network flow theory", RAND Corporation paper P-923, 1956 (no DOI; cited through CLRS ch. 24 notes).
- Floyd, "Algorithm 97: Shortest path", Communications of the ACM 5(6):345, 1962, doi:10.1145/367766.368168.
  Warshall, "A theorem on Boolean matrices", Journal of the ACM 9(1):11–12, 1962, doi:10.1145/321105.321107.
- Kahn, "Topological sorting of large networks", Communications of the ACM 5(11):558–562, 1962,
  doi:10.1145/368996.369025.
- Kruskal, "On the shortest spanning subtree of a graph and the traveling salesman problem", Proceedings of
  the American Mathematical Society 7(1):48–50, 1956, doi:10.1090/S0002-9939-1956-0078686-7.
- Prim, "Shortest connection networks and some generalizations", Bell System Technical Journal
  36(6):1389–1401, 1957, doi:10.1002/j.1538-7305.1957.tb01515.x.
- Tarjan, "Depth-first search and linear graph algorithms", SIAM Journal on Computing 1(2):146–160, 1972,
  doi:10.1137/0201010.
- Smith, "Tint fill", Proceedings of SIGGRAPH '79 (6th annual conference on computer graphics and
  interactive techniques), 1979, pp. 276–283, doi:10.1145/800249.807456.
- Knuth, *The Stanford GraphBase: A Platform for Combinatorial Computing*, ACM Press / Addison-Wesley, 1993:
  GB_WORDS (adjacent words differ in one position; built with "five hash tables, one for each of the strings
  of four letters obtained by suppressing a single letter", section "The output phase") and LADDERS (shortest
  word ladders). Checked against the CWEB source gb_words.w on CTAN, 2026-11-12.
- The Go Programming Language Specification, "For statements with range clause" (range over function
  iterators, Go 1.23+), go.dev/ref/spec; package iter (iter.Seq, iter.Seq2).

## Decisions (lead to log)
- **Weights: int64**, exact. Every path sum is checked: a sum that does not fit, or that equals
  `Unreachable`, is `ErrOverflow`. **Lost**: float64 with NaN/Inf rejected (sums lose exactness above 2⁵³,
  results depend on summation order, and a "no path" sentinel of +Inf mixes with real values). Callers with
  fractional weights scale to integers (cents, millimetres).
- **One Graph type, two constructors** (`NewDirected`, `NewUndirected`); unweighted graphs use weight 0
  (BFS-based functions ignore weights). Edges are validated before anything is stored (all or nothing).
  Parallel edges and self-loops are allowed (an undirected self-loop is stored once). The adjacency rows
  share one backing array, each clipped to its degree, so AddEdge never aliases another row. **Lost**: a
  `directed bool` parameter (unreadable at the call site); separate weighted/unweighted types (doubles every
  algorithm's API).
- **Size guards (L41, review H3)**: every allocation sized from an argument is checked against
  `maxAllocBytes` (2⁴⁷ bytes on 64-bit, 2³¹ on 32-bit) and returns `ErrTooLarge` instead of a makeslice
  panic: NewDirected/NewUndirected (n list headers of 24 bytes; 16 bytes per stored arc), FloydWarshall
  (V² int64 entries: V > 2²²), and the grid functions (one visited byte per cell).
- **Overflow (review M1)**: in Dijkstra and BellmanFord a relaxation whose sum overflows upward is "no
  improvement" (skipped), never an error by itself; ErrOverflow is returned only when some reachable vertex's
  shortest distance is ≥ math.MaxInt64 (or, in BellmanFord, a distance falls below math.MinInt64). Test:
  NewDirected(2, {0→1 1, 1→0 MaxInt64}) from 0 → [0 1].
- **Nil graph** → `ErrNilGraph` from every algorithm; methods on a nil *Graph return zero values.
- **Directed vs undirected**: TopologicalSort needs directed (`ErrUndirected`); IsBipartite, Kruskal and Prim
  need undirected (`ErrDirected`); ConnectedComponents on a directed graph returns WEAKLY connected
  components (edge direction ignored), stated in the doc.
- **ConnectedComponents calls unionfind** (one Union per edge, O(V + E α(V))). **Lost**: BFS per component
  (O(V + E) but needs reverse arcs for the weak components of a directed graph; union-find handles both
  kinds with one code path and is the CLRS §21.1 method).
- **Grids**: `[][]T` must be rectangular (`ErrNotRectangular`); rows×cols that overflow int is `ErrTooLarge`.
  **FloodFillInPlace and SurroundedRegionsInPlace modify the grid** (names say so). **Lost**: dst copies
  (O(rows×cols) extra per call on images that are usually edited in place; a caller who wants a copy clones
  first). Both return the number of cells changed.
- **Traversals return iter.Seq** built lazily: each range over the sequence runs a fresh traversal, so a
  sequence may be ranged twice; breaking out stops the work. Argument errors are returned up front, before
  any iteration.
- Dijkstra and Prim use `algo/heap` with lazy deletion (push the new key, skip stale entries on pop):
  O((V + E) log E) = O((V + E) log V) time, O(V + E) heap entries. **Lost**: decrease-key heap (heap has
  none; an indexed heap would give O((V+E) log V) with O(V) heap space — same bound); Fibonacci heap
  (O(E + V log V), Fredman–Tarjan 1987, not cited here, high constants).

### ALG-GRAPH-001…008 Graph (adjacency list)
- **Source**: CLRS §22.1 (adjacency-list representation; Θ(V + E) memory; undirected edges stored in both lists).
- **Best known**: O(V + E) build and memory (reading the input is Ω(V + E)).
- **Worked example**: CLRS §22.1 Figure 22.1(a,b): the undirected graph on 1..5 (here 0..4) with edges 1–2,
  1–5, 2–5, 2–3, 2–4, 3–4, 4–5: degrees 2 4 2 3 3, 7 edges. Test TestNewUndirected_clrsFigure22_1.
- **Errors**: `ErrNegativeVertexCount` (n < 0), `ErrVertexOutOfRange` (an endpoint outside 0..n−1, wrapped
  with the edge index), `ErrTooLarge` (L41 guard on list headers and arcs), `ErrNilGraph` (AddEdge on nil). Input slices not modified.
- **Tests**: Figure 22.1; directed/undirected, self-loop, parallel edges, n 0; Edges round trip; AddEdge
  does not alias other rows; errors. Benchmarks NewDirected/NewUndirected/AddEdge/Neighbors/Edges at
  1e3/1e4/1e5 edges; VertexCount/EdgeCount/Directed single case.
- **Tags**: graph, adjacency-list.

### ALG-GRAPH-009 BFS · 010 DFS
- **Source**: CLRS §22.2 (BFS), §22.3 (DFS); Tarjan 1972 (DFS, linear time). Go spec range-over-func.
- **Best known**: Θ(V + E) time, O(V) extra.
- **Approach**: BFS with a slice queue (head index); DFS iterative with an explicit stack of (vertex, next
  arc) frames — the same preorder as the recursive DFS-VISIT visiting arcs in list order, with no recursion
  depth limit. **Lost**: recursive DFS (stack depth V; a path of 10⁷ vertices would grow the goroutine
  stack to GBs); the "push all neighbours" stack DFS (different order than DFS-VISIT, O(E) stack).
- **Worked example**: CLRS Figure 22.3 (BFS from s on the 8-vertex undirected graph r s t u v w x y):
  distances r1 s0 t2 u3 v2 w1 x2 y3. CLRS Figure 22.4 (DFS on the directed graph u v w x y z, edges
  u→v u→x v→y w→y w→z x→v y→x z→z): DFS-VISIT(u) discovers u v y x (then the book restarts at w: w z);
  our single-source DFS from u yields u v y x, from w yields w y x v z. Tests TestBFS_clrsFigure22_3, TestDFS_clrsFigure22_4.
- **Errors**: `ErrNilGraph`, `ErrVertexOutOfRange`.
- **Tests**: figures; one vertex; disconnected; early break; re-ranging; oracle against a recursive DFS
  and a level-by-level BFS on 300 random graphs. Benchmarks at V = 1e3/1e4/1e5 (E = 4V).
- **Tags**: bfs, dfs, traversal, iterator. LeetCode: 797 (DFS paths, pattern only).

### ALG-GRAPH-011 ShortestPathUnweighted
- **Source**: CLRS §22.2 (BFS computes shortest-path distances, Theorem 22.5; PRINT-PATH reconstruction).
- **Best known**: O(V + E), O(V) extra. Bidirectional BFS has the same worst case (lost: two frontiers, more code).
- **Output**: vertex path source…target; source == target → [source]; unreachable → nil, nil.
- **Worked example**: CLRS Figure 22.3: the shortest s → y path is s w x y (3 edges, δ(s, y) = 3; it is unique).
- **Errors**: `ErrNilGraph`, `ErrVertexOutOfRange`. **Tags**: bfs, shortest-path, unweighted. LeetCode: 1091.

### ALG-GRAPH-012 Dijkstra
- **Source**: Dijkstra 1959 (problem 2); CLRS §24.3.
- **Best known**: O(E + V log V) with a Fibonacci heap; binary heap O((V + E) log V). Lower bound Ω(V log V)
  for comparison-based (it sorts). Chosen: binary heap with lazy deletion, O((V + E) log E), O(V + E) extra.
- **Worked example**: CLRS Figure 24.6 (s t x y z; edges s→t10 s→y5 t→x1 t→y2 y→t3 y→x9 y→z2 x→z4 z→x6
  z→s7): δ = s0 t8 x9 y5 z7. Test TestDijkstra_clrsFigure24_6.
- **Errors**: `ErrNilGraph`, `ErrVertexOutOfRange`, `ErrNegativeWeight` (any arc in the graph), `ErrOverflow`
  (a reachable vertex whose shortest distance is ≥ math.MaxInt64; an overflowing non-shortest path is skipped).
- **Tests**: figure; unreachable; zero weights; parallel edges; oracle vs BellmanFord on 300 random graphs.
  Benchmark V = 1e3/1e4/1e5, E = 4V.
- **Tags**: dijkstra, shortest-path, weighted, heap. LeetCode: 743, 1514.

### ALG-GRAPH-013 BellmanFord
- **Source**: Bellman 1958; Ford 1956; CLRS §24.1 (negative-cycle check after V−1 passes).
- **Best known**: O(V·E) for general weights; chosen O(V·(V + E)) with early exit when a pass relaxes
  nothing. **Lost**: SPFA queue variant (same worst case, no cleaner source).
- **Worked example**: CLRS Figure 24.4 (s t x y z; s→t6 s→y7 t→x5 t→y8 t→z−4 x→t−2 y→x−3 y→z9 z→s2 z→x7):
  δ = s0 t2 x4 y7 z−2. Test TestBellmanFord_clrsFigure24_4.
- **Errors**: `ErrNilGraph`, `ErrVertexOutOfRange`, `ErrNegativeCycle` (reachable from source), `ErrOverflow`
  (a distance below math.MinInt64 — can precede detection of a negative cycle — or a reachable vertex whose
  shortest distance is ≥ math.MaxInt64; an upward-overflowing non-shortest path is skipped).
- **Tags**: bellman-ford, shortest-path, negative-weights, negative-cycle. LeetCode: 787 (pattern only).

### ALG-GRAPH-014 FloydWarshall
- **Source**: Floyd 1962; Warshall 1962; CLRS §25.2.
- **Best known**: Θ(V³) for dense all-pairs (Johnson O(V E + V² log V) on sparse — lost: needs Bellman-Ford
  + V Dijkstras, more memory traffic; a separate function later if needed). Space Θ(V²) result, O(1) extra.
- **Worked example**: CLRS Figure 25.4 (1..5: 1→2 3, 1→3 8, 1→5 −4, 2→4 1, 2→5 7, 3→2 4, 4→1 2, 4→3 −5,
  5→4 6): D⁽⁵⁾ rows 0 1 −3 2 −4 / 3 0 −4 1 −1 / 7 4 0 5 3 / 2 −1 −5 0 −2 / 8 5 1 6 0. Test TestFloydWarshall_clrsFigure25_4.
- **Errors**: `ErrNilGraph`, `ErrNegativeCycle` (any d[i][i] < 0, detected the moment it appears), `ErrOverflow`.
  `ErrTooLarge` when the V² table exceeds 2⁴⁷ bytes (V > 2²²; L41), tested through distanceTableFits and an
  end-to-end V = 2²²+1 case (skipped under -short). V = 10⁴ allocates 800 MB.
- **Tags**: floyd-warshall, all-pairs, shortest-path, dense. LeetCode: 1334.

### ALG-GRAPH-015 TopologicalSort
- **Source**: Kahn 1962; CLRS §22.4 (the DFS method — lost: Kahn detects a cycle without a 3-colour DFS
  and its output order is the BFS-like "sources first" order).
- **Best known**: Θ(V + E), O(V) extra (the result slice doubles as the queue).
- **Order**: initial sources in increasing vertex order, then FIFO; deterministic.
- **Worked example**: CLRS Figure 22.7 (Professor Bumstead's clothes: undershorts→pants, undershorts→shoes,
  pants→belt, pants→shoes, belt→jacket, shirt→belt, shirt→tie, tie→jacket, socks→shoes; watch isolated):
  any output must respect every edge; test checks validity and Kahn's exact order.
- **CourseSchedule** (prerequisite pairs → can all be finished) is a TEST of TestTopologicalSort_courseSchedule,
  not a function.
- **Errors**: `ErrNilGraph`, `ErrUndirected`, `ErrCycle` (including a self-loop).
- **Tags**: topological-sort, kahn, dag, cycle-detection, scheduling. LeetCode: 207, 210.

### ALG-GRAPH-016 ConnectedComponents
- **Source**: CLRS §21.1 CONNECTED-COMPONENTS (disjoint sets); calls algo/unionfind.
- **Best known**: Θ(V + E) (BFS); chosen O(V + E α(V)), O(V) extra. Labels 0..count−1 in order of each
  component's smallest vertex.
- **Worked example**: CLRS Figure 21.1 (a..j, 4 components). **Errors**: `ErrNilGraph`.
- **Tags**: connected-components, union-find. LeetCode: 323, 547.

### ALG-GRAPH-017 IsBipartite
- **Source**: CLRS Exercise 22.2-7 (split wrestlers into two kinds with every rivalry across: BFS 2-colouring,
  O(V + E)); Sedgewick & Wayne, *Algorithms*, 4th ed., 2011, §4.1 (two-colourability).
- **Best known**: Θ(V + E), O(V) extra. side[v] false/true for the two parts (each component's smallest
  vertex is false); side nil when not bipartite.
- **Worked example**: hand: a 4-cycle is bipartite {0,2}/{1,3}; a triangle is not; a self-loop is not.
- **Errors**: `ErrNilGraph`, `ErrDirected`. **Tags**: bipartite, two-coloring, bfs. LeetCode: 785, 886.

### ALG-GRAPH-018 Kruskal · 019 Prim
- **Source**: Kruskal 1956; Prim 1957; CLRS §23.2 (MST-KRUSKAL with disjoint sets; MST-PRIM).
- **Best known**: O(E log V) for both with binary heaps / sorting; O(E α(V)) after sorting. Lower bound
  Ω(E). (Randomised linear-time Karger–Klein–Tarjan, lost: complex, not needed.)
- **Approach**: Kruskal sorts a copy of the edges stably by weight and unions with algo/unionfind, stopping
  at V − 1 tree edges. Prim restarts from each unvisited vertex, lazy heap. Both return a minimum spanning
  FOREST when the graph is disconnected; negative weights are allowed.
- **Worked example**: CLRS Figure 23.1/23.4 (a..i; a–b4 a–h8 b–c8 b–h11 c–d7 c–f4 c–i2 d–e9 d–f14 e–f10
  f–g2 g–h1 g–i6 h–i7): total weight 37, 8 edges. Tests TestKruskal_clrsFigure23_4, TestPrim_clrsFigure23_5.
- **Errors**: `ErrNilGraph`, `ErrDirected`, `ErrOverflow` (total weight).
- **Tests**: oracle Kruskal total == Prim total == brute force over all edge subsets on small random graphs.
- **Tags**: mst, kruskal, prim, union-find, greedy. LeetCode: 1584.

### ALG-GRAPH-020 CloneGraph
- **Source**: CLRS §22.1 (adjacency lists; copying each list). Deep copy: the clone shares no memory with g.
- **Best known**: Θ(V + E). Two allocations plus the struct. nil → nil.
- **Tags**: graph, clone, deep-copy. LeetCode: 133.

### ALG-GRAPH-021 WordLadder
- **Source**: Knuth, *The Stanford GraphBase*, 1993 (GB_WORDS adjacency by "suppressing a single letter"
  hash tables; LADDERS shortest ladders); BFS on the implicit graph, CLRS §22.2.
- **Best known**: O(N·L²) time and space with the masked-key buckets (N words of length L: N·L keys of L−1
  bytes); each bucket is scanned once then dropped. **Lost**: trying all 26 (or 256) letters per position
  (O(N·L·Σ·L), alphabet-dependent); bidirectional BFS (same worst case; optional in the brief, not built).
- **Output**: one shortest ladder start…target (deterministic for the same input); nil when none; [start]
  when start == target. target must be in dictionary; start need not be. Letters are bytes (multi-byte UTF-8
  letters are not one-letter changes). Words of another length are ignored.
- **Worked example**: Knuth's demo pair is not reproduced; hand example: cold → warm differ in all 4
  letters, so ≥ 5 words; with dictionary cord card ward warm word worm wore core (in that order) BFS returns
  cold cord word ward warm.
- **Tags**: word-ladder, bfs, implicit-graph, strings. LeetCode: 127, 126.

### ALG-GRAPH-022 GridBFS · 023 GridDFS · 024 NumberOfIslands · 025 FloodFillInPlace · 026 SurroundedRegionsInPlace
- **Source**: BFS/DFS on the grid graph (vertices = cells, edges = 4- or 8-neighbours), CLRS §22.2, §22.3;
  Smith 1979 (flood "tint fill" of a 4- or 8-connected region of one value). NumberOfIslands = connected
  components of the passable subgraph (CLRS §22.3 / §21.1). SurroundedRegions = cells of `region` whose
  component touches no border cell: BFS from all border cells of that value at once (a BFS from a virtual
  source joined to every border cell, CLRS §22.2), then fill the rest.
- **Best known**: Θ(R·C) time, O(R·C) extra (the visited bitmap / queue). FloodFillInPlace needs no visited
  array: changed cells no longer match.
- **Worked examples** (hand): a 4×5 grid with 4 islands under 4-neighbours that join into 1 under
  8-neighbours; GridBFS/GridDFS orders on a 3×3 maze (4 and 8 neighbours); a 3×3 flood fill (5 cells from
  the centre under 4, 3 cells from a corner under 8); a 4×4 O/X board where the inner O block (3 cells) is
  filled under 4-neighbours and nothing under 8 (a diagonal reaches the border O).
- **Errors**: `ErrNotRectangular`, `ErrTooLarge` (rows×cols visited bytes above 2⁴⁷, L41), `ErrCellOutOfRange`,
  `ErrNilPredicate`, `ErrNeighborhood`. Benchmarks GridBFS/GridDFS/FloodFillInPlace use all-open grids so the
  work grows with n (review B1).
- **Tests**: hand cases, 1×1, empty grid, all water, all land, 8 vs 4, oracle vs unionfind.IslandsOnline's
  last count for NumberOfIslands (4-neighbours) on 300 random grids. Benchmarks on square grids of
  1e4/1e5/1e6 cells (the n of the 1e3/1e4/1e5 rule, ×10 because one cell is cheap).
- **Tags**: grid, bfs, dfs, flood-fill, islands, connected-components. LeetCode: 200, 733, 130, 1020.
