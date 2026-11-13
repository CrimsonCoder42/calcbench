package meteorology

import (
	"errors"
	"math"
	"testing"
)

// The NWS wind chill chart (FCM-R19-2003, from the WCT equation) prints 10 °F at 15 mph as −7 °F and −20 °F at
// 40 mph as −57 °F; hand computation with 15^0.16 = 1.542353... gives −6.5895 °F. Tolerance 1e-3 °F against
// the hand value (4-digit hand arithmetic), 0.5 °F against the chart's integers.
func TestWindChillF_cases(t *testing.T) {
	got, err := WindChillF(10, 15)
	if err != nil || math.Abs(got-(-6.5895)) > 1e-3 {
		t.Errorf("10 °F 15 mph: %v %v", got, err)
	}
	cases := []struct{ t, v, chart float64 }{
		{10, 15, -7},
		{-20, 40, -57},
		{40, 5, 36},
		{0, 30, -26},
	}
	for _, c := range cases {
		got, err := WindChillF(c.t, c.v)
		if err != nil || math.Abs(got-c.chart) > 0.5 {
			t.Errorf("%v °F %v mph: %v %v, chart %v", c.t, c.v, got, err, c.chart)
		}
	}
}

func TestWindChillF_errors(t *testing.T) {
	cases := []struct {
		t, v float64
		want error
	}{
		{51, 10, ErrOutOfRange},
		{10, 2.9, ErrOutOfRange},
		{10, 201, ErrOutOfRange},
		{-500, 10, ErrOutOfRange},
		{math.NaN(), 10, ErrNotFinite},
		{10, math.Inf(1), ErrNotFinite},
	}
	for _, c := range cases {
		if _, err := WindChillF(c.t, c.v); !errors.Is(err, c.want) {
			t.Errorf("(%v, %v): %v, want %v", c.t, c.v, err, c.want)
		}
	}
	if _, err := WindChillF(50, 3); err != nil {
		t.Errorf("range edge: %v", err)
	}
}

func BenchmarkWindChillF(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = WindChillF(10, 15)
	}
}
