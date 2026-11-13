package tree

import (
	"errors"
	"math"
	"slices"
	"testing"
)

// sizesPerfectlyBalanced reports Wirth's condition: at every node the subtree sizes differ by at most 1.
func sizesPerfectlyBalanced(root *Node[int]) bool {
	_, ok := foldPostorder(root, 0, func(_ *Node[int], l, r int) (int, bool) {
		return l + r + 1, l-r <= 1 && r-l <= 1
	})
	return ok
}

// Wirth §4.4.1 rule (left gets n DIV 2 nodes) applied by hand to 1…7: root 4, (2: 1, 3), (6: 5, 7).
func TestSortedToBalancedBST_wirthRule(t *testing.T) {
	b, err := SortedToBalancedBST([]int{1, 2, 3, 4, 5, 6, 7})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := slices.Collect(Preorder(b.Root())), []int{4, 2, 1, 3, 6, 5, 7}; !slices.Equal(got, want) {
		t.Errorf("preorder = %v, want %v", got, want)
	}
	// Six keys: the left gets 3, the right 2.
	b, _ = SortedToBalancedBST([]int{1, 2, 3, 4, 5, 6})
	if got, want := slices.Collect(Preorder(b.Root())), []int{4, 2, 1, 3, 6, 5}; !slices.Equal(got, want) {
		t.Errorf("preorder(6 keys) = %v, want %v", got, want)
	}
}

func TestSortedToBalancedBST_cases(t *testing.T) {
	for _, n := range []int{0, 1, 2, 3, 8, 100, 1023, 1024, deepN} {
		keys := make([]int, n)
		for i := range keys {
			keys[i] = i*2 - n // negatives, zero, positives
		}
		b, err := SortedToBalancedBST(keys)
		if err != nil {
			t.Fatalf("n=%d: %v", n, err)
		}
		wantHeight := -1
		if n > 0 {
			wantHeight = 0
			for 1<<(wantHeight+1) <= n {
				wantHeight++
			}
		}
		if b.Len() != n || Height(b.Root()) != wantHeight || !sizesPerfectlyBalanced(b.Root()) ||
			!slices.Equal(slices.Collect(Inorder(b.Root())), keys) || ValidateBST(b.Root()) != nil {
			t.Errorf("n=%d: Len %d, Height %d (want ⌊log2 n⌋ = %d), not perfectly balanced or not in order", n, b.Len(), Height(b.Root()), wantHeight)
		}
	}
}

func TestSortedToBalancedBST_extremesAndUsable(t *testing.T) {
	b, err := SortedToBalancedBST([]int{math.MinInt, 0, math.MaxInt})
	if err != nil || b.Root().Value != 0 {
		t.Fatalf("got root %v, err %v", b.Root(), err)
	}
	if !b.Insert(5) || !b.Delete(math.MinInt) || b.Len() != 3 || ValidateBST(b.Root()) != nil {
		t.Error("the returned BST is not usable for Insert/Delete")
	}
}

func TestSortedToBalancedBST_errors(t *testing.T) {
	for _, keys := range [][]int{{1, 1}, {2, 1}, {1, 2, 3, 3}, {1, 3, 2}} {
		if b, err := SortedToBalancedBST(keys); !errors.Is(err, ErrNotSorted) || b != nil {
			t.Errorf("SortedToBalancedBST(%v) = (%v, %v), want ErrNotSorted", keys, b, err)
		}
	}
	if _, err := SortedToBalancedBST([]float64{math.NaN(), math.NaN()}); !errors.Is(err, ErrNotSorted) {
		t.Errorf("two NaNs (equal under cmp.Compare): err %v, want ErrNotSorted", err)
	}
}

func TestSortedToBalancedBST_doesNotModifyInput(t *testing.T) {
	keys := []int{1, 2, 3, 4, 5}
	b, _ := SortedToBalancedBST(keys)
	b.Delete(3)
	if !slices.Equal(keys, []int{1, 2, 3, 4, 5}) {
		t.Errorf("input modified: %v", keys)
	}
}

func BenchmarkSortedToBalancedBST(b *testing.B) {
	for _, size := range benchSizes {
		keys := make([]int, size.n)
		for i := range keys {
			keys[i] = i
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = SortedToBalancedBST(keys)
			}
		})
	}
}
