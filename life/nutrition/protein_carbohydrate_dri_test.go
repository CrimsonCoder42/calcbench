package nutrition

import "testing"

// Values are the table's printed numbers, compared exactly (copied, not computed).

// NASEM 2019 Appendix J p. 572: Males 14–18 y protein 52, carbohydrate 130 (RDAs); Infants 7–12 mo protein
// 11.0 (RDA), carbohydrate 95* (AI).
func TestProteinCarbohydrateDRIFor_example(t *testing.T) {
	if got, err := ProteinCarbohydrateDRIFor(16*12, false, LifeStageNone); err != nil || got != (ProteinCarbohydrateDRI{52, 130, false, false}) {
		t.Errorf("male 16 y: got (%+v, %v)", got, err)
	}
	if got, err := ProteinCarbohydrateDRIFor(9, true, LifeStageNone); err != nil || got != (ProteinCarbohydrateDRI{11.0, 95, false, true}) {
		t.Errorf("9 months: got (%+v, %v)", got, err)
	}
}

// parity: ledger dri_protein_carbohydrate, source cases 1–10 (protein g, carbohydrate g). Case 1 (6 months):
// the source returned 11 / 95 from DGA Table A1-1 "6 Through 11 Months"; the primary DRI table's row
// "Infants 0–6 mo" holds 6 completed months, so the primary's 9.1* / 60* is asserted.
func TestProteinCarbohydrateDRIFor_parity(t *testing.T) {
	tests := []struct {
		name           string
		months         float64
		female         bool
		protein, carbs float64
	}{
		{"source case 1 (primary value; source 11 / 95)", 6, true, 9.1, 60},
		{"source case 2", 11, false, 11, 95},
		{"source case 3", 18, false, 13, 130},
		{"source case 4", 36, true, 13, 130},
		{"source case 5", 6 * 12, true, 19, 130},
		{"source case 6", 13 * 12, false, 34, 130},
		{"source case 7", 16 * 12, false, 52, 130},
		{"source case 8", 18 * 12, true, 46, 130},
		{"source case 9", 19 * 12, false, 56, 130},
		{"source case 10", 70 * 12, true, 46, 130},
	}
	for _, tc := range tests {
		got, err := ProteinCarbohydrateDRIFor(tc.months, tc.female, LifeStageNone)
		if err != nil || got.ProteinGPerDay != tc.protein || got.CarbohydrateGPerDay != tc.carbs {
			t.Errorf("%s: got (%+v, %v), want %v / %v", tc.name, got, err, tc.protein, tc.carbs)
		}
	}
}

// parity: ledger dri_protein_carbohydrate, source refusal case (4 months). The source refused under 6 months
// (DGA Appendix 1 starts at 6 months); the primary's row "Infants 0–6 mo" gives 9.1* / 60*, asserted here.
func TestProteinCarbohydrateDRIFor_parityFourMonths(t *testing.T) {
	if got, err := ProteinCarbohydrateDRIFor(4, true, LifeStageNone); err != nil || got != (ProteinCarbohydrateDRI{9.1, 60, true, true}) {
		t.Errorf("got (%+v, %v), want {9.1 60 true true}", got, err)
	}
}

// Every row of the table, with both sides of each age edge.
func TestProteinCarbohydrateDRIFor_rows(t *testing.T) {
	tests := []struct {
		months float64
		female bool
		stage  LifeStage
		want   ProteinCarbohydrateDRI
	}{
		{0, false, LifeStageNone, ProteinCarbohydrateDRI{9.1, 60, true, true}},
		{6.99, false, LifeStageNone, ProteinCarbohydrateDRI{9.1, 60, true, true}},
		{7, false, LifeStageNone, ProteinCarbohydrateDRI{11.0, 95, false, true}},
		{11.99, false, LifeStageNone, ProteinCarbohydrateDRI{11.0, 95, false, true}},
		{12, false, LifeStageNone, ProteinCarbohydrateDRI{13, 130, false, false}},
		{47.99, false, LifeStageNone, ProteinCarbohydrateDRI{13, 130, false, false}},
		{48, false, LifeStageNone, ProteinCarbohydrateDRI{19, 130, false, false}},
		{107.99, true, LifeStageNone, ProteinCarbohydrateDRI{19, 130, false, false}},
		{108, true, LifeStageNone, ProteinCarbohydrateDRI{34, 130, false, false}},
		{167.99, false, LifeStageNone, ProteinCarbohydrateDRI{34, 130, false, false}},
		{168, false, LifeStageNone, ProteinCarbohydrateDRI{52, 130, false, false}},
		{168, true, LifeStageNone, ProteinCarbohydrateDRI{46, 130, false, false}},
		{227.99, false, LifeStageNone, ProteinCarbohydrateDRI{52, 130, false, false}},
		{228, false, LifeStageNone, ProteinCarbohydrateDRI{56, 130, false, false}},
		{228, true, LifeStageNone, ProteinCarbohydrateDRI{46, 130, false, false}},
		{900, false, LifeStageNone, ProteinCarbohydrateDRI{56, 130, false, false}},
		{168, true, LifeStagePregnancy, ProteinCarbohydrateDRI{71, 175, false, false}},
		{611.99, true, LifeStagePregnancy, ProteinCarbohydrateDRI{71, 175, false, false}},
		{300, true, LifeStageLactation, ProteinCarbohydrateDRI{71, 210, false, false}},
	}
	for _, tc := range tests {
		if got, err := ProteinCarbohydrateDRIFor(tc.months, tc.female, tc.stage); err != nil || got != tc.want {
			t.Errorf("(%v, %v, %v): got (%+v, %v), want %+v", tc.months, tc.female, tc.stage, got, err, tc.want)
		}
	}
}

func TestProteinCarbohydrateDRIFor_errors(t *testing.T) {
	f := func(months float64, female bool, stage LifeStage) func() error {
		return func() error { _, err := ProteinCarbohydrateDRIFor(months, female, stage); return err }
	}
	runErrCases(t, []errCase{
		{"male lactation", f(300, false, LifeStageLactation), ErrDomain},
		{"pregnancy 167.99 months", f(167.99, true, LifeStagePregnancy), ErrOutOfRange},
		{"pregnancy 612 months", f(612, true, LifeStagePregnancy), ErrOutOfRange},
		{"unknown stage", f(300, true, LifeStage(7)), ErrOutOfRange},
		{"months negative", f(-1, false, LifeStageNone), ErrOutOfRange},
		{"months NaN", f(nan, false, LifeStageNone), ErrNotFinite},
		{"months +Inf", f(pInf, false, LifeStageNone), ErrNotFinite},
		{"months -Inf", f(nInf, false, LifeStageNone), ErrNotFinite},
	})
}

func TestProteinCarbohydrateDRIFor_noAllocation(t *testing.T) {
	var d ProteinCarbohydrateDRI
	if n := testing.AllocsPerRun(100, func() { d, _ = ProteinCarbohydrateDRIFor(360, true, LifeStageNone) }); n != 0 {
		t.Errorf("ProteinCarbohydrateDRIFor allocates %v times, want 0", n)
	}
	sink = d.ProteinGPerDay
}

func BenchmarkProteinCarbohydrateDRIFor(b *testing.B) {
	b.ReportAllocs()
	months := 360.0
	var d ProteinCarbohydrateDRI
	for b.Loop() {
		d, _ = ProteinCarbohydrateDRIFor(months, true, LifeStageNone)
	}
	sink = d.ProteinGPerDay
}
