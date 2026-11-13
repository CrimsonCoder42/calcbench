package trig

import (
	"math"
	"testing"
)

// Worked example, hand-derived from DLMF 4.37.8 (arcsech z = arccosh(1/z)): arcsech 1 = 0;
// arcsech ½ = arccosh 2 = ln(2 + √3) = 1.3169578969248166.
func TestAsech_dlmfExample(t *testing.T) {
	for _, tc := range []struct{ x, want, tol float64 }{{1, 0, 0}, {0.5, 1.3169578969248166, ulps4}} {
		if got, err := Asech(tc.x); err != nil || !near(got, tc.want, tc.tol) {
			t.Errorf("Asech(%v) = (%v, %v), want %v", tc.x, got, err, tc.want)
		}
	}
}

func TestAsech_cases(t *testing.T) {
	tests := []struct {
		name    string
		x, want float64
		tol     float64
	}{
		// 0.9 is 0.900000000000000022…; ln((1 + √(1 − x²))/x) to 50 digits (Python decimal) is
		// 0.46714530810326196152…; acosh(1/x) would give …216 (3 ulps off).
		{"x = 0.9, one ulp", 0.9, 0.46714530810326196, 0x1p-52},
		// Below 2^-27, A&S 4.6.24 is ln(2/x) = ln 2 − ln x to within x², far below an ulp.
		{"smallest subnormal", 5e-324, math.Ln2 - math.Log(5e-324), ulps4},
		{"tiny", 1e-300, math.Ln2 - math.Log(1e-300), ulps4},
		{"just below 1", math.Nextafter(1, 0), math.Sqrt(2 * 0x1p-53), 1e-7},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := Asech(tc.x); err != nil || !relNear(got, tc.want, tc.tol) {
				t.Errorf("Asech(%v) = (%v, %v), want %v", tc.x, got, err, tc.want)
			}
		})
	}
}

func TestAsech_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { _, err := Asech(math.NaN()); return err }, ErrNotFinite},
		{"+Inf", func() error { _, err := Asech(math.Inf(1)); return err }, ErrNotFinite},
		{"zero", func() error { _, err := Asech(0); return err }, ErrDomain},
		{"negative", func() error { _, err := Asech(-0.5); return err }, ErrDomain},
		{"above 1", func() error { _, err := Asech(math.Nextafter(1, 2)); return err }, ErrDomain},
	})
}

// Oracle: sech(Asech x) = x for x in [0.01, 0.99], where the round trip is well conditioned:
// tolerance 1e-14 relative. Both branches (x ≤ ½ and x > ½) are covered.
func TestAsech_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		x := 0.01 + rng.Float64()*0.98
		y, err := Asech(x)
		if err != nil || y < 0 || !relNear(1/math.Cosh(y), x, 1e-14) {
			t.Fatalf("Asech(%v) = (%v, %v)", x, y, err)
		}
	}
}

func TestAsech_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = Asech(0.7) }); n != 0 {
		t.Errorf("Asech allocates %v times, want 0", n)
	}
}

func BenchmarkAsech(b *testing.B) {
	b.ReportAllocs()
	x := 0.7
	for b.Loop() {
		sink, _ = Asech(x)
	}
}
