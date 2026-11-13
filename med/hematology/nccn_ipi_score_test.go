package hematology

import (
	"errors"
	"math"
	"testing"
)

var (
	sinkInt    int
	sinkFloat  float64
	sinkString string
)

// Zhou 2014 Table 3: each age and LDH edge on both sides; the binary items one point each; maximum 8.
func TestNCCNIPIScore_example(t *testing.T) {
	tests := []struct {
		age, ldh     float64
		stage, extra bool
		ecog, want   int
	}{
		{40, 1, false, false, 1, 0}, {40.1, 1, false, false, 0, 1}, {60, 1, false, false, 0, 1},
		{60.1, 1, false, false, 0, 2}, {75, 1, false, false, 0, 2}, {75.1, 1, false, false, 0, 3},
		{30, 1.01, false, false, 0, 1}, {30, 3, false, false, 0, 1}, {30, 3.01, false, false, 0, 2},
		{30, 0.5, true, false, 0, 1}, {30, 0.5, false, true, 0, 1}, {30, 0.5, false, false, 2, 1},
		{65, 2, true, true, 1, 5}, {80, 4, true, true, 4, 8},
	}
	for _, tc := range tests {
		if got, err := NCCNIPIScore(tc.age, tc.ldh, tc.stage, tc.extra, tc.ecog); err != nil || got != tc.want {
			t.Errorf("NCCNIPIScore(%v, %v, %v, %v, %d) = (%d, %v), want %d", tc.age, tc.ldh, tc.stage, tc.extra, tc.ecog, got, err, tc.want)
		}
	}
}

func TestNCCNIPIScore_errors(t *testing.T) {
	cases := []struct {
		name     string
		age, ldh float64
		ecog     int
		want     error
	}{
		{"age NaN", math.NaN(), 1, 0, ErrNotFinite}, {"ldh +Inf", 50, math.Inf(1), 0, ErrNotFinite},
		{"ldh -Inf", 50, math.Inf(-1), 0, ErrNotFinite}, {"age -1", -1, 1, 0, ErrOutOfRange},
		{"age 131", 131, 1, 0, ErrOutOfRange}, {"ldh -0.1", 50, -0.1, 0, ErrOutOfRange},
		{"ldh 101", 50, 101, 0, ErrOutOfRange}, {"ecog -1", 50, 1, -1, ErrOutOfRange}, {"ecog 5", 50, 1, 5, ErrOutOfRange},
	}
	for _, c := range cases {
		if _, err := NCCNIPIScore(c.age, c.ldh, false, false, c.ecog); !errors.Is(err, c.want) {
			t.Errorf("%s: error = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestNCCNIPIScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = NCCNIPIScore(65, 2, true, true, 1) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkNCCNIPIScore(b *testing.B) {
	b.ReportAllocs()
	age, ldh, ecog := 65.0, 2.0, 1
	for b.Loop() {
		sinkInt, _ = NCCNIPIScore(age, ldh, true, true, ecog)
	}
}
