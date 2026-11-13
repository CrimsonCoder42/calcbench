package renal

import (
	"math"
	"testing"
)

// Levey 2006 abstract equation, HAND-DERIVED: 175 × Scr^−1.154 × age^−0.203 × 1.212 [black] × 0.742 [female].
// Tolerance 1e-12 relative (same arithmetic written out).
func TestEGFRMDRD_equation(t *testing.T) {
	tests := []struct {
		scr, age      float64
		female, black bool
	}{
		{1.0, 60, true, false}, {1.0, 60, false, false}, {2.3, 45, false, true}, {0.6, 18, true, true}, {5, 90, false, false},
	}
	for _, tc := range tests {
		want := 175 * math.Pow(tc.scr, -1.154) * math.Pow(tc.age, -0.203)
		if tc.black {
			want *= 1.212
		}
		if tc.female {
			want *= 0.742
		}
		got, err := EGFRMDRD(tc.scr, tc.age, tc.female, tc.black)
		if err != nil || !relNear(got, want, 1e-12) {
			t.Errorf("%+v: (%v, %v), want %v", tc, got, err, want)
		}
		si, err := EGFRMDRDSI(umolOf(tc.scr), tc.age, tc.female, tc.black)
		if err != nil || !relNear(si, got, 1e-12) {
			t.Errorf("SI %+v: (%v, %v), want %v", tc, si, err, got)
		}
	}
}

// Doc comment example: Scr 1.0, 60 y, female → 56.5558 (HAND-DERIVED, 175 × 60^−0.203 × 0.742).
func TestEGFRMDRD_example(t *testing.T) {
	if got, err := EGFRMDRD(1.0, 60, true, false); err != nil || math.Abs(got-56.5558) > 1e-3 {
		t.Errorf("got (%v, %v), want 56.5558", got, err)
	}
}

// parity: ledger mdrd_gfr, source case 1: the app uses 175 and the same exponents and multipliers (male, 50 y,
// Scr 1.2, not black → 175 × 1.2^−1.154 × 50^−0.203).
func TestEGFRMDRD_parity(t *testing.T) {
	want := 175 * math.Pow(1.2, -1.154) * math.Pow(50, -0.203)
	if got, err := EGFRMDRD(1.2, 50, false, false); err != nil || !relNear(got, want, 1e-12) {
		t.Errorf("got (%v, %v), want %v", got, err, want)
	}
}

func TestEGFRMDRD_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"scr NaN", func() error { _, e := EGFRMDRD(nan, 50, false, false); return e }, ErrNotFinite},
		{"age +Inf", func() error { _, e := EGFRMDRD(1, pInf, false, false); return e }, ErrNotFinite},
		{"age -Inf", func() error { _, e := EGFRMDRD(1, nInf, false, false); return e }, ErrNotFinite},
		{"scr 0", func() error { _, e := EGFRMDRD(0, 50, false, false); return e }, ErrOutOfRange},
		{"age 17.9", func() error { _, e := EGFRMDRD(1, 17.9, false, false); return e }, ErrOutOfRange},
		{"overflow", func() error { _, e := EGFRMDRD(1e-300, 18, false, false); return e }, ErrDomain},
		{"SI NaN", func() error { _, e := EGFRMDRDSI(nan, 50, false, false); return e }, ErrNotFinite},
		{"SI -1", func() error { _, e := EGFRMDRDSI(-1, 50, false, false); return e }, ErrOutOfRange},
	})
	if _, err := EGFRMDRD(1, 18, false, false); err != nil {
		t.Errorf("age 18: %v", err)
	}
}

func TestEGFRMDRD_noAlloc(t *testing.T) {
	n := testing.AllocsPerRun(100, func() {
		sink, _ = EGFRMDRD(1.2, 50, true, false)
		sink, _ = EGFRMDRDSI(106, 50, true, false)
	})
	if n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkEGFRMDRD(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sink, _ = EGFRMDRD(1.2, 50, true, false)
	}
}

func BenchmarkEGFRMDRDSI(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sink, _ = EGFRMDRDSI(106, 50, true, false)
	}
}
