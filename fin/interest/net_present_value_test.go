package interest

import (
	"fmt"
	"math"
	"math/big"
	"testing"
)

// bigNPV returns Σ cₜ/(1+r)ᵗ at 256 bits.
func bigNPV(r float64, flows []int64) *big.Float {
	sum := bigI(0)
	for t, c := range flows {
		sum.Add(sum, new(big.Float).SetPrec(bigPrec).Mul(bigI(c), bigCompound(r, -t)))
	}
	return sum
}

func TestNetPresentValue_SourceExample(t *testing.T) {
	// Principles of Finance (OpenStax 2022) §16.2 Table 16.5: −30,000; 6,000 … 18,000 at 9 % → $26,946.90.
	got, err := NetPresentValue(0.09, []int64{-3000000, 600000, 800000, 1000000, 1200000, 1400000, 1600000, 1800000})
	wantCents(t, "§16.2 Table 16.5", got, err, 2694690)
	// §16.2 Table 16.3 (embroidery machine): the book shows $2,835.63, the sum of its rounded row
	// values; the unrounded sum is $2,835.6194, so the half-even cent is $2,835.62.
	got, err = NetPresentValue(0.09, []int64{-1600000, 200000, 400000, 500000, 500000, 500000, 500000})
	wantCents(t, "§16.2 Table 16.3 (unrounded rows)", got, err, 283562)
}

func TestNetPresentValue_Table(t *testing.T) {
	cases := []struct {
		name  string
		rate  float64
		flows []int64
		want  int64
	}{
		{"one flow at t = 0 is not discounted", 0.5, []int64{12345}, 12345},
		{"zero rate is the plain sum", 0, []int64{-100, 40, 40, 40}, 20},
		{"two flows", 0.1, []int64{-1000, 1100}, 0},
		{"negative rate", -0.5, []int64{0, 100}, 200},
		{"all zero", 0.07, []int64{0, 0, 0}, 0},
		{"max exact flow", 0, []int64{maxExactCents}, maxExactCents},
	}
	for _, c := range cases {
		got, err := NetPresentValue(c.rate, c.flows)
		wantCents(t, c.name, got, err, c.want)
	}
}

func TestNetPresentValue_Errors(t *testing.T) {
	for _, r := range badRates {
		_, err := NetPresentValue(r, []int64{1, 2})
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err := NetPresentValue(0.05, nil)
	wantErr(t, "empty", err, ErrDomain)
	_, err = NetPresentValue(0.05, []int64{1, maxExactCents + 1})
	wantErr(t, "flow 2^53+1", err, ErrOutOfRange)
	_, err = NetPresentValue(0, []int64{maxExactCents, maxExactCents})
	wantErr(t, "sum above 2^53", err, ErrOverflow)
	flows := make([]int64, 400)
	flows[399] = 1
	_, err = NetPresentValue(-1+1e-12, flows)
	wantErr(t, "beyond float64", err, ErrOverflow)
}

func TestNetPresentValue_Oracle(t *testing.T) {
	rng := testRNG()
	for k := 0; k < 200; k++ {
		n := 1 + rng.IntN(60)
		flows := make([]int64, n)
		for t := range flows {
			flows[t] = rng.Int64N(200_000_000) - 100_000_000
		}
		r := rng.Float64()*0.3 - 0.05
		got, err := NetPresentValue(r, flows)
		if err != nil {
			t.Fatal(err)
		}
		// Horner's rounding is about n·ε·Σ|cₜ|vᵗ, which a cancelling sum can make larger than the
		// distance to a half cent; allow that much beyond the half-even cent.
		exact, _ := bigNPV(r, flows).Float64()
		mag, vt := 0.0, 1.0
		for _, c := range flows {
			mag += math.Abs(float64(c)) * vt
			vt /= 1 + r
		}
		if d := math.Abs(float64(got) - exact); d > 0.5+1e-14*float64(n)*mag {
			t.Errorf("rate %v n %d: got %d, exact %.6f", r, n, got, exact)
		}
	}
}

func TestNetPresentValue_NoAllocs(t *testing.T) {
	flows := []int64{-1600000, 200000, 400000, 500000, 500000, 500000, 500000}
	wantNoAllocs(t, "NetPresentValue", func() { sinkCents, _ = NetPresentValue(0.09, flows) })
}

func BenchmarkNetPresentValue(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		flows := make([]int64, n)
		rng := testRNG()
		for t := range flows {
			flows[t] = rng.Int64N(1_000_000)
		}
		flows[0] = -int64(n) * 400_000
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkCents, _ = NetPresentValue(0.0005, flows)
			}
		})
	}
}
