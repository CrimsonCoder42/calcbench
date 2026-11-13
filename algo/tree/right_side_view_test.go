package tree

import (
	"math/rand/v2"
	"slices"
	"testing"
)

// The shared tree of Laaksonen §14.4 p. 139, levels {1} {2 3} {4 5 7} {6}: the last of each is 1 3 7 6.
func TestRightSideView_sharedTree(t *testing.T) {
	if got := RightSideView(nil, sharedTree()); !slices.Equal(got, []int{1, 3, 7, 6}) {
		t.Errorf("RightSideView = %v, want [1 3 7 6]", got)
	}
}

func TestRightSideView_cases(t *testing.T) {
	if got := RightSideView[int](nil, nil); len(got) != 0 {
		t.Errorf("empty: %v", got)
	}
	if got := RightSideView([]int{9}, &Node[int]{Value: 1, Left: &Node[int]{Value: 2}}); !slices.Equal(got, []int{9, 1, 2}) {
		t.Errorf("appends to dst, left-only child seen: %v", got)
	}
	if got := RightSideView(nil, chain(deepN, true)); len(got) != deepN || got[deepN-1] != deepN-1 {
		t.Errorf("chain 1e5: len %d", len(got))
	}
}

func TestRightSideView_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		root := randomTree(rng, rng.IntN(60))
		var want []int
		for _, level := range refLevels(root, 0, nil) {
			want = append(want, level[len(level)-1])
		}
		if got := RightSideView(nil, root); !slices.Equal(got, want) {
			t.Fatalf("trial %d: RightSideView = %v, want %v", trial, got, want)
		}
	}
}

func BenchmarkRightSideView(b *testing.B) {
	for _, size := range benchSizes {
		root := balancedTree(size.n)
		dst := make([]int, 0, 64)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				dst = RightSideView(dst[:0], root)
			}
		})
	}
}
