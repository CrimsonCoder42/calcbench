package nutrition

import "testing"

// Values are the table's printed decimals, compared exactly (they are copied, not computed).

// IOM 2005 Table S-2, rows "> 19 y" (male) and "Pregnancy 14 through 50 y".
func TestTotalWaterAI_example(t *testing.T) {
	if got, err := TotalWaterAI(360, false, LifeStageNone); err != nil || got != (WaterAI{3.7, 3.0}) {
		t.Errorf("man 30 y: got (%+v, %v), want {3.7 3.0}", got, err)
	}
	if got, err := TotalWaterAI(336, true, LifeStagePregnancy); err != nil || got != (WaterAI{3.0, 2.3}) {
		t.Errorf("pregnancy 28 y: got (%+v, %v), want {3.0 2.3}", got, err)
	}
}

// parity: ledger fluid_intake_ai, source cases 1–12. The source returned mL/day (L/d × 1000); the table's own
// unit L/day is asserted here, so 700 mL → 0.7 L.
func TestTotalWaterAI_parity(t *testing.T) {
	tests := []struct {
		name   string
		months float64
		female bool
		stage  LifeStage
		want   WaterAI
	}{
		{"source case 1", 3, true, LifeStageNone, WaterAI{0.7, 0.7}},
		{"source case 2", 9, false, LifeStageNone, WaterAI{0.8, 0.6}},
		{"source case 3", 24, false, LifeStageNone, WaterAI{1.3, 0.9}},
		{"source case 4", 72, true, LifeStageNone, WaterAI{1.7, 1.2}},
		{"source case 5", 120, false, LifeStageNone, WaterAI{2.4, 1.8}},
		{"source case 6", 120, true, LifeStageNone, WaterAI{2.1, 1.6}},
		{"source case 7", 192, false, LifeStageNone, WaterAI{3.3, 2.6}},
		{"source case 8", 192, true, LifeStageNone, WaterAI{2.3, 1.8}},
		{"source case 9", 360, false, LifeStageNone, WaterAI{3.7, 3.0}},
		{"source case 10", 360, true, LifeStageNone, WaterAI{2.7, 2.2}},
		{"source case 11", 336, true, LifeStagePregnancy, WaterAI{3.0, 2.3}},
		{"source case 12", 420, true, LifeStageLactation, WaterAI{3.8, 3.1}},
	}
	for _, tc := range tests {
		if got, err := TotalWaterAI(tc.months, tc.female, tc.stage); err != nil || got != tc.want {
			t.Errorf("%s: got (%+v, %v), want %+v", tc.name, got, err, tc.want)
		}
	}
}

// Every group edge, each side: 6.99 | 7, 11.99 | 12, 47.99 | 48, 107.99 | 108, 167.99 | 168, 227.99 | 228 months.
func TestTotalWaterAI_groupEdges(t *testing.T) {
	tests := []struct {
		months float64
		want   WaterAI
	}{
		{0, WaterAI{0.7, 0.7}}, {6.99, WaterAI{0.7, 0.7}},
		{7, WaterAI{0.8, 0.6}}, {11.99, WaterAI{0.8, 0.6}},
		{12, WaterAI{1.3, 0.9}}, {47.99, WaterAI{1.3, 0.9}},
		{48, WaterAI{1.7, 1.2}}, {107.99, WaterAI{1.7, 1.2}},
		{108, WaterAI{2.4, 1.8}}, {167.99, WaterAI{2.4, 1.8}},
		{168, WaterAI{3.3, 2.6}}, {227.99, WaterAI{3.3, 2.6}},
		{228, WaterAI{3.7, 3.0}}, {1200, WaterAI{3.7, 3.0}},
	}
	for _, tc := range tests {
		if got, err := TotalWaterAI(tc.months, false, LifeStageNone); err != nil || got != tc.want {
			t.Errorf("%v months: got (%+v, %v), want %+v", tc.months, got, err, tc.want)
		}
	}
	// Pregnancy rows: 14 through 50 y, so 168 and 611.99 months accepted.
	for _, m := range []float64{168, 611.99} {
		if _, err := TotalWaterAI(m, true, LifeStageLactation); err != nil {
			t.Errorf("lactation %v months refused: %v", m, err)
		}
	}
}

func TestTotalWaterAI_errors(t *testing.T) {
	f := func(months float64, female bool, stage LifeStage) func() error {
		return func() error { _, err := TotalWaterAI(months, female, stage); return err }
	}
	runErrCases(t, []errCase{
		// parity: ledger fluid_intake_ai, source refusals: male pregnancy 30 y; pregnancy 13 y; lactation 52 y.
		{"male pregnancy", f(360, false, LifeStagePregnancy), ErrDomain},
		{"pregnancy 13 y", f(156, true, LifeStagePregnancy), ErrOutOfRange},
		{"lactation 52 y", f(624, true, LifeStageLactation), ErrOutOfRange},
		{"pregnancy 167.99 months", f(167.99, true, LifeStagePregnancy), ErrOutOfRange},
		{"pregnancy 612 months", f(612, true, LifeStagePregnancy), ErrOutOfRange},
		{"unknown stage", f(360, true, LifeStage(3)), ErrOutOfRange},
		{"negative stage", f(360, true, LifeStage(-1)), ErrOutOfRange},
		{"months negative", f(-0.01, true, LifeStageNone), ErrOutOfRange},
		{"months NaN", f(nan, true, LifeStageNone), ErrNotFinite},
		{"months +Inf", f(pInf, true, LifeStageNone), ErrNotFinite},
		{"months -Inf", f(nInf, true, LifeStageNone), ErrNotFinite},
	})
}

func TestTotalWaterAI_noAllocation(t *testing.T) {
	var w WaterAI
	if n := testing.AllocsPerRun(100, func() { w, _ = TotalWaterAI(360, true, LifeStageNone) }); n != 0 {
		t.Errorf("TotalWaterAI allocates %v times, want 0", n)
	}
	sink = w.TotalLPerDay
}

func BenchmarkTotalWaterAI(b *testing.B) {
	b.ReportAllocs()
	months := 360.0
	var w WaterAI
	for b.Loop() {
		w, _ = TotalWaterAI(months, true, LifeStageNone)
	}
	sink = w.TotalLPerDay
}
