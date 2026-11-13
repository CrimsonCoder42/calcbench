package trauma

import "testing"

// Kim 2021 Table 1: weight and SBP edges, grades. HAND-DERIVED (all categorical items at +2 → base 8).
func TestPediatricTraumaScore_edges(t *testing.T) {
	tests := []struct {
		w, sbp float64
		want   int
	}{
		{25, 100, 12}, {20.1, 90.1, 12}, {20, 90, 10}, {10, 50, 10}, {9.9, 49.9, 6}, {0, 0, 6},
	}
	for _, tc := range tests {
		if got, err := PediatricTraumaScore(tc.w, tc.sbp, PTSNone, PTSNone, PTSNone, PTSNone); err != nil || got != tc.want {
			t.Errorf("w %v sbp %v: (%d, %v), want %d", tc.w, tc.sbp, got, err, tc.want)
		}
	}
	if got, _ := PediatricTraumaScore(5, 40, PTSMajor, PTSMajor, PTSMajor, PTSMajor); got != -6 {
		t.Errorf("minimum = %d, want -6", got)
	}
	if got, _ := PediatricTraumaScore(30, 100, PTSMinor, PTSMinor, PTSMinor, PTSMinor); got != 8 {
		t.Errorf("minor grades = %d, want 8", got)
	}
}

// Doc comment example: 1+2+2+2+1+1 = 9 (HAND-DERIVED).
func TestPediatricTraumaScore_example(t *testing.T) {
	if got, err := PediatricTraumaScore(15, 95, PTSNone, PTSNone, PTSMinor, PTSMinor); err != nil || got != 9 {
		t.Errorf("got (%d, %v), want 9", got, err)
	}
}

// parity: ledger pediatric_trauma_score, source items (−1/+1/+2 each, range −6 to 12): 8 kg, SBP 60, airway
// maintainable, obtunded, open fracture, major wound → −1+1+1+1−1−1 = 0.
func TestPediatricTraumaScore_parity(t *testing.T) {
	if got, err := PediatricTraumaScore(8, 60, PTSMinor, PTSMinor, PTSMajor, PTSMajor); err != nil || got != 0 {
		t.Errorf("got (%d, %v), want 0", got, err)
	}
}

func TestPediatricTraumaScore_errors(t *testing.T) {
	for _, v := range nonFinite {
		_, err := PediatricTraumaScore(v, 90, PTSNone, PTSNone, PTSNone, PTSNone)
		checkErr(t, "weight", err, ErrNotFinite)
		_, err = PediatricTraumaScore(15, v, PTSNone, PTSNone, PTSNone, PTSNone)
		checkErr(t, "sbp", err, ErrNotFinite)
	}
	_, err := PediatricTraumaScore(-1, 90, PTSNone, PTSNone, PTSNone, PTSNone)
	checkErr(t, "weight -1", err, ErrOutOfRange)
	_, err = PediatricTraumaScore(15, -1, PTSNone, PTSNone, PTSNone, PTSNone)
	checkErr(t, "sbp -1", err, ErrOutOfRange)
	for _, g := range []PTSGrade{0, 3, -2} {
		_, err = PediatricTraumaScore(15, 90, g, PTSNone, PTSNone, PTSNone)
		checkErr(t, "grade", err, ErrOutOfRange)
		_, err = PediatricTraumaScore(15, 90, PTSNone, PTSNone, PTSNone, g)
		checkErr(t, "grade", err, ErrOutOfRange)
	}
}

func TestPediatricTraumaScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() {
		sinkInt, _ = PediatricTraumaScore(15, 95, PTSNone, PTSNone, PTSMinor, PTSMinor)
	}); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkPediatricTraumaScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = PediatricTraumaScore(15, 95, PTSNone, PTSNone, PTSMinor, PTSMinor)
	}
}
