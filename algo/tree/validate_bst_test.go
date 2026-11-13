package tree

import (
	"errors"
	"math"
	"testing"
)

// CLRS §12.1: the Figure 12.2 tree satisfies the binary-search-tree property.
func TestValidateBST_clrsTree(t *testing.T) {
	if err := ValidateBST(clrsBST(t).Root()); err != nil {
		t.Errorf("ValidateBST(CLRS tree) = %v, want nil", err)
	}
}

func TestValidateBST_cases(t *testing.T) {
	n := func(v int, l, r *Node[int]) *Node[int] { return &Node[int]{Value: v, Left: l, Right: r} }
	tests := []struct {
		name  string
		root  *Node[int]
		valid bool
	}{
		{"empty", nil, true},
		{"one", n(1, nil, nil), true},
		{"two ordered", n(2, n(1, nil, nil), nil), true},
		{"two reversed", n(1, n(2, nil, nil), nil), false},
		{"duplicate keys", n(2, n(2, nil, nil), nil), false},
		// Every parent-child pair is ordered, but 6 sits in the left subtree of 5.
		{"grandchild on the wrong side", n(5, n(3, nil, n(6, nil, nil)), n(8, nil, nil)), false},
		{"extremes", n(0, n(math.MinInt, nil, nil), n(math.MaxInt, nil, nil)), true},
		{"right chain 1e5", chain(deepN, false), true},
		{"left chain 1e5 (descending)", chain(deepN, true), false},
		{"shared tree (not a BST)", sharedTree(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBST(tt.root)
			if tt.valid && err != nil {
				t.Errorf("ValidateBST = %v, want nil", err)
			}
			if !tt.valid && !errors.Is(err, ErrNotBST) {
				t.Errorf("ValidateBST = %v, want ErrNotBST", err)
			}
		})
	}
}

func TestValidateBST_floats(t *testing.T) {
	root := &Node[float64]{Value: 0, Left: &Node[float64]{Value: math.NaN()}, Right: &Node[float64]{Value: math.Inf(1)}}
	if err := ValidateBST(root); err != nil {
		t.Errorf("NaN < 0 < +Inf under cmp.Compare: got %v", err)
	}
	root.Left.Value = math.Copysign(0, -1) // -0 equals 0: a duplicate key
	if err := ValidateBST(root); !errors.Is(err, ErrNotBST) {
		t.Errorf("-0 left of 0: got %v, want ErrNotBST", err)
	}
}

func BenchmarkValidateBST(b *testing.B) {
	for _, size := range benchSizes {
		root := balancedTree(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = ValidateBST(root)
			}
		})
	}
}
