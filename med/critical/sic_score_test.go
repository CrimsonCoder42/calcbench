package critical

import "testing"

// Iba 2017 Table 3 and footnotes: each item's edges on both sides; the SIC criterion. HAND-DERIVED.
func TestSICScore_table3(t *testing.T) {
	tests := []struct {
		inr, plt      float64
		sofa          int
		pi, pp, ps, t int
		sic           bool
	}{
		{1.2, 150, 0, 0, 0, 0, 0, false},
		{1.21, 150, 0, 1, 0, 0, 1, false},
		{1.4, 150, 0, 1, 0, 0, 1, false},
		{1.41, 150, 0, 2, 0, 0, 2, false},
		{1.0, 149.9, 0, 0, 1, 0, 1, false},
		{1.0, 100, 0, 0, 1, 0, 1, false},
		{1.0, 99.9, 0, 0, 2, 0, 2, false},
		{1.0, 200, 1, 0, 0, 1, 1, false},
		{1.0, 200, 2, 0, 0, 2, 2, false},
		{1.0, 200, 16, 0, 0, 2, 2, false},
		{1.5, 90, 3, 2, 2, 2, 6, true},
		{1.3, 120, 2, 1, 1, 2, 4, false}, // total 4 but PT + platelets = 2, not exceeding 2
		{1.5, 120, 1, 2, 1, 1, 4, true},  // PT + platelets = 3
		{1.5, 120, 0, 2, 1, 0, 3, false}, // total 3
	}
	for i, tc := range tests {
		r, err := SICScore(tc.inr, tc.plt, tc.sofa)
		if err != nil || r.PTINRPoints != tc.pi || r.PlateletPoints != tc.pp || r.SOFAPoints != tc.ps || r.Total != tc.t || r.SIC != tc.sic {
			t.Errorf("case %d: (%+v, %v), want %d %d %d %d %v", i, r, err, tc.pi, tc.pp, tc.ps, tc.t, tc.sic)
		}
	}
}

// parity: ledger sic_score, source case 1 (app takes INR, platelets and a SOFA input 0–24; here the four-item SOFA
// sum 0–16 that the table defines). INR 1.5, platelets 90, SOFA 3 → 6.
func TestSICScore_parity(t *testing.T) {
	if r, _ := SICScore(1.5, 90, 3); r.Total != 6 || !r.SIC {
		t.Errorf("got %+v, want 6, SIC", r)
	}
}

func TestSICScore_errors(t *testing.T) {
	for _, v := range nonFiniteB6 {
		_, err := SICScore(v, 100, 1)
		wantErrB6(t, "inr", err, ErrNotFinite)
		_, err = SICScore(1, v, 1)
		wantErrB6(t, "plt", err, ErrNotFinite)
	}
	_, err := SICScore(0, 100, 1)
	wantErrB6(t, "inr 0", err, ErrOutOfRange)
	_, err = SICScore(1, -1, 1)
	wantErrB6(t, "plt < 0", err, ErrOutOfRange)
	for _, s := range []int{-1, 17} {
		_, err = SICScore(1, 100, s)
		wantErrB6(t, "sofa", err, ErrOutOfRange)
	}
}

func TestSICScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { r, _ := SICScore(1.3, 120, 2); sinkB6Int = r.Total }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkSICScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		r, _ := SICScore(1.3, 120, 2)
		sinkB6Int = r.Total
	}
}
