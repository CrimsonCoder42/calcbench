package critical

import "testing"

var sinkSOFA2 SOFA2Scores

var sofa2Zero = SOFA2Input{GCS: 15, OxygenRatioMmHg: 400, MAPMmHg: 80, BilirubinMgDL: 0.8, CreatinineMgDL: 0.9, PlateletsKPerUL: 200}

// Ranzani 2025 Table 2 and footnotes f, h, l, m, p: both sides of every edge. HAND-DERIVED.
func TestSOFA2Score_table2(t *testing.T) {
	brain := func(s SOFA2Scores) int { return s.Brain }
	respi := func(s SOFA2Scores) int { return s.Respiratory }
	cardio := func(s SOFA2Scores) int { return s.Cardiovascular }
	liv := func(s SOFA2Scores) int { return s.Liver }
	kid := func(s SOFA2Scores) int { return s.Kidney }
	hem := func(s SOFA2Scores) int { return s.Hemostasis }
	tests := []struct {
		name string
		mod  func(*SOFA2Input)
		pick func(SOFA2Scores) int
		want int
	}{
		{"gcs 15", func(*SOFA2Input) {}, brain, 0},
		{"gcs 15 delirium drugs", func(i *SOFA2Input) { i.DeliriumDrugs = true }, brain, 1},
		{"gcs 13", func(i *SOFA2Input) { i.GCS = 13 }, brain, 1},
		{"gcs 12", func(i *SOFA2Input) { i.GCS = 12 }, brain, 2},
		{"gcs 9", func(i *SOFA2Input) { i.GCS = 9 }, brain, 2},
		{"gcs 8", func(i *SOFA2Input) { i.GCS = 8 }, brain, 3},
		{"gcs 6", func(i *SOFA2Input) { i.GCS = 6 }, brain, 3},
		{"gcs 5", func(i *SOFA2Input) { i.GCS = 5 }, brain, 4},
		{"gcs 5 delirium", func(i *SOFA2Input) { i.GCS = 5; i.DeliriumDrugs = true }, brain, 4},
		{"pf 300.1", func(i *SOFA2Input) { i.OxygenRatioMmHg = 300.1 }, respi, 0},
		{"pf 300", func(i *SOFA2Input) { i.OxygenRatioMmHg = 300 }, respi, 1},
		{"pf 225.1", func(i *SOFA2Input) { i.OxygenRatioMmHg = 225.1 }, respi, 1},
		{"pf 225", func(i *SOFA2Input) { i.OxygenRatioMmHg = 225 }, respi, 2},
		{"pf 150 no support capped", func(i *SOFA2Input) { i.OxygenRatioMmHg = 150 }, respi, 2},
		{"pf 60 no support capped", func(i *SOFA2Input) { i.OxygenRatioMmHg = 60 }, respi, 2},
		{"pf 150.1 support", func(i *SOFA2Input) { i.OxygenRatioMmHg = 150.1; i.AdvancedVentSupport = true }, respi, 2},
		{"pf 150 support", func(i *SOFA2Input) { i.OxygenRatioMmHg = 150; i.AdvancedVentSupport = true }, respi, 3},
		{"pf 75.1 support", func(i *SOFA2Input) { i.OxygenRatioMmHg = 75.1; i.AdvancedVentSupport = true }, respi, 3},
		{"pf 75 support", func(i *SOFA2Input) { i.OxygenRatioMmHg = 75; i.AdvancedVentSupport = true }, respi, 4},
		{"pf 75 precluded", func(i *SOFA2Input) { i.OxygenRatioMmHg = 75; i.VentSupportPrecluded = true }, respi, 4},
		{"ecmo", func(i *SOFA2Input) { i.RespiratoryECMO = true }, respi, 4},
		{"sf 250.1", func(i *SOFA2Input) { i.RatioIsSpO2FiO2 = true; i.OxygenRatioMmHg = 250.1 }, respi, 1},
		{"sf 250", func(i *SOFA2Input) { i.RatioIsSpO2FiO2 = true; i.OxygenRatioMmHg = 250 }, respi, 2},
		{"sf 200 support", func(i *SOFA2Input) {
			i.RatioIsSpO2FiO2, i.OxygenRatioMmHg, i.AdvancedVentSupport = true, 200, true
		}, respi, 3},
		{"sf 120.1 support", func(i *SOFA2Input) {
			i.RatioIsSpO2FiO2, i.OxygenRatioMmHg, i.AdvancedVentSupport = true, 120.1, true
		}, respi, 3},
		{"sf 120 support", func(i *SOFA2Input) {
			i.RatioIsSpO2FiO2, i.OxygenRatioMmHg, i.AdvancedVentSupport = true, 120, true
		}, respi, 4},
		{"map 70", func(i *SOFA2Input) { i.MAPMmHg = 70 }, cardio, 0},
		{"map 69.9", func(i *SOFA2Input) { i.MAPMmHg = 69.9 }, cardio, 1},
		{"ne 0.2", func(i *SOFA2Input) { i.NorepiPlusEpiUgKgMin = 0.2 }, cardio, 2},
		{"ne 0.21", func(i *SOFA2Input) { i.NorepiPlusEpiUgKgMin = 0.21 }, cardio, 3},
		{"ne 0.4", func(i *SOFA2Input) { i.NorepiPlusEpiUgKgMin = 0.4 }, cardio, 3},
		{"ne 0.41", func(i *SOFA2Input) { i.NorepiPlusEpiUgKgMin = 0.41 }, cardio, 4},
		{"other only", func(i *SOFA2Input) { i.OtherVasopressorOrInotrope = true }, cardio, 2},
		{"low + other", func(i *SOFA2Input) { i.NorepiPlusEpiUgKgMin = 0.1; i.OtherVasopressorOrInotrope = true }, cardio, 3},
		{"medium + other", func(i *SOFA2Input) { i.NorepiPlusEpiUgKgMin = 0.3; i.OtherVasopressorOrInotrope = true }, cardio, 4},
		{"mechanical", func(i *SOFA2Input) { i.MechanicalSupport = true }, cardio, 4},
		{"dopa 20", func(i *SOFA2Input) { i.DopamineUgKgMin = 20 }, cardio, 2},
		{"dopa 20.1", func(i *SOFA2Input) { i.DopamineUgKgMin = 20.1 }, cardio, 3},
		{"dopa 40", func(i *SOFA2Input) { i.DopamineUgKgMin = 40 }, cardio, 3},
		{"dopa 40.1", func(i *SOFA2Input) { i.DopamineUgKgMin = 40.1 }, cardio, 4},
		{"dopa + low ne", func(i *SOFA2Input) { i.DopamineUgKgMin = 5; i.NorepiPlusEpiUgKgMin = 0.1 }, cardio, 3},
		{"precluded map 70", func(i *SOFA2Input) { i.VasoactivesPrecluded = true; i.MAPMmHg = 70 }, cardio, 0},
		{"precluded map 60", func(i *SOFA2Input) { i.VasoactivesPrecluded = true; i.MAPMmHg = 60 }, cardio, 1},
		{"precluded map 59", func(i *SOFA2Input) { i.VasoactivesPrecluded = true; i.MAPMmHg = 59 }, cardio, 2},
		{"precluded map 50", func(i *SOFA2Input) { i.VasoactivesPrecluded = true; i.MAPMmHg = 50 }, cardio, 2},
		{"precluded map 49", func(i *SOFA2Input) { i.VasoactivesPrecluded = true; i.MAPMmHg = 49 }, cardio, 3},
		{"precluded map 40", func(i *SOFA2Input) { i.VasoactivesPrecluded = true; i.MAPMmHg = 40 }, cardio, 3},
		{"precluded map 39", func(i *SOFA2Input) { i.VasoactivesPrecluded = true; i.MAPMmHg = 39 }, cardio, 4},
		{"bili 1.2", func(i *SOFA2Input) { i.BilirubinMgDL = 1.2 }, liv, 0},
		{"bili 1.21", func(i *SOFA2Input) { i.BilirubinMgDL = 1.21 }, liv, 1},
		{"bili 3.0", func(i *SOFA2Input) { i.BilirubinMgDL = 3.0 }, liv, 1},
		{"bili 3.01", func(i *SOFA2Input) { i.BilirubinMgDL = 3.01 }, liv, 2},
		{"bili 6.0", func(i *SOFA2Input) { i.BilirubinMgDL = 6.0 }, liv, 2},
		{"bili 6.01", func(i *SOFA2Input) { i.BilirubinMgDL = 6.01 }, liv, 3},
		{"bili 12.0", func(i *SOFA2Input) { i.BilirubinMgDL = 12.0 }, liv, 3},
		{"bili 12.01", func(i *SOFA2Input) { i.BilirubinMgDL = 12.01 }, liv, 4},
		{"cr 1.2", func(i *SOFA2Input) { i.CreatinineMgDL = 1.2 }, kid, 0},
		{"cr 1.21", func(i *SOFA2Input) { i.CreatinineMgDL = 1.21 }, kid, 1},
		{"cr 2.0", func(i *SOFA2Input) { i.CreatinineMgDL = 2.0 }, kid, 1},
		{"cr 2.01", func(i *SOFA2Input) { i.CreatinineMgDL = 2.01 }, kid, 2},
		{"cr 3.5", func(i *SOFA2Input) { i.CreatinineMgDL = 3.5 }, kid, 2},
		{"cr 3.51", func(i *SOFA2Input) { i.CreatinineMgDL = 3.51 }, kid, 3},
		{"uo 6-12h", func(i *SOFA2Input) { i.UrineOutput = SOFA2UrineBelow05For6To12h }, kid, 1},
		{"uo 12h", func(i *SOFA2Input) { i.UrineOutput = SOFA2UrineBelow05For12h }, kid, 2},
		{"uo 24h", func(i *SOFA2Input) { i.UrineOutput = SOFA2UrineBelow03For24hOrAnuria12h }, kid, 3},
		{"uo lower than cr", func(i *SOFA2Input) { i.CreatinineMgDL = 3; i.UrineOutput = SOFA2UrineBelow05For6To12h }, kid, 2},
		{"rrt", func(i *SOFA2Input) { i.RRT = true }, kid, 4},
		{"plt 150.1", func(i *SOFA2Input) { i.PlateletsKPerUL = 150.1 }, hem, 0},
		{"plt 150", func(i *SOFA2Input) { i.PlateletsKPerUL = 150 }, hem, 1},
		{"plt 100", func(i *SOFA2Input) { i.PlateletsKPerUL = 100 }, hem, 2},
		{"plt 80.1", func(i *SOFA2Input) { i.PlateletsKPerUL = 80.1 }, hem, 2},
		{"plt 80", func(i *SOFA2Input) { i.PlateletsKPerUL = 80 }, hem, 3},
		{"plt 50.1", func(i *SOFA2Input) { i.PlateletsKPerUL = 50.1 }, hem, 3},
		{"plt 50", func(i *SOFA2Input) { i.PlateletsKPerUL = 50 }, hem, 4},
	}
	for _, tc := range tests {
		in := sofa2Zero
		tc.mod(&in)
		s, err := SOFA2Score(in)
		if err != nil || tc.pick(s) != tc.want {
			t.Errorf("%s: (%+v, %v), want %d", tc.name, s, err, tc.want)
		}
	}
}

// Doc comment example: 1+2+1+1+1+2 = 8 (HAND-DERIVED); and the maximum 24.
func TestSOFA2Score_example(t *testing.T) {
	in := SOFA2Input{GCS: 14, OxygenRatioMmHg: 200, MAPMmHg: 65, BilirubinMgDL: 2.0, CreatinineMgDL: 1.5, PlateletsKPerUL: 90}
	if s, err := SOFA2Score(in); err != nil || s.Total != 8 {
		t.Errorf("got (%+v, %v), want 8", s, err)
	}
	max := SOFA2Input{GCS: 3, OxygenRatioMmHg: 50, AdvancedVentSupport: true, MAPMmHg: 50, NorepiPlusEpiUgKgMin: 0.6,
		BilirubinMgDL: 15, CreatinineMgDL: 4, RRT: true, PlateletsKPerUL: 20}
	if s, _ := SOFA2Score(max); s.Total != 24 {
		t.Errorf("maximum: %+v, want 24", s)
	}
}

// parity: ledger sofa2_score, source case 1. The app sums six picked organ scores whose option labels are the
// Table 2 rows; the same rows from raw values: brain 2 (GCS 10), respiratory 3 (225 → 2 without support; 150 with
// support → 3), cardiovascular 2 (NE 0.1), liver 1, kidney 2, hemostasis 1 → 11.
func TestSOFA2Score_parity(t *testing.T) {
	in := SOFA2Input{GCS: 10, OxygenRatioMmHg: 150, AdvancedVentSupport: true, MAPMmHg: 75, NorepiPlusEpiUgKgMin: 0.1,
		BilirubinMgDL: 2.5, CreatinineMgDL: 3.0, PlateletsKPerUL: 140}
	if s, err := SOFA2Score(in); err != nil || s.Total != 11 {
		t.Errorf("got (%+v, %v), want 11", s, err)
	}
}

func TestSOFA2Score_errors(t *testing.T) {
	fields := []func(*SOFA2Input, float64){
		func(i *SOFA2Input, v float64) { i.OxygenRatioMmHg = v },
		func(i *SOFA2Input, v float64) { i.MAPMmHg = v },
		func(i *SOFA2Input, v float64) { i.NorepiPlusEpiUgKgMin = v },
		func(i *SOFA2Input, v float64) { i.DopamineUgKgMin = v },
		func(i *SOFA2Input, v float64) { i.BilirubinMgDL = v },
		func(i *SOFA2Input, v float64) { i.CreatinineMgDL = v },
		func(i *SOFA2Input, v float64) { i.PlateletsKPerUL = v },
	}
	for _, set := range fields {
		for _, v := range nonFiniteB6 {
			in := sofa2Zero
			set(&in, v)
			_, err := SOFA2Score(in)
			wantErrB6(t, "nonfinite", err, ErrNotFinite)
		}
		in := sofa2Zero
		set(&in, -1)
		_, err := SOFA2Score(in)
		wantErrB6(t, "-1", err, ErrOutOfRange)
	}
	in := sofa2Zero
	in.GCS = 2
	_, err := SOFA2Score(in)
	wantErrB6(t, "gcs 2", err, ErrOutOfRange)
	in = sofa2Zero
	in.UrineOutput = 4
	_, err = SOFA2Score(in)
	wantErrB6(t, "urine 4", err, ErrOutOfRange)
	in.UrineOutput = -1
	_, err = SOFA2Score(in)
	wantErrB6(t, "urine -1", err, ErrOutOfRange)
}

func TestSOFA2Score_noAlloc(t *testing.T) {
	in := SOFA2Input{GCS: 14, OxygenRatioMmHg: 200, MAPMmHg: 65, BilirubinMgDL: 2.0, CreatinineMgDL: 1.5, PlateletsKPerUL: 90}
	if n := testing.AllocsPerRun(100, func() { sinkSOFA2, _ = SOFA2Score(in) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkSOFA2Score(b *testing.B) {
	in := SOFA2Input{GCS: 14, OxygenRatioMmHg: 200, MAPMmHg: 65, BilirubinMgDL: 2.0, CreatinineMgDL: 1.5, PlateletsKPerUL: 90}
	b.ReportAllocs()
	for b.Loop() {
		sinkSOFA2, _ = SOFA2Score(in)
	}
}
