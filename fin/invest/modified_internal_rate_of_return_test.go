package invest

import (
	"math"
	"math/rand/v2"
	"testing"
)

// embroidery is OpenStax Principles of Finance §16.4's embroidery machine: −$16,000 now, then $2,000,
// $4,000 and $5,000 a year for four years.
var embroidery = []int64{-1600000, 200000, 400000, 500000, 500000, 500000, 500000}

// OpenStax §16.4 at a 9 % cost of capital: MIRR = 12 % (whole percent, so ±5e-3). The section's terminal
// value is printed as $31,595.22 but its own column sums to $31,589.23; hand: (3158922.93/1600000)^(1/6) − 1
// = 0.1200476138 (±1e-9).
func TestModifiedInternalRateOfReturn_openStax(t *testing.T) {
	got, err := ModifiedInternalRateOfReturn(embroidery, 0.09, 0.09)
	if err != nil || !near(got, 0.12, 5e-3) || !near(got, 0.1200476138, 1e-9) {
		t.Errorf("MIRR = %v, %v; want 0.1200476138", got, err)
	}
}

func TestModifiedInternalRateOfReturn_cases(t *testing.T) {
	// Two flows: the MIRR is the simple growth rate, whatever the rates.
	got, err := ModifiedInternalRateOfReturn([]int64{-1000, 1100}, 0.3, 0.7)
	if err != nil || !near(got, 0.1, 1e-15) {
		t.Errorf("two flows: %v, %v; want 0.1", got, err)
	}
	// Zero rates: (Σ inflows / Σ outflows)^(1/n) − 1.
	got, err = ModifiedInternalRateOfReturn([]int64{-100, 0, 0, 800}, 0, 0)
	if err != nil || !near(got, 1, 1e-15) {
		t.Errorf("zero rates: %v, %v; want 1", got, err)
	}
	// A loss: inflows below outflows give a negative MIRR.
	got, err = ModifiedInternalRateOfReturn([]int64{-1000, 250}, 0, 0)
	if err != nil || !near(got, -0.75, 1e-15) {
		t.Errorf("loss: %v, %v; want −0.75", got, err)
	}
}

// Oracle: when the finance and reinvestment rates equal the IRR, the MIRR equals the IRR; here a
// project built to have IRR i (−P, then level flows) is checked on 100 random cases.
func TestModifiedInternalRateOfReturn_equalsIRRAtIRR(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for range 100 {
		n := r.IntN(30) + 1
		i := float64(r.IntN(2000)+1) / 10000
		pmt := r.Int64N(10_000_000) + 10_000_000
		pv := float64(pmt) * -math.Expm1(-float64(n)*math.Log1p(i)) / i
		flows := make([]int64, n+1)
		flows[0] = -int64(math.Round(pv))
		for k := 1; k <= n; k++ {
			flows[k] = pmt
		}
		// The rounded price moves the true IRR by at most ~0.5/pv relative; compare at 1e-6.
		got, err := ModifiedInternalRateOfReturn(flows, i, i)
		noErr(t, "MIRR", err)
		if !near(got, i, 1e-6) {
			t.Errorf("n %d i %v pmt %d: MIRR %v", n, i, pmt, got)
		}
	}
}

func TestModifiedInternalRateOfReturn_errors(t *testing.T) {
	_, err := ModifiedInternalRateOfReturn([]int64{-100}, 0.1, 0.1)
	wantErr(t, "one flow", err, ErrDomain)
	_, err = ModifiedInternalRateOfReturn(nil, 0.1, 0.1)
	wantErr(t, "empty", err, ErrDomain)
	_, err = ModifiedInternalRateOfReturn([]int64{100, 100}, 0.1, 0.1)
	wantErr(t, "no outflow", err, ErrNoSolution)
	_, err = ModifiedInternalRateOfReturn([]int64{-100, -100}, 0.1, 0.1)
	wantErr(t, "no inflow", err, ErrNoSolution)
	_, err = ModifiedInternalRateOfReturn([]int64{-100, maxCents + 1}, 0.1, 0.1)
	wantErr(t, "flow > 2^53", err, ErrOutOfRange)
	for _, bad := range []float64{math.NaN(), math.Inf(1), -1} {
		_, err = ModifiedInternalRateOfReturn(embroidery, bad, 0.1)
		wantErr(t, "bad finance rate", err, ErrRate)
		_, err = ModifiedInternalRateOfReturn(embroidery, 0.1, bad)
		wantErr(t, "bad reinvest rate", err, ErrRate)
	}
	big := make([]int64, 2000)
	big[0], big[1] = -1, maxCents
	_, err = ModifiedInternalRateOfReturn(big, 0, 1e6) // (1e6)^1998 overflows
	wantErr(t, "overflow", err, ErrOverflow)
}

func TestModifiedInternalRateOfReturn_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = ModifiedInternalRateOfReturn(embroidery, 0.09, 0.09) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

// benchFlows returns −1e6 then n−1 inflows of 1e5 cents (fixed, no randomness needed).
func benchFlows(n int) []int64 {
	f := make([]int64, n)
	f[0] = -100_000_000
	for i := 1; i < n; i++ {
		f[i] = 100_000
	}
	return f
}

func BenchmarkModifiedInternalRateOfReturn(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		flows := benchFlows(n)
		b.Run(itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, sinkErr = ModifiedInternalRateOfReturn(flows, 0.0001, 0.0001)
			}
		})
	}
}
