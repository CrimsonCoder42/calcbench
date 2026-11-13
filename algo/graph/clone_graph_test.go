package graph

import (
	"math/rand/v2"
	"slices"
	"testing"
)

func TestCloneGraph_cases(t *testing.T) {
	if CloneGraph(nil) != nil {
		t.Error("CloneGraph(nil) != nil")
	}
	empty, _ := NewDirected(0, nil)
	if c := CloneGraph(empty); c == nil || c.VertexCount() != 0 || !c.Directed() {
		t.Errorf("empty clone = %v", c)
	}
	g, _ := NewUndirected(3, []Edge{{0, 1, 5}, {1, 2, -3}, {2, 2, 0}})
	c := CloneGraph(g)
	if c.Directed() || c.VertexCount() != 3 || c.EdgeCount() != 3 || !slices.Equal(c.Edges(), g.Edges()) {
		t.Fatalf("clone differs: %v vs %v", c.Edges(), g.Edges())
	}
	// Independence both ways, including rows that still have spare capacity.
	_ = c.AddEdge(0, 2, 9)
	if g.EdgeCount() != 3 || len(neighbors(g, 0)) != 1 || len(neighbors(g, 2)) != 2 {
		t.Error("AddEdge on the clone changed the original")
	}
	_ = g.AddEdge(1, 1, 4)
	if c.EdgeCount() != 4 || len(neighbors(c, 1)) != 2 {
		t.Error("AddEdge on the original changed the clone")
	}
}

// Oracle: every list of the clone equals the original's, on random graphs.
func TestCloneGraph_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for run := 0; run < 300; run++ {
		n := 1 + rng.IntN(15)
		g, _ := NewDirected(n, randomEdges(rng, n, rng.IntN(4*n), -50, 50))
		if run%2 == 1 {
			g, _ = NewUndirected(n, randomEdges(rng, n, rng.IntN(4*n), -50, 50))
		}
		c := CloneGraph(g)
		for v := range n {
			if !slices.Equal(neighbors(c, v), neighbors(g, v)) {
				t.Fatalf("run %d vertex %d: %v vs %v", run, v, neighbors(c, v), neighbors(g, v))
			}
		}
	}
}

func BenchmarkCloneGraph(b *testing.B) {
	for _, size := range benchSizes {
		g := benchGraph(size.n, false)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = CloneGraph(g)
			}
		})
	}
}
