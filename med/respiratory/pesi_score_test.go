package respiratory

import (
	"errors"
	"math"
	"testing"
)

var sinkPESIClass PESIClass

// pesiBase is a 50-year-old woman with no item present (score 50).
func pesiBase() PESIFindings {
	return PESIFindings{AgeYears: 50, PulsePerMin: 80, SystolicMmHg: 120, RespiratoryRatePerMin: 16, TemperatureC: 37, SaO2Percent: 97}
}

// HAND-DERIVED from Aujesky 2005 Table 2 (PMID 16020800): 70-year-old man, cancer, pulse 115 → 70+10+30+20 = 130.
func TestPESIScore_tableExample(t *testing.T) {
	f := pesiBase()
	f.AgeYears, f.Male, f.Cancer, f.PulsePerMin = 70, true, true, 115
	if got, err := PESIScore(f); err != nil || got != 130 {
		t.Errorf("got (%d, %v), want 130", got, err)
	}
	all := PESIFindings{AgeYears: 120, Male: true, Cancer: true, HeartFailure: true, ChronicLungDisease: true,
		PulsePerMin: 110, SystolicMmHg: 99, RespiratoryRatePerMin: 30, TemperatureC: 35.9, AlteredMentalStatus: true, SaO2Percent: 89}
	if got, err := PESIScore(all); err != nil || got != 350 {
		t.Errorf("all items: got (%d, %v), want 350", got, err)
	}
}

// Both sides of each measured threshold: pulse ≥ 110, SBP < 100, RR ≥ 30, temperature < 36, SaO₂ < 90.
func TestPESIScore_edges(t *testing.T) {
	cases := []struct {
		name string
		mod  func(*PESIFindings)
		want int
	}{
		{"pulse 109", func(f *PESIFindings) { f.PulsePerMin = 109 }, 50},
		{"pulse 110", func(f *PESIFindings) { f.PulsePerMin = 110 }, 70},
		{"SBP 100", func(f *PESIFindings) { f.SystolicMmHg = 100 }, 50},
		{"SBP 99.9", func(f *PESIFindings) { f.SystolicMmHg = 99.9 }, 80},
		{"RR 29", func(f *PESIFindings) { f.RespiratoryRatePerMin = 29 }, 50},
		{"RR 30", func(f *PESIFindings) { f.RespiratoryRatePerMin = 30 }, 70},
		{"temp 36", func(f *PESIFindings) { f.TemperatureC = 36 }, 50},
		{"temp 35.9", func(f *PESIFindings) { f.TemperatureC = 35.9 }, 70},
		{"SaO2 90", func(f *PESIFindings) { f.SaO2Percent = 90 }, 50},
		{"SaO2 89.9", func(f *PESIFindings) { f.SaO2Percent = 89.9 }, 70},
		{"altered mental status", func(f *PESIFindings) { f.AlteredMentalStatus = true }, 110},
		{"heart failure and chronic lung disease", func(f *PESIFindings) { f.HeartFailure, f.ChronicLungDisease = true, true }, 70},
	}
	for _, c := range cases {
		f := pesiBase()
		c.mod(&f)
		if got, err := PESIScore(f); err != nil || got != c.want {
			t.Errorf("%s: got (%d, %v), want %d", c.name, got, err, c.want)
		}
	}
}

// parity: ledger pesi, source case 1: age 55, female, no factor → 55, class I (the source's thresholds > 65,
// > 85, > 105, > 125 agree with the paper).
func TestPESIScore_parity(t *testing.T) {
	f := pesiBase()
	f.AgeYears = 55
	s, err := PESIScore(f)
	if err != nil || s != 55 {
		t.Fatalf("got (%d, %v), want 55", s, err)
	}
	if c, err := PESIRiskClass(s); err != nil || c != PESIClassI {
		t.Errorf("class: got (%q, %v), want %q", c, err, PESIClassI)
	}
}

func TestPESIScore_errors(t *testing.T) {
	cases := []struct {
		name string
		mod  func(*PESIFindings)
		want error
	}{
		{"age 17", func(f *PESIFindings) { f.AgeYears = 17 }, ErrOutOfRange},
		{"age 121", func(f *PESIFindings) { f.AgeYears = 121 }, ErrOutOfRange},
		{"pulse -1", func(f *PESIFindings) { f.PulsePerMin = -1 }, ErrOutOfRange},
		{"SBP 0", func(f *PESIFindings) { f.SystolicMmHg = 0 }, ErrOutOfRange},
		{"RR -1", func(f *PESIFindings) { f.RespiratoryRatePerMin = -1 }, ErrOutOfRange},
		{"temp 19", func(f *PESIFindings) { f.TemperatureC = 19 }, ErrOutOfRange},
		{"temp 46", func(f *PESIFindings) { f.TemperatureC = 46 }, ErrOutOfRange},
		{"SaO2 101", func(f *PESIFindings) { f.SaO2Percent = 101 }, ErrOutOfRange},
		{"pulse NaN", func(f *PESIFindings) { f.PulsePerMin = math.NaN() }, ErrNotFinite},
		{"SBP +Inf", func(f *PESIFindings) { f.SystolicMmHg = math.Inf(1) }, ErrNotFinite},
		{"RR -Inf", func(f *PESIFindings) { f.RespiratoryRatePerMin = math.Inf(-1) }, ErrNotFinite},
		{"temp NaN", func(f *PESIFindings) { f.TemperatureC = math.NaN() }, ErrNotFinite},
		{"SaO2 +Inf", func(f *PESIFindings) { f.SaO2Percent = math.Inf(1) }, ErrNotFinite},
	}
	for _, c := range cases {
		f := pesiBase()
		c.mod(&f)
		if _, err := PESIScore(f); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// Class edges, Aujesky 2005: ≤ 65 I, 66–85 II, 86–105 III, 106–125 IV, > 125 V.
func TestPESIRiskClass_edges(t *testing.T) {
	cases := []struct {
		total int
		want  PESIClass
	}{
		{18, PESIClassI}, {65, PESIClassI}, {66, PESIClassII}, {85, PESIClassII}, {86, PESIClassIII},
		{105, PESIClassIII}, {106, PESIClassIV}, {125, PESIClassIV}, {126, PESIClassV}, {350, PESIClassV},
	}
	for _, c := range cases {
		if got, err := PESIRiskClass(c.total); err != nil || got != c.want {
			t.Errorf("%d: got (%q, %v), want %q", c.total, got, err, c.want)
		}
	}
	for _, bad := range []int{17, 351, -1} {
		if _, err := PESIRiskClass(bad); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("%d: err = %v, want ErrOutOfRange", bad, err)
		}
	}
}

func TestPESIScore_noAlloc(t *testing.T) {
	f := pesiBase()
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = PESIScore(f) }); n != 0 {
		t.Errorf("PESIScore allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sinkPESIClass, _ = PESIRiskClass(90) }); n != 0 {
		t.Errorf("PESIRiskClass allocs = %v, want 0", n)
	}
}

func BenchmarkPESIScore(b *testing.B) {
	f := pesiBase()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkInt, _ = PESIScore(f)
	}
}

func BenchmarkPESIRiskClass(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkPESIClass, _ = PESIRiskClass(90)
	}
}
