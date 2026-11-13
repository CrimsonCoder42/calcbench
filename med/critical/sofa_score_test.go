package critical

import "testing"

var (
	sinkSOFA  SOFAScores
	sinkMSOFA MSOFAScores
)

// sofaZero scores 0 in every system of Singer 2016 Table 1.
var sofaZero = SOFAInput{PaO2FiO2MmHg: 450, PlateletsKPerUL: 200, BilirubinMgDL: 0.8, MAPMmHg: 80, GCS: 15, CreatinineMgDL: 0.9}

// Singer 2016 Table 1 / Lambden 2019 Table 1: both sides of every edge, one system at a time. HAND-DERIVED.
func TestSOFAScore_table1(t *testing.T) {
	tests := []struct {
		name string
		mod  func(*SOFAInput)
		pick func(SOFAScores) int
		want int
	}{
		{"resp 400", func(i *SOFAInput) { i.PaO2FiO2MmHg = 400 }, sofaPickResp, 0},
		{"resp 399", func(i *SOFAInput) { i.PaO2FiO2MmHg = 399 }, sofaPickResp, 1},
		{"resp 300", func(i *SOFAInput) { i.PaO2FiO2MmHg = 300 }, sofaPickResp, 1},
		{"resp 299", func(i *SOFAInput) { i.PaO2FiO2MmHg = 299 }, sofaPickResp, 2},
		{"resp 199 no support", func(i *SOFAInput) { i.PaO2FiO2MmHg = 199 }, sofaPickResp, 2},
		{"resp 200 support", func(i *SOFAInput) { i.PaO2FiO2MmHg = 200; i.RespiratorySupport = true }, sofaPickResp, 2},
		{"resp 199 support", func(i *SOFAInput) { i.PaO2FiO2MmHg = 199; i.RespiratorySupport = true }, sofaPickResp, 3},
		{"resp 100 support", func(i *SOFAInput) { i.PaO2FiO2MmHg = 100; i.RespiratorySupport = true }, sofaPickResp, 3},
		{"resp 99 support", func(i *SOFAInput) { i.PaO2FiO2MmHg = 99; i.RespiratorySupport = true }, sofaPickResp, 4},
		{"plt 150", func(i *SOFAInput) { i.PlateletsKPerUL = 150 }, sofaPickCoag, 0},
		{"plt 149", func(i *SOFAInput) { i.PlateletsKPerUL = 149 }, sofaPickCoag, 1},
		{"plt 100", func(i *SOFAInput) { i.PlateletsKPerUL = 100 }, sofaPickCoag, 1},
		{"plt 99", func(i *SOFAInput) { i.PlateletsKPerUL = 99 }, sofaPickCoag, 2},
		{"plt 50", func(i *SOFAInput) { i.PlateletsKPerUL = 50 }, sofaPickCoag, 2},
		{"plt 49", func(i *SOFAInput) { i.PlateletsKPerUL = 49 }, sofaPickCoag, 3},
		{"plt 20", func(i *SOFAInput) { i.PlateletsKPerUL = 20 }, sofaPickCoag, 3},
		{"plt 19", func(i *SOFAInput) { i.PlateletsKPerUL = 19 }, sofaPickCoag, 4},
		{"bili 1.19", func(i *SOFAInput) { i.BilirubinMgDL = 1.19 }, sofaPickLiver, 0},
		{"bili 1.2", func(i *SOFAInput) { i.BilirubinMgDL = 1.2 }, sofaPickLiver, 1},
		{"bili 1.95", func(i *SOFAInput) { i.BilirubinMgDL = 1.95 }, sofaPickLiver, 1},
		{"bili 2.0", func(i *SOFAInput) { i.BilirubinMgDL = 2.0 }, sofaPickLiver, 2},
		{"bili 5.9", func(i *SOFAInput) { i.BilirubinMgDL = 5.9 }, sofaPickLiver, 2},
		{"bili 6.0", func(i *SOFAInput) { i.BilirubinMgDL = 6.0 }, sofaPickLiver, 3},
		{"bili 11.9", func(i *SOFAInput) { i.BilirubinMgDL = 11.9 }, sofaPickLiver, 3},
		{"bili 12.0", func(i *SOFAInput) { i.BilirubinMgDL = 12.0 }, sofaPickLiver, 4},
		{"map 70", func(i *SOFAInput) { i.MAPMmHg = 70 }, sofaPickCv, 0},
		{"map 69", func(i *SOFAInput) { i.MAPMmHg = 69 }, sofaPickCv, 1},
		{"dopa 5", func(i *SOFAInput) { i.DopamineUgKgMin = 5 }, sofaPickCv, 2},
		{"dobutamine", func(i *SOFAInput) { i.Dobutamine = true }, sofaPickCv, 2},
		{"dopa 5.1", func(i *SOFAInput) { i.DopamineUgKgMin = 5.1 }, sofaPickCv, 3},
		{"dopa 15", func(i *SOFAInput) { i.DopamineUgKgMin = 15 }, sofaPickCv, 3},
		{"dopa 15.1", func(i *SOFAInput) { i.DopamineUgKgMin = 15.1 }, sofaPickCv, 4},
		{"epi 0.1", func(i *SOFAInput) { i.EpinephrineUgKgMin = 0.1 }, sofaPickCv, 3},
		{"epi 0.11", func(i *SOFAInput) { i.EpinephrineUgKgMin = 0.11 }, sofaPickCv, 4},
		{"norepi 0.1", func(i *SOFAInput) { i.NorepinephrineUgKgMin = 0.1 }, sofaPickCv, 3},
		{"norepi 0.11", func(i *SOFAInput) { i.NorepinephrineUgKgMin = 0.11 }, sofaPickCv, 4},
		{"gcs 14", func(i *SOFAInput) { i.GCS = 14 }, sofaPickCns, 1},
		{"gcs 13", func(i *SOFAInput) { i.GCS = 13 }, sofaPickCns, 1},
		{"gcs 12", func(i *SOFAInput) { i.GCS = 12 }, sofaPickCns, 2},
		{"gcs 10", func(i *SOFAInput) { i.GCS = 10 }, sofaPickCns, 2},
		{"gcs 9", func(i *SOFAInput) { i.GCS = 9 }, sofaPickCns, 3},
		{"gcs 6", func(i *SOFAInput) { i.GCS = 6 }, sofaPickCns, 3},
		{"gcs 5", func(i *SOFAInput) { i.GCS = 5 }, sofaPickCns, 4},
		{"cr 1.19", func(i *SOFAInput) { i.CreatinineMgDL = 1.19 }, sofaPickRenal, 0},
		{"cr 1.2", func(i *SOFAInput) { i.CreatinineMgDL = 1.2 }, sofaPickRenal, 1},
		{"cr 2.0", func(i *SOFAInput) { i.CreatinineMgDL = 2.0 }, sofaPickRenal, 2},
		{"cr 3.4", func(i *SOFAInput) { i.CreatinineMgDL = 3.4 }, sofaPickRenal, 2},
		{"cr 3.5", func(i *SOFAInput) { i.CreatinineMgDL = 3.5 }, sofaPickRenal, 3},
		{"cr 4.9", func(i *SOFAInput) { i.CreatinineMgDL = 4.9 }, sofaPickRenal, 3},
		{"cr 5.0", func(i *SOFAInput) { i.CreatinineMgDL = 5.0 }, sofaPickRenal, 4},
		{"uo 500", func(i *SOFAInput) { i.UrineOutputMeasured = true; i.UrineOutputMLPerDay = 500 }, sofaPickRenal, 0},
		{"uo 499", func(i *SOFAInput) { i.UrineOutputMeasured = true; i.UrineOutputMLPerDay = 499 }, sofaPickRenal, 3},
		{"uo 200", func(i *SOFAInput) { i.UrineOutputMeasured = true; i.UrineOutputMLPerDay = 200 }, sofaPickRenal, 3},
		{"uo 199", func(i *SOFAInput) { i.UrineOutputMeasured = true; i.UrineOutputMLPerDay = 199 }, sofaPickRenal, 4},
		{"uo 100 unmeasured", func(i *SOFAInput) { i.UrineOutputMLPerDay = 100 }, sofaPickRenal, 0},
		{"uo 499 cr 5", func(i *SOFAInput) {
			i.CreatinineMgDL = 5
			i.UrineOutputMeasured = true
			i.UrineOutputMLPerDay = 499
		}, sofaPickRenal, 4},
	}
	for _, tc := range tests {
		in := sofaZero
		tc.mod(&in)
		s, err := SOFAScore(in)
		if err != nil || tc.pick(s) != tc.want {
			t.Errorf("%s: (%+v, %v), want %d", tc.name, s, err, tc.want)
		}
	}
}

func sofaPickResp(s SOFAScores) int  { return s.Respiration }
func sofaPickCoag(s SOFAScores) int  { return s.Coagulation }
func sofaPickLiver(s SOFAScores) int { return s.Liver }
func sofaPickCv(s SOFAScores) int    { return s.Cardiovascular }
func sofaPickCns(s SOFAScores) int   { return s.CNS }
func sofaPickRenal(s SOFAScores) int { return s.Renal }

// Doc comment example: 2+1+1+1+1+2 = 8 (HAND-DERIVED from Table 1).
func TestSOFAScore_example(t *testing.T) {
	in := SOFAInput{PaO2FiO2MmHg: 250, PlateletsKPerUL: 120, BilirubinMgDL: 1.5, MAPMmHg: 65, GCS: 14, CreatinineMgDL: 2.1}
	s, err := SOFAScore(in)
	if err != nil || s.Total != 8 {
		t.Errorf("got (%+v, %v), want total 8", s, err)
	}
	max := SOFAInput{PaO2FiO2MmHg: 50, RespiratorySupport: true, PlateletsKPerUL: 10, BilirubinMgDL: 15, MAPMmHg: 50,
		NorepinephrineUgKgMin: 0.5, GCS: 3, CreatinineMgDL: 6}
	if s, _ := SOFAScore(max); s.Total != 24 {
		t.Errorf("maximum: %+v, want 24", s)
	}
}

// parity: ledger sofa_score, source case 1. The app computes respiration (PaO2 / (FiO2% / 100)), coagulation, CNS
// and liver as Table 1 and takes the cardiovascular and renal scores as picked values. PaO2 80 on FiO2 40%
// (ratio 200) ventilated, platelets 90, GCS 13, bilirubin 2.5, MAP < 70 (1), creatinine 2.5 (2) → 2+2+1+2+1+2 = 10.
// Ratio exactly 200 with support scores 2 in both (Table 1: 3 needs "<200").
func TestSOFAScore_parity(t *testing.T) {
	in := SOFAInput{PaO2FiO2MmHg: 200, RespiratorySupport: true, PlateletsKPerUL: 90, BilirubinMgDL: 2.5, MAPMmHg: 65,
		GCS: 13, CreatinineMgDL: 2.5}
	if s, err := SOFAScore(in); err != nil || s.Total != 10 {
		t.Errorf("got (%+v, %v), want 10", s, err)
	}
}

func TestSOFAScore_errors(t *testing.T) {
	fields := []func(*SOFAInput, float64){
		func(i *SOFAInput, v float64) { i.PaO2FiO2MmHg = v },
		func(i *SOFAInput, v float64) { i.PlateletsKPerUL = v },
		func(i *SOFAInput, v float64) { i.BilirubinMgDL = v },
		func(i *SOFAInput, v float64) { i.MAPMmHg = v },
		func(i *SOFAInput, v float64) { i.DopamineUgKgMin = v },
		func(i *SOFAInput, v float64) { i.EpinephrineUgKgMin = v },
		func(i *SOFAInput, v float64) { i.NorepinephrineUgKgMin = v },
		func(i *SOFAInput, v float64) { i.CreatinineMgDL = v },
		func(i *SOFAInput, v float64) { i.UrineOutputMLPerDay = v },
	}
	for k, set := range fields {
		for _, v := range nonFiniteB6 {
			in := sofaZero
			set(&in, v)
			_, err := SOFAScore(in)
			wantErrB6(t, "field nonfinite", err, ErrNotFinite)
		}
		in := sofaZero
		set(&in, -1)
		if _, err := SOFAScore(in); err == nil {
			t.Errorf("field %d = -1: no error", k)
		} else {
			wantErrB6(t, "field -1", err, ErrOutOfRange)
		}
	}
	for _, g := range []int{2, 16} {
		in := sofaZero
		in.GCS = g
		_, err := SOFAScore(in)
		wantErrB6(t, "gcs", err, ErrOutOfRange)
	}
}

func TestSOFAScore_noAlloc(t *testing.T) {
	in := SOFAInput{PaO2FiO2MmHg: 250, PlateletsKPerUL: 120, BilirubinMgDL: 1.5, MAPMmHg: 65, GCS: 14, CreatinineMgDL: 2.1}
	if n := testing.AllocsPerRun(100, func() { sinkSOFA, _ = SOFAScore(in) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkSOFAScore(b *testing.B) {
	in := SOFAInput{PaO2FiO2MmHg: 250, PlateletsKPerUL: 120, BilirubinMgDL: 1.5, MAPMmHg: 65, GCS: 14, CreatinineMgDL: 2.1}
	b.ReportAllocs()
	for b.Loop() {
		sinkSOFA, _ = SOFAScore(in)
	}
}
