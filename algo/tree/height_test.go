package tree

import (
	"math/rand/v2"
	"testing"
)

// The shared tree of Laaksonen (2018) §14.4 p. 139 has longest path 1 2 5 6: 3 edges (hand-checked).
func TestHeight_laaksonenExample(t *testing.T) {
	if got := Height(sharedTree()); got != 3 {
		t.Errorf("Height(shared tree) = %d, want 3", got)
	}
}

func TestHeight_cases(t *testing.T) {
	tests := []struct {
		name string
		root *Node[int]
		want int
	}{
		{"empty", nil, -1},
		{"one node", &Node[int]{}, 0},
		{"two nodes", chain(2, true), 1},
		{"left chain 1e5", chain(deepN, true), deepN - 1},
		{"right chain 1e5", chain(deepN, false), deepN - 1},
		{"balanced 7", balancedTree(7), 2},
		{"balanced 8", balancedTree(8), 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Height(tt.root); got != tt.want {
				t.Errorf("Height = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestHeight_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		root := randomTree(rng, rng.IntN(80))
		if got, want := Height(root), refHeight(root); got != want {
			t.Fatalf("trial %d: Height = %d, want %d", trial, got, want)
		}
	}
}

func BenchmarkHeight(b *testing.B) {
	for _, size := range benchSizes {
		root := balancedTree(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = Height(root)
			}
		})
	}
}
