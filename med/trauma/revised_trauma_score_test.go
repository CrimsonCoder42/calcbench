package trauma

import (
	"errors"
	"math"
	"testing"
)

var sinkRTS RevisedTraumaScoreResult

// Yousefi 2024 Table 1 and Yu 2021 Table 1 (both citing Champion 1989): every code band, both sides of each edge.
func TestRevisedTraumaScore_codes(t *testing.T) {
	gcs := []struct{ v, code int }{{3, 0}, {4, 1}, {5, 1}, {6, 2}, {8, 2}, {9, 3}, {12, 3}, {13, 4}, {15, 4}}
	for _, tc := range gcs {
		r, err := RevisedTraumaScore(tc.v, 120, 16)
		if err != nil || r.GCSCode != tc.code {
			t.Errorf("GCS %d → %d (%v), want %d", tc.v, r.GCSCode, err, tc.code)
		}
	}
	sbp := []struct{ v, code int }{{0, 0}, {1, 1}, {49, 1}, {50, 2}, {75, 2}, {76, 3}, {89, 3}, {90, 4}, {300, 4}}
	for _, tc := range sbp {
		r, err := RevisedTraumaScore(15, tc.v, 16)
		if err != nil || r.SBPCode != tc.code {
			t.Errorf("SBP %d → %d (%v), want %d", tc.v, r.SBPCode, err, tc.code)
		}
	}
	rr := []struct{ v, code int }{{0, 0}, {1, 1}, {5, 1}, {6, 2}, {9, 2}, {10, 4}, {29, 4}, {30, 3}, {80, 3}}
	for _, tc := range rr {
		r, err := RevisedTraumaScore(15, 120, tc.v)
		if err != nil || r.RRCode != tc.code {
			t.Errorf("RR %d → %d (%v), want %d", tc.v, r.RRCode, err, tc.code)
		}
	}
}

// HAND-DERIVED from "RTS = (0.9368 × GCS value) + (0.7326 × SBP value) + (0.2908 × RR value)" (Yousefi 2024):
// all codes 4 → T-RTS 12, RTS 7.8408 (the 0–7.8408 maximum); all codes 0 → 0; GCS 10 (3), SBP 80 (3), RR 8 (2):
// T-RTS 8, RTS 2.8104 + 2.1978 + 0.5816 = 5.5898. Tolerance 1e-12 absolute: exact arithmetic on printed coefficients.
func TestRevisedTraumaScore_source(t *testing.T) {
	tests := []struct {
		gcs, sbp, rr int
		trts         int
		rts          float64
	}{
		{15, 120, 16, 12, 7.8408},
		{3, 0, 0, 0, 0},
		{10, 80, 8, 8, 5.5898},
	}
	for _, tc := range tests {
		r, err := RevisedTraumaScore(tc.gcs, tc.sbp, tc.rr)
		if err != nil || r.TriageRTS != tc.trts || math.Abs(r.RTS-tc.rts) > 1e-12 {
			t.Errorf("%+v → %+v (%v), want T-RTS %d RTS %v", tc, r, err, tc.trts, tc.rts)
		}
	}
}

// parity: ledger revised_trauma_score, source case 1. The entry carries no _test.go and no JSON example; its
// app code uses the same codes and coefficients, so the defaults (GCS 15, SBP 120, RR 16) give T-RTS 12.
func TestRevisedTraumaScore_parity(t *testing.T) {
	r, err := RevisedTraumaScore(15, 120, 16)
	if err != nil || r.TriageRTS != 12 {
		t.Errorf("T-RTS = %d (%v), want 12", r.TriageRTS, err)
	}
}

func TestRevisedTraumaScore_errors(t *testing.T) {
	cases := []struct {
		name         string
		gcs, sbp, rr int
	}{
		{"gcs 2", 2, 120, 16}, {"gcs 16", 16, 120, 16}, {"sbp -1", 15, -1, 16}, {"rr -1", 15, 120, -1},
	}
	for _, tc := range cases {
		if _, err := RevisedTraumaScore(tc.gcs, tc.sbp, tc.rr); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("%s: err = %v, want ErrOutOfRange", tc.name, err)
		}
	}
}

func TestRevisedTraumaScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkRTS, _ = RevisedTraumaScore(12, 80, 8) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkRevisedTraumaScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkRTS, _ = RevisedTraumaScore(12, 80, 8)
	}
}
