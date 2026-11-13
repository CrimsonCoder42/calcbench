package trauma

import "testing"

// Stiell 1995 abstract: five variables, "age 55 years or older". HAND-DERIVED.
func TestOttawaKneeRule_source(t *testing.T) {
	tests := []struct {
		age  float64
		f    OttawaKneeFindings
		want int
	}{
		{54.9, OttawaKneeFindings{}, 0},
		{55, OttawaKneeFindings{}, 1},
		{30, OttawaKneeFindings{FibularHeadTenderness: true}, 1},
		{30, OttawaKneeFindings{IsolatedPatellarTenderness: true, UnableToFlex90: true}, 2},
		{60, OttawaKneeFindings{true, true, true, true}, 5},
		{30, OttawaKneeFindings{UnableToBearWeightFourSteps: true}, 1},
	}
	for i, tc := range tests {
		if got, err := OttawaKneeRule(tc.age, tc.f); err != nil || got != tc.want {
			t.Errorf("case %d: (%d, %v), want %d", i, got, err, tc.want)
		}
	}
}

// parity: ledger ottawa_knee_rule, source case 1 (app anyPositive over five 0/1 items; age_55 = 1 is age 60 here).
func TestOttawaKneeRule_parity(t *testing.T) {
	if got, _ := OttawaKneeRule(60, OttawaKneeFindings{}); got != 1 {
		t.Errorf("got %d, want 1", got)
	}
}

func TestOttawaKneeRule_errors(t *testing.T) {
	for _, v := range nonFinite {
		_, err := OttawaKneeRule(v, OttawaKneeFindings{})
		checkErr(t, "non-finite", err, ErrNotFinite)
	}
	_, err := OttawaKneeRule(130.5, OttawaKneeFindings{})
	checkErr(t, "range", err, ErrOutOfRange)
}

func TestOttawaKneeRule_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = OttawaKneeRule(60, OttawaKneeFindings{}) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkOttawaKneeRule(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = OttawaKneeRule(60, OttawaKneeFindings{})
	}
}
