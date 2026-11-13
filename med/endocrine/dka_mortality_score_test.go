package endocrine

import "testing"

// below is a findings value with every numeric predictor on the zero-point side of its edge.
var dkaBelow = DKAMortalityFindings{PHAtPresentation: 7.2, RegularInsulinIUAt12h: 30, GlucoseMmolLAt12h: 12}

// HAND-DERIVED from Efstathiou 2002 (abstract): each predictor alone gives its points; all give 25.
func TestDKAMortalityScore_example(t *testing.T) {
	all := DKAMortalityFindings{true, 6.9, 60, 18, true, true}
	scd := dkaBelow
	scd.SevereCoexistingDiseases = true
	ph := dkaBelow
	ph.PHAtPresentation = 6.95
	ins := dkaBelow
	ins.RegularInsulinIUAt12h = 51
	glu := dkaBelow
	glu.GlucoseMmolLAt12h = 16.8
	mental := dkaBelow
	mental.DepressedMentalStateAt24h = true
	fever := dkaBelow
	fever.FeverAt24h = true
	tests := []struct {
		name string
		f    DKAMortalityFindings
		want int
	}{
		{"none", dkaBelow, 0},
		{"all", all, 25},
		{"SCD", scd, 6},
		{"pH", ph, 4},
		{"insulin", ins, 4},
		{"glucose", glu, 4},
		{"mental", mental, 4},
		{"fever", fever, 3},
	}
	for _, tc := range tests {
		if got, err := DKAMortalityScore(tc.f); err != nil || got != tc.want {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// The source's edges are strict: pH < 7.0, insulin > 50 IU, glucose > 16.7 mmol/L.
func TestDKAMortalityScore_edges(t *testing.T) {
	tests := []struct {
		name         string
		ph, ins, glu float64
		want         int
	}{
		{"pH 7.0 scores 0", 7.0, 30, 12, 0},
		{"pH 6.99 scores 4", 6.99, 30, 12, 4},
		{"insulin 50 scores 0", 7.2, 50, 12, 0},
		{"insulin 50.5 scores 4", 7.2, 50.5, 12, 4},
		{"glucose 16.7 scores 0", 7.2, 30, 16.7, 0},
		{"glucose 16.71 scores 4", 7.2, 30, 16.71, 4},
		{"pH 14 accepted", 14, 0, 12, 0},
	}
	for _, tc := range tests {
		f := DKAMortalityFindings{PHAtPresentation: tc.ph, RegularInsulinIUAt12h: tc.ins, GlucoseMmolLAt12h: tc.glu}
		if got, err := DKAMortalityScore(f); err != nil || got != tc.want {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// parity: ledger dka_mpm, the app's six yes/no items with points 6/4/4/4/4/3 (cases: all no → 0, all yes → 25,
// comorbidity + pH → 10). The app's glucose item is "> 300 mg/dL"; the primary states > 16.7 mmol/L, used here.
func TestDKAMortalityScore_parity(t *testing.T) {
	two := dkaBelow
	two.SevereCoexistingDiseases = true
	two.PHAtPresentation = 6.8
	if got, err := DKAMortalityScore(two); err != nil || got != 10 {
		t.Errorf("SCD + pH: got (%v, %v), want 10", got, err)
	}
}

func TestDKAMortalityScore_errors(t *testing.T) {
	with := func(mut func(*DKAMortalityFindings)) func() error {
		return func() error { f := dkaBelow; mut(&f); _, err := DKAMortalityScore(f); return err }
	}
	runErrCases(t, []errCase{
		{"pH NaN", with(func(f *DKAMortalityFindings) { f.PHAtPresentation = nan }), ErrNotFinite},
		{"insulin +Inf", with(func(f *DKAMortalityFindings) { f.RegularInsulinIUAt12h = pInf }), ErrNotFinite},
		{"glucose -Inf", with(func(f *DKAMortalityFindings) { f.GlucoseMmolLAt12h = nInf }), ErrNotFinite},
		{"pH 0", with(func(f *DKAMortalityFindings) { f.PHAtPresentation = 0 }), ErrOutOfRange},
		{"pH above 14", with(func(f *DKAMortalityFindings) { f.PHAtPresentation = 14.01 }), ErrOutOfRange},
		{"insulin negative", with(func(f *DKAMortalityFindings) { f.RegularInsulinIUAt12h = -1 }), ErrOutOfRange},
		{"glucose 0", with(func(f *DKAMortalityFindings) { f.GlucoseMmolLAt12h = 0 }), ErrOutOfRange},
	})
}

func TestDKAMortalityScore_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = DKAMortalityScore(dkaBelow) }); n != 0 {
		t.Errorf("DKAMortalityScore allocates %v times, want 0", n)
	}
}

func BenchmarkDKAMortalityScore(b *testing.B) {
	b.ReportAllocs()
	f := DKAMortalityFindings{true, 6.9, 60, 18, true, true}
	for b.Loop() {
		sinkInt, _ = DKAMortalityScore(f)
	}
}
