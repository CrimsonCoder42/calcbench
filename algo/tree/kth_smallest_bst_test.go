package tree

import (
	"errors"
	"math/rand/v2"
	"slices"
	"testing"
)

// CLRS §12.1 Theorem 12.1: the inorder walk of the Figure 12.2 tree is 2 3 4 6 7 9 13 15 17 18 20.
func TestKthSmallestBST_clrsTree(t *testing.T) {
	root := clrsBST(t).Root()
	for k, want := range []int{2, 3, 4, 6, 7, 9, 13, 15, 17, 18, 20} {
		if got, err := KthSmallestBST(root, k+1); err != nil || got != want {
			t.Errorf("KthSmallestBST(k=%d) = (%d, %v), want %d", k+1, got, err, want)
		}
	}
}

func TestKthSmallestBST_errors(t *testing.T) {
	root := clrsBST(t).Root()
	for _, k := range []int{0, -1, 12} {
		if _, err := KthSmallestBST(root, k); !errors.Is(err, ErrRange) {
			t.Errorf("KthSmallestBST(k=%d) error = %v, want ErrRange", k, err)
		}
	}
	if _, err := KthSmallestBST[int](nil, 1); !errors.Is(err, ErrRange) {
		t.Errorf("empty tree: error = %v, want ErrRange", err)
	}
}

func TestKthSmallestBST_deepChain(t *testing.T) {
	root := chain(deepN, false) // 0 → 1 → … all right children: a valid BST
	if got, err := KthSmallestBST(root, deepN); err != nil || got != deepN-1 {
		t.Errorf("KthSmallestBST(chain, n) = (%d, %v), want %d", got, err, deepN-1)
	}
}

func TestKthSmallestBST_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		var b BST[int]
		for range 1 + rng.IntN(80) {
			b.Insert(rng.IntN(200) - 100)
		}
		sorted := slices.Collect(Inorder(b.Root()))
		k := 1 + rng.IntN(len(sorted))
		if got, err := KthSmallestBST(b.Root(), k); err != nil || got != sorted[k-1] {
			t.Fatalf("trial %d: KthSmallestBST(k=%d) = (%d, %v), want %d", trial, k, got, err, sorted[k-1])
		}
	}
}

// k = n/2 on a balanced tree: O(h + k) grows with n.
func BenchmarkKthSmallestBST(b *testing.B) {
	for _, size := range benchSizes {
		root := balancedTree(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = KthSmallestBST(root, size.n/2)
			}
		})
	}
}
