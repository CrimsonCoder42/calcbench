package cardio

import (
	"math"
	"testing"
)

// HAND-DERIVED from Ranucci 2018's printed rule. Tolerance 1e-12 absolute.
func TestACEFII_values(t *testing.T) {
	tests := []struct {
		name        string
		age, ef, cr float64
		emergency   bool
		hct, want   float64
	}{
		{"example: 70/50 + 0.2·3", 70, 50, 1.0, false, 33, 2.0},
		{"HCT 36 edge: no anaemia points", 70, 50, 1.0, false, 36, 1.4},
		{"creatinine 2.0 edge: no points", 60, 60, 2.0, false, 40, 1},
		{"creatinine 2.01: +2", 60, 60, 2.01, false, 40, 3},
		{"emergency: +3", 60, 60, 1, true, 40, 4},
	}
	for _, tc := range tests {
		got, err := ACEFII(tc.age, tc.ef, tc.cr, tc.emergency, tc.hct)
		if err != nil || math.Abs(got-tc.want) > 1e-12 {
			t.Errorf("%s: (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

func TestACEFII_errors(t *testing.T) {
	_, err := ACEFII(math.NaN(), 50, 1, false, 40)
	checkErr(t, "age NaN", err, ErrNotFinite)
	_, err = ACEFII(70, math.Inf(1), 1, false, 40)
	checkErr(t, "ef +Inf", err, ErrNotFinite)
	_, err = ACEFII(70, 0.5, 1, false, 40)
	checkErr(t, "ef 0.5", err, ErrOutOfRange)
	_, err = ACEFII(70, 50, 0, false, 40)
	checkErr(t, "creatinine 0", err, ErrOutOfRange)
	_, err = ACEFII(70, 50, 1, false, 101)
	checkErr(t, "hct 101", err, ErrOutOfRange)
	_, err = ACEFII(17, 50, 1, false, 40)
	checkErr(t, "age 17", err, ErrOutOfRange)
}

func TestACEFII_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = ACEFII(70, 50, 1, false, 33) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkACEFII(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = ACEFII(70, 50, 1, false, 33)
	}
}
