package gastro

import (
	"math"
	"testing"
)

// Jin 2024 Table 1 (PMC11214112): each item's points, both sides of every category edge, and the range ends
// −3 and 14. Integer points, exact comparison.
func TestAcutePancreatitisDiagnosisScore_example(t *testing.T) {
	base := AcutePancreatitisFindings{LipaseTimesULN: 3}
	with := func(mod func(*AcutePancreatitisFindings)) AcutePancreatitisFindings { f := base; mod(&f); return f }
	tests := []struct {
		name string
		f    AcutePancreatitisFindings
		want int
	}{
		{"worked", AcutePancreatitisFindings{PriorEpisodes: 1, Cholelithiasis: true, EpigastricPain: true, WorseningSeverity: true, PainSeverity: 8, LipaseTimesULN: 25}, 11},
		{"all zero", base, 0},
		{"minimum −3", with(func(f *AcutePancreatitisFindings) {
			f.AbdominalSurgeryPrior2Months = true
			f.PainDurationAtLeast5Days = true
		}), -3},
		{"maximum 14", AcutePancreatitisFindings{PriorEpisodes: 4, Cholelithiasis: true, EpigastricPain: true, WorseningSeverity: true, PainSeverity: 10, LipaseTimesULN: 20}, 14},
		{"3 prior episodes", with(func(f *AcutePancreatitisFindings) { f.PriorEpisodes = 3 }), 3},
		{"4 prior episodes", with(func(f *AcutePancreatitisFindings) { f.PriorEpisodes = 4 }), 4},
		{"9 prior episodes capped at 4", with(func(f *AcutePancreatitisFindings) { f.PriorEpisodes = 9 }), 4},
		{"pain 3", with(func(f *AcutePancreatitisFindings) { f.PainSeverity = 3 }), 0},
		{"pain 4", with(func(f *AcutePancreatitisFindings) { f.PainSeverity = 4 }), 2},
		{"pain 6", with(func(f *AcutePancreatitisFindings) { f.PainSeverity = 6 }), 2},
		{"pain 7", with(func(f *AcutePancreatitisFindings) { f.PainSeverity = 7 }), 3},
		{"lipase 9.99", with(func(f *AcutePancreatitisFindings) { f.LipaseTimesULN = 9.99 }), 0},
		{"lipase 10", with(func(f *AcutePancreatitisFindings) { f.LipaseTimesULN = 10 }), 1},
		{"lipase 19.99", with(func(f *AcutePancreatitisFindings) { f.LipaseTimesULN = 19.99 }), 1},
		{"lipase 20", with(func(f *AcutePancreatitisFindings) { f.LipaseTimesULN = 20 }), 2},
	}
	for _, tc := range tests {
		if got, err := AcutePancreatitisDiagnosisScore(tc.f); err != nil || got != tc.want {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// parity: ledger acute_pancreatitis_prediction. The app takes category indices: pain_severity 0/1/2 (mild,
// moderate, severe) and lipase 0/1/2 (3–10, 10–20, ≥ 20 × ULN); here they are passed as a pain of 2, 5 or 8 and
// a lipase of 5, 15 or 25 × ULN. Its defaults (all 0) → 0.
func TestAcutePancreatitisDiagnosisScore_parity(t *testing.T) {
	pain := [3]int{2, 5, 8}
	lipase := [3]float64{5, 15, 25}
	cases := []struct {
		episodes                         int
		chole, surgery, epi, worse, long bool
		painIdx, lipIdx, want            int
	}{
		{0, false, false, false, false, false, 0, 0, 0}, // parity: ledger acute_pancreatitis_prediction, source case 1 (defaults)
		{4, true, false, true, true, false, 2, 2, 14},   // source case 2: 4+2+2+1+3+2
		{0, false, true, false, false, true, 0, 0, -3},  // source case 3: −2 −1
		{2, true, false, true, false, false, 1, 1, 9},   // source case 4: 2+2+2+2+1
	}
	for i, c := range cases {
		f := AcutePancreatitisFindings{c.episodes, c.chole, c.surgery, c.epi, c.worse, c.long, pain[c.painIdx], lipase[c.lipIdx]}
		if got, err := AcutePancreatitisDiagnosisScore(f); err != nil || got != c.want {
			t.Errorf("case %d: got (%v, %v), want %v", i+1, got, err, c.want)
		}
	}
}

func TestAcutePancreatitisDiagnosisScore_errors(t *testing.T) {
	call := func(f AcutePancreatitisFindings) func() error {
		return func() error { _, err := AcutePancreatitisDiagnosisScore(f); return err }
	}
	runErrCases(t, []errCase{
		{"lipase NaN", call(AcutePancreatitisFindings{LipaseTimesULN: nan}), ErrNotFinite},
		{"lipase +Inf", call(AcutePancreatitisFindings{LipaseTimesULN: pInf}), ErrNotFinite},
		{"lipase -Inf", call(AcutePancreatitisFindings{LipaseTimesULN: nInf}), ErrNotFinite},
		{"NaN before range", call(AcutePancreatitisFindings{PriorEpisodes: -1, LipaseTimesULN: math.NaN()}), ErrNotFinite},
		{"lipase 2.99", call(AcutePancreatitisFindings{LipaseTimesULN: 2.99}), ErrOutOfRange},
		{"episodes −1", call(AcutePancreatitisFindings{PriorEpisodes: -1, LipaseTimesULN: 5}), ErrOutOfRange},
		{"pain −1", call(AcutePancreatitisFindings{PainSeverity: -1, LipaseTimesULN: 5}), ErrOutOfRange},
		{"pain 11", call(AcutePancreatitisFindings{PainSeverity: 11, LipaseTimesULN: 5}), ErrOutOfRange},
	})
}

func TestAcutePancreatitisDiagnosisScore_noAllocation(t *testing.T) {
	f := AcutePancreatitisFindings{PriorEpisodes: 1, Cholelithiasis: true, PainSeverity: 8, LipaseTimesULN: 25}
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = AcutePancreatitisDiagnosisScore(f) }); n != 0 {
		t.Errorf("AcutePancreatitisDiagnosisScore allocates %v times, want 0", n)
	}
}

func BenchmarkAcutePancreatitisDiagnosisScore(b *testing.B) {
	b.ReportAllocs()
	f := AcutePancreatitisFindings{PriorEpisodes: 1, Cholelithiasis: true, PainSeverity: 8, LipaseTimesULN: 25}
	for b.Loop() {
		sinkInt, _ = AcutePancreatitisDiagnosisScore(f)
	}
}
