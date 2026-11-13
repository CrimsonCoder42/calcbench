package tree

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// clrsKeys inserted in this order give the tree of CLRS 3rd ed. §12.2 Figure 12.2 (root 15): built by hand,
// 15 → (6, 18); 6 → (3, 7); 3 → (2, 4); 7 → (nil, 13); 13 → (9, nil); 18 → (17, 20).
var clrsKeys = []int{15, 6, 18, 3, 7, 17, 20, 2, 4, 13, 9}

func clrsBST(t testing.TB) *BST[int] {
	var b BST[int]
	for _, k := range clrsKeys {
		if !b.Insert(k) {
			t.Fatalf("Insert(%d) = false on a fresh key", k)
		}
	}
	return &b
}

func TestBSTInsert_clrsExample(t *testing.T) {
	b := clrsBST(t)
	if got, want := slices.Collect(Preorder(b.Root())), []int{15, 6, 3, 2, 4, 7, 13, 9, 18, 17, 20}; !slices.Equal(got, want) {
		t.Errorf("preorder = %v, want %v", got, want)
	}
	if got := slices.Collect(Inorder(b.Root())); !slices.IsSorted(got) || len(got) != len(clrsKeys) {
		t.Errorf("inorder = %v, want the sorted keys", got)
	}
	if b.Len() != len(clrsKeys) {
		t.Errorf("Len = %d, want %d", b.Len(), len(clrsKeys))
	}
}

func TestBSTInsert_duplicate(t *testing.T) {
	b := clrsBST(t)
	if b.Insert(7) || b.Len() != len(clrsKeys) {
		t.Errorf("Insert(duplicate 7) = true or Len changed to %d", b.Len())
	}
}

func TestBSTSearch_cases(t *testing.T) {
	b := clrsBST(t)
	for _, k := range clrsKeys {
		if !b.Search(k) {
			t.Errorf("Search(%d) = false, want true", k)
		}
	}
	for _, k := range []int{0, 1, 5, 14, 16, 21, math.MinInt, math.MaxInt} {
		if b.Search(k) {
			t.Errorf("Search(%d) = true, want false", k)
		}
	}
	var empty BST[int]
	if empty.Search(1) {
		t.Error("Search on an empty BST = true")
	}
}

// CLRS §12.3 Figure 12.4 cases: a leaf, a node with one child, a node with two children whose successor
// is its right child, and one whose successor lies deeper.
func TestBSTDelete_clrsCases(t *testing.T) {
	tests := []struct {
		name string
		key  int
		pre  []int
	}{
		{"leaf", 2, []int{15, 6, 3, 4, 7, 13, 9, 18, 17, 20}},
		{"one child", 13, []int{15, 6, 3, 2, 4, 7, 9, 18, 17, 20}},
		{"two children, successor is right child", 18, []int{15, 6, 3, 2, 4, 7, 13, 9, 20, 17}},
		{"two children, successor deeper", 15, []int{17, 6, 3, 2, 4, 7, 13, 9, 18, 20}},
		{"two children, successor 7", 6, []int{15, 7, 3, 2, 4, 13, 9, 18, 17, 20}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := clrsBST(t)
			if !b.Delete(tt.key) {
				t.Fatalf("Delete(%d) = false", tt.key)
			}
			if got := slices.Collect(Preorder(b.Root())); !slices.Equal(got, tt.pre) {
				t.Errorf("preorder after Delete(%d) = %v, want %v", tt.key, got, tt.pre)
			}
			if b.Len() != len(clrsKeys)-1 || b.Search(tt.key) {
				t.Errorf("Len = %d or key still found", b.Len())
			}
			if err := ValidateBST(b.Root()); err != nil {
				t.Error(err)
			}
		})
	}
}

// Deletion relinks the successor node instead of copying its key, so node identity is kept.
func TestBSTDelete_keepsNodeIdentity(t *testing.T) {
	b := clrsBST(t)
	succ := b.Root().Right.Left // the node holding 17
	b.Delete(15)
	if b.Root() != succ || succ.Value != 17 {
		t.Errorf("root after Delete(15) is %p holding %d, want the old 17 node %p", b.Root(), b.Root().Value, succ)
	}
}

func TestBSTDelete_absentAndEmpty(t *testing.T) {
	b := clrsBST(t)
	if b.Delete(100) || b.Len() != len(clrsKeys) {
		t.Error("Delete(absent) = true or Len changed")
	}
	var empty BST[int]
	if empty.Delete(1) || empty.Len() != 0 {
		t.Error("Delete on empty BST")
	}
	one := BST[int]{}
	one.Insert(5)
	if !one.Delete(5) || one.Len() != 0 || one.Root() != nil {
		t.Error("deleting the only key did not leave an empty tree")
	}
}

func TestBST_nilReceiver(t *testing.T) {
	var b *BST[int]
	if b.Insert(1) || b.Search(1) || b.Delete(1) || b.Len() != 0 || b.Root() != nil {
		t.Error("nil *BST methods did not report empty/false")
	}
}

func TestBSTLen_counts(t *testing.T) {
	var b BST[int]
	for i, k := range []int{5, 3, 5, 8, 3, 1} {
		b.Insert(k)
		want := []int{1, 2, 2, 3, 3, 4}[i]
		if b.Len() != want {
			t.Fatalf("after inserting %d: Len = %d, want %d", k, b.Len(), want)
		}
	}
}

func TestBSTRoot_emptyAndOne(t *testing.T) {
	var b BST[string]
	if b.Root() != nil {
		t.Error("Root of empty BST is not nil")
	}
	b.Insert("k")
	if r := b.Root(); r == nil || r.Value != "k" || r.Left != nil || r.Right != nil {
		t.Errorf("Root = %+v, want a single node k", r)
	}
}

func TestBSTInsert_floatsNaNAndZero(t *testing.T) {
	var b BST[float64]
	for _, v := range []float64{1, math.NaN(), math.Inf(-1), 0, math.Copysign(0, -1), math.NaN(), math.Inf(1)} {
		b.Insert(v)
	}
	// cmp.Compare: NaN equals NaN and sorts first; -0 equals +0.
	got := slices.Collect(Inorder(b.Root()))
	if len(got) != 5 || !math.IsNaN(got[0]) || !math.IsInf(got[1], -1) || got[2] != 0 || got[3] != 1 || !math.IsInf(got[4], 1) {
		t.Errorf("inorder = %v, want [NaN -Inf 0 1 +Inf]", got)
	}
	if !b.Search(math.NaN()) || !b.Delete(math.Copysign(0, -1)) || b.Search(0) {
		t.Error("NaN not found, or -0 not deleted as 0")
	}
}

// A chain from sorted insertion is handled iteratively (no recursion on depth).
func TestBST_sortedInsertionChain(t *testing.T) {
	const n = 10_000
	var b BST[int]
	for k := range n {
		b.Insert(k)
	}
	if Height(b.Root()) != n-1 || !b.Search(n-1) {
		t.Fatalf("Height = %d, want %d", Height(b.Root()), n-1)
	}
	for k := range n {
		if !b.Delete(k) {
			t.Fatalf("Delete(%d) = false", k)
		}
	}
	if b.Len() != 0 || b.Root() != nil {
		t.Error("tree not empty after deleting every key")
	}
}

// Oracle: random inserts, searches and deletes against a map, checking order and size throughout.
func TestBST_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 200 {
		var b BST[int]
		model := map[int]bool{}
		for op := range 200 {
			k := rng.IntN(60) - 30
			switch rng.IntN(3) {
			case 0:
				if got, want := b.Insert(k), !model[k]; got != want {
					t.Fatalf("trial %d op %d: Insert(%d) = %v, want %v", trial, op, k, got, want)
				}
				model[k] = true
			case 1:
				if got, want := b.Delete(k), model[k]; got != want {
					t.Fatalf("trial %d op %d: Delete(%d) = %v, want %v", trial, op, k, got, want)
				}
				delete(model, k)
			default:
				if got, want := b.Search(k), model[k]; got != want {
					t.Fatalf("trial %d op %d: Search(%d) = %v, want %v", trial, op, k, got, want)
				}
			}
		}
		if b.Len() != len(model) || ValidateBST(b.Root()) != nil {
			t.Fatalf("trial %d: Len = %d, want %d; valid: %v", trial, b.Len(), len(model), ValidateBST(b.Root()))
		}
		keys := make([]int, 0, len(model))
		for k := range model {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		if got := slices.Collect(Inorder(b.Root())); !slices.Equal(got, keys) {
			t.Fatalf("trial %d: inorder = %v, want %v", trial, got, keys)
		}
	}
}

func randomKeys(n int) []int {
	rng := rand.New(rand.NewPCG(1, 2))
	keys := make([]int, n)
	for i := range keys {
		keys[i] = rng.Int()
	}
	return keys
}

func BenchmarkBSTInsert(b *testing.B) {
	for _, size := range benchSizes {
		keys := randomKeys(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				var t BST[int]
				for _, k := range keys {
					t.Insert(k)
				}
			}
		})
	}
}

func BenchmarkBSTSearch(b *testing.B) {
	for _, size := range benchSizes {
		keys := randomKeys(size.n)
		var t BST[int]
		for _, k := range keys {
			t.Insert(k)
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			i := 0
			for b.Loop() {
				_ = t.Search(keys[i%len(keys)])
				i++
			}
		})
	}
}

func BenchmarkBSTDelete(b *testing.B) {
	for _, size := range benchSizes {
		keys := randomKeys(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				b.StopTimer()
				var t BST[int]
				for _, k := range keys {
					t.Insert(k)
				}
				b.StartTimer()
				for _, k := range keys {
					t.Delete(k)
				}
			}
		})
	}
}

func BenchmarkBSTLen(b *testing.B) {
	t := clrsBST(b)
	b.ReportAllocs()
	for b.Loop() {
		_ = t.Len()
	}
}

func BenchmarkBSTRoot(b *testing.B) {
	t := clrsBST(b)
	b.ReportAllocs()
	for b.Loop() {
		_ = t.Root()
	}
}
