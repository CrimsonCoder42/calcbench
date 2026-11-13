package gastro

import "testing"

var sinkISEE ISEECategory

// HAND-DERIVED from Dellon 2022 Table 1: each cell alone, the worked example, and the eosinophil edges.
func TestISEEScore_example(t *testing.T) {
	tests := []struct {
		name string
		f    ISEEFindings
		want int
	}{
		{"nothing", ISEEFindings{}, 0},
		{"weekly", ISEEFindings{Symptoms: ISEESymptomsWeekly}, 1},
		{"daily", ISEEFindings{Symptoms: ISEESymptomsDaily}, 2},
		{"multiple daily", ISEEFindings{Symptoms: ISEESymptomsMultiple}, 4},
		{"food impaction adult", ISEEFindings{FoodImpactionERVisit: true}, 2},
		{"food impaction child", ISEEFindings{FoodImpactionERVisit: true, Under18: true}, 4},
		{"under 18 alone", ISEEFindings{Under18: true}, 0},
		{"hospitalization", ISEEFindings{HospitalizationDueToEoE: true}, 4},
		{"perforation", ISEEFindings{EsophagealPerforation: true}, 15},
		{"three 15-point complications add", ISEEFindings{EsophagealPerforation: true, Malnutrition: true, PersistentInflammation: true}, 45},
		{"localized", ISEEFindings{EndoscopicInflammation: ISEEExtentLocalized}, 1},
		{"diffuse", ISEEFindings{EndoscopicInflammation: ISEEExtentDiffuse}, 2},
		{"eos 14.9", ISEEFindings{PeakEosinophilsPerHPF: 14.9}, 0},
		{"eos 15", ISEEFindings{PeakEosinophilsPerHPF: 15}, 1},
		{"eos 60", ISEEFindings{PeakEosinophilsPerHPF: 60}, 1},
		{"eos 60.1", ISEEFindings{PeakEosinophilsPerHPF: 60.1}, 2},
		{"rings pass easily", ISEEFindings{EndoscopicFibrostenosis: ISEEFibrostenosisPassesEasily}, 1},
		{"rings dilation or snug", ISEEFindings{EndoscopicFibrostenosis: ISEEFibrostenosisDilationOrSnug}, 2},
		{"cannot pass", ISEEFindings{EndoscopicFibrostenosis: ISEEFibrostenosisCannotPass}, 15},
		{"BZH/LPF", ISEEFindings{FibrostenoticHistology: true}, 2},
		{"worked", ISEEFindings{Symptoms: ISEESymptomsDaily, FoodImpactionERVisit: true, EndoscopicInflammation: ISEEExtentDiffuse,
			PeakEosinophilsPerHPF: 70, EndoscopicFibrostenosis: ISEEFibrostenosisPassesEasily}, 9},
	}
	for _, tc := range tests {
		if got, err := ISEEScore(tc.f); err != nil || got != tc.want {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// parity: ledger isee_eoe, the app's point arrays (symptoms [0,1,2,4], food impaction [0,2,4], hospitalization 4,
// perforation/malnutrition/inflammation 15 each, edema [0,1,2], eos [0,1,2], rings [0,1,2,15], BZH 2) agree with
// Table 1. All-maximum picker values → 4 + 4 + 4 + 45 + 2 + 2 + 15 + 2 = 78.
func TestISEEScore_parity(t *testing.T) {
	f := ISEEFindings{ISEESymptomsMultiple, true, true, true, true, true, true, ISEEExtentDiffuse, 100, ISEEFibrostenosisCannotPass, true}
	if got, err := ISEEScore(f); err != nil || got != 78 {
		t.Errorf("all maximum: got (%v, %v), want 78", got, err)
	}
}

// Both sides of every category edge (Table 1).
func TestISEESeverity_edges(t *testing.T) {
	tests := []struct {
		total int
		want  ISEECategory
	}{
		{0, ISEEInactive}, {1, ISEEMildActive}, {6, ISEEMildActive}, {7, ISEEModerateActive},
		{14, ISEEModerateActive}, {15, ISEESevereActive}, {78, ISEESevereActive},
	}
	for _, tc := range tests {
		if got, err := ISEESeverity(tc.total); err != nil || got != tc.want {
			t.Errorf("ISEESeverity(%d) = (%q, %v), want %q", tc.total, got, err, tc.want)
		}
	}
}

func TestISEE_errors(t *testing.T) {
	score := func(f ISEEFindings) func() error {
		return func() error { _, err := ISEEScore(f); return err }
	}
	runErrCases(t, []errCase{
		{"eos NaN", score(ISEEFindings{PeakEosinophilsPerHPF: nan}), ErrNotFinite},
		{"eos +Inf", score(ISEEFindings{PeakEosinophilsPerHPF: pInf}), ErrNotFinite},
		{"eos -Inf", score(ISEEFindings{PeakEosinophilsPerHPF: nInf}), ErrNotFinite},
		{"eos negative", score(ISEEFindings{PeakEosinophilsPerHPF: -1}), ErrOutOfRange},
		{"symptoms unknown", score(ISEEFindings{Symptoms: 4}), ErrOutOfRange},
		{"inflammation unknown", score(ISEEFindings{EndoscopicInflammation: 3}), ErrOutOfRange},
		{"fibrostenosis unknown", score(ISEEFindings{EndoscopicFibrostenosis: -1}), ErrOutOfRange},
		{"severity negative", func() error { _, err := ISEESeverity(-1); return err }, ErrOutOfRange},
	})
}

func TestISEE_noAllocation(t *testing.T) {
	f := ISEEFindings{Symptoms: ISEESymptomsDaily, PeakEosinophilsPerHPF: 40}
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = ISEEScore(f) }); n != 0 {
		t.Errorf("ISEEScore allocates %v times, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sinkISEE, _ = ISEESeverity(9) }); n != 0 {
		t.Errorf("ISEESeverity allocates %v times, want 0", n)
	}
}

func BenchmarkISEEScore(b *testing.B) {
	b.ReportAllocs()
	f := ISEEFindings{Symptoms: ISEESymptomsDaily, FoodImpactionERVisit: true, PeakEosinophilsPerHPF: 70}
	for b.Loop() {
		sinkInt, _ = ISEEScore(f)
	}
}

func BenchmarkISEESeverity(b *testing.B) {
	b.ReportAllocs()
	v := 9
	for b.Loop() {
		sinkISEE, _ = ISEESeverity(v)
	}
}
