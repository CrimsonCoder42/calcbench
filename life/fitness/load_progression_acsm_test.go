package fitness

import (
	"errors"
	"math"
	"testing"
)

// Shared test values and helpers for the package (first file, as the sentinels are).
var (
	nan    = math.NaN()
	pInf   = math.Inf(1)
	nInf   = math.Inf(-1)
	sink   float64
	sinkI  int
	sinkLR LoadIncreaseRange
)

// relNear reports whether got is within rel (relative) of want; exact when want is 0.
func relNear(got, want, rel float64) bool {
	if want == 0 {
		return got == 0
	}
	return math.Abs(got-want) <= rel*math.Abs(want)
}

type errCase struct {
	name string
	call func() error
	want error
}

func runErrCases(t *testing.T, cases []errCase) {
	t.Helper()
	for _, c := range cases {
		if err := c.call(); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// Worked values HAND-DERIVED from the ACSM 2009 abstract's "2-10% increase in load ... one to two
// repetitions over the desired number": low = load × 1.02, high = load × 1.10, condition met at ≥ 1 rep over.
// Tolerance 1e-15 relative: one correctly rounded multiplication.
func TestLoadProgressionACSM_example(t *testing.T) {
	r, err := LoadProgressionACSM(225, 1)
	if err != nil || !relNear(r.LowLoad, 229.5, 1e-15) || !relNear(r.HighLoad, 247.5, 1e-15) || !r.ConditionMet {
		t.Errorf("LoadProgressionACSM(225, 1) = (%+v, %v), want {229.5 247.5 true}", r, err)
	}
}

// parity: ledger acsm_load_progression, source cases 1–5 (expected values printed to 2 decimals in the source;
// compared here to 1e-12 relative). "keep" = condition not met, "increase" = met.
func TestLoadProgressionACSM_parity(t *testing.T) {
	cases := []struct {
		load      float64
		over      int
		low, high float64
		met       bool
	}{
		{100, 0, 102.00, 110.00, false}, // parity: ledger acsm_load_progression, source case 1
		{225, 1, 229.50, 247.50, true},  // parity: ledger acsm_load_progression, source case 2
		{45, 2, 45.90, 49.50, true},     // parity: ledger acsm_load_progression, source case 3
		{135, 3, 137.70, 148.50, true},  // parity: ledger acsm_load_progression, source case 4
		{315, 0, 321.30, 346.50, false}, // parity: ledger acsm_load_progression, source case 5
	}
	for _, c := range cases {
		r, err := LoadProgressionACSM(c.load, c.over)
		if err != nil || !relNear(r.LowLoad, c.low, 1e-12) || !relNear(r.HighLoad, c.high, 1e-12) || r.ConditionMet != c.met {
			t.Errorf("LoadProgressionACSM(%v, %d) = (%+v, %v), want {%v %v %v}", c.load, c.over, r, err, c.low, c.high, c.met)
		}
	}
}

// Both sides of the condition edge (0 / 1 rep over) and of the load edge (0 / smallest positive).
func TestLoadProgressionACSM_edges(t *testing.T) {
	if r, _ := LoadProgressionACSM(100, 0); r.ConditionMet {
		t.Error("0 reps over: condition met, want not met")
	}
	if r, _ := LoadProgressionACSM(100, 1); !r.ConditionMet {
		t.Error("1 rep over: condition not met, want met")
	}
	if r, err := LoadProgressionACSM(5e-324, 0); err != nil || r.HighLoad < 0 {
		t.Errorf("smallest load: (%+v, %v), want no error", r, err)
	}
}

// parity: ledger acsm_load_progression, source refusal cases: negative reps over is refused (case "-1");
// fractional reps (1.5) cannot be passed (repsOverTarget is int); the kg-unit and age-15 refusals have no
// counterpart (the function takes kg or lb and has no age input).
func TestLoadProgressionACSM_errors(t *testing.T) {
	call := func(l float64, o int) func() error {
		return func() error { _, err := LoadProgressionACSM(l, o); return err }
	}
	runErrCases(t, []errCase{
		{"load NaN", call(nan, 1), ErrNotFinite},
		{"load +Inf", call(pInf, 1), ErrNotFinite},
		{"load -Inf", call(nInf, 1), ErrNotFinite},
		{"load 0", call(0, 1), ErrOutOfRange},
		{"load negative", call(-1, 1), ErrOutOfRange},
		{"reps over -1 (parity)", call(100, -1), ErrOutOfRange},
		{"high end overflows", call(math.MaxFloat64, 1), ErrDomain},
	})
}

func TestLoadProgressionACSM_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkLR, _ = LoadProgressionACSM(225, 1) }); n != 0 {
		t.Errorf("LoadProgressionACSM allocates %v times, want 0", n)
	}
}

func BenchmarkLoadProgressionACSM(b *testing.B) {
	b.ReportAllocs()
	l, o := 225.0, 1
	for b.Loop() {
		sinkLR, _ = LoadProgressionACSM(l, o)
	}
}
