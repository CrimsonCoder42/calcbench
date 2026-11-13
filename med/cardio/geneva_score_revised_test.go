package cardio

import (
	"math"
	"testing"
)

// Le Gal 2006 abstract points; heart-rate and age edges on both sides; maximum 1 + 5 + 3+2+2+3+2+4 = 22.
func TestGenevaScoreRevised_items(t *testing.T) {
	tests := []struct {
		name    string
		age, hr float64
		f       GenevaRevisedFindings
		want    int
	}{
		{"65, HR 74", 65, 74, GenevaRevisedFindings{}, 0},
		{"65.5, HR 75", 65.5, 75, GenevaRevisedFindings{}, 4},
		{"HR 94", 40, 94, GenevaRevisedFindings{}, 3},
		{"HR 95", 40, 95, GenevaRevisedFindings{}, 5},
		{"example", 70, 100, GenevaRevisedFindings{Hemoptysis: true}, 8},
		{"maximum", 80, 120, GenevaRevisedFindings{true, true, true, true, true, true}, 22},
	}
	for _, tc := range tests {
		if got, err := GenevaScoreRevised(tc.age, tc.hr, tc.f); err != nil || got != tc.want {
			t.Errorf("%s: (%d, %v), want %d", tc.name, got, err, tc.want)
		}
	}
	_, err := GenevaScoreRevised(math.NaN(), 80, GenevaRevisedFindings{})
	checkErr(t, "age NaN", err, ErrNotFinite)
	_, err = GenevaScoreRevised(50, math.Inf(1), GenevaRevisedFindings{})
	checkErr(t, "hr +Inf", err, ErrNotFinite)
	_, err = GenevaScoreRevised(50, 0, GenevaRevisedFindings{})
	checkErr(t, "hr 0", err, ErrOutOfRange)
	_, err = GenevaScoreRevised(17, 80, GenevaRevisedFindings{})
	checkErr(t, "age 17", err, ErrOutOfRange)
}

func TestGenevaScoreRevised_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = GenevaScoreRevised(70, 100, GenevaRevisedFindings{}) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkGenevaScoreRevised(b *testing.B) {
	b.ReportAllocs()
	f := GenevaRevisedFindings{Hemoptysis: true}
	for b.Loop() {
		sinkInt, _ = GenevaScoreRevised(70, 100, f)
	}
}
