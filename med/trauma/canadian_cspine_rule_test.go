package trauma

import "testing"

// Stiell 2001 abstract: the three questions. HAND-DERIVED cases, one per path.
func TestCanadianCSpineRule_paths(t *testing.T) {
	tests := []struct {
		name               string
		age                float64
		f                  CanadianCSpineFindings
		high, low, radiogr bool
	}{
		{"age 65 is high risk", 65, CanadianCSpineFindings{SittingPositionInED: true, AbleToRotateNeck45: true}, true, true, true},
		{"age 64.9, low risk, rotates", 64.9, CanadianCSpineFindings{SittingPositionInED: true, AbleToRotateNeck45: true}, false, true, false},
		{"dangerous mechanism", 30, CanadianCSpineFindings{DangerousMechanism: true, AmbulatorySinceInjury: true, AbleToRotateNeck45: true}, true, true, true},
		{"paresthesias", 30, CanadianCSpineFindings{ParesthesiasExtremities: true}, true, false, true},
		{"no low-risk factor", 30, CanadianCSpineFindings{AbleToRotateNeck45: true}, false, false, true},
		{"low risk, cannot rotate", 30, CanadianCSpineFindings{DelayedOnsetNeckPain: true}, false, true, true},
		{"rear-end, rotates", 30, CanadianCSpineFindings{SimpleRearEndCollision: true, AbleToRotateNeck45: true}, false, true, false},
		{"no midline tenderness, rotates", 30, CanadianCSpineFindings{NoMidlineTenderness: true, AbleToRotateNeck45: true}, false, true, false},
	}
	for _, tc := range tests {
		r, err := CanadianCSpineRule(tc.age, tc.f)
		if err != nil || r.HighRiskFactor != tc.high || r.LowRiskFactor != tc.low || r.Radiography != tc.radiogr {
			t.Errorf("%s: (%+v, %v), want high %v low %v radiography %v", tc.name, r, err, tc.high, tc.low, tc.radiogr)
		}
	}
}

// parity: ledger canadian_cspine, source cases (app compute: hasHighRisk / hasLowRisk / canRotate; no source _test.go).
func TestCanadianCSpineRule_parity(t *testing.T) {
	r, _ := CanadianCSpineRule(70, CanadianCSpineFindings{}) // app age65 = 1
	if !r.Radiography {
		t.Errorf("age65 → %+v, want radiography", r)
	}
}

func TestCanadianCSpineRule_errors(t *testing.T) {
	for _, v := range nonFinite {
		_, err := CanadianCSpineRule(v, CanadianCSpineFindings{})
		checkErr(t, "non-finite", err, ErrNotFinite)
	}
	for _, v := range []float64{-0.1, 130.1} {
		_, err := CanadianCSpineRule(v, CanadianCSpineFindings{})
		checkErr(t, "range", err, ErrOutOfRange)
	}
	for _, v := range []float64{0, 130} {
		if _, err := CanadianCSpineRule(v, CanadianCSpineFindings{}); err != nil {
			t.Errorf("age %v: %v", v, err)
		}
	}
}

func TestCanadianCSpineRule_noAlloc(t *testing.T) {
	f := CanadianCSpineFindings{SittingPositionInED: true, AbleToRotateNeck45: true}
	if n := testing.AllocsPerRun(100, func() { r, _ := CanadianCSpineRule(30, f); sinkBool = r.Radiography }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkCanadianCSpineRule(b *testing.B) {
	b.ReportAllocs()
	f := CanadianCSpineFindings{SittingPositionInED: true, AbleToRotateNeck45: true}
	for b.Loop() {
		r, _ := CanadianCSpineRule(30, f)
		sinkBool = r.Radiography
	}
}
