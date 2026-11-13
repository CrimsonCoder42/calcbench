package interest

import (
	"fmt"
	"math"
	"testing"
)

// npvTolerance bounds |NPV(irr)| for a root within 1e-12·max(1,|r|): the slope Σ t·|cₜ|·v^(t+1)
// times four times the rate tolerance, plus Horner's rounding, about n·ε·Σ|cₜ|·vᵗ (with margin).
func npvTolerance(rate float64, flows []int64) float64 {
	v := 1 / (1 + rate)
	s, st, vt := 0.0, 0.0, 1.0
	for t, c := range flows {
		a := math.Abs(float64(c))
		s += a * vt
		st += float64(t) * a * vt * v
		vt *= v
	}
	return 4e-12*math.Max(1, math.Abs(rate))*st + 1e-14*float64(len(flows))*s + 1e-9
}

func TestInternalRateOfReturn_SourceExample(t *testing.T) {
	// Principles of Finance (OpenStax 2022) §16.3, Table 16.7: the embroidery machine's flows
	// (−16,000; 2,000; 4,000; 5,000 ×4) have IRR 14.09 %.
	got, err := InternalRateOfReturn([]int64{-1600000, 200000, 400000, 500000, 500000, 500000, 500000})
	wantFloat(t, "§16.3", got, err, 0.1409, 5e-5)
}

func TestInternalRateOfReturn_Table(t *testing.T) {
	got, err := InternalRateOfReturn([]int64{-1000, 1100})
	wantFloat(t, "one period 10 %", got, err, 0.1, 1e-11)
	got, err = InternalRateOfReturn([]int64{-1000, 1000})
	wantFloat(t, "zero IRR is exact", got, err, 0, 0)
	got, err = InternalRateOfReturn([]int64{-1000, 500})
	wantFloat(t, "negative IRR", got, err, -0.5, 1e-11)
	got, err = InternalRateOfReturn([]int64{1000, -1100})
	wantFloat(t, "borrowing: inflow first", got, err, 0.1, 1e-11)
	got, err = InternalRateOfReturn([]int64{-100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1000000000})
	wantFloat(t, "very high IRR", got, err, math.Pow(1e7, 0.1)-1, 1e-9)
	got, err = InternalRateOfReturn([]int64{-1000000, 1})
	wantFloat(t, "IRR near −1", got, err, -0.999999, 1e-11)
}

func TestInternalRateOfReturn_MultipleRoots(t *testing.T) {
	// Principles of Finance §16.3, Table 16.8 (Olympic facility): −350, +950, −620 changes sign twice
	// and has two IRRs. The roots of 620v² − 950v + 350 = 0 with v = 1/(1+r) are r ≈ 9.18 % and
	// 62.25 %; the root bracketed nearest 0 is returned.
	got, err := InternalRateOfReturn([]int64{-35000000000, 95000000000, -62000000000})
	disc := math.Sqrt(950*950 - 4*620*350)
	vNear := (950 + disc) / (2 * 620)
	wantFloat(t, "nearest root", got, err, 1/vNear-1, 1e-10)
}

func TestInternalRateOfReturn_Errors(t *testing.T) {
	_, err := InternalRateOfReturn(nil)
	wantErr(t, "empty", err, ErrDomain)
	_, err = InternalRateOfReturn([]int64{-5})
	wantErr(t, "one flow", err, ErrDomain)
	_, err = InternalRateOfReturn([]int64{-5, maxExactCents + 1})
	wantErr(t, "flow 2^53+1", err, ErrOutOfRange)
	_, err = InternalRateOfReturn([]int64{5, 6, 7})
	wantErr(t, "all inflows", err, ErrNoSolution)
	_, err = InternalRateOfReturn([]int64{0, 0})
	wantErr(t, "all zero", err, ErrNoSolution)
	// −1, +3, −3 never crosses zero: NPV = −(1 − 3v + 3v²) < 0 for every v.
	_, err = InternalRateOfReturn([]int64{-100, 300, -300})
	wantErr(t, "no real root", err, ErrNoSolution)
}

func TestInternalRateOfReturn_NoConvergence(t *testing.T) {
	// The public functions' brackets are at most 5e5 wide, so bisection always converges within 200
	// steps; the limit is exercised on the solver itself with a bracket of width 2e300.
	_, err := rtsafe(func(r float64) (float64, float64) { return r - 0.123, 0 }, -1e300, -1e300, 1e300)
	wantErr(t, "solver limit", err, ErrNoConvergence)
}

func TestInternalRateOfReturn_Oracle(t *testing.T) {
	// NPV(IRR) ≈ 0 on conventional projects (one outflow, then inflows: a unique IRR).
	rng := testRNG()
	for k := 0; k < 300; k++ {
		n := 2 + rng.IntN(60)
		flows := make([]int64, n)
		sum := int64(0)
		for t := 1; t < n; t++ {
			flows[t] = rng.Int64N(10_000_000)
			sum += flows[t]
		}
		flows[0] = -1 - int64(float64(sum)*(0.2+rng.Float64()*1.5))
		irr, err := InternalRateOfReturn(flows)
		if err != nil {
			t.Fatalf("%v: %v", flows, err)
		}
		if v := npvFloat(irr, flows); math.Abs(v) > npvTolerance(irr, flows) {
			t.Errorf("NPV(IRR = %v) = %v cents, tolerance %v", irr, v, npvTolerance(irr, flows))
		}
	}
}

func TestInternalRateOfReturn_NoAllocs(t *testing.T) {
	flows := []int64{-1600000, 200000, 400000, 500000, 500000, 500000, 500000}
	wantNoAllocs(t, "InternalRateOfReturn", func() { sinkFloat, _ = InternalRateOfReturn(flows) })
}

func BenchmarkInternalRateOfReturn(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		flows := make([]int64, n)
		rng := testRNG()
		for t := 1; t < n; t++ {
			flows[t] = rng.Int64N(1_000_000)
		}
		flows[0] = -int64(n) * 400_000
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, _ = InternalRateOfReturn(flows)
			}
		})
	}
}
