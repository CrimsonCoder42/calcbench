package cardio

import (
	"math"
	"testing"
)

// Aerts 2017 abstract points (+1 × 5, −1 for reproducible pain). HAND-DERIVED cases, both sides of each age edge.
func TestINTERCHESTScore_items(t *testing.T) {
	tests := []struct {
		name   string
		age    float64
		female bool
		f      INTERCHESTFindings
		want   int
	}{
		{"man 54", 54, false, INTERCHESTFindings{}, 0},
		{"man 55", 55, false, INTERCHESTFindings{}, 1},
		{"woman 64", 64, true, INTERCHESTFindings{}, 0},
		{"woman 70, exertion, reproducible (example)", 70, true, INTERCHESTFindings{PainOnExertion: true, PainReproducible: true}, 1},
		{"minimum −1", 30, false, INTERCHESTFindings{PainReproducible: true}, -1},
		{"maximum 5", 80, false, INTERCHESTFindings{true, true, true, true, false}, 5},
	}
	for _, tc := range tests {
		got, err := INTERCHESTScore(tc.age, tc.female, tc.f)
		if err != nil || got != tc.want {
			t.Errorf("%s: got (%d, %v), want %d", tc.name, got, err, tc.want)
		}
	}
}

func TestINTERCHESTScore_errors(t *testing.T) {
	_, err := INTERCHESTScore(math.NaN(), false, INTERCHESTFindings{})
	checkErr(t, "NaN", err, ErrNotFinite)
	_, err = INTERCHESTScore(math.Inf(-1), false, INTERCHESTFindings{})
	checkErr(t, "-Inf", err, ErrNotFinite)
	_, err = INTERCHESTScore(17, false, INTERCHESTFindings{})
	checkErr(t, "17", err, ErrOutOfRange)
	_, err = INTERCHESTScore(121, false, INTERCHESTFindings{})
	checkErr(t, "121", err, ErrOutOfRange)
}

func TestINTERCHESTScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = INTERCHESTScore(60, false, INTERCHESTFindings{}) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkINTERCHESTScore(b *testing.B) {
	b.ReportAllocs()
	f := INTERCHESTFindings{PainOnExertion: true}
	for b.Loop() {
		sinkInt, _ = INTERCHESTScore(60, false, f)
	}
}
