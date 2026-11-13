package anthropometry

import "testing"

var (
	sinkMUST MUSTRisk
	sinkInt  int
)

// HAND-DERIVED from the 'MUST' tool p. 3; both sides of every edge (BMI 18.5 and 20 belong to "18.5 - 20";
// weight loss 5 and 10 belong to "5 -10").
func TestMUSTScore_edges(t *testing.T) {
	tests := []struct {
		name      string
		bmi, loss float64
		acute     bool
		want      int
	}{
		{"worked", 19, 7, false, 2},
		{"BMI 20.01", 20.01, 0, false, 0},
		{"BMI 20", 20, 0, false, 1},
		{"BMI 18.5", 18.5, 0, false, 1},
		{"BMI 18.49", 18.49, 0, false, 2},
		{"obese BMI scores 0", 35, 0, false, 0},
		{"loss 4.99", 25, 4.99, false, 0},
		{"loss 5", 25, 5, false, 1},
		{"loss 10", 25, 10, false, 1},
		{"loss 10.01", 25, 10.01, false, 2},
		{"acute", 25, 0, true, 2},
		{"maximum", 16, 15, true, 6},
		{"loss 100", 25, 100, false, 2},
	}
	for _, tc := range tests {
		if got, err := MUSTScore(tc.bmi, tc.loss, tc.acute); err != nil || got != tc.want {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// parity: ledger must_malnutrition, source cases 1–6. The source takes each step pre-scored (pickers 0–2);
// representative BMI and loss values for each picker index (BMI 22 / 19 / 17; loss 2 / 7 / 12 %).
func TestMUSTScore_parity(t *testing.T) {
	bmis := []float64{22, 19, 17}
	losses := []float64{2, 7, 12}
	tests := []struct {
		bmi, loss int
		acute     bool
		score     int
		risk      MUSTRisk
	}{
		{0, 0, false, 0, MUSTLowRisk},    // source case 1 (source label "Low Risk of Malnutrition")
		{1, 0, false, 1, MUSTMediumRisk}, // source case 2
		{0, 1, false, 1, MUSTMediumRisk}, // source case 3
		{2, 0, false, 2, MUSTHighRisk},   // source case 4
		{0, 0, true, 2, MUSTHighRisk},    // source case 5
		{2, 2, true, 6, MUSTHighRisk},    // source case 6
	}
	for i, tc := range tests {
		got, err := MUSTScore(bmis[tc.bmi], losses[tc.loss], tc.acute)
		if err != nil || got != tc.score {
			t.Errorf("source case %d: score (%v, %v), want %v", i+1, got, err, tc.score)
			continue
		}
		if r, err := MUSTRiskCategory(got); err != nil || r != tc.risk {
			t.Errorf("source case %d: risk (%q, %v), want %q", i+1, r, err, tc.risk)
		}
	}
}

func TestMUSTRiskCategory_edges(t *testing.T) {
	tests := []struct {
		score int
		want  MUSTRisk
	}{{0, MUSTLowRisk}, {1, MUSTMediumRisk}, {2, MUSTHighRisk}, {6, MUSTHighRisk}}
	for _, tc := range tests {
		if got, err := MUSTRiskCategory(tc.score); err != nil || got != tc.want {
			t.Errorf("MUSTRiskCategory(%d) = (%q, %v), want %q", tc.score, got, err, tc.want)
		}
	}
}

func TestMUST_errors(t *testing.T) {
	score := func(b, l float64) func() error {
		return func() error { _, err := MUSTScore(b, l, false); return err }
	}
	cat := func(s int) func() error {
		return func() error { _, err := MUSTRiskCategory(s); return err }
	}
	runErrCases(t, []errCase{
		{"BMI NaN", score(nan, 0), ErrNotFinite},
		{"BMI +Inf", score(pInf, 0), ErrNotFinite},
		{"loss -Inf", score(22, nInf), ErrNotFinite},
		{"BMI 0", score(0, 0), ErrOutOfRange},
		{"loss negative", score(22, -0.1), ErrOutOfRange},
		{"loss above 100", score(22, 100.1), ErrOutOfRange},
		{"risk -1", cat(-1), ErrOutOfRange},
		{"risk 7", cat(7), ErrOutOfRange},
	})
}

func TestMUST_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = MUSTScore(19, 7, false) }); n != 0 {
		t.Errorf("MUSTScore allocates %v times, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sinkMUST, _ = MUSTRiskCategory(2) }); n != 0 {
		t.Errorf("MUSTRiskCategory allocates %v times, want 0", n)
	}
}

func BenchmarkMUSTScore(b *testing.B) {
	b.ReportAllocs()
	bmi, loss := 19.0, 7.0
	for b.Loop() {
		sinkInt, _ = MUSTScore(bmi, loss, false)
	}
}

func BenchmarkMUSTRiskCategory(b *testing.B) {
	b.ReportAllocs()
	s := 2
	for b.Loop() {
		sinkMUST, _ = MUSTRiskCategory(s)
	}
}
