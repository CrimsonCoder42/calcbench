package reduce

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"strconv"
	"testing"
)

// exactSum is the oracle: the sum in 2000-bit big.Float, exact for any float64 inputs of this size, then
// rounded once to float64.
func exactSum(values []float64) float64 {
	acc := new(big.Float).SetPrec(2000)
	for _, x := range values {
		acc.Add(acc, new(big.Float).SetPrec(2000).SetFloat64(x))
	}
	f, _ := acc.Float64()
	return f
}

// Neumaier 1974 and Higham 2002 §4.3: Kahan's method fails on 1, 1e100, 1, -1e100 (it returns 0); the
// improved method returns the exact 2.
func TestSum_neumaierExample(t *testing.T) {
	got, err := Sum([]float64{1, 1e100, 1, -1e100})
	if err != nil || got != 2 {
		t.Errorf("Sum = (%v, %v), want (2, nil)", got, err)
	}
}

func TestSum_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   float64
	}{
		{"empty", nil, 0},
		{"one", []float64{3.5}, 3.5},
		{"two", []float64{1, 2}, 3},
		{"negative", []float64{-1, -2, -3}, -6},
		{"tenths", []float64{0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1}, exactSum([]float64{0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1})},
		{"max float", []float64{math.MaxFloat64, -math.MaxFloat64}, 0},
		{"smallest subnormal", []float64{math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64}, 2 * math.SmallestNonzeroFloat64},
	}
	for _, tc := range tests {
		got, err := Sum(tc.values)
		if err != nil || got != tc.want {
			t.Errorf("%s: Sum = (%v, %v), want (%v, nil)", tc.name, got, err, tc.want)
		}
	}
}

func TestSum_errors(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   error
	}{
		{"NaN", []float64{1, math.NaN()}, ErrNotFinite},
		{"+Inf", []float64{math.Inf(1)}, ErrNotFinite},
		{"-Inf", []float64{0, math.Inf(-1)}, ErrNotFinite},
		{"overflow", []float64{math.MaxFloat64, math.MaxFloat64}, ErrOverflow},
		{"partial overflow", []float64{math.MaxFloat64, math.MaxFloat64, -math.MaxFloat64}, ErrOverflow},
	}
	for _, tc := range tests {
		if _, err := Sum(tc.values); !errors.Is(err, tc.want) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.want)
		}
	}
}

// Oracle: on 300 ill-conditioned sums (values spread over 30 decades, both signs) the compensated sum is
// within 2 ulp of the exactly rounded sum when the condition number is moderate; we assert the Higham
// (4.8)-type bound |err| <= 2u·Σ|x| + 1 ulp(result), which holds for every input.
func TestSum_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	const u = 0x1p-53
	for trial := range 300 {
		n := 1 + r.IntN(200)
		values := make([]float64, n)
		var absSum float64
		for i := range values {
			values[i] = (r.Float64()*2 - 1) * math.Pow(10, float64(r.IntN(30)-15))
			absSum += math.Abs(values[i])
		}
		got, err := Sum(values)
		if err != nil {
			t.Fatalf("trial %d: %v", trial, err)
		}
		want := exactSum(values)
		tol := 2*u*absSum + math.Abs(math.Nextafter(want, math.Inf(1))-want)
		if math.Abs(got-want) > tol {
			t.Errorf("trial %d: Sum = %v, exact %v, |err| %g > %g", trial, got, want, math.Abs(got-want), tol)
		}
	}
}

func TestSum_noAlloc(t *testing.T) {
	values := []float64{1, 2, 3}
	if a := testing.AllocsPerRun(100, func() { _, _ = Sum(values) }); a != 0 {
		t.Errorf("Sum allocates %v times, want 0", a)
	}
}

func benchValues(n int) []float64 {
	r := rand.New(rand.NewPCG(3, 4))
	values := make([]float64, n)
	for i := range values {
		values[i] = r.Float64()*2 - 1
	}
	return values
}

func BenchmarkSum(b *testing.B) {
	for _, n := range []int{1e3, 1e4, 1e5} {
		values := benchValues(n)
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = Sum(values)
			}
		})
	}
}
