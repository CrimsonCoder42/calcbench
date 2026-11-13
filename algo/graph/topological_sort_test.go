package graph

import (
	"errors"
	"math/rand/v2"
	"slices"
	"testing"
)

// isTopological reports whether order lists every vertex once with each edge's tail before its head.
func isTopological(g *Graph, order []int) bool {
	n := g.VertexCount()
	if len(order) != n {
		return false
	}
	pos := make([]int, n)
	for i := range pos {
		pos[i] = -1
	}
	for i, v := range order {
		if v < 0 || v >= n || pos[v] != -1 {
			return false
		}
		pos[v] = i
	}
	for _, e := range g.Edges() {
		if pos[e.From] >= pos[e.To] {
			return false
		}
	}
	return true
}

// Cormen et al., Introduction to Algorithms, 3rd ed., §22.4, Figure 22.7 (getting dressed): undershorts→
// pants, undershorts→shoes, pants→belt, pants→shoes, belt→jacket, shirt→belt, shirt→tie, tie→jacket,
// socks→shoes; the watch has no edges. Kahn's order with sources in increasing number is derived by hand.
func TestTopologicalSort_clrsFigure22_7(t *testing.T) {
	const undershorts, pants, belt, shirt, tie, jacket, socks, shoes, watch = 0, 1, 2, 3, 4, 5, 6, 7, 8
	g, _ := NewDirected(9, []Edge{
		{undershorts, pants, 0}, {undershorts, shoes, 0}, {pants, belt, 0}, {pants, shoes, 0},
		{belt, jacket, 0}, {shirt, belt, 0}, {shirt, tie, 0}, {tie, jacket, 0}, {socks, shoes, 0},
	})
	got, err := TopologicalSort(g)
	if err != nil || !isTopological(g, got) {
		t.Fatalf("TopologicalSort = (%v, %v), not a topological order", got, err)
	}
	if want := []int{undershorts, shirt, socks, watch, pants, tie, belt, shoes, jacket}; !slices.Equal(got, want) {
		t.Errorf("order = %v, want %v", got, want)
	}
}

// Course scheduling as a test of Kahn's method: courses 0..n-1, each pair (course, needs) means needs
// must be taken first (edge needs → course). All courses can be finished exactly when there is no cycle.
func TestTopologicalSort_courseSchedule(t *testing.T) {
	tests := []struct {
		name     string
		courses  int
		needs    [][2]int // {course, prerequisite}
		feasible bool
	}{
		{"no prerequisites", 3, nil, true},
		{"a chain of three", 3, [][2]int{{1, 0}, {2, 1}}, true},
		{"two courses that need each other", 2, [][2]int{{0, 1}, {1, 0}}, false},
		{"a diamond", 4, [][2]int{{1, 0}, {2, 0}, {3, 1}, {3, 2}}, true},
		{"a course that needs itself", 2, [][2]int{{1, 1}}, false},
		{"a cycle off to the side", 5, [][2]int{{1, 0}, {3, 2}, {4, 3}, {2, 4}}, false},
	}
	for _, tt := range tests {
		edges := make([]Edge, len(tt.needs))
		for i, p := range tt.needs {
			edges[i] = Edge{From: p[1], To: p[0]}
		}
		g, _ := NewDirected(tt.courses, edges)
		order, err := TopologicalSort(g)
		if tt.feasible && (err != nil || !isTopological(g, order)) {
			t.Errorf("%s: (%v, %v), want a valid order", tt.name, order, err)
		}
		if !tt.feasible && (!errors.Is(err, ErrCycle) || order != nil) {
			t.Errorf("%s: (%v, %v), want (nil, ErrCycle)", tt.name, order, err)
		}
	}
}

func TestTopologicalSort_cases(t *testing.T) {
	empty, _ := NewDirected(0, nil)
	if got, err := TopologicalSort(empty); err != nil || len(got) != 0 {
		t.Errorf("empty: (%v, %v)", got, err)
	}
	par, _ := NewDirected(2, []Edge{{1, 0, 0}, {1, 0, 0}})
	if got, err := TopologicalSort(par); err != nil || !slices.Equal(got, []int{1, 0}) {
		t.Errorf("parallel edges: (%v, %v), want [1 0]", got, err)
	}
}

func TestTopologicalSort_errors(t *testing.T) {
	if _, err := TopologicalSort(nil); !errors.Is(err, ErrNilGraph) {
		t.Errorf("nil: %v", err)
	}
	und, _ := NewUndirected(2, []Edge{{0, 1, 0}})
	if _, err := TopologicalSort(und); !errors.Is(err, ErrUndirected) {
		t.Errorf("undirected: %v", err)
	}
}

// Oracle: random DAGs (edges low → high under a random relabelling) always sort validly; adding one
// back edge along a path makes a cycle, detected.
func TestTopologicalSort_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for run := 0; run < 300; run++ {
		n := 2 + rng.IntN(15)
		label := rng.Perm(n)
		var edges []Edge
		for range rng.IntN(3 * n) {
			a, b := rng.IntN(n), rng.IntN(n)
			if a != b {
				edges = append(edges, Edge{From: label[min(a, b)], To: label[max(a, b)]})
			}
		}
		g, _ := NewDirected(n, edges)
		order, err := TopologicalSort(g)
		if err != nil || !isTopological(g, order) {
			t.Fatalf("run %d: DAG gave (%v, %v)", run, order, err)
		}
		// Close a cycle: a chain label[0] → label[1] → … plus label[n-1] → label[0].
		for i := 0; i+1 < n; i++ {
			edges = append(edges, Edge{From: label[i], To: label[i+1]})
		}
		edges = append(edges, Edge{From: label[n-1], To: label[0]})
		cyc, _ := NewDirected(n, edges)
		if _, err := TopologicalSort(cyc); !errors.Is(err, ErrCycle) {
			t.Fatalf("run %d: cycle not detected (%v)", run, err)
		}
	}
}

func BenchmarkTopologicalSort(b *testing.B) {
	for _, size := range benchSizes {
		rng := rand.New(rand.NewPCG(1, 2))
		edges := make([]Edge, 0, 4*size.n)
		for range 4 * size.n {
			a, c := rng.IntN(size.n), rng.IntN(size.n)
			if a != c {
				edges = append(edges, Edge{From: min(a, c), To: max(a, c)})
			}
		}
		g, _ := NewDirected(size.n, edges)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = TopologicalSort(g)
			}
		})
	}
}
