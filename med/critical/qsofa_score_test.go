package critical

import "testing"

// Seymour 2016 abstract: SBP ≤ 100, RR ≥ 22, altered mentation; both sides of each edge.
func TestQSOFAScore_source(t *testing.T) {
	tests := []struct {
		rr, sbp float64
		alt     bool
		want    int
	}{
		{21, 101, false, 0},
		{22, 101, false, 1},
		{21, 100, false, 1},
		{21.9, 100.1, true, 1},
		{22, 100, true, 3},
	}
	for i, tc := range tests {
		if got, err := QSOFAScore(tc.rr, tc.sbp, tc.alt); err != nil || got != tc.want {
			t.Errorf("case %d: (%d, %v), want %d", i, got, err, tc.want)
		}
	}
}

// parity: ledgers qsofa and qsofa_calc, source cases (app: rr ≥ 22, sbp ≤ 100, gcs < 15 = altered mentation).
// GCS 14 → alteredMentation true; RR 24, SBP 95, GCS 14 → 3.
func TestQSOFAScore_parity(t *testing.T) {
	if got, _ := QSOFAScore(24, 95, true); got != 3 {
		t.Errorf("got %d, want 3", got)
	}
	if got, _ := QSOFAScore(16, 120, false); got != 0 {
		t.Errorf("got %d, want 0", got)
	}
}

func TestQSOFAScore_errors(t *testing.T) {
	for _, v := range nonFiniteB6 {
		_, err := QSOFAScore(v, 120, false)
		wantErrB6(t, "rr", err, ErrNotFinite)
		_, err = QSOFAScore(16, v, false)
		wantErrB6(t, "sbp", err, ErrNotFinite)
	}
	_, err := QSOFAScore(-1, 120, false)
	wantErrB6(t, "rr < 0", err, ErrOutOfRange)
	_, err = QSOFAScore(16, -1, false)
	wantErrB6(t, "sbp < 0", err, ErrOutOfRange)
}

func TestQSOFAScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkB6Int, _ = QSOFAScore(24, 95, true) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkQSOFAScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkB6Int, _ = QSOFAScore(24, 95, true)
	}
}
