package market

import (
	"math"
	"testing"
)

// Hand-checked: (100·1 + 102·3)/4 = 101.5 ticks × 0.5 = 50.75 (exact in binary).
func TestVWAP_cases(t *testing.T) {
	for _, c := range []struct {
		name   string
		prices []int64
		sizes  []float64
		tick   float64
		want   float64
	}{
		{"example", []int64{100, 102}, []float64{1, 3}, 0.5, 50.75},
		{"one trade", []int64{2421}, []float64{7}, 1, 2421},
		{"zero-size trade ignored", []int64{10, 99}, []float64{2, 0}, 1, 10},
	} {
		got, err := VWAP(c.prices, c.sizes, c.tick)
		if err != nil || got != c.want {
			t.Errorf("%s: %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestVWAP_errors(t *testing.T) {
	_, err := VWAP(nil, nil, 1)
	wantErr(t, "empty", err, ErrEmpty)
	_, err = VWAP([]int64{1}, []float64{1, 2}, 1)
	wantErr(t, "lengths", err, ErrLengthMismatch)
	_, err = VWAP([]int64{1}, []float64{1}, 0)
	wantErr(t, "tick 0", err, ErrOutOfRange)
	_, err = VWAP([]int64{-1}, []float64{1}, 1)
	wantErr(t, "negative price", err, ErrNegative)
	_, err = VWAP([]int64{1}, []float64{math.NaN()}, 1)
	wantErr(t, "NaN size", err, ErrNotFinite)
	_, err = VWAP([]int64{1, 2}, []float64{0, 0}, 1)
	wantErr(t, "all zero", err, ErrDomain)
	_, err = VWAP([]int64{4, 4}, []float64{math.MaxFloat64, math.MaxFloat64}, 1)
	wantErr(t, "overflow", err, ErrOverflow)
}

func BenchmarkVWAP(b *testing.B) {
	p, s := make([]int64, 1000), make([]float64, 1000)
	for i := range p {
		p[i], s[i] = int64(100+i%5), float64(i%9+1)
	}
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkErr = VWAP(p, s, 0.01)
	}
}
