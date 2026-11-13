package cardio

import (
	"math"
	"testing"
)

// Bösner 2010 score table: age/sex point at female ≥ 65, male ≥ 55; one point per finding. HAND-DERIVED cases.
func TestMarburgHeartScore_items(t *testing.T) {
	all := MarburgHeartFindings{true, true, true, true}
	tests := []struct {
		name   string
		age    float64
		female bool
		f      MarburgHeartFindings
		want   int
	}{
		{"man 54.9", 54.9, false, MarburgHeartFindings{}, 0},
		{"man 55 edge", 55, false, MarburgHeartFindings{}, 1},
		{"woman 64.9", 64.9, true, MarburgHeartFindings{}, 0},
		{"woman 65 edge", 65, true, MarburgHeartFindings{}, 1},
		{"man 60, exercise (example)", 60, false, MarburgHeartFindings{PainWorseDuringExercise: true}, 2},
		{"all, woman 93", 93, true, all, 5},
	}
	for _, tc := range tests {
		got, err := MarburgHeartScore(tc.age, tc.female, tc.f)
		if err != nil || got != tc.want {
			t.Errorf("%s: got (%d, %v), want %d", tc.name, got, err, tc.want)
		}
	}
}

func TestMarburgHeartScore_errors(t *testing.T) {
	for _, a := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		_, err := MarburgHeartScore(a, false, MarburgHeartFindings{})
		checkErr(t, "non-finite", err, ErrNotFinite)
	}
	for _, a := range []float64{34.99, 93.01} {
		_, err := MarburgHeartScore(a, false, MarburgHeartFindings{})
		checkErr(t, "range", err, ErrOutOfRange)
	}
	if _, err := MarburgHeartScore(35, true, MarburgHeartFindings{}); err != nil {
		t.Errorf("age 35: %v", err)
	}
}

func TestMarburgHeartScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = MarburgHeartScore(60, false, MarburgHeartFindings{}) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkMarburgHeartScore(b *testing.B) {
	b.ReportAllocs()
	f := MarburgHeartFindings{PainWorseDuringExercise: true}
	for b.Loop() {
		sinkInt, _ = MarburgHeartScore(60, false, f)
	}
}
