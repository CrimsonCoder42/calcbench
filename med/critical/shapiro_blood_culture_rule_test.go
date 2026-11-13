package critical

import "testing"

// Shapiro 2008 abstract: ≥ 1 major or ≥ 2 minor; age "> 65 years". HAND-DERIVED.
func TestShapiroBloodCultureRule_source(t *testing.T) {
	tests := []struct {
		name         string
		age          float64
		f            ShapiroFindings
		major, minor int
		ind          bool
	}{
		{"none", 40, ShapiroFindings{}, 0, 0, false},
		{"catheter", 40, ShapiroFindings{IndwellingCatheter: true}, 1, 0, true},
		{"endocarditis", 40, ShapiroFindings{SuspectedEndocarditis: true}, 1, 0, true},
		{"temp > 39.5", 40, ShapiroFindings{TemperatureAbove39_5C: true}, 1, 0, true},
		{"age 65 + chills", 65, ShapiroFindings{Chills: true}, 0, 1, false},
		{"age 65.1 + chills", 65.1, ShapiroFindings{Chills: true}, 0, 2, true},
		{"two labs", 40, ShapiroFindings{PlateletsBelow150k: true, CreatinineAbove2: true}, 0, 2, true},
		{"all", 70, ShapiroFindings{true, true, true, true, true, true, true, true, true, true, true, true}, 3, 10, true},
	}
	for _, tc := range tests {
		r, err := ShapiroBloodCultureRule(tc.age, tc.f)
		if err != nil || r.Major != tc.major || r.Minor != tc.minor || r.Indicated != tc.ind {
			t.Errorf("%s: (%+v, %v), want %d %d %v", tc.name, r, err, tc.major, tc.minor, tc.ind)
		}
	}
}

// parity: ledger shapiro_rule, source case 1. The app weights majors 3/3/2 points and omits neutrophil% > 80; its
// result rule (≥ 1 major or ≥ 2 minor) is the abstract's. Neutrophils > 80% and vomiting: app 1 minor (not
// indicated), primary 2 minor (indicated) — asserted here.
func TestShapiroBloodCultureRule_parity(t *testing.T) {
	r, _ := ShapiroBloodCultureRule(40, ShapiroFindings{NeutrophilsAbove80Pct: true, Vomiting: true})
	if r.Minor != 2 || !r.Indicated {
		t.Errorf("got %+v, want Minor 2, Indicated (primary; app 1 minor)", r)
	}
}

func TestShapiroBloodCultureRule_errors(t *testing.T) {
	for _, v := range nonFiniteB6 {
		_, err := ShapiroBloodCultureRule(v, ShapiroFindings{})
		wantErrB6(t, "non-finite", err, ErrNotFinite)
	}
	_, err := ShapiroBloodCultureRule(131, ShapiroFindings{})
	wantErrB6(t, "range", err, ErrOutOfRange)
}

func TestShapiroBloodCultureRule_noAlloc(t *testing.T) {
	f := ShapiroFindings{Chills: true}
	if n := testing.AllocsPerRun(100, func() { r, _ := ShapiroBloodCultureRule(70, f); sinkB6Bool = r.Indicated }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkShapiroBloodCultureRule(b *testing.B) {
	b.ReportAllocs()
	f := ShapiroFindings{Chills: true}
	for b.Loop() {
		r, _ := ShapiroBloodCultureRule(70, f)
		sinkB6Bool = r.Indicated
	}
}
