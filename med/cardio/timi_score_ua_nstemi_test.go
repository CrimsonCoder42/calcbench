package cardio

import (
	"errors"
	"math"
	"testing"
)

// Sinks and a small error checker shared by the MED-CARDIO-004… tests of this package.
var (
	sinkInt int
	sinkStr string
)

// checkErr fails unless errors.Is(err, want).
func checkErr(t *testing.T, name string, err, want error) {
	t.Helper()
	if !errors.Is(err, want) {
		t.Errorf("%s: err = %v, want %v", name, err, want)
	}
}

// Antman 2000 abstract: one point per factor, age ≥ 65 counted from ageYears.
func TestTIMIScoreUANSTEMI_items(t *testing.T) {
	all := TIMIUANSTEMIFactors{true, true, true, true, true, true}
	tests := []struct {
		name string
		age  float64
		f    TIMIUANSTEMIFactors
		want int
	}{
		{"none, age 64.99", 64.99, TIMIUANSTEMIFactors{}, 0},
		{"age edge 65", 65, TIMIUANSTEMIFactors{}, 1},
		{"age 70, ST deviation, markers (example)", 70, TIMIUANSTEMIFactors{STDeviation: true, ElevatedCardiacMarkers: true}, 3},
		{"all at 80", 80, all, 7},
		{"all six, age 50", 50, all, 6},
	}
	for _, tc := range tests {
		got, err := TIMIScoreUANSTEMI(tc.age, tc.f)
		if err != nil || got != tc.want {
			t.Errorf("%s: got (%d, %v), want %d", tc.name, got, err, tc.want)
		}
	}
}

func TestTIMIScoreUANSTEMI_errors(t *testing.T) {
	for _, age := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		_, err := TIMIScoreUANSTEMI(age, TIMIUANSTEMIFactors{})
		checkErr(t, "non-finite", err, ErrNotFinite)
	}
	for _, age := range []float64{17.999, 120.001} {
		_, err := TIMIScoreUANSTEMI(age, TIMIUANSTEMIFactors{})
		checkErr(t, "range", err, ErrOutOfRange)
	}
	if _, err := TIMIScoreUANSTEMI(18, TIMIUANSTEMIFactors{}); err != nil {
		t.Errorf("age 18: %v", err)
	}
	if _, err := TIMIScoreUANSTEMI(120, TIMIUANSTEMIFactors{}); err != nil {
		t.Errorf("age 120: %v", err)
	}
}

// Antman 2000 abstract: "4.7% for a score of 0/1; 8.3% for 2; 13.2% for 3; 19.9% for 4; 26.2% for 5; and 40.9% for 6/7".
func TestTIMIUANSTEMIEventRatePercent_table(t *testing.T) {
	want := []float64{4.7, 4.7, 8.3, 13.2, 19.9, 26.2, 40.9, 40.9}
	for s, w := range want {
		if got, err := TIMIUANSTEMIEventRatePercent(s); err != nil || got != w {
			t.Errorf("score %d: got (%v, %v), want %v", s, got, err, w)
		}
	}
	for _, s := range []int{-1, 8} {
		_, err := TIMIUANSTEMIEventRatePercent(s)
		checkErr(t, "range", err, ErrOutOfRange)
	}
}

func TestTIMIScoreUANSTEMI_noAlloc(t *testing.T) {
	f := TIMIUANSTEMIFactors{STDeviation: true}
	if n := testing.AllocsPerRun(100, func() {
		sinkInt, _ = TIMIScoreUANSTEMI(70, f)
		sinkFloat, _ = TIMIUANSTEMIEventRatePercent(3)
	}); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkTIMIScoreUANSTEMI(b *testing.B) {
	b.ReportAllocs()
	f := TIMIUANSTEMIFactors{STDeviation: true, ElevatedCardiacMarkers: true}
	for b.Loop() {
		sinkInt, _ = TIMIScoreUANSTEMI(70, f)
	}
}

func BenchmarkTIMIUANSTEMIEventRatePercent(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = TIMIUANSTEMIEventRatePercent(3)
	}
}
