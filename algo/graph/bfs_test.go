package graph

import (
	"errors"
	"math/rand/v2"
	"slices"
	"testing"
)

// Vertices of Cormen et al., Introduction to Algorithms, 3rd ed., §22.2, Figure 22.3.
const (
	fr, fs, ft, fu, fv, fw, fx, fy = 0, 1, 2, 3, 4, 5, 6, 7
)

// figure22_3 is the undirected graph of CLRS Figure 22.3 (edges r–s r–v s–w w–t w–x t–x t–u x–u x–y u–y).
func figure22_3(t testing.TB) *Graph {
	g, err := NewUndirected(8, []Edge{
		{fr, fs, 0}, {fr, fv, 0}, {fs, fw, 0}, {fw, ft, 0}, {fw, fx, 0},
		{ft, fx, 0}, {ft, fu, 0}, {fx, fu, 0}, {fx, fy, 0}, {fu, fy, 0},
	})
	if err != nil {
		t.Fatal(err)
	}
	return g
}

// levels is a slow oracle: the edge distance from source of every vertex (-1 when unreachable), by
// relaxing every edge until nothing changes.
func levels(g *Graph, source int) []int {
	n := g.VertexCount()
	d := make([]int, n)
	for i := range d {
		d[i] = -1
	}
	d[source] = 0
	for changed := true; changed; {
		changed = false
		for u := range n {
			if d[u] < 0 {
				continue
			}
			for v := range g.Neighbors(u) {
				if d[v] < 0 || d[u]+1 < d[v] {
					d[v] = d[u] + 1
					changed = true
				}
			}
		}
	}
	return d
}

// CLRS Figure 22.3: BFS from s gives δ = r1 s0 t2 u3 v2 w1 x2 y3; with our edge order the visit order is
// s r w v t x u y.
func TestBFS_clrsFigure22_3(t *testing.T) {
	g := figure22_3(t)
	seq, err := BFS(g, fs)
	if err != nil {
		t.Fatal(err)
	}
	got := slices.Collect(seq)
	if want := []int{fs, fr, fw, fv, ft, fx, fu, fy}; !slices.Equal(got, want) {
		t.Errorf("BFS order = %v, want %v", got, want)
	}
	if d := levels(g, fs); !slices.Equal(d, []int{1, 0, 2, 3, 2, 1, 2, 3}) {
		t.Errorf("oracle distances = %v, want CLRS r1 s0 t2 u3 v2 w1 x2 y3", d)
	}
}

func TestBFS_cases(t *testing.T) {
	one, _ := NewDirected(1, nil)
	if got := slices.Collect(must(BFS(one, 0))); !slices.Equal(got, []int{0}) {
		t.Errorf("one vertex: %v", got)
	}
	g, _ := NewDirected(4, []Edge{{0, 1, 0}, {2, 3, 0}, {1, 0, 0}, {0, 0, 0}})
	if got := slices.Collect(must(BFS(g, 0))); !slices.Equal(got, []int{0, 1}) {
		t.Errorf("disconnected, self-loop: %v, want [0 1]", got)
	}
	if got := slices.Collect(must(BFS(g, 3))); !slices.Equal(got, []int{3}) {
		t.Errorf("sink: %v, want [3]", got)
	}
	seq := must(BFS(figure22_3(t), fs))
	var first []int
	for v := range seq {
		first = append(first, v)
		if len(first) == 3 {
			break
		}
	}
	if !slices.Equal(first, []int{fs, fr, fw}) {
		t.Errorf("break after 3: %v", first)
	}
	if again := slices.Collect(seq); len(again) != 8 {
		t.Errorf("second range yielded %d vertices, want a fresh full search of 8", len(again))
	}
}

func TestBFS_errors(t *testing.T) {
	if _, err := BFS(nil, 0); !errors.Is(err, ErrNilGraph) {
		t.Errorf("nil graph err = %v", err)
	}
	g, _ := NewDirected(2, nil)
	for _, s := range []int{-1, 2} {
		if seq, err := BFS(g, s); !errors.Is(err, ErrVertexOutOfRange) || seq != nil {
			t.Errorf("source %d: err = %v", s, err)
		}
	}
}

// Oracle: the BFS order visits exactly the reachable vertices, by non-decreasing distance, and each
// vertex after the source has a neighbour-predecessor one level up that came earlier.
func TestBFS_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for run := 0; run < 300; run++ {
		n := 1 + rng.IntN(15)
		edges := randomEdges(rng, n, rng.IntN(3*n), 0, 0)
		g, _ := NewDirected(n, edges)
		if run%2 == 1 {
			g, _ = NewUndirected(n, edges)
		}
		s := rng.IntN(n)
		order := slices.Collect(must(BFS(g, s)))
		d := levels(g, s)
		reach := 0
		for _, x := range d {
			if x >= 0 {
				reach++
			}
		}
		if len(order) != reach || order[0] != s {
			t.Fatalf("run %d: order %v, %d reachable", run, order, reach)
		}
		pos := make([]int, n)
		for i, v := range order {
			pos[v] = i
			if i > 0 && d[v] < d[order[i-1]] {
				t.Fatalf("run %d: distance decreases at %d in %v (d %v)", run, v, order, d)
			}
		}
		for _, v := range order[1:] {
			ok := false
			for u := range n {
				for w := range g.Neighbors(u) {
					if w == v && d[u] == d[v]-1 && pos[u] < pos[v] {
						ok = true
					}
				}
			}
			if !ok {
				t.Fatalf("run %d: %d has no earlier predecessor", run, v)
			}
		}
	}
}

// must returns seq, failing loudly (a nil sequence) if err is set; used where the arguments are valid.
func must[S any](seq S, err error) S {
	if err != nil {
		panic(err) // test helper only: the arguments are valid by construction
	}
	return seq
}

func BenchmarkBFS(b *testing.B) {
	for _, size := range benchSizes {
		g := benchGraph(size.n, true)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				seq, _ := BFS(g, 0)
				for range seq {
				}
			}
		})
	}
}
