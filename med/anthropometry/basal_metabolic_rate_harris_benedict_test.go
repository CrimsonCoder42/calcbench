package anthropometry

import (
	"errors"
	"testing"
)

// HAND-DERIVED from Harris & Benedict 1918, p. 373. Tolerance 1e-12 relative (four correctly rounded
// operations on values near 10³).
func TestBasalMetabolicRateHarrisBenedict_example(t *testing.T) {
	tests := []struct {
		w, h, a float64
		female  bool
		want    float64
	}{
		{70, 175, 40, false, 1634.4625},
		{60, 165, 30, true, 1393.8155},
	}
	for _, tc := range tests {
		if got, err := BasalMetabolicRateHarrisBenedict(tc.w, tc.h, tc.a, tc.female); err != nil || !relNear(got, tc.want, 1e-12) {
			t.Errorf("(%v, %v, %v, %v) = (%v, %v), want %v", tc.w, tc.h, tc.a, tc.female, got, err, tc.want)
		}
	}
}

// parity: ledger basal_energy_expenditure, source cases 1–6 (lb × 0.45359237 → kg, in × 2.54 → cm). The source
// uses rounded constants (66.5, 13.75, 5.003, 6.775; 655.1, 9.563, 1.850, 4.676) and kg/lb 0.453592; the
// expected values here are the PRIMARY's (p. 373). Source BEE in comments.
func TestBasalMetabolicRateHarrisBenedict_parity(t *testing.T) {
	const kgPerLb, cmPerIn = 0.45359237, 2.54
	tests := []struct {
		name        string
		lb, in, age float64
		female      bool
		want        float64 // 0 means ErrOutOfRange
	}{
		{"source case 1 (source BEE 1875)", 180, 70, 30, false, 1876.18149035256},
		{"source case 2 (age 80: outside 21–70; source 997)", 100, 60, 80, true, 0},
		{"source case 3 (source BEE 1353)", 140, 64, 45, true, 1352.66841397612},
		{"source case 4 (136.1 kg: outside 25.0–124.9; source 2759)", 300, 78, 25, false, 0},
		{"source case 5 (source BEE 1725)", 160, 68, 30, false, 1726.01230964672},
		{"source case 6 (22.7 kg, 91.4 cm, age 100: outside; source 574)", 50, 36, 100, true, 0},
	}
	for _, tc := range tests {
		got, err := BasalMetabolicRateHarrisBenedict(tc.lb*kgPerLb, tc.in*cmPerIn, tc.age, tc.female)
		if tc.want == 0 {
			if !errors.Is(err, ErrOutOfRange) {
				t.Errorf("%s: err = %v, want ErrOutOfRange", tc.name, err)
			}
			continue
		}
		if err != nil || !relNear(got, tc.want, 1e-12) {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

func TestBasalMetabolicRateHarrisBenedict_range(t *testing.T) {
	for _, tc := range []struct{ w, h, a float64 }{{25.0, 151, 21}, {124.9, 200, 70}} {
		if _, err := BasalMetabolicRateHarrisBenedict(tc.w, tc.h, tc.a, false); err != nil {
			t.Errorf("edge (%v, %v, %v) refused: %v", tc.w, tc.h, tc.a, err)
		}
	}
}

func TestBasalMetabolicRateHarrisBenedict_errors(t *testing.T) {
	call := func(w, h, a float64) func() error {
		return func() error { _, err := BasalMetabolicRateHarrisBenedict(w, h, a, true); return err }
	}
	runErrCases(t, []errCase{
		{"weight NaN", call(nan, 170, 40), ErrNotFinite},
		{"height +Inf", call(70, pInf, 40), ErrNotFinite},
		{"age -Inf", call(70, 170, nInf), ErrNotFinite},
		{"NaN before range", call(0, 170, nan), ErrNotFinite},
		{"weight 24.99", call(24.99, 170, 40), ErrOutOfRange},
		{"weight 124.91", call(124.91, 170, 40), ErrOutOfRange},
		{"height 150.9", call(70, 150.9, 40), ErrOutOfRange},
		{"height 200.1", call(70, 200.1, 40), ErrOutOfRange},
		{"age 20.99", call(70, 170, 20.99), ErrOutOfRange},
		{"age 70.01", call(70, 170, 70.01), ErrOutOfRange},
	})
}

func TestBasalMetabolicRateHarrisBenedict_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = BasalMetabolicRateHarrisBenedict(70, 175, 40, false) }); n != 0 {
		t.Errorf("BasalMetabolicRateHarrisBenedict allocates %v times, want 0", n)
	}
}

func BenchmarkBasalMetabolicRateHarrisBenedict(b *testing.B) {
	b.ReportAllocs()
	w, h, a := 70.0, 175.0, 40.0
	for b.Loop() {
		sink, _ = BasalMetabolicRateHarrisBenedict(w, h, a, false)
	}
}
