package endocrine

import (
	"errors"
	"math"
	"testing"
)

type randomCase struct {
	v        float64
	symptoms bool
	want     bool
}

func checkRandom(t *testing.T, name string, f func(float64, bool) (bool, error), cases []randomCase) {
	t.Helper()
	for _, c := range cases {
		got, err := f(c.v, c.symptoms)
		if err != nil {
			t.Errorf("%s(%v, %v): unexpected error %v", name, c.v, c.symptoms, err)
			continue
		}
		if got != c.want {
			t.Errorf("%s(%v, %v) = %v, want %v", name, c.v, c.symptoms, got, c.want)
		}
	}
}

// ADA 2026 Table 2.1, fourth criterion: random PG ≥ 200 mg/dL with classic symptoms; both sides of the
// edge, with and without symptoms.
func TestRandomGlucoseCriterionADAMgDL(t *testing.T) {
	checkRandom(t, "RandomGlucoseCriterionADAMgDL", RandomGlucoseCriterionADAMgDL, []randomCase{
		{199, true, false}, {199.9, true, false}, {200, true, true}, {600, true, true},
		{199.9, false, false}, {200, false, false}, {600, false, false},
	})
}

// ADA 2026 Table 2.1, mmol/L column: ≥ 11.1 mmol/L.
func TestRandomGlucoseCriterionADAMmolL(t *testing.T) {
	checkRandom(t, "RandomGlucoseCriterionADAMmolL", RandomGlucoseCriterionADAMmolL, []randomCase{
		{11.0, true, false}, {11.09, true, false}, {11.1, true, true}, {33, true, true},
		{11.1, false, false}, {33, false, false},
	})
}

func TestRandomGlucoseCriterionADA_errors(t *testing.T) {
	fs := map[string]func(float64, bool) (bool, error){
		"MgDL": RandomGlucoseCriterionADAMgDL, "MmolL": RandomGlucoseCriterionADAMmolL,
	}
	for name, f := range fs {
		for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
			if _, err := f(v, true); !errors.Is(err, ErrNotFinite) {
				t.Errorf("%s(%v): err = %v, want ErrNotFinite", name, v, err)
			}
		}
		for _, v := range []float64{0, -5} {
			if _, err := f(v, true); !errors.Is(err, ErrOutOfRange) {
				t.Errorf("%s(%v): err = %v, want ErrOutOfRange", name, v, err)
			}
		}
		if got, err := f(math.SmallestNonzeroFloat64, true); err != nil || got {
			t.Errorf("%s(smallest positive) = %v, %v; want false, nil", name, got, err)
		}
	}
}

// parity: ledger bs_random_adult, source case 1 (bands 70–140 "Normal", 141–199 "Elevated", 200–299
// "High", ≥ 300 "Critically High"). The primary defines only the ≥ 200 mg/dL criterion, with classic
// symptoms; no band below 200 and none at 300. Asserted: the criterion.
func TestRandomGlucoseCriterionADAMgDL_parityAdult(t *testing.T) {
	checkRandom(t, "RandomGlucoseCriterionADAMgDL", RandomGlucoseCriterionADAMgDL, []randomCase{
		{140.9, true, false}, // source "Normal"
		{141, true, false},   // source "Elevated": no primary band
		{199.9, true, false}, // source "Elevated"
		{200, true, true},    // source "High"
		{200, false, false},  // source "High" without symptoms: the criterion needs them
		{300, true, true},    // source "Critically High": no primary band
	})
}

// parity: ledger bs_random_child, source case 1 (70–140 "Normal", 141–199 "Elevated", ≥ 200 "High").
// ADA 2026 section 2: the criteria "apply to children, adolescents, and adults".
func TestRandomGlucoseCriterionADAMgDL_parityChild(t *testing.T) {
	checkRandom(t, "RandomGlucoseCriterionADAMgDL", RandomGlucoseCriterionADAMgDL, []randomCase{
		{199.9, true, false}, // source "Elevated"
		{200, true, true},    // source "High"
	})
}

func TestRandomGlucoseCriterionADA_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { _, _ = RandomGlucoseCriterionADAMgDL(210, true) }); a != 0 {
		t.Errorf("MgDL: %v allocs, want 0", a)
	}
	if a := testing.AllocsPerRun(100, func() { _, _ = RandomGlucoseCriterionADAMmolL(11.7, true) }); a != 0 {
		t.Errorf("MmolL: %v allocs, want 0", a)
	}
}

func BenchmarkRandomGlucoseCriterionADAMgDL(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = RandomGlucoseCriterionADAMgDL(210, true)
	}
}

func BenchmarkRandomGlucoseCriterionADAMmolL(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = RandomGlucoseCriterionADAMmolL(11.7, true)
	}
}
