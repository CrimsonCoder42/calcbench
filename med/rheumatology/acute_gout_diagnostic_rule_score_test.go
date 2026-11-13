package rheumatology

import "testing"

var sinkFloatRheum float64

// Points (Lee 2015, PMC4444468; McConville & Cox 2023, PMC10211406): male 2, previous attack 2, onset ≤ 1 day
// 0.5, redness 1, MTP1 2.5, hypertension/CVD 1.5, uric acid > 5.88 mg/dL 3.5. Sums of halves are exact, so the
// comparison is exact.
func TestAcuteGoutDiagnosticRuleScore_table(t *testing.T) {
	tests := []struct {
		name                  string
		m, p, o, r, mtp, h, u bool
		want                  float64
	}{
		{"none", false, false, false, false, false, false, false, 0},
		{"male", true, false, false, false, false, false, false, 2},
		{"previous attack", false, true, false, false, false, false, false, 2},
		{"onset within 1 day", false, false, true, false, false, false, false, 0.5},
		{"redness", false, false, false, true, false, false, false, 1},
		{"MTP1", false, false, false, false, true, false, false, 2.5},
		{"hypertension or CVD", false, false, false, false, false, true, false, 1.5},
		{"uric acid", false, false, false, false, false, false, true, 3.5},
		{"maximum", true, true, true, true, true, true, true, 13},
		{"≤4 edge: male + onset + hypertension", true, false, true, false, false, true, false, 4},
		{"just above 4: male + MTP1", true, false, false, false, true, false, false, 4.5},
		{"≥8 edge: uric acid + MTP1 + male", true, false, false, false, true, false, true, 8},
		{"just below 8: uric acid + MTP1 + hypertension", false, false, false, false, true, true, true, 7.5},
	}
	for _, tc := range tests {
		if got := AcuteGoutDiagnosticRuleScore(tc.m, tc.p, tc.o, tc.r, tc.mtp, tc.h, tc.u); got != tc.want {
			t.Errorf("%s: got %v, want %v", tc.name, got, tc.want)
		}
	}
}

// parity: ledger acute_gout_diagnosis. The source entry has no _test.go and its JSON no example; the consumer
// sums the same seven weights, so case 1 (all seven present) is 13.
func TestAcuteGoutDiagnosticRuleScore_parity(t *testing.T) {
	if got := AcuteGoutDiagnosticRuleScore(true, true, true, true, true, true, true); got != 13 {
		t.Errorf("source case 1: got %v, want 13", got)
	}
}

func TestAcuteGoutDiagnosticRuleScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() {
		sinkFloatRheum = AcuteGoutDiagnosticRuleScore(true, false, true, false, true, false, true)
	}); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkAcuteGoutDiagnosticRuleScore(b *testing.B) {
	b.ReportAllocs()
	x := true
	for b.Loop() {
		sinkFloatRheum = AcuteGoutDiagnosticRuleScore(x, !x, x, !x, x, !x, x)
	}
}
