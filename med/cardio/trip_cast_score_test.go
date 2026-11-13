package cardio

import (
	"math"
	"testing"
)

// Nemeth 2020, Table 1; age and BMI band edges on both sides; HAND-DERIVED totals.
func TestTRiPCastScore_items(t *testing.T) {
	all := TRiPCastFactors{true, true, true, true, true, true, true, true, true}
	tests := []struct {
		name string
		tr   TRiPCastTrauma
		im   TRiPCastImmobilization
		age  float64
		male bool
		bmi  float64
		f    TRiPCastFactors
		want int
	}{
		{"minimum", TRiPCastTraumaLow, TRiPCastOtherWithPlantarSupport, 34.9, false, 24.9, TRiPCastFactors{}, 1},
		{"age 35", TRiPCastTraumaLow, TRiPCastOtherWithPlantarSupport, 35, false, 24.9, TRiPCastFactors{}, 2},
		{"age 55", TRiPCastTraumaLow, TRiPCastOtherWithPlantarSupport, 55, false, 24.9, TRiPCastFactors{}, 3},
		{"age 74.9", TRiPCastTraumaLow, TRiPCastOtherWithPlantarSupport, 74.9, false, 24.9, TRiPCastFactors{}, 3},
		{"age 75", TRiPCastTraumaLow, TRiPCastOtherWithPlantarSupport, 75, false, 24.9, TRiPCastFactors{}, 4},
		{"BMI 25", TRiPCastTraumaLow, TRiPCastOtherWithPlantarSupport, 20, false, 25, TRiPCastFactors{}, 2},
		{"BMI 34.9", TRiPCastTraumaLow, TRiPCastOtherWithPlantarSupport, 20, false, 34.9, TRiPCastFactors{}, 2},
		{"BMI 35", TRiPCastTraumaLow, TRiPCastOtherWithPlantarSupport, 20, false, 35, TRiPCastFactors{}, 3},
		{"example", TRiPCastTraumaHigh, TRiPCastLowerLegCast, 40, true, 27, TRiPCastFactors{}, 8},
		{"maximum", TRiPCastTraumaHigh, TRiPCastUpperLegCast, 80, true, 40, all, 34},
	}
	for _, tc := range tests {
		if got, err := TRiPCastScore(tc.tr, tc.im, tc.age, tc.male, tc.bmi, tc.f); err != nil || got != tc.want {
			t.Errorf("%s: (%d, %v), want %d", tc.name, got, err, tc.want)
		}
	}
	_, err := TRiPCastScore(0, TRiPCastLowerLegCast, 40, true, 25, TRiPCastFactors{})
	checkErr(t, "trauma 0", err, ErrOutOfRange)
	_, err = TRiPCastScore(TRiPCastTraumaLow, 4, 40, true, 25, TRiPCastFactors{})
	checkErr(t, "immob 4", err, ErrOutOfRange)
	_, err = TRiPCastScore(TRiPCastTraumaLow, TRiPCastLowerLegCast, math.NaN(), true, 25, TRiPCastFactors{})
	checkErr(t, "age NaN", err, ErrNotFinite)
	_, err = TRiPCastScore(TRiPCastTraumaLow, TRiPCastLowerLegCast, 17, true, 25, TRiPCastFactors{})
	checkErr(t, "age 17", err, ErrOutOfRange)
	_, err = TRiPCastScore(TRiPCastTraumaLow, TRiPCastLowerLegCast, 40, true, math.Inf(1), TRiPCastFactors{})
	checkErr(t, "bmi +Inf", err, ErrNotFinite)
}

func TestTRiPCastScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() {
		sinkInt, _ = TRiPCastScore(TRiPCastTraumaHigh, TRiPCastLowerLegCast, 40, true, 27, TRiPCastFactors{})
	}); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkTRiPCastScore(b *testing.B) {
	b.ReportAllocs()
	f := TRiPCastFactors{Surgery3Months: true}
	for b.Loop() {
		sinkInt, _ = TRiPCastScore(TRiPCastTraumaHigh, TRiPCastLowerLegCast, 40, true, 27, f)
	}
}
