package trauma

import (
	"math"
	"testing"
)

var sinkNTS NewTraumaScoreResult

// Jeong 2017 Table 1: every code band, both sides of each edge.
func TestNewTraumaScore_codes(t *testing.T) {
	sbp := []struct{ v, code int }{{0, 0}, {69, 0}, {70, 1}, {89, 1}, {90, 2}, {109, 2}, {110, 4}, {149, 4}, {150, 3}, {300, 3}}
	for _, tc := range sbp {
		r, err := NewTraumaScore(15, tc.v, 98)
		if err != nil || r.SBPCode != tc.code {
			t.Errorf("SBP %d → %d (%v), want %d", tc.v, r.SBPCode, err, tc.code)
		}
	}
	spo2 := []struct{ v, code int }{{0, 0}, {39, 0}, {40, 1}, {59, 1}, {60, 2}, {79, 2}, {80, 3}, {93, 3}, {94, 4}, {100, 4}}
	for _, tc := range spo2 {
		r, err := NewTraumaScore(15, 120, tc.v)
		if err != nil || r.SpO2Code != tc.code {
			t.Errorf("SpO2 %d → %d (%v), want %d", tc.v, r.SpO2Code, err, tc.code)
		}
	}
}

// HAND-DERIVED from Jeong 2017's equations: GCS 15, SBP 120 (code 4), SpO2 98 (code 4):
// NTS = 0.4006×15 + 0.2983×4 + 0.8709×4 = 6.009 + 1.1932 + 3.4836 = 10.6858; T-NTS = 23;
// b = −6.5406 + 10.6858 = 4.1452, Ps = 1/(1 + e^−4.1452) = 0.984407…
// GCS 3, SBP 60, SpO2 30: NTS = 1.2018, T-NTS 3, Ps = 1/(1 + e^5.3388) = 0.004779…
// Tolerance 1e-12 absolute: exact arithmetic on the printed coefficients.
func TestNewTraumaScore_equations(t *testing.T) {
	tests := []struct {
		gcs, sbp, spo2 int
		nts, tnts      float64
	}{
		{15, 120, 98, 10.6858, 23},
		{3, 60, 30, 1.2018, 3},
	}
	for _, tc := range tests {
		r, err := NewTraumaScore(tc.gcs, tc.sbp, tc.spo2)
		ps := 1 / (1 + math.Exp(-(-6.5406 + tc.nts)))
		if err != nil || math.Abs(r.NTS-tc.nts) > 1e-12 || r.TriageNTS != tc.tnts || math.Abs(r.SurvivalProbability-ps) > 1e-12 {
			t.Errorf("%+v → %+v (%v), want NTS %v T-NTS %v Ps %v", tc, r, err, tc.nts, tc.tnts, ps)
		}
	}
	r, _ := NewTraumaScore(15, 120, 98)
	if math.Abs(r.SurvivalProbability-0.984407) > 1e-6 {
		t.Errorf("Ps = %v, want 0.984407 (printed to 6 places)", r.SurvivalProbability)
	}
}

// parity: ledger nts, source case 1. The app computes GCS + 0.2983×SBP code + 0.8709×SpO2 code (GCS unweighted);
// the primary weights GCS by 0.4006. For GCS 15, SBP 120, SpO2 98 the app gives 19.6768; the primary gives 10.6858,
// asserted here.
func TestNewTraumaScore_parity(t *testing.T) {
	r, _ := NewTraumaScore(15, 120, 98)
	if math.Abs(r.NTS-10.6858) > 1e-12 {
		t.Errorf("NTS = %v, want 10.6858 (primary; app 19.6768)", r.NTS)
	}
}

func TestNewTraumaScore_errors(t *testing.T) {
	for _, in := range [][3]int{{2, 120, 98}, {16, 120, 98}, {15, -1, 98}, {15, 120, -1}, {15, 120, 101}} {
		_, err := NewTraumaScore(in[0], in[1], in[2])
		checkErr(t, "range", err, ErrOutOfRange)
	}
}

func TestNewTraumaScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkNTS, _ = NewTraumaScore(15, 120, 98) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkNewTraumaScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkNTS, _ = NewTraumaScore(15, 120, 98)
	}
}
