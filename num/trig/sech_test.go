package trig

import (
	"math"
	"testing"
)

// Worked example, hand-derived from DLMF 4.28.6 (sech z = 1/cosh z): sech 0 = 1;
// sech 1 = 2/(e + 1/e) = 0.6480542736638855.
func TestSech_dlmfExample(t *testing.T) {
	for _, tc := range []struct{ x, want, tol float64 }{{0, 1, 0}, {1, 0.6480542736638855, ulps4}, {-1, 0.6480542736638855, ulps4}} {
		if got, err := Sech(tc.x); err != nil || !relNear(got, tc.want, tc.tol) {
			t.Errorf("Sech(%v) = (%v, %v), want %v", tc.x, got, err, tc.want)
		}
	}
}

func TestSech_cases(t *testing.T) {
	tests := []struct {
		name    string
		x, want float64
	}{
		{"cosh overflows, result 0", 800, 0},
		{"largest float64", math.MaxFloat64, 0},
		{"negative large", -800, 0},
		{"tiny", 1e-300, 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := Sech(tc.x); err != nil || got != tc.want {
				t.Errorf("Sech(%v) = (%v, %v), want %v", tc.x, got, err, tc.want)
			}
		})
	}
}

func TestSech_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { _, err := Sech(math.NaN()); return err }, ErrNotFinite},
		{"+Inf", func() error { _, err := Sech(math.Inf(1)); return err }, ErrNotFinite},
		{"-Inf", func() error { _, err := Sech(math.Inf(-1)); return err }, ErrNotFinite},
	})
}

// Oracle: sech² x + tanh² x = 1 (A&S 4.5.17), for |x| ≤ 20, to a few ulps (1e-15).
func TestSech_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		x := (rng.Float64() - 0.5) * 40
		s, err := Sech(x)
		th := math.Tanh(x)
		if err != nil || s <= 0 || s > 1 || math.Abs(s*s+th*th-1) > 1e-15 {
			t.Fatalf("Sech(%v) = (%v, %v)", x, s, err)
		}
	}
}

func TestSech_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = Sech(0.7) }); n != 0 {
		t.Errorf("Sech allocates %v times, want 0", n)
	}
}

func BenchmarkSech(b *testing.B) {
	b.ReportAllocs()
	x := 0.7
	for b.Loop() {
		sink, _ = Sech(x)
	}
}
