package cardio

import (
	"math"
	"testing"
)

// Yeh 2016: range −2 to 10; age band edges on both sides.
func TestDAPTScore_items(t *testing.T) {
	all := DAPTFactors{true, true, true, true, true, true, true, true}
	tests := []struct {
		age  float64
		f    DAPTFactors
		want int
	}{
		{64.9, DAPTFactors{}, 0},
		{65, DAPTFactors{}, -1},
		{74.9, DAPTFactors{}, -1},
		{75, DAPTFactors{}, -2},
		{70, DAPTFactors{Diabetes: true, MIAtPresentation: true}, 1},
		{50, all, 10},
	}
	for _, tc := range tests {
		if got, err := DAPTScore(tc.age, tc.f); err != nil || got != tc.want {
			t.Errorf("(%v) = (%d, %v), want %d", tc.age, got, err, tc.want)
		}
	}
	_, err := DAPTScore(math.Inf(-1), DAPTFactors{})
	checkErr(t, "-Inf", err, ErrNotFinite)
	_, err = DAPTScore(17, DAPTFactors{})
	checkErr(t, "17", err, ErrOutOfRange)
}

func TestDAPTScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = DAPTScore(70, DAPTFactors{Diabetes: true}) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkDAPTScore(b *testing.B) {
	b.ReportAllocs()
	f := DAPTFactors{Diabetes: true}
	for b.Loop() {
		sinkInt, _ = DAPTScore(70, f)
	}
}
