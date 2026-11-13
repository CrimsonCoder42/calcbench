package endocrine

import (
	"math"
	"testing"
)

// HAND-DERIVED from Zisman 2016 Box 1 (bedtime minus morning); exact integer arithmetic in float64.
// parity: ledger beam_value, JSON inputs bedtimeGlucose/amGlucose (mg/dL), output bedtime − AM.
func TestBeAMValue_example(t *testing.T) {
	tests := []struct{ bed, am, want float64 }{
		{200, 130, 70},
		{130, 200, -70},
		{150, 150, 0},
		{40, 500, -460},
	}
	for _, tc := range tests {
		if got, err := BeAMValue(tc.bed, tc.am); err != nil || got != tc.want {
			t.Errorf("BeAMValue(%v, %v) = (%v, %v), want %v", tc.bed, tc.am, got, err, tc.want)
		}
	}
}

func TestBeAMValue_range(t *testing.T) {
	if got, err := BeAMValue(math.MaxFloat64, 5e-324); err != nil || got != math.MaxFloat64 {
		t.Errorf("extremes: got (%v, %v)", got, err)
	}
}

func TestBeAMValue_errors(t *testing.T) {
	call := func(b, a float64) func() error {
		return func() error { _, err := BeAMValue(b, a); return err }
	}
	runErrCases(t, []errCase{
		{"bedtime NaN", call(nan, 100), ErrNotFinite},
		{"bedtime +Inf", call(pInf, 100), ErrNotFinite},
		{"bedtime -Inf", call(nInf, 100), ErrNotFinite},
		{"morning NaN", call(100, nan), ErrNotFinite},
		{"morning +Inf", call(100, pInf), ErrNotFinite},
		{"morning -Inf", call(100, nInf), ErrNotFinite},
		{"NaN before range", call(-1, nan), ErrNotFinite},
		{"bedtime 0", call(0, 100), ErrOutOfRange},
		{"morning negative", call(100, -1), ErrOutOfRange},
	})
}

func TestBeAMValue_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = BeAMValue(200, 130) }); n != 0 {
		t.Errorf("BeAMValue allocates %v times, want 0", n)
	}
}

func BenchmarkBeAMValue(b *testing.B) {
	b.ReportAllocs()
	bed, am := 200.0, 130.0
	for b.Loop() {
		sink, _ = BeAMValue(bed, am)
	}
}
