package tree

import (
	"iter"
	"math/rand/v2"
	"slices"
	"testing"
)

var benchSizes = []struct {
	name string
	n    int
}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}

// deepN is the length of the degenerate chains that prove no walk recurses per node.
const deepN = 100_000

// sharedTree is the binary tree of Laaksonen, Competitive Programmer's Handbook (2018), §14.4, p. 139:
// 1 → (2, 3); 2 → (4, 5); 5 → (6, nil); 3 → (nil, 7).
func sharedTree() *Node[int] {
	n := func(v int, l, r *Node[int]) *Node[int] { return &Node[int]{Value: v, Left: l, Right: r} }
	return n(1, n(2, n(4, nil, nil), n(5, n(6, nil, nil), nil)), n(3, nil, n(7, nil, nil)))
}

// Orders given in the book for sharedTree (p. 139); level order read off its figure.
var (
	sharedPre   = []int{1, 2, 4, 5, 6, 3, 7}
	sharedIn    = []int{4, 2, 6, 5, 1, 3, 7}
	sharedPost  = []int{4, 6, 5, 2, 7, 3, 1}
	sharedLevel = []int{1, 2, 3, 4, 5, 7, 6}
)

// chain returns a degenerate tree 0 → 1 → … → n−1, every child on the left (or the right).
func chain(n int, left bool) *Node[int] {
	var root *Node[int]
	for v := n - 1; v >= 0; v-- {
		nd := &Node[int]{Value: v}
		if left {
			nd.Left = root
		} else {
			nd.Right = root
		}
		root = nd
	}
	return root
}

// randomTree returns a tree of n nodes holding 0…n−1 (the order of creation) in a random shape.
func randomTree(rng *rand.Rand, n int) *Node[int] {
	if n == 0 {
		return nil
	}
	nodes := []*Node[int]{{Value: 0}}
	for v := 1; v < n; v++ {
		nd := &Node[int]{Value: v}
		for {
			p := nodes[rng.IntN(len(nodes))]
			if rng.IntN(2) == 0 && p.Left == nil {
				p.Left = nd
				break
			}
			if p.Right == nil {
				p.Right = nd
				break
			}
		}
		nodes = append(nodes, nd)
	}
	return nodes[0]
}

// balancedTree returns a perfectly balanced BST of 0…n−1 (for benchmarks).
func balancedTree(n int) *Node[int] {
	keys := make([]int, n)
	for i := range keys {
		keys[i] = i
	}
	b, err := SortedToBalancedBST(keys)
	if err != nil {
		panic(err)
	}
	return b.Root()
}

// Recursive reference walks for small trees (oracles).
func refPre(n *Node[int], out []int) []int {
	if n == nil {
		return out
	}
	out = append(out, n.Value)
	out = refPre(n.Left, out)
	return refPre(n.Right, out)
}

func refIn(n *Node[int], out []int) []int {
	if n == nil {
		return out
	}
	out = refIn(n.Left, out)
	out = append(out, n.Value)
	return refIn(n.Right, out)
}

func refPost(n *Node[int], out []int) []int {
	if n == nil {
		return out
	}
	out = refPost(n.Left, out)
	out = refPost(n.Right, out)
	return append(out, n.Value)
}

func refHeight(n *Node[int]) int {
	if n == nil {
		return -1
	}
	return 1 + max(refHeight(n.Left), refHeight(n.Right))
}

// refLevels returns the values of each level, left to right.
func refLevels(n *Node[int], depth int, levels [][]int) [][]int {
	if n == nil {
		return levels
	}
	if depth == len(levels) {
		levels = append(levels, nil)
	}
	levels[depth] = append(levels[depth], n.Value)
	levels = refLevels(n.Left, depth+1, levels)
	return refLevels(n.Right, depth+1, levels)
}

var traversals = []struct {
	name string
	walk func(*Node[int]) iter.Seq[int]
	want []int
}{
	{"Inorder", Inorder[int], sharedIn},
	{"Preorder", Preorder[int], sharedPre},
	{"Postorder", Postorder[int], sharedPost},
	{"LevelOrder", LevelOrder[int], sharedLevel},
}

// Laaksonen (2018) §14.4 p. 139: preorder 1 2 4 5 6 3 7, inorder 4 2 6 5 1 3 7, postorder 4 6 5 2 7 3 1.
func checkLaaksonenExample(t *testing.T, i int) {
	tr := traversals[i]
	if got := slices.Collect(tr.walk(sharedTree())); !slices.Equal(got, tr.want) {
		t.Errorf("%s(shared tree) = %v, want %v", tr.name, got, tr.want)
	}
}

func checkEmptyAndOne(t *testing.T, i int) {
	tr := traversals[i]
	if got := slices.Collect(tr.walk(nil)); len(got) != 0 {
		t.Errorf("%s(nil) = %v, want nothing", tr.name, got)
	}
	if got := slices.Collect(tr.walk(&Node[int]{Value: 9})); !slices.Equal(got, []int{9}) {
		t.Errorf("%s(one node) = %v, want [9]", tr.name, got)
	}
}

func checkTwoNodes(t *testing.T, i int) {
	left := &Node[int]{Value: 1, Left: &Node[int]{Value: 2}}
	right := &Node[int]{Value: 1, Right: &Node[int]{Value: 2}}
	cases := []struct {
		name string
		walk func(*Node[int]) iter.Seq[int]
		l, r []int
	}{
		{"Inorder", Inorder[int], []int{2, 1}, []int{1, 2}},
		{"Preorder", Preorder[int], []int{1, 2}, []int{1, 2}},
		{"Postorder", Postorder[int], []int{2, 1}, []int{2, 1}},
		{"LevelOrder", LevelOrder[int], []int{1, 2}, []int{1, 2}},
	}
	c := cases[i]
	if got := slices.Collect(c.walk(left)); !slices.Equal(got, c.l) {
		t.Errorf("%s(left child) = %v, want %v", c.name, got, c.l)
	}
	if got := slices.Collect(c.walk(right)); !slices.Equal(got, c.r) {
		t.Errorf("%s(right child) = %v, want %v", c.name, got, c.r)
	}
}

// Breaking the range loop must stop the walk: no further yield, no panic from a second call.
func checkEarlyBreak(t *testing.T, i int) {
	for _, tr := range traversals[i : i+1] {
		for stopAfter := 1; stopAfter <= len(tr.want); stopAfter++ {
			var got []int
			for v := range tr.walk(sharedTree()) {
				got = append(got, v)
				if len(got) == stopAfter {
					break
				}
			}
			if !slices.Equal(got, tr.want[:stopAfter]) {
				t.Errorf("%s break after %d: got %v, want %v", tr.name, stopAfter, got, tr.want[:stopAfter])
			}
		}
		// Driving the iterator by hand: after yield returns false it is never called again.
		calls := 0
		tr.walk(sharedTree())(func(int) bool { calls++; return false })
		if calls != 1 {
			t.Errorf("%s: yield called %d times after returning false, want 1", tr.name, calls)
		}
	}
}

// A degenerate chain of 1e5 nodes in either direction is walked without recursion.
func checkDeepChain(t *testing.T, i int) {
	for _, left := range []bool{true, false} {
		root := chain(deepN, left)
		for _, tr := range traversals[i : i+1] {
			count, sum := 0, 0
			for v := range tr.walk(root) {
				count++
				sum += v
			}
			if count != deepN || sum != deepN*(deepN-1)/2 {
				t.Errorf("%s(chain left=%v) visited %d nodes (sum %d), want %d", tr.name, left, count, sum, deepN)
			}
		}
	}
}

func TestInorder_strings(t *testing.T) {
	root := &Node[string]{Value: "m", Left: &Node[string]{Value: "c"}, Right: &Node[string]{Value: "x"}}
	if got := slices.Collect(Inorder(root)); !slices.Equal(got, []string{"c", "m", "x"}) {
		t.Errorf("Inorder(strings) = %v", got)
	}
}

// Oracle: the iterative walks against recursive ones on random shapes (all four in one pass).
func TestLevelOrder_oracleWithOtherWalks(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		root := randomTree(rng, rng.IntN(60))
		if got, want := slices.Collect(Preorder(root)), refPre(root, nil); !slices.Equal(got, want) {
			t.Fatalf("trial %d: Preorder = %v, want %v", trial, got, want)
		}
		if got, want := slices.Collect(Inorder(root)), refIn(root, nil); !slices.Equal(got, want) {
			t.Fatalf("trial %d: Inorder = %v, want %v", trial, got, want)
		}
		if got, want := slices.Collect(Postorder(root)), refPost(root, nil); !slices.Equal(got, want) {
			t.Fatalf("trial %d: Postorder = %v, want %v", trial, got, want)
		}
		if got, want := slices.Collect(LevelOrder(root)), slices.Concat(refLevels(root, 0, nil)...); !slices.Equal(got, want) {
			t.Fatalf("trial %d: LevelOrder = %v, want %v", trial, got, want)
		}
	}
}

func TestInorder_laaksonenExample(t *testing.T)    { checkLaaksonenExample(t, 0) }
func TestInorder_emptyAndOne(t *testing.T)         { checkEmptyAndOne(t, 0) }
func TestInorder_twoNodes(t *testing.T)            { checkTwoNodes(t, 0) }
func TestInorder_earlyBreak(t *testing.T)          { checkEarlyBreak(t, 0) }
func TestInorder_deepChain(t *testing.T)           { checkDeepChain(t, 0) }
func TestPreorder_laaksonenExample(t *testing.T)   { checkLaaksonenExample(t, 1) }
func TestPreorder_emptyAndOne(t *testing.T)        { checkEmptyAndOne(t, 1) }
func TestPreorder_twoNodes(t *testing.T)           { checkTwoNodes(t, 1) }
func TestPreorder_earlyBreak(t *testing.T)         { checkEarlyBreak(t, 1) }
func TestPreorder_deepChain(t *testing.T)          { checkDeepChain(t, 1) }
func TestPostorder_laaksonenExample(t *testing.T)  { checkLaaksonenExample(t, 2) }
func TestPostorder_emptyAndOne(t *testing.T)       { checkEmptyAndOne(t, 2) }
func TestPostorder_twoNodes(t *testing.T)          { checkTwoNodes(t, 2) }
func TestPostorder_earlyBreak(t *testing.T)        { checkEarlyBreak(t, 2) }
func TestPostorder_deepChain(t *testing.T)         { checkDeepChain(t, 2) }
func TestLevelOrder_laaksonenExample(t *testing.T) { checkLaaksonenExample(t, 3) }
func TestLevelOrder_emptyAndOne(t *testing.T)      { checkEmptyAndOne(t, 3) }
func TestLevelOrder_twoNodes(t *testing.T)         { checkTwoNodes(t, 3) }
func TestLevelOrder_earlyBreak(t *testing.T)       { checkEarlyBreak(t, 3) }
func TestLevelOrder_deepChain(t *testing.T)        { checkDeepChain(t, 3) }

func benchmarkWalk(b *testing.B, walk func(*Node[int]) iter.Seq[int]) {
	for _, size := range benchSizes {
		root := balancedTree(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				s := 0
				for v := range walk(root) {
					s += v
				}
				_ = s
			}
		})
	}
}

func BenchmarkInorder(b *testing.B)    { benchmarkWalk(b, Inorder[int]) }
func BenchmarkPreorder(b *testing.B)   { benchmarkWalk(b, Preorder[int]) }
func BenchmarkPostorder(b *testing.B)  { benchmarkWalk(b, Postorder[int]) }
func BenchmarkLevelOrder(b *testing.B) { benchmarkWalk(b, LevelOrder[int]) }
