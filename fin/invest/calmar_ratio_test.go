package invest

import (
	"math"
	"testing"
)

func TestCalmarRatio_cases(t *testing.T) {
	for _, c := range []struct{ ret, mdd, want float64 }{
		{0.15, 0.30, 0.5},
		{-0.10, 0.50, -0.2},
		{0.20, 1, 0.2},
	} {
		got, err := CalmarRatio(c.ret, c.mdd)
		if err != nil || !near(got, c.want, 1e-15) {
			t.Errorf("CalmarRatio(%v, %v) = %v, %v; want %v", c.ret, c.mdd, got, err, c.want)
		}
	}
}

func TestCalmarRatio_errors(t *testing.T) {
	_, err := CalmarRatio(0.1, 0)
	wantErr(t, "zero mdd", err, ErrZeroDenominator)
	_, err = CalmarRatio(0.1, -0.1)
	wantErr(t, "negative mdd", err, ErrDomain)
	_, err = CalmarRatio(0.1, 1.1)
	wantErr(t, "mdd > 1", err, ErrDomain)
	_, err = CalmarRatio(0.1, math.NaN())
	wantErr(t, "NaN mdd", err, ErrNotFinite)
	for _, bad := range []float64{math.NaN(), math.Inf(1), -1} {
		_, err = CalmarRatio(bad, 0.2)
		wantErr(t, "bad return", err, ErrRate)
	}
	_, err = CalmarRatio(1e300, 1e-300)
	wantErr(t, "overflow", err, ErrOverflow)
}

func TestCalmarRatio_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = CalmarRatio(0.15, 0.30) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkCalmarRatio(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkErr = CalmarRatio(0.15, 0.30)
	}
}
