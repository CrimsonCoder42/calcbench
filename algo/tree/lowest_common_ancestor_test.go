package tree

import (
	"math/rand/v2"
	"testing"
)

// Hand-checked on the BST built from clrsKeys (CLRS 3rd ed. Figure 12.2 shape).
func TestLowestCommonAncestorBST_clrsTree(t *testing.T) {
	root := clrsBST(t).Root()
	tests := []struct{ p, q, want int }{
		{2, 4, 3}, {4, 2, 3}, {9, 17, 15}, {13, 9, 13}, {4, 9, 6}, {15, 15, 15}, {2, 2, 2}, {17, 20, 18},
	}
	for _, tt := range tests {
		n, ok := LowestCommonAncestorBST(root, tt.p, tt.q)
		if !ok || n == nil || n.Value != tt.want {
			t.Errorf("LowestCommonAncestorBST(%d, %d) = (%v, %v), want %d", tt.p, tt.q, n, ok, tt.want)
		}
	}
}

func TestLowestCommonAncestorBST_missing(t *testing.T) {
	root := clrsBST(t).Root()
	for _, pq := range [][2]int{{2, 5}, {5, 2}, {1, 100}, {14, 16}} {
		if n, ok := LowestCommonAncestorBST(root, pq[0], pq[1]); ok || n != nil {
			t.Errorf("LowestCommonAncestorBST(%d, %d) = (%v, true), want (nil, false)", pq[0], pq[1], n)
		}
	}
	if n, ok := LowestCommonAncestorBST[int](nil, 1, 1); ok || n != nil {
		t.Error("empty tree gave an ancestor")
	}
}

// refPath returns the root-to-target path, or nil.
func refPath(n, target *Node[int], path []*Node[int]) []*Node[int] {
	if n == nil {
		return nil
	}
	path = append(path, n)
	if n == target {
		return path
	}
	if p := refPath(n.Left, target, path); p != nil {
		return p
	}
	return refPath(n.Right, target, path)
}

// The shared tree of Laaksonen (2018) §14.4: LCA(4, 6) = 2, LCA(6, 7) = 1, LCA(5, 6) = 5 (hand-checked).
func TestLowestCommonAncestor_sharedTree(t *testing.T) {
	root := sharedTree()
	byValue := map[int]*Node[int]{}
	for _, n := range collectNodes(root) {
		byValue[n.Value] = n
	}
	for _, tt := range []struct{ p, q, want int }{{4, 6, 2}, {6, 7, 1}, {5, 6, 5}, {6, 5, 5}, {3, 3, 3}, {1, 7, 1}} {
		if got := LowestCommonAncestor(root, byValue[tt.p], byValue[tt.q]); got != byValue[tt.want] {
			t.Errorf("LowestCommonAncestor(%d, %d) = %v, want node %d", tt.p, tt.q, got, tt.want)
		}
	}
}

func TestLowestCommonAncestor_missing(t *testing.T) {
	root := sharedTree()
	stranger := &Node[int]{Value: 4}
	if got := LowestCommonAncestor(root, root.Left, stranger); got != nil {
		t.Errorf("node outside the tree: got %v, want nil", got)
	}
	if got := LowestCommonAncestor(root, nil, root); got != nil {
		t.Errorf("nil p: got %v, want nil", got)
	}
	if got := LowestCommonAncestor[int](nil, root, root); got != nil {
		t.Errorf("empty tree: got %v, want nil", got)
	}
}

func TestLowestCommonAncestor_deepChain(t *testing.T) {
	root := chain(deepN, true)
	bottom := root
	for bottom.Left != nil {
		bottom = bottom.Left
	}
	if got := LowestCommonAncestor(root, root, bottom); got != root {
		t.Errorf("LCA(root, bottom of 1e5 chain) = %v, want the root", got)
	}
}

func collectNodes(root *Node[int]) []*Node[int] {
	var out []*Node[int]
	stack := []*Node[int]{root}
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if n == nil {
			continue
		}
		out = append(out, n)
		stack = append(stack, n.Left, n.Right)
	}
	return out
}

// Oracle: the last common node of the two root paths.
func TestLowestCommonAncestor_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		root := randomTree(rng, 1+rng.IntN(60))
		nodes := collectNodes(root)
		p, q := nodes[rng.IntN(len(nodes))], nodes[rng.IntN(len(nodes))]
		pp, qp := refPath(root, p, nil), refPath(root, q, nil)
		var want *Node[int]
		for i := 0; i < len(pp) && i < len(qp) && pp[i] == qp[i]; i++ {
			want = pp[i]
		}
		if got := LowestCommonAncestor(root, p, q); got != want {
			t.Fatalf("trial %d: LCA(%d, %d) = %v, want %v", trial, p.Value, q.Value, got, want)
		}
	}
}

// Oracle for the BST version: the general version on the same nodes.
func TestLowestCommonAncestorBST_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		var b BST[int]
		for range 1 + rng.IntN(60) {
			b.Insert(rng.IntN(100))
		}
		keys := slicesOf(b.Root())
		p, q := keys[rng.IntN(len(keys))], keys[rng.IntN(len(keys))]
		got, ok := LowestCommonAncestorBST(b.Root(), p, q)
		want := LowestCommonAncestor(b.Root(), searchNode(b.Root(), p), searchNode(b.Root(), q))
		if !ok || got != want {
			t.Fatalf("trial %d: LowestCommonAncestorBST(%d, %d) = (%v, %v), want %v", trial, p, q, got, ok, want)
		}
	}
}

func slicesOf(root *Node[int]) []int {
	var out []int
	for v := range Inorder(root) {
		out = append(out, v)
	}
	return out
}

func BenchmarkLowestCommonAncestorBST(b *testing.B) {
	for _, size := range benchSizes {
		root := balancedTree(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = LowestCommonAncestorBST(root, 0, size.n-1)
			}
		})
	}
}

func BenchmarkLowestCommonAncestor(b *testing.B) {
	for _, size := range benchSizes {
		root := balancedTree(size.n)
		p, q := searchNode(root, 0), searchNode(root, size.n-1) // the two ends: the walk visits every node
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = LowestCommonAncestor(root, p, q)
			}
		})
	}
}
