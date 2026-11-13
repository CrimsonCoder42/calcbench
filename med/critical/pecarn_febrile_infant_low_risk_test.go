package critical

import "testing"

// Kuppermann 2019 abstract: negative UA, ANC ≤ 4090/µL, PCT ≤ 1.71 ng/mL; both sides of each edge.
func TestPECARNFebrileInfantLowRisk_source(t *testing.T) {
	tests := []struct {
		pos      bool
		anc, pct float64
		want     bool
	}{
		{false, 4090, 1.71, true},
		{false, 4091, 1.71, false},
		{false, 4090, 1.72, false},
		{true, 1000, 0.1, false},
		{false, 0, 0, true},
	}
	for i, tc := range tests {
		if got, err := PECARNFebrileInfantLowRisk(30, tc.pos, tc.anc, tc.pct); err != nil || got != tc.want {
			t.Errorf("case %d: (%v, %v), want %v", i, got, err, tc.want)
		}
	}
}

// parity: ledger pecarn_febrile_infant, source case 1. The app uses ANC ≤ 4000 and PCT ≤ 0.5 (not the 2019
// derivation's values); ANC 4050, PCT 1.0, negative UA: app not low risk, primary low risk — asserted here.
func TestPECARNFebrileInfantLowRisk_parity(t *testing.T) {
	if got, _ := PECARNFebrileInfantLowRisk(30, false, 4050, 1.0); !got {
		t.Error("got false, want true (primary 4090 / 1.71; app false)")
	}
}

func TestPECARNFebrileInfantLowRisk_errors(t *testing.T) {
	for _, v := range nonFiniteB6 {
		_, err := PECARNFebrileInfantLowRisk(v, false, 1000, 0.1)
		wantErrB6(t, "age", err, ErrNotFinite)
		_, err = PECARNFebrileInfantLowRisk(30, false, v, 0.1)
		wantErrB6(t, "anc", err, ErrNotFinite)
		_, err = PECARNFebrileInfantLowRisk(30, false, 1000, v)
		wantErrB6(t, "pct", err, ErrNotFinite)
	}
	for _, a := range []float64{-0.1, 60.1} {
		_, err := PECARNFebrileInfantLowRisk(a, false, 1000, 0.1)
		wantErrB6(t, "age range", err, ErrOutOfRange)
	}
	for _, a := range []float64{0, 60} {
		if _, err := PECARNFebrileInfantLowRisk(a, false, 1000, 0.1); err != nil {
			t.Errorf("age %v: %v", a, err)
		}
	}
	_, err := PECARNFebrileInfantLowRisk(30, false, -1, 0.1)
	wantErrB6(t, "anc < 0", err, ErrOutOfRange)
	_, err = PECARNFebrileInfantLowRisk(30, false, 1000, -0.1)
	wantErrB6(t, "pct < 0", err, ErrOutOfRange)
}

func TestPECARNFebrileInfantLowRisk_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkB6Bool, _ = PECARNFebrileInfantLowRisk(30, false, 3000, 0.2) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkPECARNFebrileInfantLowRisk(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkB6Bool, _ = PECARNFebrileInfantLowRisk(30, false, 3000, 0.2)
	}
}
