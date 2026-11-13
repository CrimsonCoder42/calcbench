package critical

import (
	"errors"
	"math"
	"testing"
)

var sinkSAPS3 SAPS3Result

// saps3Base is an admission in every 0-point band: its score is the offset 16 alone (Table 2 footnote 12).
func saps3Base() SAPS3Input {
	return SAPS3Input{
		AgeYears: 30, SurgicalStatus: SAPS3ScheduledSurgery, SurgicalSite: SAPS3SiteAllOthers,
		GCS: 15, TemperatureC: 37, HeartRatePerMin: 80, LeukocytesGL: 8, PH: 7.4, PlateletsGL: 200,
		SystolicMmHg: 130, PaO2MmHg: 90, FiO2: 0.21,
	}
}

func saps3Points(t *testing.T, in SAPS3Input, biliMgDL, creatMgDL float64) int {
	t.Helper()
	r, err := SAPS3Score(in, biliMgDL, creatMgDL)
	if err != nil {
		t.Fatalf("SAPS3Score(%+v): %v", in, err)
	}
	return r.Score - 16
}

// Moreno 2005 (PMC1315315) Table 1, Box I, Box II, Box III and Table 2, read cell by cell from the XML:
// every printed band, both sides of each edge, one item at a time over the 16-point base.
func TestSAPS3Score_points(t *testing.T) {
	if got := saps3Points(t, saps3Base(), 0.5, 0.8); got != 0 {
		t.Fatalf("base → %d over 16, want 0", got)
	}
	type mod func(*SAPS3Input)
	tests := []struct {
		name string
		f    mod
		want int
	}{
		{"age 39.9", func(in *SAPS3Input) { in.AgeYears = 39.9 }, 0},
		{"age 40", func(in *SAPS3Input) { in.AgeYears = 40 }, 5},
		{"age 59.9", func(in *SAPS3Input) { in.AgeYears = 59.9 }, 5},
		{"age 60", func(in *SAPS3Input) { in.AgeYears = 60 }, 9},
		{"age 69.9", func(in *SAPS3Input) { in.AgeYears = 69.9 }, 9},
		{"age 70", func(in *SAPS3Input) { in.AgeYears = 70 }, 13},
		{"age 74.9", func(in *SAPS3Input) { in.AgeYears = 74.9 }, 13},
		{"age 75", func(in *SAPS3Input) { in.AgeYears = 75 }, 15},
		{"age 79.9", func(in *SAPS3Input) { in.AgeYears = 79.9 }, 15},
		{"age 80", func(in *SAPS3Input) { in.AgeYears = 80 }, 18},
		{"cancer therapy", func(in *SAPS3Input) { in.CancerTherapy = true }, 3},
		{"chronic HF", func(in *SAPS3Input) { in.ChronicHeartFailureNYHAIV = true }, 6},
		{"haematological cancer", func(in *SAPS3Input) { in.HaematologicalCancer = true }, 6},
		{"HF + haematological (footnote 3: double)", func(in *SAPS3Input) { in.ChronicHeartFailureNYHAIV, in.HaematologicalCancer = true, true }, 12},
		{"cirrhosis", func(in *SAPS3Input) { in.Cirrhosis = true }, 8},
		{"AIDS", func(in *SAPS3Input) { in.AIDS = true }, 8},
		{"cancer", func(in *SAPS3Input) { in.Cancer = true }, 11},
		{"days 13.9", func(in *SAPS3Input) { in.HospitalDaysBeforeICU = 13.9 }, 0},
		{"days 14", func(in *SAPS3Input) { in.HospitalDaysBeforeICU = 14 }, 6},
		{"days 27.9", func(in *SAPS3Input) { in.HospitalDaysBeforeICU = 27.9 }, 6},
		{"days 28", func(in *SAPS3Input) { in.HospitalDaysBeforeICU = 28 }, 7},
		{"emergency room", func(in *SAPS3Input) { in.Location = SAPS3LocationEmergencyRoom }, 5},
		{"other ICU", func(in *SAPS3Input) { in.Location = SAPS3LocationOtherICU }, 7},
		{"other location", func(in *SAPS3Input) { in.Location = SAPS3LocationOther }, 8},
		{"vasoactive", func(in *SAPS3Input) { in.VasoactiveDrugsBeforeAdmission = true }, 3},
		{"unplanned", func(in *SAPS3Input) { in.UnplannedAdmission = true }, 3},
		{"rhythm", func(in *SAPS3Input) { in.RhythmDisturbances = true }, -5},
		{"seizures", func(in *SAPS3Input) { in.Seizures = true }, -4},
		{"rhythm + seizures (footnote 13)", func(in *SAPS3Input) { in.RhythmDisturbances, in.Seizures = true, true }, -4},
		{"hypovolemic shock", func(in *SAPS3Input) { in.HypovolemicShock = true }, 3},
		{"acute abdomen", func(in *SAPS3Input) { in.DigestiveAcuteAbdomenOrOther = true }, 3},
		{"coma", func(in *SAPS3Input) { in.ComaStuporConfusionAgitation = true }, 4},
		{"septic shock", func(in *SAPS3Input) { in.SepticShock = true }, 5},
		{"anaphylactic shock", func(in *SAPS3Input) { in.AnaphylacticMixedUndefinedShock = true }, 5},
		{"liver failure", func(in *SAPS3Input) { in.LiverFailure = true }, 6},
		{"focal deficit", func(in *SAPS3Input) { in.FocalNeurologicDeficit = true }, 7},
		{"pancreatitis", func(in *SAPS3Input) { in.SeverePancreatitis = true }, 9},
		{"mass effect", func(in *SAPS3Input) { in.IntracranialMassEffect = true }, 10},
		{"no surgery", func(in *SAPS3Input) { in.SurgicalStatus = SAPS3NoSurgery }, 5},
		{"emergency surgery", func(in *SAPS3Input) { in.SurgicalStatus = SAPS3EmergencySurgery }, 6},
		{"transplantation", func(in *SAPS3Input) { in.SurgicalSite = SAPS3SiteTransplantation }, -11},
		{"trauma", func(in *SAPS3Input) { in.SurgicalSite = SAPS3SiteTrauma }, -8},
		{"CABG", func(in *SAPS3Input) { in.SurgicalSite = SAPS3SiteCABGWithoutValve }, -6},
		{"neurosurgery CVA", func(in *SAPS3Input) { in.SurgicalSite = SAPS3SiteNeurosurgeryCVA }, 5},
		{"nosocomial", func(in *SAPS3Input) { in.NosocomialInfection = true }, 4},
		{"respiratory infection", func(in *SAPS3Input) { in.RespiratoryInfection = true }, 5},
		{"GCS 3", func(in *SAPS3Input) { in.GCS = 3 }, 15},
		{"GCS 4", func(in *SAPS3Input) { in.GCS = 4 }, 15},
		{"GCS 5", func(in *SAPS3Input) { in.GCS = 5 }, 10},
		{"GCS 6", func(in *SAPS3Input) { in.GCS = 6 }, 7},
		{"GCS 7", func(in *SAPS3Input) { in.GCS = 7 }, 2},
		{"GCS 12", func(in *SAPS3Input) { in.GCS = 12 }, 2},
		{"GCS 13", func(in *SAPS3Input) { in.GCS = 13 }, 0},
		{"temp 34.9", func(in *SAPS3Input) { in.TemperatureC = 34.9 }, 7},
		{"temp 35", func(in *SAPS3Input) { in.TemperatureC = 35 }, 0},
		{"HR 119", func(in *SAPS3Input) { in.HeartRatePerMin = 119 }, 0},
		{"HR 120", func(in *SAPS3Input) { in.HeartRatePerMin = 120 }, 5},
		{"HR 159", func(in *SAPS3Input) { in.HeartRatePerMin = 159 }, 5},
		{"HR 160", func(in *SAPS3Input) { in.HeartRatePerMin = 160 }, 7},
		{"leukocytes 14.9", func(in *SAPS3Input) { in.LeukocytesGL = 14.9 }, 0},
		{"leukocytes 15", func(in *SAPS3Input) { in.LeukocytesGL = 15 }, 2},
		{"pH 7.25", func(in *SAPS3Input) { in.PH = 7.25 }, 3},
		{"pH 7.26", func(in *SAPS3Input) { in.PH = 7.26 }, 0},
		{"platelets 19", func(in *SAPS3Input) { in.PlateletsGL = 19 }, 13},
		{"platelets 20", func(in *SAPS3Input) { in.PlateletsGL = 20 }, 8},
		{"platelets 49", func(in *SAPS3Input) { in.PlateletsGL = 49 }, 8},
		{"platelets 50", func(in *SAPS3Input) { in.PlateletsGL = 50 }, 5},
		{"platelets 99", func(in *SAPS3Input) { in.PlateletsGL = 99 }, 5},
		{"platelets 100", func(in *SAPS3Input) { in.PlateletsGL = 100 }, 0},
		{"SBP 39", func(in *SAPS3Input) { in.SystolicMmHg = 39 }, 11},
		{"SBP 40", func(in *SAPS3Input) { in.SystolicMmHg = 40 }, 8},
		{"SBP 69", func(in *SAPS3Input) { in.SystolicMmHg = 69 }, 8},
		{"SBP 70", func(in *SAPS3Input) { in.SystolicMmHg = 70 }, 3},
		{"SBP 119", func(in *SAPS3Input) { in.SystolicMmHg = 119 }, 3},
		{"SBP 120", func(in *SAPS3Input) { in.SystolicMmHg = 120 }, 0},
		{"PaO2 59 no MV", func(in *SAPS3Input) { in.PaO2MmHg = 59 }, 5},
		{"PaO2 60 no MV", func(in *SAPS3Input) { in.PaO2MmHg = 60 }, 0},
		{"P/F 99 MV", func(in *SAPS3Input) { in.MechanicalVentilation, in.PaO2MmHg, in.FiO2 = true, 49.5, 0.5 }, 11},
		{"P/F 100 MV", func(in *SAPS3Input) { in.MechanicalVentilation, in.PaO2MmHg, in.FiO2 = true, 50, 0.5 }, 7},
	}
	for _, tc := range tests {
		in := saps3Base()
		tc.f(&in)
		if got := saps3Points(t, in, 0.5, 0.8); got != tc.want {
			t.Errorf("%s: %d points, want %d", tc.name, got, tc.want)
		}
	}
	labs := []struct {
		bili, creat float64
		want        int
	}{
		{1.99, 0.8, 0}, {2, 0.8, 4}, {5.99, 0.8, 4}, {6, 0.8, 5},
		{0.5, 1.19, 0}, {0.5, 1.2, 2}, {0.5, 1.99, 2}, {0.5, 2, 7}, {0.5, 3.49, 7}, {0.5, 3.5, 8},
	}
	for _, tc := range labs {
		if got := saps3Points(t, saps3Base(), tc.bili, tc.creat); got != tc.want {
			t.Errorf("bilirubin %v creatinine %v: %d points, want %d", tc.bili, tc.creat, got, tc.want)
		}
	}
}

// The µmol/L rows of Table 1, both sides of each printed edge.
func TestSAPS3ScoreSI_points(t *testing.T) {
	labs := []struct {
		bili, creat float64
		want        int
	}{
		{34.1, 70, 0}, {34.2, 70, 4}, {102.5, 70, 4}, {102.6, 70, 5},
		{10, 106, 0}, {10, 106.1, 2}, {10, 176.7, 2}, {10, 176.8, 7}, {10, 309.3, 7}, {10, 309.4, 8},
	}
	for _, tc := range labs {
		r, err := SAPS3ScoreSI(saps3Base(), tc.bili, tc.creat)
		if err != nil || r.Score-16 != tc.want {
			t.Errorf("bilirubin %v creatinine %v µmol/L: %d (%v), want %d", tc.bili, tc.creat, r.Score-16, err, tc.want)
		}
	}
}

// HAND-DERIVED from Tables 1–2 and "Logit = −32.6659 +ln(SAPS 3 score +20.5958) ×7.3068":
// age 65 (9) + 16 + unplanned (3) + no surgery (5) + septic shock (5) + respiratory infection (5) + HR 125 (5)
// + SBP 85 (3) + PaO2 55 without ventilation (5) = 56; logit = −0.9650395245547…, probability 0.2758703335….
// The paper's minimum, 0: 16 − 5 (rhythm disturbances) − 11 (transplantation). Tolerance 1e-12 absolute:
// double-precision evaluation of the printed equation.
func TestSAPS3Score_source(t *testing.T) {
	in := saps3Base()
	in.AgeYears, in.UnplannedAdmission, in.SurgicalStatus = 65, true, SAPS3NoSurgery
	in.SepticShock, in.RespiratoryInfection = true, true
	in.HeartRatePerMin, in.SystolicMmHg, in.PaO2MmHg = 125, 85, 55
	r, err := SAPS3Score(in, 1, 1)
	if err != nil || r.Score != 56 || math.Abs(r.Logit-(-0.9650395245547436)) > 1e-12 || math.Abs(r.ProbabilityOfDeath-0.2758703335012793) > 1e-12 {
		t.Errorf("got %+v (%v), want 56, −0.96504, 0.27587", r, err)
	}
	rs, err := SAPS3ScoreSI(in, 20, 80)
	if err != nil || rs != r {
		t.Errorf("SI got %+v (%v), want %+v", rs, err, r)
	}
	lo := saps3Base()
	lo.RhythmDisturbances, lo.SurgicalSite = true, SAPS3SiteTransplantation
	r, err = SAPS3Score(lo, 0.5, 0.8)
	if err != nil || r.Score != 0 || math.Abs(r.Logit-(-10.562193055186441)) > 1e-12 {
		t.Errorf("minimum: got %+v (%v), want 0", r, err)
	}
}

// Every additive item at its maximum. Table 1's footnote ("points are assigned for all applicable
// combinations" for comorbidities, reasons for admission and acute infection) gives 251; the paper's text says
// "in theory, vary from a minimum of 0 points to a maximum of 217 points". DECISION NEEDED (ledger saps3_score):
// this test pins the footnote's reading as coded.
func TestSAPS3Score_allItems(t *testing.T) {
	in := SAPS3Input{
		AgeYears: 85, CancerTherapy: true, ChronicHeartFailureNYHAIV: true, HaematologicalCancer: true,
		Cirrhosis: true, AIDS: true, Cancer: true, HospitalDaysBeforeICU: 30, Location: SAPS3LocationOther,
		VasoactiveDrugsBeforeAdmission: true, UnplannedAdmission: true, HypovolemicShock: true,
		DigestiveAcuteAbdomenOrOther: true, ComaStuporConfusionAgitation: true, SepticShock: true,
		AnaphylacticMixedUndefinedShock: true, LiverFailure: true, FocalNeurologicDeficit: true,
		SeverePancreatitis: true, IntracranialMassEffect: true, SurgicalStatus: SAPS3EmergencySurgery,
		SurgicalSite: SAPS3SiteNeurosurgeryCVA, NosocomialInfection: true, RespiratoryInfection: true,
		GCS: 3, TemperatureC: 30, HeartRatePerMin: 170, LeukocytesGL: 30, PH: 7.0, PlateletsGL: 10,
		SystolicMmHg: 30, MechanicalVentilation: true, PaO2MmHg: 40, FiO2: 1,
	}
	if r, err := SAPS3Score(in, 8, 4); err != nil || r.Score != 251 {
		t.Errorf("got %d (%v), want 251", r.Score, err)
	}
}

// parity: ledger saps3_score, source case 1. No _test.go or JSON example. The entry takes pre-banded picker
// indexes whose points differ from Table 1 (e.g. "AIDS (14 pts)", "ER (8 pts)", "GCS 6–8 (13 pts)"); with every
// index 0 its points are 0 + 16, which this function also gives for an admission in every 0-point band.
func TestSAPS3Score_parity(t *testing.T) {
	r, _ := SAPS3Score(saps3Base(), 0.5, 0.8)
	if r.Score != 16 {
		t.Errorf("got %d, want 16", r.Score)
	}
}

func TestSAPS3Score_errors(t *testing.T) {
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := SAPS3Score(saps3Base(), v, 1); !errors.Is(err, ErrNotFinite) {
			t.Errorf("bilirubin %v: err = %v", v, err)
		}
		if _, err := SAPS3Score(saps3Base(), 1, v); !errors.Is(err, ErrNotFinite) {
			t.Errorf("creatinine %v: err = %v", v, err)
		}
		if _, err := SAPS3ScoreSI(saps3Base(), v, 1); !errors.Is(err, ErrNotFinite) {
			t.Errorf("SI bilirubin %v: err = %v", v, err)
		}
		if _, err := SAPS3ScoreSI(saps3Base(), 1, v); !errors.Is(err, ErrNotFinite) {
			t.Errorf("SI creatinine %v: err = %v", v, err)
		}
		for i, f := range []func(*SAPS3Input){
			func(in *SAPS3Input) { in.AgeYears = v }, func(in *SAPS3Input) { in.HospitalDaysBeforeICU = v },
			func(in *SAPS3Input) { in.TemperatureC = v }, func(in *SAPS3Input) { in.HeartRatePerMin = v },
			func(in *SAPS3Input) { in.LeukocytesGL = v }, func(in *SAPS3Input) { in.PH = v },
			func(in *SAPS3Input) { in.PlateletsGL = v }, func(in *SAPS3Input) { in.SystolicMmHg = v },
			func(in *SAPS3Input) { in.PaO2MmHg = v }, func(in *SAPS3Input) { in.FiO2 = v },
		} {
			in := saps3Base()
			f(&in)
			if _, err := SAPS3Score(in, 1, 1); !errors.Is(err, ErrNotFinite) {
				t.Errorf("field %d = %v: err = %v, want ErrNotFinite", i, v, err)
			}
		}
	}
	for i, f := range []func(*SAPS3Input){
		func(in *SAPS3Input) { in.AgeYears = -1 }, func(in *SAPS3Input) { in.AgeYears = 131 },
		func(in *SAPS3Input) { in.HospitalDaysBeforeICU = -1 }, func(in *SAPS3Input) { in.TemperatureC = 9.9 },
		func(in *SAPS3Input) { in.TemperatureC = 50.1 }, func(in *SAPS3Input) { in.HeartRatePerMin = -1 },
		func(in *SAPS3Input) { in.LeukocytesGL = -1 }, func(in *SAPS3Input) { in.PH = 5.9 },
		func(in *SAPS3Input) { in.PH = 8.1 }, func(in *SAPS3Input) { in.PlateletsGL = -1 },
		func(in *SAPS3Input) { in.SystolicMmHg = -1 }, func(in *SAPS3Input) { in.PaO2MmHg = -1 },
		func(in *SAPS3Input) { in.GCS = 2 }, func(in *SAPS3Input) { in.GCS = 16 },
		func(in *SAPS3Input) { in.MechanicalVentilation, in.FiO2 = true, 0.2 },
		func(in *SAPS3Input) { in.MechanicalVentilation, in.FiO2 = true, 1.01 },
		func(in *SAPS3Input) { in.Location = SAPS3Location(4) },
		func(in *SAPS3Input) { in.SurgicalStatus = SAPS3SurgicalStatus(3) },
		func(in *SAPS3Input) { in.SurgicalSite = SAPS3SurgicalSite(-1) },
	} {
		in := saps3Base()
		f(&in)
		if _, err := SAPS3Score(in, 1, 1); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("case %d: err = %v, want ErrOutOfRange", i, err)
		}
	}
	if _, err := SAPS3Score(saps3Base(), -0.1, 1); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("bilirubin −0.1: err = %v", err)
	}
	if _, err := SAPS3ScoreSI(saps3Base(), 1, -0.1); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("SI creatinine −0.1: err = %v", err)
	}
	// FiO2 is not checked without ventilation (Table 1 uses it only with MV).
	in := saps3Base()
	in.FiO2 = 0
	if _, err := SAPS3Score(in, 1, 1); err != nil {
		t.Errorf("FiO2 0 without MV: err = %v, want nil", err)
	}
}

func TestSAPS3Score_noAlloc(t *testing.T) {
	in := saps3Base()
	in.SepticShock, in.MechanicalVentilation, in.PaO2MmHg, in.FiO2 = true, true, 70, 0.6
	if n := testing.AllocsPerRun(100, func() { sinkSAPS3, _ = SAPS3Score(in, 2.5, 1.5) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sinkSAPS3, _ = SAPS3ScoreSI(in, 40, 130) }); n != 0 {
		t.Errorf("SI allocs = %v, want 0", n)
	}
}

func BenchmarkSAPS3Score(b *testing.B) {
	in := saps3Base()
	in.SepticShock, in.MechanicalVentilation, in.PaO2MmHg, in.FiO2 = true, true, 70, 0.6
	b.ReportAllocs()
	for b.Loop() {
		sinkSAPS3, _ = SAPS3Score(in, 2.5, 1.5)
	}
}

func BenchmarkSAPS3ScoreSI(b *testing.B) {
	in := saps3Base()
	in.SepticShock, in.MechanicalVentilation, in.PaO2MmHg, in.FiO2 = true, true, 70, 0.6
	b.ReportAllocs()
	for b.Loop() {
		sinkSAPS3, _ = SAPS3ScoreSI(in, 40, 130)
	}
}
