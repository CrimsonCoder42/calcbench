package unionfind

import (
	"errors"
	"math"
	"math/bits"
	"math/rand/v2"
	"strconv"
	"testing"
)

var benchSizes = []struct {
	name string
	n    int
}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}

// Cormen et al., Introduction to Algorithms, 3rd ed., §21.1, Figure 21.1: vertices a..j (here 0..9),
// edges processed in the order (b,d) (e,g) (a,c) (h,i) (a,b) (e,f) (b,c); the components are
// {a,b,c,d} {e,f,g} {h,i} {j}.
func TestDisjointSet_clrsFigure21_1(t *testing.T) {
	const a, b, c, d, e, f, g, h, i, j = 0, 1, 2, 3, 4, 5, 6, 7, 8, 9
	s, err := NewDisjointSet(10)
	if err != nil {
		t.Fatal(err)
	}
	edges := [][2]int{{b, d}, {e, g}, {a, c}, {h, i}, {a, b}, {e, f}, {b, c}}
	wantMerged := []bool{true, true, true, true, true, true, false}
	for k, edge := range edges {
		merged, err := s.Union(edge[0], edge[1])
		if err != nil || merged != wantMerged[k] {
			t.Errorf("Union%v = (%v, %v), want (%v, nil)", edge, merged, err, wantMerged[k])
		}
	}
	if s.Count() != 4 {
		t.Errorf("Count = %d, want 4", s.Count())
	}
	for _, tc := range []struct {
		x, y int
		want bool
	}{{a, d, true}, {c, b, true}, {e, f, true}, {g, f, true}, {h, i, true}, {e, h, false}, {a, j, false}, {j, j, true}} {
		if got, err := s.Connected(tc.x, tc.y); err != nil || got != tc.want {
			t.Errorf("Connected(%d, %d) = (%v, %v), want (%v, nil)", tc.x, tc.y, got, err, tc.want)
		}
	}
}

func TestDisjointSet_cases(t *testing.T) {
	t.Run("n 0", func(t *testing.T) {
		s, err := NewDisjointSet(0)
		if err != nil || s.Count() != 0 || s.Len() != 0 {
			t.Errorf("NewDisjointSet(0) = (count %d, len %d, %v), want (0, 0, nil)", s.Count(), s.Len(), err)
		}
	})
	t.Run("one element", func(t *testing.T) {
		s, _ := NewDisjointSet(1)
		if r, err := s.Find(0); r != 0 || err != nil {
			t.Errorf("Find(0) = (%d, %v), want (0, nil)", r, err)
		}
		if merged, err := s.Union(0, 0); merged || err != nil {
			t.Errorf("Union(0, 0) = (%v, %v), want (false, nil)", merged, err)
		}
		if s.Count() != 1 {
			t.Errorf("Count = %d, want 1", s.Count())
		}
	})
	t.Run("two elements", func(t *testing.T) {
		s, _ := NewDisjointSet(2)
		if ok, _ := s.Connected(0, 1); ok {
			t.Error("Connected(0, 1) before Union = true")
		}
		if merged, _ := s.Union(1, 0); !merged {
			t.Error("Union(1, 0) = false, want true")
		}
		if merged, _ := s.Union(0, 1); merged {
			t.Error("repeated Union(0, 1) = true, want false")
		}
		r0, _ := s.Find(0)
		r1, _ := s.Find(1)
		if r0 != r1 || s.Count() != 1 {
			t.Errorf("roots %d %d, count %d; want equal roots, count 1", r0, r1, s.Count())
		}
	})
	t.Run("zero value and MakeSet", func(t *testing.T) {
		var s DisjointSet
		if s.Len() != 0 || s.Count() != 0 {
			t.Fatal("zero DisjointSet not empty")
		}
		for want := 0; want < 3; want++ {
			if got := s.MakeSet(); got != want {
				t.Errorf("MakeSet = %d, want %d", got, want)
			}
		}
		_, _ = s.Union(0, 2)
		if s.Len() != 3 || s.Count() != 2 {
			t.Errorf("Len %d Count %d, want 3 2", s.Len(), s.Count())
		}
		if ok, _ := s.Connected(2, 0); !ok {
			t.Error("Connected(2, 0) = false after Union(0, 2)")
		}
	})
	t.Run("chain stays shallow", func(t *testing.T) {
		const n = 1 << 12
		s, _ := NewDisjointSet(n)
		for i := 1; i < n; i++ {
			_, _ = s.Union(i-1, i)
		}
		maxRank := 0
		for _, r := range s.rank {
			maxRank = max(maxRank, int(r))
		}
		if limit := bits.Len(uint(n)) - 1; maxRank > limit {
			t.Errorf("max rank %d > log2 n = %d", maxRank, limit)
		}
		if s.Count() != 1 {
			t.Errorf("Count = %d, want 1", s.Count())
		}
	})
}

func TestDisjointSet_nil(t *testing.T) {
	var s *DisjointSet
	if s.Count() != 0 || s.Len() != 0 || s.MakeSet() != -1 {
		t.Error("nil DisjointSet: want Count 0, Len 0, MakeSet -1")
	}
	if _, err := s.Find(0); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("nil Find err = %v, want ErrOutOfRange", err)
	}
	if _, err := s.Union(0, 0); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("nil Union err = %v, want ErrOutOfRange", err)
	}
	if _, err := s.Connected(0, 0); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("nil Connected err = %v, want ErrOutOfRange", err)
	}
}

func TestDisjointSet_errors(t *testing.T) {
	for _, n := range []int{-1, math.MinInt} {
		if s, err := NewDisjointSet(n); !errors.Is(err, ErrNegativeSize) || s != nil {
			t.Errorf("NewDisjointSet(%d) = (%v, %v), want (nil, ErrNegativeSize)", n, s, err)
		}
	}
	// L41: a size whose arrays exceed 2⁴⁷ bytes (2³¹ on 32-bit) is an error, never a makeslice panic.
	for _, n := range []int{math.MaxInt, maxAllocBytes/(strconv.IntSize/8) + 1} {
		if s, err := NewDisjointSet(n); !errors.Is(err, ErrTooLarge) || s != nil {
			t.Errorf("NewDisjointSet(%d) = (%v, %v), want (nil, ErrTooLarge)", n, s, err)
		}
	}
	s, _ := NewDisjointSet(3)
	for _, x := range []int{-1, 3, math.MaxInt, math.MinInt} {
		if _, err := s.Find(x); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("Find(%d) err = %v, want ErrOutOfRange", x, err)
		}
		if _, err := s.Union(x, 0); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("Union(%d, 0) err = %v, want ErrOutOfRange", x, err)
		}
		if _, err := s.Union(0, x); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("Union(0, %d) err = %v, want ErrOutOfRange", x, err)
		}
		if _, err := s.Connected(x, 0); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("Connected(%d, 0) err = %v, want ErrOutOfRange", x, err)
		}
		if _, err := s.Connected(0, x); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("Connected(0, %d) err = %v, want ErrOutOfRange", x, err)
		}
	}
	if s.Count() != 3 {
		t.Errorf("Count after failed calls = %d, want 3", s.Count())
	}
}

// naiveSets is the oracle: a label per element, relabelled in full on every union.
type naiveSets struct{ label []int }

func (ns *naiveSets) union(a, b int) bool {
	la, lb := ns.label[a], ns.label[b]
	if la == lb {
		return false
	}
	for i, l := range ns.label {
		if l == lb {
			ns.label[i] = la
		}
	}
	return true
}

func (ns *naiveSets) count() int {
	seen := map[int]bool{}
	for _, l := range ns.label {
		seen[l] = true
	}
	return len(seen)
}

func TestDisjointSet_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for run := 0; run < 300; run++ {
		n := 1 + rng.IntN(40)
		s, _ := NewDisjointSet(n)
		ns := &naiveSets{label: make([]int, n)}
		for i := range ns.label {
			ns.label[i] = i
		}
		for op := 0; op < 3*n; op++ {
			a, b := rng.IntN(n), rng.IntN(n)
			if rng.IntN(2) == 0 {
				got, _ := s.Union(a, b)
				if want := ns.union(a, b); got != want {
					t.Fatalf("run %d: Union(%d, %d) = %v, want %v", run, a, b, got, want)
				}
			} else {
				got, _ := s.Connected(a, b)
				if want := ns.label[a] == ns.label[b]; got != want {
					t.Fatalf("run %d: Connected(%d, %d) = %v, want %v", run, a, b, got, want)
				}
			}
			if s.Count() != ns.count() {
				t.Fatalf("run %d: Count = %d, want %d", run, s.Count(), ns.count())
			}
		}
	}
}

func TestDisjointSet_noAllocs(t *testing.T) {
	s, _ := NewDisjointSet(1000)
	i := 0
	allocs := testing.AllocsPerRun(100, func() {
		_, _ = s.Union(i%1000, (i*7+3)%1000)
		_, _ = s.Find(i % 1000)
		_, _ = s.Connected(i%1000, (i+1)%1000)
		_ = s.Count()
		_ = s.Len()
		i++
	})
	if allocs != 0 {
		t.Errorf("Union/Find/Connected allocated %v times per run, want 0", allocs)
	}
}

// randomPairs returns n pairs of elements of 0..n-1 from the fixed seed.
func randomPairs(n int) [][2]int {
	rng := rand.New(rand.NewPCG(1, 2))
	pairs := make([][2]int, n)
	for i := range pairs {
		pairs[i] = [2]int{rng.IntN(n), rng.IntN(n)}
	}
	return pairs
}

func BenchmarkNewDisjointSet(b *testing.B) {
	for _, size := range benchSizes {
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = NewDisjointSet(size.n)
			}
		})
	}
}

func BenchmarkDisjointSetUnion(b *testing.B) {
	for _, size := range benchSizes {
		pairs := randomPairs(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				s, _ := NewDisjointSet(size.n)
				for _, p := range pairs {
					_, _ = s.Union(p[0], p[1])
				}
			}
		})
	}
}

func BenchmarkDisjointSetFind(b *testing.B) {
	for _, size := range benchSizes {
		pairs := randomPairs(size.n)
		s, _ := NewDisjointSet(size.n)
		for _, p := range pairs[:size.n/2] {
			_, _ = s.Union(p[0], p[1])
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				for x := range size.n {
					_, _ = s.Find(x)
				}
			}
		})
	}
}

func BenchmarkDisjointSetConnected(b *testing.B) {
	for _, size := range benchSizes {
		pairs := randomPairs(size.n)
		s, _ := NewDisjointSet(size.n)
		for _, p := range pairs[:size.n/2] {
			_, _ = s.Union(p[0], p[1])
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				for _, p := range pairs {
					_, _ = s.Connected(p[0], p[1])
				}
			}
		})
	}
}

func BenchmarkDisjointSetMakeSet(b *testing.B) {
	for _, size := range benchSizes {
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				var s DisjointSet
				for range size.n {
					s.MakeSet()
				}
			}
		})
	}
}

func BenchmarkDisjointSetCount(b *testing.B) {
	s, _ := NewDisjointSet(1000)
	b.ReportAllocs()
	for b.Loop() {
		_ = s.Count()
	}
}

func BenchmarkDisjointSetLen(b *testing.B) {
	s, _ := NewDisjointSet(1000)
	b.ReportAllocs()
	for b.Loop() {
		_ = s.Len()
	}
}
