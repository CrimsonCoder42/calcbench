package cardio

import (
	"math"
	"testing"
)

// Colivicchi 2003 abstract: four items, age > 65 strictly.
func TestOESILScore_items(t *testing.T) {
	tests := []struct {
		age  float64
		f    OESILFindings
		want int
	}{
		{65, OESILFindings{}, 0},
		{65.01, OESILFindings{}, 1},
		{70, OESILFindings{AbnormalECG: true}, 2},
		{90, OESILFindings{true, true, true}, 4},
	}
	for _, tc := range tests {
		got, err := OESILScore(tc.age, tc.f)
		if err != nil || got != tc.want {
			t.Errorf("OESILScore(%v, %+v) = (%d, %v), want %d", tc.age, tc.f, got, err, tc.want)
		}
	}
	_, err := OESILScore(math.NaN(), OESILFindings{})
	checkErr(t, "NaN", err, ErrNotFinite)
	_, err = OESILScore(math.Inf(1), OESILFindings{})
	checkErr(t, "+Inf", err, ErrNotFinite)
	_, err = OESILScore(17, OESILFindings{})
	checkErr(t, "17", err, ErrOutOfRange)
}

// Colivicchi 2003 abstract: "0% for a score of 0, 0.8% for 1 point; 19.6% for 2 points; 34.7% for 3 points; 57.1% for 4 points".
func TestOESILMortalityPercent_table(t *testing.T) {
	for s, w := range []float64{0, 0.8, 19.6, 34.7, 57.1} {
		if got, err := OESILMortalityPercent(s); err != nil || got != w {
			t.Errorf("score %d: (%v, %v), want %v", s, got, err, w)
		}
	}
	for _, s := range []int{-1, 5} {
		_, err := OESILMortalityPercent(s)
		checkErr(t, "range", err, ErrOutOfRange)
	}
}

func TestOESILScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() {
		sinkInt, _ = OESILScore(70, OESILFindings{})
		sinkFloat, _ = OESILMortalityPercent(2)
	}); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkOESILScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = OESILScore(70, OESILFindings{AbnormalECG: true})
	}
}

func BenchmarkOESILMortalityPercent(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = OESILMortalityPercent(2)
	}
}
