package tree

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

func nd[T any](v T, l, r *Node[T]) *Node[T] { return &Node[T]{Value: v, Left: l, Right: r} }

// refMaxPathSum tries every pair of endpoints: the path between u and v is unique in a tree.
func refMaxPathSum(root *Node[int]) int {
	nodes := collectNodes(root)
	best := math.MinInt
	for _, u := range nodes {
		for _, v := range nodes {
			pu, pv := refPath(root, u, nil), refPath(root, v, nil)
			i := 0
			for i < len(pu) && i < len(pv) && pu[i] == pv[i] {
				i++
			}
			s := pu[i-1].Value // the highest node of the path
			for _, x := range pu[i:] {
				s += x.Value
			}
			for _, x := range pv[i:] {
				s += x.Value
			}
			best = max(best, s)
		}
	}
	return best
}

// Laaksonen §14.2 Algorithm 1 (best path through each highest node), weighted, worked by hand:
// −10 → (9, 20 → (15, 7)); the best path is 15 → 20 → 7 = 42, not through the root.
func TestMaxPathSum_handExample(t *testing.T) {
	root := nd(-10, nd(9, nil, nil), nd(20, nd(15, nil, nil), nd(7, nil, nil)))
	if got, err := MaxPathSum(root); err != nil || got != 42 {
		t.Errorf("MaxPathSum = (%d, %v), want 42", got, err)
	}
}

func TestMaxPathSum_cases(t *testing.T) {
	tests := []struct {
		name string
		root *Node[int]
		want int
	}{
		{"one node", nd(5, nil, nil), 5},
		{"one negative node", nd(-3, nil, nil), -3},
		{"all negative: the largest value", nd(-3, nd(-1, nil, nil), nd(-2, nil, nil)), -1},
		{"two nodes", nd(1, nd(2, nil, nil), nil), 3},
		{"through the root", nd(1, nd(2, nil, nil), nd(3, nil, nil)), 6},
		{"zeros", nd(0, nd(0, nil, nil), nil), 0},
		{"skip a negative child", nd(2, nd(-1, nil, nil), nd(4, nd(-5, nil, nil), nil)), 6},
		{"extremes cancel", nd(math.MaxInt, nd(math.MinInt, nil, nil), nil), math.MaxInt},
		{"shared tree", sharedTree(), 6 + 5 + 2 + 1 + 3 + 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := MaxPathSum(tt.root); err != nil || got != tt.want {
				t.Errorf("MaxPathSum = (%d, %v), want %d", got, err, tt.want)
			}
		})
	}
}

func TestMaxPathSum_errors(t *testing.T) {
	if _, err := MaxPathSum[int](nil); !errors.Is(err, ErrEmptyTree) {
		t.Errorf("nil: %v, want ErrEmptyTree", err)
	}
	if _, err := MaxPathSum(nd(math.MaxInt, nd(1, nil, nil), nil)); !errors.Is(err, ErrOverflow) {
		t.Errorf("MaxInt + 1: %v, want ErrOverflow", err)
	}
	if _, err := MaxPathSum(nd[int8](100, nd[int8](20, nil, nil), nd[int8](10, nil, nil))); !errors.Is(err, ErrOverflow) {
		t.Errorf("int8 130: %v, want ErrOverflow", err)
	}
	for _, bad := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := MaxPathSum(nd(1.0, nd(bad, nil, nil), nil)); !errors.Is(err, ErrNotFinite) {
			t.Errorf("value %v: %v, want ErrNotFinite", bad, err)
		}
	}
	if _, err := MaxPathSum(nd(math.MaxFloat64, nd(math.MaxFloat64, nil, nil), nil)); !errors.Is(err, ErrOverflow) {
		t.Errorf("float sum to +Inf: %v, want ErrOverflow", err)
	}
}

func TestMaxPathSum_typesAndDepth(t *testing.T) {
	if got, err := MaxPathSum(nd[uint8](200, nd[uint8](55, nil, nil), nil)); err != nil || got != 255 {
		t.Errorf("uint8 = (%d, %v), want 255", got, err)
	}
	// 1.5 + 2.25 is exact in binary floating point.
	if got, err := MaxPathSum(nd(1.5, nd(2.25, nil, nil), nil)); err != nil || got != 3.75 {
		t.Errorf("float = (%v, %v), want 3.75", got, err)
	}
	if got, err := MaxPathSum(chain(deepN, true)); err != nil || got != deepN*(deepN-1)/2 {
		t.Errorf("chain 1e5 = (%d, %v), want %d", got, err, deepN*(deepN-1)/2)
	}
}

func TestMaxPathSum_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		root := randomTree(rng, 1+rng.IntN(25))
		for _, n := range collectNodes(root) {
			n.Value = rng.IntN(41) - 20
		}
		if got, err := MaxPathSum(root); err != nil || got != refMaxPathSum(root) {
			t.Fatalf("trial %d: MaxPathSum = (%d, %v), want %d", trial, got, err, refMaxPathSum(root))
		}
	}
}

func BenchmarkMaxPathSum(b *testing.B) {
	for _, size := range benchSizes {
		root := balancedTree(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = MaxPathSum(root)
			}
		})
	}
}
