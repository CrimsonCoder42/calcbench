package cardio

import (
	"math"
	"testing"
)

// Gage 2001 abstract: points and the stroke rate per 100 person-years for each score.
func TestCHADS2_sourceTable(t *testing.T) {
	want := []float64{1.9, 2.8, 4.0, 5.9, 8.5, 12.5, 18.2}
	for s, w := range want {
		if got, err := CHADS2StrokeRatePer100PersonYears(s); err != nil || got != w {
			t.Errorf("rate(%d) = (%v, %v), want %v", s, got, err, w)
		}
	}
	for _, s := range []int{-1, 7} {
		_, err := CHADS2StrokeRatePer100PersonYears(s)
		checkErr(t, "rate range", err, ErrOutOfRange)
	}
}

// parity: ledger chads2_score, source cases 1–6. The source takes "age ≥ 75" as a 0/1 input (with a caller-supplied age
// of 40, outside Gage's 65–95 cohort); here age75 = 1 is age 80 and age75 = 0 is age 70.
func TestCHADS2Score_parity(t *testing.T) {
	tests := []struct {
		age   float64
		f     CHADS2Factors
		score int
		rate  float64
	}{
		{70, CHADS2Factors{}, 0, 1.9},                                       // source case 1
		{70, CHADS2Factors{CongestiveHeartFailure: true}, 1, 2.8},           // source case 2
		{80, CHADS2Factors{Hypertension: true}, 2, 4.0},                     // source case 3
		{80, CHADS2Factors{Hypertension: true, Diabetes: true}, 3, 5.9},     // source case 4
		{70, CHADS2Factors{Diabetes: true, PriorStrokeOrTIA: true}, 3, 5.9}, // source case 5
		{80, CHADS2Factors{true, true, true, true}, 6, 18.2},                // source case 6
	}
	for i, tc := range tests {
		s, err := CHADS2Score(tc.age, tc.f)
		r, _ := CHADS2StrokeRatePer100PersonYears(s)
		if err != nil || s != tc.score || r != tc.rate {
			t.Errorf("case %d: (%d, %v, %v), want (%d, %v)", i+1, s, r, err, tc.score, tc.rate)
		}
	}
}

func TestCHADS2Score_edges(t *testing.T) {
	if s, _ := CHADS2Score(74.99, CHADS2Factors{}); s != 0 {
		t.Errorf("74.99 → %d, want 0", s)
	}
	if s, _ := CHADS2Score(75, CHADS2Factors{}); s != 1 {
		t.Errorf("75 → %d, want 1", s)
	}
	for _, a := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		_, err := CHADS2Score(a, CHADS2Factors{})
		checkErr(t, "non-finite", err, ErrNotFinite)
	}
	for _, a := range []float64{64.99, 95.01} {
		_, err := CHADS2Score(a, CHADS2Factors{})
		checkErr(t, "range", err, ErrOutOfRange)
	}
	for _, a := range []float64{65, 95} {
		if _, err := CHADS2Score(a, CHADS2Factors{}); err != nil {
			t.Errorf("age %v: %v", a, err)
		}
	}
}

func TestCHADS2Score_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() {
		sinkInt, _ = CHADS2Score(80, CHADS2Factors{Hypertension: true})
		sinkFloat, _ = CHADS2StrokeRatePer100PersonYears(2)
	}); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkCHADS2Score(b *testing.B) {
	b.ReportAllocs()
	f := CHADS2Factors{Hypertension: true}
	for b.Loop() {
		sinkInt, _ = CHADS2Score(80, f)
	}
}

func BenchmarkCHADS2StrokeRatePer100PersonYears(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = CHADS2StrokeRatePer100PersonYears(2)
	}
}
