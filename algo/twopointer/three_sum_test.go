package twopointer

import (
	"cmp"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// naiveThreeSum is the oracle: every triple of positions, exact int sums, deduplicated and sorted.
func naiveThreeSum(values []int8, target int) [][3]int8 {
	var out [][3]int8
	n := len(values)
	for a := 0; a < n; a++ {
		for b := a + 1; b < n; b++ {
			for c := b + 1; c < n; c++ {
				if int(values[a])+int(values[b])+int(values[c]) == target {
					tr := []int8{values[a], values[b], values[c]}
					slices.Sort(tr)
					out = append(out, [3]int8{tr[0], tr[1], tr[2]})
				}
			}
		}
	}
	slices.SortFunc(out, compareTriple)
	return slices.Compact(out)
}

func compareTriple[T cmp.Ordered](x, y [3]T) int {
	for k := range 3 {
		if c := cmp.Compare(x[k], y[k]); c != 0 {
			return c
		}
	}
	return 0
}

// The 3SUM problem of Gajentaan & Overmars (1995) on a hand-made set of eight integers; the four
// zero-sum triples are hand-derived.
func TestThreeSum_worked(t *testing.T) {
	got := ThreeSum(nil, []int{30, -40, -20, -10, 40, 0, 10, 5}, 0)
	want := [][3]int{{-40, 0, 40}, {-40, 10, 30}, {-20, -10, 30}, {-10, 0, 10}}
	if !slices.Equal(got, want) {
		t.Errorf("ThreeSum(eight ints, 0) = %v, want %v", got, want)
	}
}

func TestThreeSum_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int64
		target int64
		want   [][3]int64
	}{
		{"empty", nil, 0, nil},
		{"two elements", []int64{1, -1}, 0, nil},
		{"exactly three", []int64{3, 1, 2}, 6, [][3]int64{{1, 2, 3}}},
		{"all zeros: once", []int64{0, 0, 0, 0}, 0, [][3]int64{{0, 0, 0}}},
		{"duplicates collapse", []int64{2, -1, -1, 0, 1, -3}, 0, [][3]int64{{-3, 1, 2}, {-1, -1, 2}, {-1, 0, 1}}},
		{"none", []int64{1, 2, 3, 4}, 100, nil},
		{"non-zero target", []int64{5, 1, 4, 2, 3}, 9, [][3]int64{{1, 3, 5}, {2, 3, 4}}},
		{"extremes, exact", []int64{math.MaxInt64, math.MaxInt64, math.MinInt64}, math.MaxInt64 - 1, [][3]int64{{math.MinInt64, math.MaxInt64, math.MaxInt64}}},
		{"overflow does not invent", []int64{math.MaxInt64, math.MaxInt64, 2}, 0, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ThreeSum(nil, tt.values, tt.target)
			if !slices.Equal(got, tt.want) {
				t.Errorf("ThreeSum(%v, %d) = %v, want %v", tt.values, tt.target, got, tt.want)
			}
		})
	}
	if got := ThreeSum(nil, []int8{100, 100, 100}, 44); len(got) != 0 {
		t.Errorf("int8 300 wraps to 44 but must not match: %v", got)
	}
}

func TestThreeSum_appendsToDst(t *testing.T) {
	got := ThreeSum([][3]int{{9, 9, 9}}, []int{1, 2, 3}, 6)
	if want := [][3]int{{9, 9, 9}, {1, 2, 3}}; !slices.Equal(got, want) {
		t.Errorf("ThreeSum(dst, 1 2 3, 6) = %v, want %v", got, want)
	}
}

// Oracle: random int8 data with many duplicates against the O(n³) enumeration; input never modified.
func TestThreeSum_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		values := make([]int8, rng.IntN(14))
		spread := 1 + rng.IntN(256)
		for i := range values {
			values[i] = int8(rng.IntN(spread) - spread/2)
		}
		before := slices.Clone(values)
		target := int8(rng.IntN(21) - 10)
		got := ThreeSum(nil, values, target)
		if want := naiveThreeSum(values, int(target)); !slices.Equal(got, want) {
			t.Fatalf("trial %d: ThreeSum(%v, %d) = %v, want %v", trial, values, target, got, want)
		}
		if !slices.Equal(values, before) {
			t.Fatalf("trial %d: input modified", trial)
		}
	}
}

// O(n²): sizes 1e2, 1e3, 1e4 finish within a second each.
func BenchmarkThreeSum(b *testing.B) {
	for _, size := range []struct {
		name string
		n    int
	}{{"1e2", 100}, {"1e3", 1_000}, {"1e4", 10_000}} {
		rng := rand.New(rand.NewPCG(1, 2))
		values := make([]int64, size.n)
		for i := range values {
			values[i] = int64(rng.IntN(1 << 30))
		}
		var dst [][3]int64
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				dst = ThreeSum(dst[:0], values, 1<<30)
			}
			benchSink = len(dst)
		})
	}
}
