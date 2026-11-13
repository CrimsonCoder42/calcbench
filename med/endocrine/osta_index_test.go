package endocrine

import "testing"

var sinkOSTA OSTARisk

// HAND-DERIVED from "(body weight in kg − age in years) × 0.2", truncated (Subramaniam 2018 PMC6068473).
func TestOSTAIndex_formula(t *testing.T) {
	cases := []struct {
		kg, age float64
		want    int
	}{
		{50, 65, -3},   // −3.0
		{60, 50, 2},    // 2.0
		{55, 61, -1},   // −1.2 truncates toward zero
		{70, 40, 6},    // 6.0
		{45, 80, -7},   // −7.0
		{64.9, 60, 0},  // 0.98
		{54.9, 60, -1}, // −1.02
	}
	for _, c := range cases {
		if got, err := OSTAIndex(c.kg, c.age); err != nil || got != c.want {
			t.Errorf("(%v kg, %v y): got %d, %v; want %d", c.kg, c.age, got, err, c.want)
		}
	}
}

// parity: ledger ost_osteoporosis, source cases 1–7 (ost_osteoporosis_test.go TestOST_MatchesTSX): the consumer
// converts lb with ÷ 2.205; this function uses the exact 0.45359237 kg/lb, and every case agrees.
func TestOSTAIndexImperial_parity(t *testing.T) {
	cases := []struct {
		age, lb float64
		want    int
	}{
		{65, 130, -1}, {50, 150, 3}, {80, 110, -6}, {70, 170, 1}, {45, 220, 10}, {90, 120, -7}, {41, 70, -1},
	}
	for i, c := range cases {
		if got, err := OSTAIndexImperial(c.lb, c.age); err != nil || got != c.want {
			t.Errorf("case %d: got %d, %v; want %d", i+1, got, err, c.want)
		}
	}
}

// Imperial and metric twins agree: 100 lb = 45.359237 kg exactly.
func TestOSTAIndexImperial_agreesWithMetric(t *testing.T) {
	for _, age := range []float64{40, 55, 70, 85} {
		a, _ := OSTAIndexImperial(100, age)
		b, _ := OSTAIndex(45.359237, age)
		if a != b {
			t.Errorf("age %v: imperial %d, metric %d", age, a, b)
		}
	}
}

func TestOSTAIndex_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"kg NaN", func() error { _, err := OSTAIndex(nan, 60); return err }, ErrNotFinite},
		{"age +Inf", func() error { _, err := OSTAIndex(60, pInf); return err }, ErrNotFinite},
		{"lb -Inf", func() error { _, err := OSTAIndexImperial(nInf, 60); return err }, ErrNotFinite},
		{"kg 0", func() error { _, err := OSTAIndex(0, 60); return err }, ErrOutOfRange},
		{"age negative", func() error { _, err := OSTAIndex(60, -1); return err }, ErrOutOfRange},
		{"lb 0", func() error { _, err := OSTAIndexImperial(0, 60); return err }, ErrOutOfRange},
	})
}

// Subramaniam 2018: both sides of each edge.
func TestOSTARiskCategory_edges(t *testing.T) {
	cases := []struct {
		idx  int
		want OSTARisk
	}{{5, OSTALowRisk}, {0, OSTALowRisk}, {-1, OSTAModerateRisk}, {-4, OSTAModerateRisk}, {-5, OSTAHighRisk}, {-12, OSTAHighRisk}}
	for _, c := range cases {
		if got := OSTARiskCategory(c.idx); got != c.want {
			t.Errorf("%d: got %q, want %q", c.idx, got, c.want)
		}
	}
}

func TestOSTA_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = OSTAIndex(55, 60) }); n != 0 {
		t.Errorf("OSTAIndex allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = OSTAIndexImperial(130, 60) }); n != 0 {
		t.Errorf("OSTAIndexImperial allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sinkOSTA = OSTARiskCategory(-2) }); n != 0 {
		t.Errorf("OSTARiskCategory allocs = %v, want 0", n)
	}
}

func BenchmarkOSTAIndex(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkInt, _ = OSTAIndex(55, 60)
	}
}

func BenchmarkOSTAIndexImperial(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkInt, _ = OSTAIndexImperial(130, 60)
	}
}

func BenchmarkOSTARiskCategory(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkOSTA = OSTARiskCategory(i%13 - 8)
	}
}
