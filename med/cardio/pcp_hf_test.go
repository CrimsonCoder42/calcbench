package cardio

import (
	"math"
	"testing"
)

// HAND-DERIVED from Khan 2019 Table 2 (Python, float64, the same formula): two input sets per race-sex group.
// Set A: 55 y, SBP 120 untreated, non-smoker, glucose 100 untreated, TC 200, HDL 50, BMI 27, QRS 100.
// Set B: 65 y, SBP 150 treated, smoker, glucose 130 treated, TC 220, HDL 40, BMI 32, QRS 110.
// Tolerance 1e-9 relative: a handful of logs and one pow, no cancellation beyond 1 − S0^x.
// Check on MeanCV: set A's IndX is within 1 of each group's MeanCV (171.5, 99.73, 28.73, 233.9), as expected for
// near-mean inputs.
func TestPCPHF10YearRiskPercent_table2(t *testing.T) {
	a := PCPHFInputs{AgeYears: 55, SystolicMmHg: 120, GlucoseMgDL: 100, TotalCholesterolMgDL: 200, HDLMgDL: 50, BMIKgM2: 27, QRSDurationMs: 100}
	b := PCPHFInputs{AgeYears: 65, SystolicMmHg: 150, TreatedBP: true, CurrentSmoker: true, GlucoseMgDL: 130, TreatedGlucose: true, TotalCholesterolMgDL: 220, HDLMgDL: 40, BMIKgM2: 32, QRSDurationMs: 110}
	tests := []struct {
		female, black bool
		wantA, wantB  float64
	}{
		{false, false, 1.43908984908242, 39.251805572711504},
		{true, false, 0.8012683458698322, 30.259136756552095},
		{false, true, 1.524492702243918, 30.75501937116021},
		{true, true, 1.1207120483402, 34.62659581747902},
	}
	for _, tc := range tests {
		for _, c := range []struct {
			in   PCPHFInputs
			want float64
		}{{a, tc.wantA}, {b, tc.wantB}} {
			in := c.in
			in.Female, in.Black = tc.female, tc.black
			got, err := PCPHF10YearRiskPercent(in)
			if err != nil || math.Abs(got-c.want) > 1e-9*c.want {
				t.Errorf("female=%v black=%v %+v = (%v, %v), want %v", tc.female, tc.black, c.in, got, err, c.want)
			}
		}
	}
}

// parity: ledger pcp_hf — no source test file or JSON example; source case 1 is set A for a White man,
// which the app's coefficient table also gives (its table matches Table 2).
func TestPCPHF10YearRiskPercent_parity(t *testing.T) {
	in := PCPHFInputs{AgeYears: 55, SystolicMmHg: 120, GlucoseMgDL: 100, TotalCholesterolMgDL: 200, HDLMgDL: 50, BMIKgM2: 27, QRSDurationMs: 100}
	got, _ := PCPHF10YearRiskPercent(in) // parity: ledger pcp_hf, source case 1
	if math.Abs(got-1.43908984908242) > 1e-9 {
		t.Errorf("got %v", got)
	}
}

func TestPCPHF10YearRiskPercent_errors(t *testing.T) {
	base := PCPHFInputs{AgeYears: 55, SystolicMmHg: 120, GlucoseMgDL: 100, TotalCholesterolMgDL: 200, HDLMgDL: 50, BMIKgM2: 27, QRSDurationMs: 100}
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		in := base
		in.QRSDurationMs = v
		_, err := PCPHF10YearRiskPercent(in)
		checkErr(t, "non-finite", err, ErrNotFinite)
	}
	for _, age := range []float64{29.9, 80.1} {
		in := base
		in.AgeYears = age
		_, err := PCPHF10YearRiskPercent(in)
		checkErr(t, "age", err, ErrOutOfRange)
	}
	for _, age := range []float64{30, 80} {
		in := base
		in.AgeYears = age
		if _, err := PCPHF10YearRiskPercent(in); err != nil {
			t.Errorf("age %v: %v", age, err)
		}
	}
	in := base
	in.HDLMgDL = 0
	_, err := PCPHF10YearRiskPercent(in)
	checkErr(t, "hdl", err, ErrOutOfRange)
}

func TestPCPHF10YearRiskPercent_noAlloc(t *testing.T) {
	in := PCPHFInputs{AgeYears: 55, SystolicMmHg: 120, GlucoseMgDL: 100, TotalCholesterolMgDL: 200, HDLMgDL: 50, BMIKgM2: 27, QRSDurationMs: 100}
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = PCPHF10YearRiskPercent(in) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkPCPHF10YearRiskPercent(b *testing.B) {
	b.ReportAllocs()
	in := PCPHFInputs{AgeYears: 55, SystolicMmHg: 120, GlucoseMgDL: 100, TotalCholesterolMgDL: 200, HDLMgDL: 50, BMIKgM2: 27, QRSDurationMs: 100}
	for b.Loop() {
		sinkFloat, _ = PCPHF10YearRiskPercent(in)
	}
}
