package trig

import (
	"math"
	"testing"
)

// Worked example: A&S 4.6.25, arccoth x = ½ ln((x+1)/(x−1)); at x = 2 that is ½ ln 3 = 0.5493061443340549.
func TestAcoth_asExample(t *testing.T) {
	for _, tc := range []struct{ x, want float64 }{{2, 0.5493061443340549}, {-2, -0.5493061443340549}} {
		if got, err := Acoth(tc.x); err != nil || !relNear(got, tc.want, ulps4) {
			t.Errorf("Acoth(%v) = (%v, %v), want %v", tc.x, got, err, tc.want)
		}
	}
}

func TestAcoth_cases(t *testing.T) {
	tests := []struct {
		name    string
		x, want float64
	}{
		// x = 1 + 2^-52: (x+1)/(x−1) = 2^53 + 1, so the result is ½ ln(2^53 + 1) = 18.368400284838550755…
		// (Python decimal, 50 digits): exact rational reference for the claim that log1p keeps accuracy at 1⁺.
		{"just above 1", 1 + 0x1p-52, 18.36840028483855},
		{"just below −1", -1 - 0x1p-52, -18.36840028483855},
		// x = 3: ½ ln 2.
		{"three", 3, math.Ln2 / 2},
		{"largest float64", math.MaxFloat64, 1 / math.MaxFloat64},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := Acoth(tc.x); err != nil || !relNear(got, tc.want, ulps4) {
				t.Errorf("Acoth(%v) = (%v, %v), want %v", tc.x, got, err, tc.want)
			}
		})
	}
}

func TestAcoth_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { _, err := Acoth(math.NaN()); return err }, ErrNotFinite},
		{"+Inf", func() error { _, err := Acoth(math.Inf(1)); return err }, ErrNotFinite},
		{"pole at 1", func() error { _, err := Acoth(1); return err }, ErrDomain},
		{"pole at −1", func() error { _, err := Acoth(-1); return err }, ErrDomain},
		{"zero", func() error { _, err := Acoth(0); return err }, ErrDomain},
		{"inside (−1, 1)", func() error { _, err := Acoth(0.5); return err }, ErrDomain},
	})
}

// Oracle: coth(Acoth x) = x for 1.01 ≤ |x| ≤ 100, to 1e-12 relative (the round trip amplifies by up to
// about 1/(x²−1)·x); and Acoth x = atanh(1/x) where 1/x is exact (x a power of two), to a few ulps.
func TestAcoth_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		x := 1.01 + rng.Float64()*98.99
		if rng.IntN(2) == 0 {
			x = -x
		}
		y, err := Acoth(x)
		if err != nil || math.Signbit(y) != math.Signbit(x) || !relNear(1/math.Tanh(y), x, 1e-12) {
			t.Fatalf("Acoth(%v) = (%v, %v)", x, y, err)
		}
	}
	for e := 1; e < 60; e++ {
		x := math.Ldexp(1, e)
		if y, _ := Acoth(x); !relNear(y, math.Atanh(1/x), ulps4) {
			t.Fatalf("Acoth(2^%d) = %v, want %v", e, y, math.Atanh(1/x))
		}
	}
}

func TestAcoth_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = Acoth(2.5) }); n != 0 {
		t.Errorf("Acoth allocates %v times, want 0", n)
	}
}

func BenchmarkAcoth(b *testing.B) {
	b.ReportAllocs()
	x := 2.5
	for b.Loop() {
		sink, _ = Acoth(x)
	}
}
