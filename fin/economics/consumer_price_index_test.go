package economics

import (
	"errors"
	"math"
	"testing"
)

// OpenStax, Principles of Economics 3e, §22.1 "Tracking Inflation": base year Period 3, basket $107.
// The book prints one decimal (109.8, 93.4, 99.5), truncated, not rounded (93.4579… is printed 93.4);
// the exact ratios are compared to 1e-12 relative (one division and one multiplication, each within
// half an ulp), and the truncation to one decimal must equal the printed value.
func TestConsumerPriceIndex_openStaxExample(t *testing.T) {
	cases := []struct {
		name          string
		current, base int64
		want          float64
	}{
		{"period4", 11750, 10700, 11750.0 / 10700.0 * 100},
		{"period1", 10000, 10700, 10000.0 / 10700.0 * 100},
		{"period2", 10650, 10700, 10650.0 / 10700.0 * 100},
	}
	printed := map[string]float64{"period4": 109.8, "period1": 93.4, "period2": 99.5}
	for _, c := range cases {
		got, err := ConsumerPriceIndex(c.current, c.base)
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if math.Abs(got-c.want) > 1e-12*c.want {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
		}
		if math.Trunc(got*10)/10 != printed[c.name] {
			t.Errorf("%s: %v does not truncate to the book's %v", c.name, got, printed[c.name])
		}
	}
}

func TestConsumerPriceIndex_cases(t *testing.T) {
	cases := []struct {
		name          string
		current, base int64
		want          float64
	}{
		{"baseIsHundred", 10700, 10700, 100},
		{"zeroCurrent", 0, 10700, 0},
		{"double", 200, 100, 200},
		{"maxExact", maxExactCents, maxExactCents, 100},
		{"oneCent", 1, 1, 100},
	}
	for _, c := range cases {
		got, err := ConsumerPriceIndex(c.current, c.base)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestConsumerPriceIndex_errors(t *testing.T) {
	cases := []struct {
		name          string
		current, base int64
		want          error
	}{
		{"currentAboveMax", maxExactCents + 1, 100, ErrOutOfRange},
		{"baseAboveMax", 100, maxExactCents + 1, ErrOutOfRange},
		{"baseBelowMin", 100, math.MinInt64, ErrOutOfRange},
		{"negativeCurrent", -1, 100, ErrNegative},
		{"negativeBase", 100, -1, ErrNegative},
		{"zeroBase", 100, 0, ErrZeroDenominator},
	}
	for _, c := range cases {
		if _, err := ConsumerPriceIndex(c.current, c.base); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestConsumerPriceIndex_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { _, _ = ConsumerPriceIndex(11750, 10700) }); a != 0 {
		t.Fatalf("allocs = %v, want 0", a)
	}
}

var sinkFloat float64

func BenchmarkConsumerPriceIndex(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = ConsumerPriceIndex(11750, 10700)
	}
}
