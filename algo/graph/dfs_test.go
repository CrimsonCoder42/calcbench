package graph

import (
	"errors"
	"math/rand/v2"
	"slices"
	"testing"
)

// recursiveDFS is the oracle: CLRS DFS-VISIT, recursive, appending each vertex when discovered.
func recursiveDFS(g *Graph, v int, seen []bool, out []int) []int {
	seen[v] = true
	out = append(out, v)
	for w := range g.Neighbors(v) {
		if !seen[w] {
			out = recursiveDFS(g, w, seen, out)
		}
	}
	return out
}

// Cormen et al., Introduction to Algorithms, 3rd ed., §22.3, Figure 22.4: the directed graph u v w x y z
// with edges u→v u→x v→y w→y w→z x→v y→x z→z. DFS-VISIT(u) discovers u v y x (times 1 2 3 4); from w,
// our single-source DFS discovers w y x v z.
func TestDFS_clrsFigure22_4(t *testing.T) {
	const u, v, w, x, y, z = 0, 1, 2, 3, 4, 5
	g, _ := NewDirected(6, []Edge{{u, v, 0}, {u, x, 0}, {v, y, 0}, {w, y, 0}, {w, z, 0}, {x, v, 0}, {y, x, 0}, {z, z, 0}})
	if got := slices.Collect(must(DFS(g, u))); !slices.Equal(got, []int{u, v, y, x}) {
		t.Errorf("DFS(u) = %v, want [u v y x]", got)
	}
	if got := slices.Collect(must(DFS(g, w))); !slices.Equal(got, []int{w, y, x, v, z}) {
		t.Errorf("DFS(w) = %v, want [w y x v z]", got)
	}
}

func TestDFS_cases(t *testing.T) {
	one, _ := NewUndirected(1, []Edge{{0, 0, 0}})
	if got := slices.Collect(must(DFS(one, 0))); !slices.Equal(got, []int{0}) {
		t.Errorf("one vertex with a self-loop: %v", got)
	}
	// A long path: an iterative DFS must not depend on call-stack depth.
	const n = 200_000
	path := make([]Edge, n-1)
	for i := range path {
		path[i] = Edge{From: i, To: i + 1}
	}
	g, _ := NewDirected(n, path)
	count, last := 0, -1
	for v := range must(DFS(g, 0)) {
		count++
		last = v
	}
	if count != n || last != n-1 {
		t.Errorf("path of %d: visited %d, last %d", n, count, last)
	}
	seq := must(DFS(g, 0))
	for v := range seq {
		if v == 5 {
			break
		}
	}
	if again := slices.Collect(seq); len(again) != n {
		t.Errorf("re-ranging after break visited %d, want %d", len(again), n)
	}
}

func TestDFS_errors(t *testing.T) {
	if _, err := DFS(nil, 0); !errors.Is(err, ErrNilGraph) {
		t.Errorf("nil graph err = %v", err)
	}
	g, _ := NewDirected(2, nil)
	for _, s := range []int{-1, 2} {
		if seq, err := DFS(g, s); !errors.Is(err, ErrVertexOutOfRange) || seq != nil {
			t.Errorf("source %d: err = %v", s, err)
		}
	}
}

func TestDFS_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for run := 0; run < 300; run++ {
		n := 1 + rng.IntN(15)
		edges := randomEdges(rng, n, rng.IntN(3*n), 0, 0)
		g, _ := NewDirected(n, edges)
		if run%2 == 1 {
			g, _ = NewUndirected(n, edges)
		}
		s := rng.IntN(n)
		got := slices.Collect(must(DFS(g, s)))
		want := recursiveDFS(g, s, make([]bool, n), nil)
		if !slices.Equal(got, want) {
			t.Fatalf("run %d: DFS = %v, recursive = %v", run, got, want)
		}
	}
}

func BenchmarkDFS(b *testing.B) {
	for _, size := range benchSizes {
		g := benchGraph(size.n, true)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				seq, _ := DFS(g, 0)
				for range seq {
				}
			}
		})
	}
}
