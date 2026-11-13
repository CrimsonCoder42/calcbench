package float

import (
	"errors"
	"math"
	"testing"
)

var sinkB bool

func TestIsClose_Cases(t *testing.T) {
	inf, nan := math.Inf(1), math.NaN()
	cases := []struct {
		name           string
		a, b, rel, abs float64
		want           bool
	}{
		{"PEP 485 example: within 1e-9 relative", 1, 1 + 1e-10, 1e-9, 0, true},
		{"PEP 485: near zero needs absTol", 1e-12, 0, 1e-9, 0, false},
		{"near zero with absTol", 1e-12, 0, 1e-9, 1e-10, true},
		{"identical", 3.5, 3.5, 0, 0, true},
		{"zero tolerances, different", 1, math.Nextafter(1, 2), 0, 0, false},
		{"one ulp with eps tolerance", 1, math.Nextafter(1, 2), MachineEpsilon, 0, true},
		{"symmetric D3: scale is max(|a|,|b|)", 100, 90, 0.1, 0, true},
		{"symmetric D3 swapped", 90, 100, 0.1, 0, true},
		{"just outside", 100, 89, 0.1, 0, false},
		{"opposite signs", -1, 1, 0.5, 0, false},
		{"+Inf equals +Inf", inf, inf, 0, 0, true},
		{"+Inf vs -Inf", inf, -inf, 1, 1, false},
		{"Inf vs finite", inf, math.MaxFloat64, 1, 1, false},
		{"NaN vs NaN", nan, nan, 1, 1, false},
		{"NaN vs number", nan, 1, 1, 1, false},
		{"huge opposite (difference overflows)", math.MaxFloat64, -math.MaxFloat64, 0.5, 0, false},
		{"subnormals with absTol", 5e-324, 1e-323, 0, 1e-300, true},
	}
	for _, c := range cases {
		got, err := IsClose(c.a, c.b, c.rel, c.abs)
		if err != nil || got != c.want {
			t.Errorf("%s: IsClose(%v, %v, %v, %v) = %v, %v; want %v", c.name, c.a, c.b, c.rel, c.abs, got, err, c.want)
		}
	}
}

func TestIsClose_Errors(t *testing.T) {
	for _, tol := range [][2]float64{{-1e-9, 0}, {0, -1}, {math.NaN(), 0}, {0, math.Inf(1)}} {
		if _, err := IsClose(1, 1, tol[0], tol[1]); !errors.Is(err, ErrBadTolerance) {
			t.Errorf("IsClose tol %v: err = %v; want ErrBadTolerance", tol, err)
		}
		if _, err := AllClose(nil, nil, tol[0], tol[1]); !errors.Is(err, ErrBadTolerance) {
			t.Errorf("AllClose tol %v: err = %v; want ErrBadTolerance", tol, err)
		}
	}
}

func TestAllClose_Cases(t *testing.T) {
	cases := []struct {
		name string
		a, b []float64
		want bool
	}{
		{"empty", nil, []float64{}, true},
		{"one", []float64{1}, []float64{1 + 1e-12}, true},
		{"all close", []float64{1, 2, 3}, []float64{1, 2 + 1e-12, 3}, true},
		{"last differs", []float64{1, 2, 3}, []float64{1, 2, 3.1}, false},
		{"NaN never close", []float64{math.NaN()}, []float64{math.NaN()}, false},
	}
	for _, c := range cases {
		got, err := AllClose(c.a, c.b, 1e-9, 0)
		if err != nil || got != c.want {
			t.Errorf("%s: AllClose = %v, %v; want %v", c.name, got, err, c.want)
		}
	}
	if _, err := AllClose([]float64{1}, []float64{1, 2}, 1e-9, 0); !errors.Is(err, ErrLengthMismatch) {
		t.Errorf("length mismatch: err = %v; want ErrLengthMismatch", err)
	}
}

func TestAllClose_NoAlloc(t *testing.T) {
	a, b := []float64{1, 2, 3}, []float64{1, 2, 3}
	if n := testing.AllocsPerRun(100, func() { sinkB, _ = AllClose(a, b, 1e-9, 0) }); n != 0 {
		t.Errorf("allocs = %v; want 0", n)
	}
}

func BenchmarkIsClose(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkB, _ = IsClose(1, 1+1e-10, 1e-9, 0)
	}
}

func BenchmarkAllClose(b *testing.B) {
	a := make([]float64, 1000)
	c := make([]float64, 1000)
	for i := range a {
		a[i] = float64(i)
		c[i] = float64(i) * (1 + 1e-12)
	}
	b.ReportAllocs()
	for b.Loop() {
		sinkB, _ = AllClose(a, c, 1e-9, 0)
	}
}
