package economics

import (
	"errors"
	"math"
	"testing"
)

// OpenStax, Principles of Economics 3e, §22.1: "109.8 – 100/100 = 0.098 = 9.8%". Tolerance 1e-15:
// 109.8 has no exact binary form, so 109.8/100 − 1 differs from 0.098 by a few ulps.
func TestInflationRate_openStaxExample(t *testing.T) {
	got, err := InflationRate(100, 109.8)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(got-0.098) > 1e-15 {
		t.Fatalf("got %v, want 0.098", got)
	}
}

func TestInflationRate_cases(t *testing.T) {
	cases := []struct {
		name         string
		prev, cur    float64
		want, tolAbs float64
	}{
		{"noChange", 250, 250, 0, 0},
		{"deflation", 200, 190, -0.05, 1e-16},
		{"toZero", 100, 0, -1, 0},
		{"doubling", 100, 200, 1, 0},
		{"maxFloat", math.MaxFloat64, math.MaxFloat64, 0, 0},
		{"tinyPrevious", 1e-300, 2e-300, 1, 1e-15},
	}
	for _, c := range cases {
		got, err := InflationRate(c.prev, c.cur)
		if err != nil || math.Abs(got-c.want) > c.tolAbs {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestInflationRate_errors(t *testing.T) {
	cases := []struct {
		name      string
		prev, cur float64
		want      error
	}{
		{"prevNaN", math.NaN(), 100, ErrNotFinite},
		{"curNaN", 100, math.NaN(), ErrNotFinite},
		{"prevInf", math.Inf(1), 100, ErrNotFinite},
		{"curNegInf", 100, math.Inf(-1), ErrNotFinite},
		{"prevNegative", -1, 100, ErrNegative},
		{"curNegative", 100, -1, ErrNegative},
		{"prevZero", 0, 100, ErrZeroDenominator},
		{"resultOverflow", 1e-300, 1e308, ErrOverflow},
	}
	for _, c := range cases {
		if _, err := InflationRate(c.prev, c.cur); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestInflationRate_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { _, _ = InflationRate(100, 109.8) }); a != 0 {
		t.Fatalf("allocs = %v, want 0", a)
	}
}

func BenchmarkInflationRate(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = InflationRate(100, 109.8)
	}
}
