package computing

import (
	"errors"
	"math"
	"testing"
)

// Hand computation from L/R (Kurose & Ross, 7th ed., §1.4.1): 500 MB = 500 × 10⁶ × 8 = 4 × 10⁹ bit;
// 4e9 / 5e6 bit/s = 800 s. Exact in float64.
func TestTransmissionTime_cases(t *testing.T) {
	if got, err := TransmissionTime(4e9, 5e6); err != nil || got != 800 {
		t.Errorf("got %v %v", got, err)
	}
	if got, err := TransmissionTime(0, 5e6); err != nil || got != 0 {
		t.Errorf("zero size: %v %v", got, err)
	}
	for _, c := range []struct {
		l, r float64
		want error
	}{
		{1, 0, ErrZero},
		{-1, 1, ErrNegative},
		{1, -1, ErrNegative},
		{math.NaN(), 1, ErrNotFinite},
		{1, math.Inf(1), ErrNotFinite},
		{math.MaxFloat64, 0.5, ErrOverflow},
	} {
		if _, err := TransmissionTime(c.l, c.r); !errors.Is(err, c.want) {
			t.Errorf("(%v, %v): %v, want %v", c.l, c.r, err, c.want)
		}
	}
}

func BenchmarkTransmissionTime(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = TransmissionTime(4e9, 5e6)
	}
}
