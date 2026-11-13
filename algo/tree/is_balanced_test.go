package tree

import (
	"math/rand/v2"
	"testing"
)

func refBalanced(n *Node[int]) bool {
	if n == nil {
		return true
	}
	d := refHeight(n.Left) - refHeight(n.Right)
	return d >= -1 && d <= 1 && refBalanced(n.Left) && refBalanced(n.Right)
}

// Knuth Vol. 3 §6.2.3 condition checked by hand on the shared tree of Laaksonen §14.4: subtree heights at
// 1: 2 vs 1, at 2: 0 vs 1, at 5: 0 vs −1, at 3: −1 vs 0, so it is balanced.
func TestIsBalanced_sharedTree(t *testing.T) {
	if !IsBalanced(sharedTree()) {
		t.Error("IsBalanced(shared tree) = false, want true")
	}
}

func TestIsBalanced_cases(t *testing.T) {
	unbalancedDeep := &Node[int]{Value: -1, Left: balancedTree(4), Right: chain(1, true)} // 2 vs 0 at the root only
	tests := []struct {
		name string
		root *Node[int]
		want bool
	}{
		{"empty", nil, true},
		{"one node", &Node[int]{}, true},
		{"two nodes", chain(2, false), true},
		{"chain of 3", chain(3, true), false},
		{"chain of 1e5", chain(deepN, false), false},
		{"root unbalanced, subtrees fine", unbalancedDeep, false},
		{"perfect 15", balancedTree(15), true},
		{"perfectly balanced 1e5", balancedTree(deepN), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsBalanced(tt.root); got != tt.want {
				t.Errorf("IsBalanced = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsBalanced_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		root := randomTree(rng, rng.IntN(12))
		if got, want := IsBalanced(root), refBalanced(root); got != want {
			t.Fatalf("trial %d: IsBalanced = %v, want %v", trial, got, want)
		}
	}
}

func BenchmarkIsBalanced(b *testing.B) {
	for _, size := range benchSizes {
		root := balancedTree(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = IsBalanced(root)
			}
		})
	}
}
