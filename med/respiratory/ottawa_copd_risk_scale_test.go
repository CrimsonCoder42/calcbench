package respiratory

import (
	"errors"
	"math"
	"testing"
)

var sinkOttawaCOPD OttawaCOPDRisk

// ottawaBase has no item present: heart rate 90, hemoglobin 130 g/L, urea 6, CO₂ 26.
func ottawaBase() OttawaCOPDFindings {
	return OttawaCOPDFindings{HeartRateOnArrivalPerMin: 90, HemoglobinGL: 130, UreaMmolL: 6, SerumCO2MmolL: 26}
}

// HAND-DERIVED from Stiell 2014 Table 7 (PMID 24549125): intubation 2 + HR 115 2 + Hb 95 g/L 3 = 7; every item
// present = 16, the table's "possible range 0–16".
func TestOttawaCOPDRiskScale_tableExample(t *testing.T) {
	f := ottawaBase()
	f.IntubationForRespiratory, f.HeartRateOnArrivalPerMin, f.HemoglobinGL = true, 115, 95
	if got, err := OttawaCOPDRiskScale(f); err != nil || got != 7 {
		t.Errorf("got (%d, %v), want 7", got, err)
	}
	all := OttawaCOPDFindings{true, true, true, 110, true, true, true, 99, 12, 35}
	if got, err := OttawaCOPDRiskScale(all); err != nil || got != 16 {
		t.Errorf("all: got (%d, %v), want 16", got, err)
	}
	if got, err := OttawaCOPDRiskScale(ottawaBase()); err != nil || got != 0 {
		t.Errorf("none: got (%d, %v), want 0", got, err)
	}
}

// Both sides of each measured threshold: HR ≥ 110, Hb < 100 g/L, urea ≥ 12, CO₂ ≥ 35.
func TestOttawaCOPDRiskScale_edges(t *testing.T) {
	cases := []struct {
		name string
		mod  func(*OttawaCOPDFindings)
		want int
	}{
		{"HR 109", func(f *OttawaCOPDFindings) { f.HeartRateOnArrivalPerMin = 109 }, 0},
		{"HR 110", func(f *OttawaCOPDFindings) { f.HeartRateOnArrivalPerMin = 110 }, 2},
		{"Hb 100", func(f *OttawaCOPDFindings) { f.HemoglobinGL = 100 }, 0},
		{"Hb 99.9", func(f *OttawaCOPDFindings) { f.HemoglobinGL = 99.9 }, 3},
		{"urea 11.9", func(f *OttawaCOPDFindings) { f.UreaMmolL = 11.9 }, 0},
		{"urea 12", func(f *OttawaCOPDFindings) { f.UreaMmolL = 12 }, 1},
		{"CO2 34.9", func(f *OttawaCOPDFindings) { f.SerumCO2MmolL = 34.9 }, 0},
		{"CO2 35", func(f *OttawaCOPDFindings) { f.SerumCO2MmolL = 35 }, 1},
		{"CABG and PVD", func(f *OttawaCOPDFindings) { f.CoronaryBypassGraft, f.PeripheralVascularIntervention = true, true }, 2},
		{"walk test and ECG", func(f *OttawaCOPDFindings) { f.TooIllForWalkTest, f.AcuteIschemicECG = true, true }, 4},
		{"congestion", func(f *OttawaCOPDFindings) { f.PulmonaryCongestionCXR = true }, 1},
	}
	for _, c := range cases {
		f := ottawaBase()
		c.mod(&f)
		if got, err := OttawaCOPDRiskScale(f); err != nil || got != c.want {
			t.Errorf("%s: got (%d, %v), want %d", c.name, got, err, c.want)
		}
	}
}

// parity: ledger ottawa_copd, source case 1: the source's weights (1,1,2,2,2,2,1,3,1,1) agree with Table 7;
// intubation + HR ≥ 110 + Hb < 100 g/L = 7, as the source computes.
func TestOttawaCOPDRiskScale_parity(t *testing.T) {
	f := ottawaBase()
	f.IntubationForRespiratory, f.HeartRateOnArrivalPerMin, f.HemoglobinGL = true, 120, 90
	if got, err := OttawaCOPDRiskScale(f); err != nil || got != 7 {
		t.Errorf("got (%d, %v), want 7", got, err)
	}
}

func TestOttawaCOPDRiskScale_errors(t *testing.T) {
	cases := []struct {
		name string
		mod  func(*OttawaCOPDFindings)
		want error
	}{
		{"HR -1", func(f *OttawaCOPDFindings) { f.HeartRateOnArrivalPerMin = -1 }, ErrOutOfRange},
		{"Hb 0", func(f *OttawaCOPDFindings) { f.HemoglobinGL = 0 }, ErrOutOfRange},
		{"urea -1", func(f *OttawaCOPDFindings) { f.UreaMmolL = -1 }, ErrOutOfRange},
		{"CO2 -1", func(f *OttawaCOPDFindings) { f.SerumCO2MmolL = -1 }, ErrOutOfRange},
		{"HR NaN", func(f *OttawaCOPDFindings) { f.HeartRateOnArrivalPerMin = math.NaN() }, ErrNotFinite},
		{"Hb +Inf", func(f *OttawaCOPDFindings) { f.HemoglobinGL = math.Inf(1) }, ErrNotFinite},
		{"urea -Inf", func(f *OttawaCOPDFindings) { f.UreaMmolL = math.Inf(-1) }, ErrNotFinite},
		{"CO2 NaN", func(f *OttawaCOPDFindings) { f.SerumCO2MmolL = math.NaN() }, ErrNotFinite},
	}
	for _, c := range cases {
		f := ottawaBase()
		c.mod(&f)
		if _, err := OttawaCOPDRiskScale(f); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// Table 8: 0 Low; 1, 2 Medium; 3, 4 High; 5 … 10 and "> 10" Very high.
func TestOttawaCOPDRiskCategory_edges(t *testing.T) {
	cases := []struct {
		total int
		want  OttawaCOPDRisk
	}{{0, OttawaCOPDLow}, {1, OttawaCOPDMedium}, {2, OttawaCOPDMedium}, {3, OttawaCOPDHigh}, {4, OttawaCOPDHigh},
		{5, OttawaCOPDVeryHigh}, {10, OttawaCOPDVeryHigh}, {11, OttawaCOPDVeryHigh}, {16, OttawaCOPDVeryHigh}}
	for _, c := range cases {
		if got, err := OttawaCOPDRiskCategory(c.total); err != nil || got != c.want {
			t.Errorf("%d: got (%q, %v), want %q", c.total, got, err, c.want)
		}
	}
	for _, bad := range []int{-1, 17} {
		if _, err := OttawaCOPDRiskCategory(bad); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("%d: err = %v, want ErrOutOfRange", bad, err)
		}
	}
}

func TestOttawaCOPDRiskScale_noAlloc(t *testing.T) {
	f := ottawaBase()
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = OttawaCOPDRiskScale(f) }); n != 0 {
		t.Errorf("OttawaCOPDRiskScale allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sinkOttawaCOPD, _ = OttawaCOPDRiskCategory(3) }); n != 0 {
		t.Errorf("OttawaCOPDRiskCategory allocs = %v, want 0", n)
	}
}

func BenchmarkOttawaCOPDRiskScale(b *testing.B) {
	f := ottawaBase()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkInt, _ = OttawaCOPDRiskScale(f)
	}
}

func BenchmarkOttawaCOPDRiskCategory(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkOttawaCOPD, _ = OttawaCOPDRiskCategory(3)
	}
}
