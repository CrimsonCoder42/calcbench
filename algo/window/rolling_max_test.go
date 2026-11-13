package window

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Worked example on the data of NIST/SEMATECH e-Handbook §6.4.2.1 (supplier data), N = 3. Expected
// maxima are hand-derived from the source's data.
func TestRollingMax_nistExample(t *testing.T) {
	want := []float64{9, 12, 12, 12, 12, 12, 13, 13, 13, 11}
	got, err := RollingMax(nil, nistSupplier, 3)
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("RollingMax(nistSupplier, 3) = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestRollingMax_cases(t *testing.T) {
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
		{"window equals the length", []int{3, -1, 2}, 3, []int{3}},
		{"two elements", []int{5, 2}, 2, []int{5}},
		{"duplicates", []int{1, 1, 2, 2, 1, 1}, 2, []int{1, 2, 2, 2, 1}},
		{"already sorted", []int{1, 2, 3, 4, 5}, 3, []int{3, 4, 5}},
		{"reversed", []int{5, 4, 3, 2, 1}, 3, []int{5, 4, 3}},
		{"negative and zero", []int{0, -3, 0, -1}, 2, []int{0, 0, 0}},
		{"extremes of the type", []int{math.MinInt, math.MaxInt, math.MinInt}, 2, []int{math.MaxInt, math.MaxInt}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RollingMax(nil, tt.values, tt.windowLen)
			if err != nil || len(got) != len(tt.want) || !slices.Equal(got, tt.want) {
				t.Errorf("RollingMax(%v, %d) = (%v, %v), want (%v, nil)", tt.values, tt.windowLen, got, err, tt.want)
			}
		})
	}
}

func TestRollingMax_floats(t *testing.T) {
	values := []float64{math.NaN(), math.NaN(), 1, math.Inf(1), -2}
	got, err := RollingMax(nil, values, 2)
	// cmp.Compare order: NaN is the least value, so it is the maximum only of an all-NaN window.
	want := []float64{math.NaN(), 1, math.Inf(1), math.Inf(1)}
	if err != nil || len(got) != len(want) {
		t.Fatalf("RollingMax(floats, 2) = (%v, %v), want %v", got, err, want)
	}
	for i := range want {
		if !(got[i] == want[i] || (math.IsNaN(got[i]) && math.IsNaN(want[i]))) {
			t.Errorf("RollingMax(floats, 2)[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}

func TestRollingMax_errors(t *testing.T) {
	for _, w := range []int{0, -1, math.MinInt} {
		got, err := RollingMax([]int{9}, []int{1, 2}, w)
		if !errors.Is(err, ErrWindowLen) || !slices.Equal(got, []int{9}) {
			t.Errorf("RollingMax([9], [1 2], %d) = (%v, %v), want ([9], ErrWindowLen)", w, got, err)
		}
	}
}

func TestRollingMax_noAllocationWithBuffer(t *testing.T) {
	values := benchFloats(1000)
	buf := make([]float64, 0, len(values))
	allocs := testing.AllocsPerRun(20, func() { buf, _ = RollingMax(buf[:0], values, benchWindow) })
	if allocs != 0 {
		t.Errorf("RollingMax with a reused buffer allocates %v times per call, want 0", allocs)
	}
}

func TestRollingMax_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		n := rng.IntN(300)
		values := make([]int, n)
		for i := range values {
			values[i] = rng.IntN(1 + rng.IntN(50))
		}
		w := 1 + rng.IntN(n+2)
		got, err := RollingMax(nil, values, w)
		want := naiveWindowExtremes(values, w, true)
		if err != nil || !slices.Equal(got, want) {
			t.Fatalf("trial %d: RollingMax(%v, %d) = (%v, %v), want %v", trial, values, w, got, err, want)
		}
	}
}

func BenchmarkRollingMax(b *testing.B) {
	for _, size := range benchSizes {
		values := benchFloats(size.n)
		buf := make([]float64, 0, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				buf, _ = RollingMax(buf[:0], values, benchWindow)
			}
		})
	}
}
