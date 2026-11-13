package window

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// naiveWindowExtremes is the slow obvious oracle shared by the RollingMin and RollingMax tests: rescan
// each window, O(n·w).
func naiveWindowExtremes(values []int, w int, wantMax bool) []int {
	var out []int
	for end := w; end <= len(values); end++ {
		window := values[end-w : end]
		if wantMax {
			out = append(out, slices.Max(window))
		} else {
			out = append(out, slices.Min(window))
		}
	}
	return out
}

// Worked example on the data of NIST/SEMATECH e-Handbook §6.4.2.1 (supplier data), N = 3. Expected
// minima are hand-derived from the source's data.
func TestRollingMin_nistExample(t *testing.T) {
	want := []float64{8, 8, 9, 9, 9, 7, 7, 7, 9, 9}
	got, err := RollingMin(nil, nistSupplier, 3)
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("RollingMin(nistSupplier, 3) = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestRollingMin_cases(t *testing.T) {
	tests := []struct {
		name      string
		values    []int
		windowLen int
		want      []int
	}{
		{"empty", nil, 1, nil},
		{"window longer than the data", []int{1, 2}, 3, nil},
		{"one element", []int{4}, 1, []int{4}},
		{"window 1 is the identity", []int{3, -1, 2}, 1, []int{3, -1, 2}},
		{"window equals the length", []int{3, -1, 2}, 3, []int{-1}},
		{"two elements", []int{5, 2}, 2, []int{2}},
		{"duplicates", []int{2, 2, 1, 1, 2, 2}, 2, []int{2, 1, 1, 1, 2}},
		{"already sorted", []int{1, 2, 3, 4, 5}, 3, []int{1, 2, 3}},
		{"reversed", []int{5, 4, 3, 2, 1}, 3, []int{3, 2, 1}},
		{"negative and zero", []int{0, -3, 0, -1}, 2, []int{-3, -3, -1}},
		{"extremes of the type", []int{math.MaxInt, math.MinInt, math.MaxInt}, 2, []int{math.MinInt, math.MinInt}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RollingMin(nil, tt.values, tt.windowLen)
			if err != nil || len(got) != len(tt.want) || !slices.Equal(got, tt.want) {
				t.Errorf("RollingMin(%v, %d) = (%v, %v), want (%v, nil)", tt.values, tt.windowLen, got, err, tt.want)
			}
		})
	}
}

func TestRollingMin_floats(t *testing.T) {
	values := []float64{1, math.NaN(), 2, math.Inf(-1), 3, math.Inf(1)}
	got, err := RollingMin(nil, values, 2)
	// cmp.Compare order: NaN is the least value, so it is the minimum of every window holding it.
	want := []float64{math.NaN(), math.NaN(), math.Inf(-1), math.Inf(-1), 3}
	if err != nil || len(got) != len(want) {
		t.Fatalf("RollingMin(floats, 2) = (%v, %v), want %v", got, err, want)
	}
	for i := range want {
		if !(got[i] == want[i] || (math.IsNaN(got[i]) && math.IsNaN(want[i]))) {
			t.Errorf("RollingMin(floats, 2)[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}

func TestRollingMin_strings(t *testing.T) {
	got, err := RollingMin(nil, []string{"pear", "fig", "plum", "kiwi"}, 2)
	if want := []string{"fig", "fig", "kiwi"}; err != nil || !slices.Equal(got, want) {
		t.Errorf("RollingMin(fruit, 2) = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestRollingMin_errors(t *testing.T) {
	for _, w := range []int{0, -1, math.MinInt} {
		got, err := RollingMin([]int{9}, []int{1, 2}, w)
		if !errors.Is(err, ErrWindowLen) || !slices.Equal(got, []int{9}) {
			t.Errorf("RollingMin([9], [1 2], %d) = (%v, %v), want ([9], ErrWindowLen)", w, got, err)
		}
	}
}

// The deque is a stack array up to stackRingMax and one heap allocation above it.
func TestRollingMin_allocations(t *testing.T) {
	values := make([]int, 1000)
	for i := range values {
		values[i] = (i * 7919) % 1000
	}
	buf := make([]int, 0, len(values))
	for _, tc := range []struct {
		w    int
		want float64
	}{{1, 0}, {stackRingMax, 0}, {stackRingMax + 1, 1}} {
		allocs := testing.AllocsPerRun(20, func() { buf, _ = RollingMin(buf[:0], values, tc.w) })
		if allocs != tc.want {
			t.Errorf("RollingMin(w=%d) with a reused buffer allocates %v times per call, want %v", tc.w, allocs, tc.want)
		}
	}
}

// Oracle: rescan every window. Lengths and windows cross stackRingMax so both deque storages and the
// ring wrap-around are exercised; small value ranges make ties common.
func TestRollingMin_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		n := rng.IntN(300)
		values := make([]int, n)
		for i := range values {
			values[i] = rng.IntN(1 + rng.IntN(50))
		}
		w := 1 + rng.IntN(n+2)
		before := slices.Clone(values)
		got, err := RollingMin(nil, values, w)
		want := naiveWindowExtremes(values, w, false)
		if err != nil || !slices.Equal(got, want) {
			t.Fatalf("trial %d: RollingMin(%v, %d) = (%v, %v), want %v", trial, values, w, got, err, want)
		}
		if !slices.Equal(values, before) {
			t.Fatalf("trial %d: RollingMin modified its input", trial)
		}
	}
}

func BenchmarkRollingMin(b *testing.B) {
	for _, size := range benchSizes {
		values := benchFloats(size.n)
		buf := make([]float64, 0, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				buf, _ = RollingMin(buf[:0], values, benchWindow)
			}
		})
	}
}
