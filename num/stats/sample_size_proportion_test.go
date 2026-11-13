package stats

import (
	"errors"
	"math"
	"testing"
)

// z₀.₉₇₅ = 1.959963984540054 (the 0.975 standard normal quantile). Hand computations from Cochran §4.4:
// 95 %, p = 0.5, d = 0.05: n₀ = z²·0.25/0.0025 = 384.1458820694124; N = 1000: n = n₀/(1 + (n₀ − 1)/1000)
// = 277.7335…; p = 0 needs no sample.
func TestSampleSizeProportion_cases(t *testing.T) {
	const z = 1.959963984540054
	n0 := z * z * 0.25 / 0.0025
	tests := []struct {
		name       string
		conf, p, d float64
		pop        int64
		want       float64
	}{
		{"infinite", 0.95, 0.5, 0.05, 0, n0},
		{"finite 1000", 0.95, 0.5, 0.05, 1000, n0 / (1 + (n0-1)/1000)},
		{"p = 0", 0.95, 0, 0.05, 0, 0},
		{"p = 0.2, d = 0.03", 0.95, 0.2, 0.03, 0, z * z * 0.16 / 0.0009},
	}
	for _, tc := range tests {
		got, err := SampleSizeProportion(tc.conf, tc.p, tc.d, tc.pop)
		if err != nil || math.Abs(got-tc.want) > 1e-9*math.Max(1, tc.want) {
			t.Errorf("%s: got %v, %v; want %v", tc.name, got, err, tc.want)
		}
	}
	if got, _ := SampleSizeProportion(0.95, 0.5, 0.05, 0); math.Ceil(got) != 385 {
		t.Errorf("ceil = %v, want 385", math.Ceil(got))
	}
}

// Hand: 95 %, p = 0.6, n = 100: z·√(0.0024) = 0.0960182…; N = 1000 multiplies by √(900/999).
// Round trip: the margin at the finite-population sample size is the requested d up to Cochran's
// N − 1 ≈ N simplification (within 0.1 % here).
func TestMarginOfErrorProportion_cases(t *testing.T) {
	const z = 1.959963984540054
	want := z * math.Sqrt(0.24/100)
	if got, err := MarginOfErrorProportion(0.95, 0.6, 100, 0); err != nil || math.Abs(got-want) > 1e-12 {
		t.Errorf("infinite: %v, %v; want %v", got, err, want)
	}
	if got, err := MarginOfErrorProportion(0.95, 0.6, 100, 1000); err != nil || math.Abs(got-want*math.Sqrt(900.0/999)) > 1e-12 {
		t.Errorf("finite: %v, %v", got, err)
	}
	if got, err := MarginOfErrorProportion(0.95, 0.5, 50, 50); err != nil || got != 0 {
		t.Errorf("census: %v, %v; want 0", got, err)
	}
	n, _ := SampleSizeProportion(0.95, 0.5, 0.05, 10000)
	if d, err := MarginOfErrorProportion(0.95, 0.5, int64(math.Round(n)), 10000); err != nil || math.Abs(d-0.05) > 5e-5 {
		t.Errorf("round trip: %v, %v", d, err)
	}
}

func TestSampleSizeProportion_errors(t *testing.T) {
	e := func(_ float64, err error) error { return err }
	tests := []struct {
		name      string
		err, want error
	}{
		{"conf NaN", e(SampleSizeProportion(math.NaN(), 0.5, 0.05, 0)), ErrNotFinite},
		{"conf 1", e(SampleSizeProportion(1, 0.5, 0.05, 0)), ErrOutOfRange},
		{"conf 0", e(SampleSizeProportion(0, 0.5, 0.05, 0)), ErrOutOfRange},
		{"p NaN", e(SampleSizeProportion(0.95, math.NaN(), 0.05, 0)), ErrNotFinite},
		{"p > 1", e(SampleSizeProportion(0.95, 1.5, 0.05, 0)), ErrOutOfRange},
		{"d 0", e(SampleSizeProportion(0.95, 0.5, 0, 0)), ErrOutOfRange},
		{"d Inf", e(SampleSizeProportion(0.95, 0.5, math.Inf(1), 0)), ErrNotFinite},
		{"population < 0", e(SampleSizeProportion(0.95, 0.5, 0.05, -1)), ErrOutOfRange},
		{"MOE conf Inf", e(MarginOfErrorProportion(math.Inf(1), 0.5, 10, 0)), ErrNotFinite},
		{"MOE p NaN", e(MarginOfErrorProportion(0.95, math.NaN(), 10, 0)), ErrNotFinite},
		{"MOE n 0", e(MarginOfErrorProportion(0.95, 0.5, 0, 0)), ErrOutOfRange},
		{"MOE n > N", e(MarginOfErrorProportion(0.95, 0.5, 11, 10)), ErrOutOfRange},
		{"MOE N 1", e(MarginOfErrorProportion(0.95, 0.5, 1, 1)), ErrOutOfRange},
	}
	for _, tc := range tests {
		if !errors.Is(tc.err, tc.want) {
			t.Errorf("%s: err = %v, want %v", tc.name, tc.err, tc.want)
		}
	}
}

func BenchmarkSampleSizeProportion(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = SampleSizeProportion(0.95, 0.5, 0.05, 1000)
	}
}

func BenchmarkMarginOfErrorProportion(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = MarginOfErrorProportion(0.95, 0.6, 100, 1000)
	}
}
