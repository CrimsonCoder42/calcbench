package invest

import (
	"math"
	"math/rand/v2"
	"testing"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

// OpenStax Principles of Finance §16.4: embroidery machine at 9 %, PI = $18,836/$16,000 = 1.18
// (two decimals: ±5e-3); hand: 1883561.938/1600000 = 1.1772262112 (±1e-9).
func TestProfitabilityIndex_openStax(t *testing.T) {
	got, err := ProfitabilityIndex(0.09, embroidery)
	if err != nil || !near(got, 1.18, 5e-3) || !near(got, 1.1772262112, 1e-9) {
		t.Errorf("PI = %v, %v; want 1.1772262112", got, err)
	}
}

func TestProfitabilityIndex_cases(t *testing.T) {
	got, err := ProfitabilityIndex(0, []int64{-100, 50, 50})
	if err != nil || got != 1 {
		t.Errorf("break-even: %v, %v; want 1", got, err)
	}
	got, err = ProfitabilityIndex(0.1, []int64{-100})
	if err != nil || got != 0 {
		t.Errorf("no inflow: %v, %v; want 0", got, err)
	}
	// Outflows later in the project count at their present value: −100 now, +330 at t=1, −121 at t=2.
	got, err = ProfitabilityIndex(0.1, []int64{-100, 330, -121})
	if err != nil || !near(got, 300.0/200.0, 1e-14) {
		t.Errorf("mixed: %v, %v; want 1.5", got, err)
	}
}

// Oracle: for a single initial outflow, PI = (NPV + I)/I with NPV from interest.NetPresentValue
// (rounded to the cent, so ±1e-6 relative for I ≥ 1e6 cents).
func TestProfitabilityIndex_npvIdentity(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for range 200 {
		n := r.IntN(40) + 2
		flows := make([]int64, n)
		flows[0] = -(r.Int64N(1e9) + 1e6)
		for k := 1; k < n; k++ {
			flows[k] = r.Int64N(1e8)
		}
		rate := float64(r.IntN(3000)) / 10000
		got, err := ProfitabilityIndex(rate, flows)
		noErr(t, "PI", err)
		npv, err := interest.NetPresentValue(rate, flows)
		noErr(t, "NPV", err)
		want := float64(npv-flows[0]) / float64(-flows[0])
		if math.Abs(got-want) > 1e-6 {
			t.Errorf("rate %v: PI %v, (NPV+I)/I %v", rate, got, want)
		}
	}
}

func TestProfitabilityIndex_errors(t *testing.T) {
	_, err := ProfitabilityIndex(0.1, nil)
	wantErr(t, "empty", err, ErrDomain)
	_, err = ProfitabilityIndex(0.1, []int64{100, 200})
	wantErr(t, "no outflow", err, ErrZeroDenominator)
	_, err = ProfitabilityIndex(0.1, []int64{-maxCents - 1, 200})
	wantErr(t, "flow < −2^53", err, ErrOutOfRange)
	for _, bad := range []float64{math.NaN(), math.Inf(-1), -1} {
		_, err = ProfitabilityIndex(bad, embroidery)
		wantErr(t, "bad rate", err, ErrRate)
	}
	big := make([]int64, 2000)
	big[0], big[1999] = -1, maxCents
	_, err = ProfitabilityIndex(-0.999, big) // v = 1000, v^1999 overflows
	wantErr(t, "overflow", err, ErrOverflow)
}

func TestProfitabilityIndex_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = ProfitabilityIndex(0.09, embroidery) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkProfitabilityIndex(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		flows := benchFlows(n)
		b.Run(itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, sinkErr = ProfitabilityIndex(0.0001, flows)
			}
		})
	}
}
