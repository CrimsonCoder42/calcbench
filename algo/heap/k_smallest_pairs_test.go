package heap

import (
	"cmp"
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

func ip(i, j int) IndexPair { return IndexPair{i, j} }

// allPairsSorted is the slow oracle: every pair, sorted by (sum, i, j).
func allPairsSorted(a, b []int) []IndexPair {
	out := []IndexPair{}
	for i := range a {
		for j := range b {
			out = append(out, ip(i, j))
		}
	}
	slices.SortFunc(out, func(p, q IndexPair) int {
		return cmp.Or(cmp.Compare(a[p.I]+b[p.J], a[q.I]+b[q.J]), cmp.Compare(p.I, q.I), cmp.Compare(p.J, q.J))
	})
	return out
}

// Worked example written for this repo (the sources print none). a = 0 3 10, b = 1 2 20. The nine sums
// by hand: (0,0)=1 (0,1)=2 (1,0)=4 (1,1)=5 (2,0)=11 (2,1)=12 (0,2)=20 (1,2)=23 (2,2)=30.
func TestKSmallestPairs_workedExample(t *testing.T) {
	got, err := KSmallestPairs(nil, []int{0, 3, 10}, []int{1, 2, 20}, 5)
	want := []IndexPair{ip(0, 0), ip(0, 1), ip(1, 0), ip(1, 1), ip(2, 0)}
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("KSmallestPairs = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestKSmallestPairs_cases(t *testing.T) {
	tests := []struct {
		name string
		a, b []int
		k    int
		want []IndexPair
	}{
		{"empty a", nil, []int{1}, 3, nil},
		{"empty b", []int{1}, nil, 3, nil},
		{"k 0", []int{1}, []int{1}, 0, nil},
		{"one by one", []int{4}, []int{5}, 1, []IndexPair{ip(0, 0)}},
		{"k above all pairs", []int{1, 2}, []int{1}, 9, []IndexPair{ip(0, 0), ip(1, 0)}},
		{"ties ordered by i then j", []int{1, 1}, []int{2, 2}, 4, []IndexPair{ip(0, 0), ip(0, 1), ip(1, 0), ip(1, 1)}},
		{"cross tie", []int{0, 5}, []int{5, 10}, 3, []IndexPair{ip(0, 0), ip(0, 1), ip(1, 0)}},
		{"one long row wins", []int{0, 100}, []int{0, 1, 2, 3}, 4, []IndexPair{ip(0, 0), ip(0, 1), ip(0, 2), ip(0, 3)}},
		{"negative and zero", []int{-5, 0}, []int{-1, 0}, 3, []IndexPair{ip(0, 0), ip(0, 1), ip(1, 0)}},
		{"extremes cancel", []int{math.MinInt}, []int{math.MaxInt}, 1, []IndexPair{ip(0, 0)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := KSmallestPairs(nil, tt.a, tt.b, tt.k)
			if err != nil || !slices.Equal(got, tt.want) {
				t.Errorf("KSmallestPairs(%v, %v, %d) = (%v, %v), want (%v, nil)", tt.a, tt.b, tt.k, got, err, tt.want)
			}
		})
	}
}

// Problem-site examples, in our own words: two short ascending lists and a small k; a list with a
// repeated value, where both copies pair with the same partner; and k larger than the number of pairs.
func TestKSmallestPairs_problemSiteShapes(t *testing.T) {
	got, _ := KSmallestPairs(nil, []int{2, 8, 9}, []int{1, 3, 5}, 3)
	if want := []IndexPair{ip(0, 0), ip(0, 1), ip(0, 2)}; !slices.Equal(got, want) {
		t.Errorf("short lists = %v, want %v", got, want)
	}
	got, _ = KSmallestPairs(nil, []int{1, 1, 4}, []int{1, 2}, 2)
	if want := []IndexPair{ip(0, 0), ip(1, 0)}; !slices.Equal(got, want) {
		t.Errorf("repeated value = %v, want %v", got, want)
	}
	got, _ = KSmallestPairs(nil, []int{1, 2}, []int{3}, 10)
	if want := []IndexPair{ip(0, 0), ip(1, 0)}; !slices.Equal(got, want) {
		t.Errorf("k too large = %v, want %v", got, want)
	}
}

func TestKSmallestPairs_floats(t *testing.T) {
	got, err := KSmallestPairs(nil, []float64{-0.5, 0.25}, []float64{0.5, 1}, 4)
	if want := []IndexPair{ip(0, 0), ip(0, 1), ip(1, 0), ip(1, 1)}; err != nil || !slices.Equal(got, want) {
		t.Errorf("floats = (%v, %v), want %v", got, err, want)
	}
}

func TestKSmallestPairs_appendsToDst(t *testing.T) {
	dst := []IndexPair{ip(9, 9)}
	got, err := KSmallestPairs(dst, []int{1}, []int{1, 2}, 2)
	if want := []IndexPair{ip(9, 9), ip(0, 0), ip(0, 1)}; err != nil || !slices.Equal(got, want) {
		t.Errorf("append = (%v, %v), want %v", got, err, want)
	}
}

func TestKSmallestPairs_errors(t *testing.T) {
	dst := []IndexPair{ip(7, 7)}
	if got, err := KSmallestPairs(dst, []int{1}, []int{1}, -1); !errors.Is(err, ErrNegativeK) || len(got) != 1 {
		t.Errorf("k -1 = (%v, %v), want ErrNegativeK and dst unchanged", got, err)
	}
	for _, bad := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if got, err := KSmallestPairs(dst, []float64{bad}, []float64{1}, 1); !errors.Is(err, ErrNotFinite) || len(got) != 1 {
			t.Errorf("a = [%v]: (%v, %v), want ErrNotFinite", bad, got, err)
		}
		if got, err := KSmallestPairs(dst, []float64{0, 1}, []float64{1, bad}, 4); !errors.Is(err, ErrNotFinite) || len(got) != 1 {
			t.Errorf("b = [1 %v]: (%v, %v), want ErrNotFinite and dst unchanged", bad, got, err)
		}
	}
	if _, err := KSmallestPairs(nil, []int8{100}, []int8{100}, 1); !errors.Is(err, ErrOverflow) {
		t.Errorf("int8 100+100: err = %v, want ErrOverflow", err)
	}
	if _, err := KSmallestPairs(nil, []int{math.MinInt}, []int{-1}, 1); !errors.Is(err, ErrOverflow) {
		t.Errorf("MinInt-1: err = %v, want ErrOverflow", err)
	}
	if _, err := KSmallestPairs(nil, []uint{math.MaxUint}, []uint{1}, 1); !errors.Is(err, ErrOverflow) {
		t.Errorf("MaxUint+1: err = %v, want ErrOverflow", err)
	}
}

func TestKSmallestPairs_doesNotModifyInput(t *testing.T) {
	a, b := []int{1, 4, 9}, []int{2, 3}
	if _, err := KSmallestPairs(nil, a, b, 6); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(a, []int{1, 4, 9}) || !slices.Equal(b, []int{2, 3}) {
		t.Errorf("inputs modified: %v %v", a, b)
	}
}

func TestKSmallestPairs_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		a := make([]int, rng.IntN(12))
		b := make([]int, rng.IntN(12))
		for i := range a {
			a[i] = rng.IntN(30) - 10
		}
		for i := range b {
			b[i] = rng.IntN(30) - 10
		}
		slices.Sort(a)
		slices.Sort(b)
		k := rng.IntN(len(a)*len(b) + 4)
		all := allPairsSorted(a, b)
		want := all[:min(k, len(all))]
		got, err := KSmallestPairs(nil, a, b, k)
		if err != nil || !slices.Equal(got, want) && !(len(got) == 0 && len(want) == 0) {
			t.Fatalf("trial %d: KSmallestPairs(%v, %v, %d) = (%v, %v), want %v", trial, a, b, k, got, err, want)
		}
	}
}

func BenchmarkKSmallestPairs(b *testing.B) {
	for _, size := range benchSizes {
		x := randomInts(size.n)
		y := randomInts(size.n)
		slices.Sort(x)
		slices.Sort(y)
		dst := make([]IndexPair, 0, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				dst, _ = KSmallestPairs(dst[:0], x, y, size.n)
			}
		})
	}
}
