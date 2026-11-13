package stats

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Worked example: e-Handbook §7.2.6.2, the 12 resistivity values, p = 0.90. Printed: "Y(90) = 95.1981"
// (type 6, p(N+1) = 11.7), R7 "95.1957", R8 "95.1972". Re-derived exactly: 95.19807, 95.19568, 95.197243….
func TestPercentile_eHandbook(t *testing.T) {
	for _, tt := range []struct {
		method QuantileMethod
		exact  float64
		print  float64
	}{
		{HyndmanFan6, 95.19807, 95.1981},
		{HyndmanFan7, 95.19568, 95.1957},
		{HyndmanFan8, 95.19724333333333, 95.1972},
	} {
		got, err := Percentile(resistivity, 0.9, tt.method)
		// Tolerance 1e-14 relative against the exact re-derivation (rounding of stored decimals), and
		// 5e-5 absolute against the printed 4-decimal value.
		if err != nil || !near(got, tt.exact, 1e-14) || math.Abs(got-tt.print) > 5e-5 {
			t.Errorf("Percentile(resistivity, 0.9, type %d) = %v, %v; want %v (printed %v)", tt.method, got, err, tt.exact, tt.print)
		}
	}
}

// hfTable is 3, 1, 4, 1, 5, 9, 2, 6, 5, 3 (n = 10) under all nine types, computed by hand from
// Hyndman and Fan's eq. (1) in exact rational arithmetic. p = 0.3 makes np = 3 exactly, which float64
// computes as 3.0000000000000004: types 1–3 must still treat g as 0.
var hfTable = []struct {
	p    float64
	want [9]float64
}{
	{0, [9]float64{1, 1, 1, 1, 1, 1, 1, 1, 1}},
	{0.1, [9]float64{1, 1, 1, 1, 1, 1, 1, 1, 1}},
	{0.25, [9]float64{2, 2, 1, 1.5, 2, 1.75, 2.25, 1.9166666666666667, 1.9375}},
	{0.3, [9]float64{2, 2.5, 2, 2, 2.5, 2.3, 2.7, 2.433333333333333, 2.45}},
	{0.5, [9]float64{3, 3.5, 3, 3, 3.5, 3.5, 3.5, 3.5, 3.5}},
	{0.75, [9]float64{5, 5, 5, 5, 5, 5.25, 5, 5.083333333333333, 5.0625}},
	{0.9, [9]float64{6, 7.5, 6, 6, 7.5, 8.7, 6.3, 7.9, 7.8}},
	{1, [9]float64{9, 9, 9, 9, 9, 9, 9, 9, 9}},
}

func TestPercentile_allTypes(t *testing.T) {
	data := []float64{3, 1, 4, 1, 5, 9, 2, 6, 5, 3}
	for _, row := range hfTable {
		for i, want := range row.want {
			method := QuantileMethod(i + 1)
			// Tolerance 1e-15: one interpolation.
			if got, err := Percentile(data, row.p, method); err != nil || !near(got, want, 1e-15) {
				t.Errorf("Percentile(p=%v, type %d) = %v, %v; want %v", row.p, method, got, err, want)
			}
		}
	}
}

func TestPercentile_cases(t *testing.T) {
	m := math.MaxFloat64
	for _, tt := range []struct {
		name   string
		values []float64
		p      float64
		want   float64
	}{
		{"one element", []float64{4}, 0.37, 4},
		{"two elements, type 7 midpoint", []float64{1, 3}, 0.5, 2},
		{"ties straddle the rank", []float64{1, 2, 2, 2, 9}, 0.6, 2},
		{"negative", []float64{-4, -2}, 0.25, -3.5},
		{"extremes, no overflow", []float64{-m, m}, 0.5, 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Percentile(tt.values, tt.p, DefaultQuantileMethod); err != nil || got != tt.want {
				t.Errorf("Percentile(%v, %v) = %v, %v; want %v", tt.values, tt.p, got, err, tt.want)
			}
		})
	}
}

func TestPercentile_errors(t *testing.T) {
	for _, tt := range []struct {
		name   string
		values []float64
		p      float64
		method QuantileMethod
		want   error
	}{
		{"empty", nil, 0.5, HyndmanFan7, ErrEmpty},
		{"NaN value", []float64{math.NaN()}, 0.5, HyndmanFan7, ErrNotFinite},
		{"p below 0", []float64{1}, -0.01, HyndmanFan7, ErrOutOfRange},
		{"p above 1", []float64{1}, 1.01, HyndmanFan7, ErrOutOfRange},
		{"p NaN", []float64{1}, math.NaN(), HyndmanFan7, ErrOutOfRange},
		{"method 0", []float64{1}, 0.5, 0, ErrOutOfRange},
		{"method 10", []float64{1}, 0.5, 10, ErrOutOfRange},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Percentile(tt.values, tt.p, tt.method); !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("Percentile = %v, %v; want 0, %v", got, err, tt.want)
			}
		})
	}
}

// The oracle sorts a copy and applies type 7 (R's default) directly: x(1 + (n−1)p) interpolated.
func TestPercentile_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		n := 1 + rng.IntN(300)
		values := randomFloats(rng, n)
		for i := range values { // force ties
			if rng.IntN(4) == 0 {
				values[i] = values[rng.IntN(n)]
			}
		}
		p := rng.Float64()
		s := slices.Sorted(slices.Values(values))
		h := float64(n-1) * p
		lo := int(math.Floor(h))
		want := s[lo]
		if lo+1 < n {
			want += (h - float64(lo)) * (s[lo+1] - s[lo])
		}
		// Tolerance 1e-14: one interpolation; the position may round differently.
		if got, err := Percentile(values, p, HyndmanFan7); err != nil || !near(got, want, 1e-14) {
			t.Fatalf("trial %d: Percentile(n=%d, p=%v) = %v, %v; want %v", trial, n, p, got, err, want)
		}
	}
}

func TestPercentile_doesNotModifyInput(t *testing.T) {
	values := slices.Clone(resistivity)
	if _, err := Percentile(values, 0.9, HyndmanFan6); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(values, resistivity) {
		t.Error("Percentile modified its input")
	}
}

func TestPercentile_allocations(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { benchSink, _ = Percentile(michelso, 0.9, HyndmanFan7) }); a != 1 {
		t.Errorf("Percentile allocates %v times per call, want 1 (the copy)", a)
	}
}

func BenchmarkPercentile(b *testing.B) {
	for _, size := range benchSizes {
		inputs := benchInputs(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; b.Loop(); i++ {
				benchSink, _ = Percentile(inputs[i&7], 0.9, HyndmanFan7)
			}
		})
	}
}
