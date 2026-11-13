package cardio

import (
	"math"
	"testing"
)

// Sanfilippo 2019 score table; HAND-DERIVED totals including the negative extremes.
func TestIMPEDEVTEScore_items(t *testing.T) {
	all := IMPEDEFactors{true, true, true, true, false, true, true}
	tests := []struct {
		name string
		bmi  float64
		dex  IMPEDEDexamethasone
		pro  IMPEDEThromboprophylaxis
		f    IMPEDEFactors
		want int
	}{
		{"nothing, BMI 24.9", 24.9, IMPEDEDexNone, IMPEDEProphylaxisNone, IMPEDEFactors{}, 0},
		{"BMI 25", 25, IMPEDEDexNone, IMPEDEProphylaxisNone, IMPEDEFactors{}, 1},
		{"example", 27, IMPEDEDexHigh, IMPEDEProphylaxisProphylactic, IMPEDEFactors{ImmunomodulatoryDrug: true}, 6},
		{"low dex", 20, IMPEDEDexLow, IMPEDEProphylaxisNone, IMPEDEFactors{}, 2},
		{"minimum", 20, IMPEDEDexNone, IMPEDEProphylaxisTherapeutic, IMPEDEFactors{AsianPacificIslander: true}, -7},
		{"maximum", 30, IMPEDEDexHigh, IMPEDEProphylaxisNone, all, 24},
	}
	for _, tc := range tests {
		if got, err := IMPEDEVTEScore(tc.bmi, tc.dex, tc.pro, tc.f); err != nil || got != tc.want {
			t.Errorf("%s: (%d, %v), want %d", tc.name, got, err, tc.want)
		}
	}
	_, err := IMPEDEVTEScore(math.NaN(), IMPEDEDexNone, IMPEDEProphylaxisNone, IMPEDEFactors{})
	checkErr(t, "NaN", err, ErrNotFinite)
	_, err = IMPEDEVTEScore(-1, IMPEDEDexNone, IMPEDEProphylaxisNone, IMPEDEFactors{})
	checkErr(t, "bmi < 0", err, ErrOutOfRange)
	_, err = IMPEDEVTEScore(25, 3, IMPEDEProphylaxisNone, IMPEDEFactors{})
	checkErr(t, "dex", err, ErrOutOfRange)
	_, err = IMPEDEVTEScore(25, IMPEDEDexNone, -1, IMPEDEFactors{})
	checkErr(t, "prophylaxis", err, ErrOutOfRange)
}

func TestIMPEDEVTEScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() {
		sinkInt, _ = IMPEDEVTEScore(27, IMPEDEDexHigh, IMPEDEProphylaxisProphylactic, IMPEDEFactors{})
	}); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkIMPEDEVTEScore(b *testing.B) {
	b.ReportAllocs()
	f := IMPEDEFactors{ImmunomodulatoryDrug: true}
	for b.Loop() {
		sinkInt, _ = IMPEDEVTEScore(27, IMPEDEDexHigh, IMPEDEProphylaxisProphylactic, f)
	}
}
