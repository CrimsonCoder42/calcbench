package graph

import (
	"errors"
	"math"
	"testing"
)

// CLRS Figure 23.5: MST-PRIM on Figure 23.1 from a gives 8 edges of total weight 37; the first edge
// added is a–b (4).
func TestPrim_clrsFigure23_5(t *testing.T) {
	g := figure23_1(t)
	tree, total, err := Prim(g)
	if err != nil || total != 37 || len(tree) != 8 {
		t.Fatalf("Prim = (%v, %d, %v), want 8 edges weighing 37", tree, total, err)
	}
	checkSpanningForest(t, g, tree, total)
	if tree[0] != (Edge{From: 0, To: 1, Weight: 4}) {
		t.Errorf("first edge %v, want a–b 4", tree[0])
	}
}

func TestPrim_cases(t *testing.T) {
	tests := []struct {
		name  string
		n     int
		edges []Edge
		total int64
	}{
		{"empty", 0, nil, 0},
		{"one vertex, self-loop", 1, []Edge{{0, 0, -5}}, 0},
		{"parallel edges take the lighter", 2, []Edge{{0, 1, 9}, {1, 0, 3}}, 3},
		{"negative weights", 3, []Edge{{0, 1, -4}, {1, 2, -1}, {0, 2, -2}}, -6},
		{"forest: restarts at the next unreached vertex", 5, []Edge{{0, 1, 1}, {1, 2, 2}, {0, 2, 3}, {3, 4, 4}}, 7},
		{"isolated vertices", 4, []Edge{{1, 3, 5}}, 5},
	}
	for _, tt := range tests {
		g, _ := NewUndirected(tt.n, tt.edges)
		tree, total, err := Prim(g)
		if err != nil || total != tt.total {
			t.Errorf("%s: (%v, %d, %v), want total %d", tt.name, tree, total, err, tt.total)
			continue
		}
		checkSpanningForest(t, g, tree, total)
	}
}

func TestPrim_errors(t *testing.T) {
	if _, _, err := Prim(nil); !errors.Is(err, ErrNilGraph) {
		t.Errorf("nil: %v", err)
	}
	d, _ := NewDirected(2, []Edge{{0, 1, 1}})
	if _, _, err := Prim(d); !errors.Is(err, ErrDirected) {
		t.Errorf("directed: %v", err)
	}
	over, _ := NewUndirected(3, []Edge{{0, 1, math.MaxInt64}, {1, 2, math.MaxInt64}})
	if _, _, err := Prim(over); !errors.Is(err, ErrOverflow) {
		t.Errorf("overflow: %v", err)
	}
}

func BenchmarkPrim(b *testing.B) {
	for _, size := range benchSizes {
		g := benchGraph(size.n, false)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _ = Prim(g)
			}
		})
	}
}
