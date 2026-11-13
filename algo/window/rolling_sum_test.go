package window

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"slices"
	"testing"
)

// nistSupplier is the example of NIST/SEMATECH e-Handbook of Statistical Methods, §6.4.2 and §6.4.2.1:
// what each of 12 suppliers delivers, in 1000-dollar units. Shared by every test file in this package.
var nistSupplier = []float64{9, 8, 9, 12, 9, 12, 11, 7, 13, 9, 11, 10}

// benchSizes are the collection sizes of .claude/rules/go-complexity.md. Shared by the benchmarks.
var benchSizes = []struct {
	name string
	n    int
}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}

// benchWindow is the window length used by every benchmark: within stackRingMax, so the hot path
// with a reused dst allocates nothing.
const benchWindow = 50

// benchFloats returns n fixed-seed values in [-1000, 1000).
func benchFloats(n int) []float64 {
	rng := rand.New(rand.NewPCG(1, 2))
	s := make([]float64, n)
	for i := range s {
		s[i] = 2000*rng.Float64() - 1000
	}
	return s
}

// exactWindowSum is the oracle: the exact rational sum of the window, rounded once to float64.
func exactWindowSum(window []float64) float64 {
	var total big.Rat
	for _, v := range window {
		var r big.Rat
		total.Add(&total, r.SetFloat64(v))
	}
	f, _ := total.Float64()
	return f
}

// NIST/SEMATECH e-Handbook §6.4.2.1, moving average with N = 3 over the supplier data. The handbook
// prints the averages; the sums (average × 3) are hand-derived and exact in float64.
func TestRollingSum_nistExample(t *testing.T) {
	want := []float64{26, 29, 30, 33, 32, 30, 31, 29, 33, 30}
	got, err := RollingSum(nil, nistSupplier, 3)
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("RollingSum(nistSupplier, 3) = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestRollingSum_cases(t *testing.T) {
	tests := []struct {
		name      string
		values    []float64
		windowLen int
		want      []float64
	}{
		{"empty", nil, 1, nil},
		{"window longer than the data", []float64{1, 2}, 3, nil},
		{"one element, window 1", []float64{4}, 1, []float64{4}},
		{"window 1 is the identity", []float64{3, -1, 2}, 1, []float64{3, -1, 2}},
		{"window equals the length", []float64{3, -1, 2}, 3, []float64{4}},
		{"two windows", []float64{1, 2, 3}, 2, []float64{3, 5}},
		{"negative and zero", []float64{-2, 0, -3, 5}, 2, []float64{-2, -3, 2}},
		{"signed zeros", []float64{math.Copysign(0, -1), 0}, 1, []float64{0, 0}},
		// After 1e16 leaves, an uncompensated running sum would have lost the 1s it absorbed.
		{"large value leaves the window", []float64{1e16, 1, 1, 1, 1}, 2, []float64{1e16, 2, 2, 2}},
		{"largest finite values, window 1", []float64{math.MaxFloat64, math.MaxFloat64}, 1, []float64{math.MaxFloat64, math.MaxFloat64}},
		{"smallest subnormal", []float64{5e-324, 5e-324}, 2, []float64{1e-323}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RollingSum(nil, tt.values, tt.windowLen)
			if err != nil || len(got) != len(tt.want) || !slices.Equal(got, tt.want) {
				t.Errorf("RollingSum(%v, %d) = (%v, %v), want (%v, nil)", tt.values, tt.windowLen, got, err, tt.want)
			}
		})
	}
}

func TestRollingSum_errors(t *testing.T) {
	tests := []struct {
		name      string
		values    []float64
		windowLen int
		wantErr   error
	}{
		{"window zero", []float64{1, 2}, 0, ErrWindowLen},
		{"window negative", []float64{1, 2}, -1, ErrWindowLen},
		{"window at the minimum of int", []float64{1, 2}, math.MinInt, ErrWindowLen},
		{"NaN", []float64{1, math.NaN(), 2}, 2, ErrNotFinite},
		{"+Inf", []float64{1, math.Inf(1)}, 1, ErrNotFinite},
		{"-Inf outside every window still rejected", []float64{1, 2, math.Inf(-1)}, 5, ErrNotFinite},
		{"sum overflows", []float64{math.MaxFloat64, math.MaxFloat64}, 2, ErrOverflow},
		{"sum overflows to -Inf", []float64{-math.MaxFloat64, -math.MaxFloat64}, 2, ErrOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := []float64{7}
			got, err := RollingSum(dst, tt.values, tt.windowLen)
			if !errors.Is(err, tt.wantErr) || !slices.Equal(got, []float64{7}) {
				t.Errorf("RollingSum(%v, %d) = (%v, %v), want ([7], %v)", tt.values, tt.windowLen, got, err, tt.wantErr)
			}
		})
	}
}

func TestRollingSum_appendsToDst(t *testing.T) {
	got, err := RollingSum([]float64{-1}, []float64{1, 2, 3}, 2)
	if err != nil || !slices.Equal(got, []float64{-1, 3, 5}) {
		t.Errorf("RollingSum([-1], [1 2 3], 2) = (%v, %v), want ([-1 3 5], nil)", got, err)
	}
}

func TestRollingSum_noAllocationWithBuffer(t *testing.T) {
	values := benchFloats(1000)
	buf := make([]float64, 0, len(values))
	allocs := testing.AllocsPerRun(20, func() { buf, _ = RollingSum(buf[:0], values, benchWindow) })
	if allocs != 0 {
		t.Errorf("RollingSum with a reused buffer allocates %v times per call, want 0", allocs)
	}
}

func TestRollingSum_doesNotModifyInput(t *testing.T) {
	values := slices.Clone(nistSupplier)
	if _, err := RollingSum(nil, values, 3); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(values, nistSupplier) {
		t.Errorf("RollingSum modified its input: %v", values)
	}
}

// eps is 2⁻⁵², above the unit roundoff 2⁻⁵³, so tolerances built from it are safe.
const eps = 0x1p-52

// sumTolerance bounds the error of the sum of the window starting at index start. Since the last
// restart, at most 2w-1 values have entered the running sum, all inside values[start-w+1 : start+w],
// in at most m = 3w Neumaier steps. sum + comp then carries their total exactly except for rounding
// inside comp, at most m²·ε²·Σ|those values| (Higham §4.3); the final rounding of sum + comp adds
// ε·|exact|, doubled for margin. Nothing outside that span may enter the bound: that is what proves
// the error does not depend on earlier history.
func sumTolerance(values []float64, start, w int, exact float64) float64 {
	spanAbs := 0.0
	for _, v := range values[max(0, start-w+1):min(len(values), start+w)] {
		spanAbs += math.Abs(v)
	}
	m := float64(3 * w)
	return 2*eps*math.Abs(exact) + m*m*eps*eps*spanAbs
}

// mixedMagnitudes returns n fixed-seed values with signs mixed and magnitudes spread over 10⁻¹⁶…10¹⁶,
// so cancellation and absorption both happen.
func mixedMagnitudes(rng *rand.Rand, n int) []float64 {
	values := make([]float64, n)
	for i := range values {
		values[i] = (2*rng.Float64() - 1) * math.Pow(10, float64(rng.IntN(33)-16))
	}
	return values
}

// Oracle: the exact sum of each window, within sumTolerance.
func TestRollingSum_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		n := 1 + rng.IntN(200)
		values := mixedMagnitudes(rng, n)
		w := 1 + rng.IntN(n)
		got, err := RollingSum(nil, values, w)
		if err != nil || len(got) != n-w+1 {
			t.Fatalf("trial %d: RollingSum(n=%d, w=%d) = %d sums, %v", trial, n, w, len(got), err)
		}
		for i, g := range got {
			want := exactWindowSum(values[i : i+w])
			if tol := sumTolerance(values, i, w, want); math.Abs(g-want) > tol {
				t.Fatalf("trial %d: window %d (w=%d) sum = %v, exact %v, error %g > tolerance %g", trial, i, w, g, want, math.Abs(g-want), tol)
			}
		}
	}
}

// A long history of mixed magnitudes with short windows: the error of each window must stay within the
// local bound however many values came before (a sum that is never restarted fails this past n ≈ 1e4).
func TestRollingSum_errorDoesNotGrowWithHistory(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	values := mixedMagnitudes(rng, 20_000)
	for _, w := range []int{2, 7, 16} {
		got, err := RollingSum(nil, values, w)
		if err != nil {
			t.Fatal(err)
		}
		for i, g := range got {
			want := exactWindowSum(values[i : i+w])
			if tol := sumTolerance(values, i, w, want); math.Abs(g-want) > tol {
				t.Fatalf("w=%d: window %d sum = %v, exact %v, error %g > tolerance %g", w, i, g, want, math.Abs(g-want), tol)
			}
		}
	}
}

// Reviewer's counterexample (2026-11-12): huge values leave, tiny ones remain. The last windows must be
// the exact sums to within a few ulps, with the right sign.
func TestRollingSum_tinyWindowsAfterHugeValues(t *testing.T) {
	values := []float64{-9.11e15, 0.941, -0.818, -7.26e15, -6.96e-17, -6.14e-17, 2.4e-18, -1.3e-18}
	got, err := RollingSum(nil, values, 2)
	if err != nil {
		t.Fatal(err)
	}
	for i := 4; i < len(got); i++ {
		want := exactWindowSum(values[i : i+2])
		if math.Abs(got[i]-want) > 2*eps*math.Abs(want) {
			t.Errorf("window %d sum = %v, exact %v", i, got[i], want)
		}
	}
}

func BenchmarkRollingSum(b *testing.B) {
	for _, size := range benchSizes {
		values := benchFloats(size.n)
		buf := make([]float64, 0, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				buf, _ = RollingSum(buf[:0], values, benchWindow)
			}
		})
	}
}
