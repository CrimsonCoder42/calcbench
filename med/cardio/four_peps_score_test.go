package cardio

import (
	"math"
	"testing"
)

var sinkFourPEPS FourPEPSProbability

var fourPEPSBase = FourPEPSFindings{AgeYears: 70, HeartRateBpm: 90, SpO2Percent: 98}

// Roy 2021 Table 3 (PMC7931139): each item's points, both sides of each threshold. HAND-DERIVED totals.
func TestFourPEPSScore_table(t *testing.T) {
	cases := []struct {
		name string
		mod  func(*FourPEPSFindings)
		want int
	}{
		{"base", func(f *FourPEPSFindings) {}, 0},
		{"age 49", func(f *FourPEPSFindings) { f.AgeYears = 49 }, -2},
		{"age 50", func(f *FourPEPSFindings) { f.AgeYears = 50 }, -1},
		{"age 64", func(f *FourPEPSFindings) { f.AgeYears = 64 }, -1},
		{"age 65", func(f *FourPEPSFindings) { f.AgeYears = 65 }, 0},
		{"male", func(f *FourPEPSFindings) { f.Male = true }, 2},
		{"chronic resp", func(f *FourPEPSFindings) { f.ChronicRespiratoryDisease = true }, -1},
		{"HR 79", func(f *FourPEPSFindings) { f.HeartRateBpm = 79 }, -1},
		{"HR 80", func(f *FourPEPSFindings) { f.HeartRateBpm = 80 }, 0},
		{"chest pain + dyspnea", func(f *FourPEPSFindings) { f.ChestPainAndAcuteDyspnea = true }, 1},
		{"estrogen", func(f *FourPEPSFindings) { f.EstrogenicTreatment = true }, 2},
		{"VTE", func(f *FourPEPSFindings) { f.PersonalHistoryVTE = true }, 2},
		{"syncope", func(f *FourPEPSFindings) { f.Syncope = true }, 2},
		{"immobility", func(f *FourPEPSFindings) { f.ImmobilityWithin4Weeks = true }, 2},
		{"SpO2 94.9", func(f *FourPEPSFindings) { f.SpO2Percent = 94.9 }, 3},
		{"SpO2 95", func(f *FourPEPSFindings) { f.SpO2Percent = 95 }, 0},
		{"calf", func(f *FourPEPSFindings) { f.CalfPainOrUnilateralEdema = true }, 3},
		{"PE most likely", func(f *FourPEPSFindings) { f.PEMostLikelyDiagnosis = true }, 5},
		{"minimum", func(f *FourPEPSFindings) { f.AgeYears = 30; f.ChronicRespiratoryDisease = true; f.HeartRateBpm = 60 }, -4},
		{"maximum", func(f *FourPEPSFindings) {
			*f = FourPEPSFindings{AgeYears: 70, Male: true, HeartRateBpm: 110, ChestPainAndAcuteDyspnea: true, EstrogenicTreatment: true, PersonalHistoryVTE: true, Syncope: true, ImmobilityWithin4Weeks: true, SpO2Percent: 90, CalfPainOrUnilateralEdema: true, PEMostLikelyDiagnosis: true}
		}, 22},
		{"doc example", func(f *FourPEPSFindings) {
			f.Male = true
			f.HeartRateBpm = 95
			f.SpO2Percent = 93
			f.PEMostLikelyDiagnosis = true
		}, 10},
	}
	for _, c := range cases {
		f := fourPEPSBase
		c.mod(&f)
		got, err := FourPEPSScore(f)
		if err != nil || got != c.want {
			t.Errorf("%s: got %d, %v; want %d", c.name, got, err, c.want)
		}
	}
}

// parity: ledger four_peps, source case 1: the consumer sums age (−2 < 50, −1 ≤ 64), male 2, chronic resp −1,
// HR < 80 −1, chest pain + dyspnea 1, estrogen 2, VTE 2, syncope 2, immobility 2, SpO2 < 95 3, calf 3, PE 5;
// woman 45, HR < 80, estrogen, calf → −2 − 1 + 2 + 3 = 2 → Low CPP.
func TestFourPEPSScore_parity(t *testing.T) {
	f := FourPEPSFindings{AgeYears: 45, HeartRateBpm: 70, SpO2Percent: 97, EstrogenicTreatment: true, CalfPainOrUnilateralEdema: true}
	got, err := FourPEPSScore(f)
	if err != nil || got != 2 {
		t.Fatalf("got %d, %v; want 2", got, err)
	}
	if lv, _ := FourPEPSLevel(got); lv != FourPEPSLow {
		t.Errorf("level %q, want %q", lv, FourPEPSLow)
	}
}

func TestFourPEPSScore_errors(t *testing.T) {
	cases := []struct {
		name string
		mod  func(*FourPEPSFindings)
		want error
	}{
		{"age NaN", func(f *FourPEPSFindings) { f.AgeYears = math.NaN() }, ErrNotFinite},
		{"HR +Inf", func(f *FourPEPSFindings) { f.HeartRateBpm = math.Inf(1) }, ErrNotFinite},
		{"SpO2 -Inf", func(f *FourPEPSFindings) { f.SpO2Percent = math.Inf(-1) }, ErrNotFinite},
		{"age 17", func(f *FourPEPSFindings) { f.AgeYears = 17 }, ErrOutOfRange},
		{"HR negative", func(f *FourPEPSFindings) { f.HeartRateBpm = -1 }, ErrOutOfRange},
		{"SpO2 101", func(f *FourPEPSFindings) { f.SpO2Percent = 101 }, ErrOutOfRange},
	}
	for _, c := range cases {
		f := fourPEPSBase
		c.mod(&f)
		_, err := FourPEPSScore(f)
		checkErr(t, c.name, err, c.want)
	}
}

// Roy 2021 Table 3: both sides of each CPP edge.
func TestFourPEPSLevel_edges(t *testing.T) {
	cases := []struct {
		s    int
		want FourPEPSProbability
	}{{-4, FourPEPSVeryLow}, {-1, FourPEPSVeryLow}, {0, FourPEPSLow}, {5, FourPEPSLow}, {6, FourPEPSModerate}, {12, FourPEPSModerate}, {13, FourPEPSHigh}, {22, FourPEPSHigh}}
	for _, c := range cases {
		if got, err := FourPEPSLevel(c.s); err != nil || got != c.want {
			t.Errorf("%d: got %q, %v; want %q", c.s, got, err, c.want)
		}
	}
	for _, s := range []int{-5, 23} {
		_, err := FourPEPSLevel(s)
		checkErr(t, "out of range", err, ErrOutOfRange)
	}
}

func TestFourPEPS_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = FourPEPSScore(fourPEPSBase) }); n != 0 {
		t.Errorf("FourPEPSScore allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sinkFourPEPS, _ = FourPEPSLevel(7) }); n != 0 {
		t.Errorf("FourPEPSLevel allocs = %v, want 0", n)
	}
}

func BenchmarkFourPEPSScore(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkInt, _ = FourPEPSScore(fourPEPSBase)
	}
}

func BenchmarkFourPEPSLevel(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkFourPEPS, _ = FourPEPSLevel(i%27 - 4)
	}
}
