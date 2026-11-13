package stats

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"slices"
	"testing"
)

// Shared fixtures. NIST Statistical Reference Datasets (StRD), Univariate Summary Statistics,
// https://www.itl.nist.gov/div898/strd/univ/ (read 2026-11-12); certified values quoted where used. The
// data were transcribed and checked against the certified mean and SD in exact rational arithmetic.

// numAcc1 is StRD NumAcc1: certified mean 10000002, certified sample SD 1.
var numAcc1 = []float64{10000001, 10000003, 10000002}

// numAcc builds StRD NumAcc2–4: first the mean, then (low, high) 500 times; certified SD 0.1.
func numAcc(mean, low, high float64) []float64 {
	v := []float64{mean}
	for range 500 {
		v = append(v, low, high)
	}
	return v
}

var (
	numAcc2 = numAcc(1.2, 1.1, 1.3)                      // certified mean 1.2, SD 0.1
	numAcc3 = numAcc(1000000.2, 1000000.1, 1000000.3)    // certified mean 1000000.2, SD 0.1
	numAcc4 = numAcc(10000000.2, 10000000.1, 10000000.3) // certified mean 10000000.2, SD 0.1
)

// StRD Michelso (n = 100): certified mean 299.852400000000, SD 0.0790105478190518.
// StRD Mavro (n = 50): certified mean 2.00185600000000, SD 0.000429123454003053.
// StRD Lew (n = 200): certified mean -177.435000000000, SD 277.332168044316.
// StRD Linear Least Squares Norris (36 points): https://www.itl.nist.gov/div898/strd/lls/data/LINKS/DATA/Norris.dat,
// certified B0 = -0.262323073774029, B1 = 1.00211681802045, R² = 0.999993745883712.
var michelso = []float64{
	299.85, 299.74, 299.90, 300.07, 299.93, 299.85, 299.95, 299.98, 299.98, 299.88,
	300.00, 299.98, 299.93, 299.65, 299.76, 299.81, 300.00, 300.00, 299.96, 299.96,
	299.96, 299.94, 299.96, 299.94, 299.88, 299.80, 299.85, 299.88, 299.90, 299.84,
	299.83, 299.79, 299.81, 299.88, 299.88, 299.83, 299.80, 299.79, 299.76, 299.80,
	299.88, 299.88, 299.88, 299.86, 299.72, 299.72, 299.62, 299.86, 299.97, 299.95,
	299.88, 299.91, 299.85, 299.87, 299.84, 299.84, 299.85, 299.84, 299.84, 299.84,
	299.89, 299.81, 299.81, 299.82, 299.80, 299.77, 299.76, 299.74, 299.75, 299.76,
	299.91, 299.92, 299.89, 299.86, 299.88, 299.72, 299.84, 299.85, 299.85, 299.78,
	299.89, 299.84, 299.78, 299.81, 299.76, 299.81, 299.79, 299.81, 299.82, 299.85,
	299.87, 299.87, 299.81, 299.74, 299.81, 299.94, 299.95, 299.80, 299.81, 299.87,
}
var mavro = []float64{
	2.00180, 2.00170, 2.00180, 2.00190, 2.00180, 2.00170, 2.00150, 2.00140, 2.00150, 2.00150,
	2.00170, 2.00180, 2.00180, 2.00190, 2.00190, 2.00210, 2.00200, 2.00160, 2.00140, 2.00130,
	2.00130, 2.00150, 2.00150, 2.00160, 2.00150, 2.00140, 2.00130, 2.00140, 2.00150, 2.00140,
	2.00150, 2.00160, 2.00150, 2.00160, 2.00190, 2.00200, 2.00200, 2.00210, 2.00220, 2.00230,
	2.00240, 2.00250, 2.00270, 2.00260, 2.00260, 2.00260, 2.00270, 2.00260, 2.00250, 2.00240,
}
var lew = []float64{
	-213, -564, -35, -15, 141, 115, -420, -360, 203, -338,
	-431, 194, -220, -513, 154, -125, -559, 92, -21, -579,
	-52, 99, -543, -175, 162, -457, -346, 204, -300, -474,
	164, -107, -572, -8, 83, -541, -224, 180, -420, -374,
	201, -236, -531, 83, 27, -564, -112, 131, -507, -254,
	199, -311, -495, 143, -46, -579, -90, 136, -472, -338,
	202, -287, -477, 169, -124, -568, 17, 48, -568, -135,
	162, -430, -422, 172, -74, -577, -13, 92, -534, -243,
	194, -355, -465, 156, -81, -578, -64, 139, -449, -384,
	193, -198, -538, 110, -44, -577, -6, 66, -552, -164,
	161, -460, -344, 205, -281, -504, 134, -28, -576, -118,
	156, -437, -381, 200, -220, -540, 83, 11, -568, -160,
	172, -414, -408, 188, -125, -572, -32, 139, -492, -321,
	205, -262, -504, 142, -83, -574, 0, 48, -571, -106,
	137, -501, -266, 190, -391, -406, 194, -186, -553, 83,
	-13, -577, -49, 103, -515, -280, 201, 300, -506, 131,
	-45, -578, -80, 138, -462, -361, 201, -211, -554, 32,
	74, -533, -235, 187, -372, -442, 182, -147, -566, 25,
	68, -535, -244, 194, -351, -463, 174, -125, -570, 15,
	72, -550, -190, 172, -424, -385, 198, -218, -536, 96,
}
var norrisY = []float64{
	0.1, 338.8, 118.1, 888.0, 9.2, 228.1, 668.5, 998.5, 449.1, 778.9,
	559.2, 0.3, 0.1, 778.1, 668.8, 339.3, 448.9, 10.8, 557.7, 228.3,
	998.0, 888.8, 119.6, 0.3, 0.6, 557.6, 339.3, 888.0, 998.5, 778.9,
	10.2, 117.6, 228.9, 668.4, 449.2, 0.2,
}
var norrisX = []float64{
	0.2, 337.4, 118.2, 884.6, 10.1, 226.5, 666.3, 996.3, 448.6, 777.0,
	558.2, 0.4, 0.6, 775.5, 666.9, 338.0, 447.5, 11.6, 556.0, 228.1,
	995.8, 887.6, 120.2, 0.3, 0.3, 556.8, 339.1, 887.2, 999.0, 779.0,
	11.1, 118.3, 229.2, 669.1, 448.9, 0.5,
}

// resistivity are the 12 resistivity measurements (ohm·cm) printed in NIST/SEMATECH e-Handbook §7.2.6.2
// "Percentiles", in printed order.
var resistivity = []float64{95.1772, 95.1567, 95.1937, 95.1959, 95.1442, 95.0610, 95.1591, 95.1195, 95.1065, 95.0925, 95.1990, 95.1682}

// benchSizes are the collection sizes of .claude/rules/go-complexity.md. Shared by the benchmarks.
var benchSizes = []struct {
	name string
	n    int
}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}

// benchSink keeps benchmark results observable.
var benchSink float64

// benchInputs returns 8 fixed-seed random inputs of length n (mean 100, SD 10). Benchmarks rotate
// through them so that no single input stays in cache.
func benchInputs(n int) [8][]float64 {
	rng := rand.New(rand.NewPCG(1, 2))
	var in [8][]float64
	for k := range in {
		in[k] = randomFloats(rng, n)
	}
	return in
}

// randomFloats returns n values from a normal distribution with mean 100 and SD 10.
func randomFloats(rng *rand.Rand, n int) []float64 {
	v := make([]float64, n)
	for i := range v {
		v[i] = 100 + 10*rng.NormFloat64()
	}
	return v
}

// near reports whether got is within relTol of want, relative to |want| (absolute when want is 0).
func near(got, want, relTol float64) bool {
	if want == 0 {
		return math.Abs(got) <= relTol
	}
	return math.Abs(got-want) <= relTol*math.Abs(want)
}

// bigMean is the oracle mean: the sum in 512-bit floating point, divided by n.
func bigMean(values []float64) *big.Float {
	s := new(big.Float).SetPrec(512)
	for _, v := range values {
		s.Add(s, new(big.Float).SetPrec(512).SetFloat64(v))
	}
	return s.Quo(s, new(big.Float).SetPrec(512).SetInt64(int64(len(values))))
}

// bigSumOfSquares is the oracle Σ(x − x̄)², naive two-pass in 512-bit floating point.
func bigSumOfSquares(values []float64) float64 {
	m := bigMean(values)
	s := new(big.Float).SetPrec(512)
	for _, v := range values {
		d := new(big.Float).SetPrec(512).SetFloat64(v)
		d.Sub(d, m)
		s.Add(s, d.Mul(d, d))
	}
	f, _ := s.Float64()
	return f
}

// Worked example: StRD NumAcc1, certified mean 10000002 (exact).
func TestArithmeticMean_strdNumAcc(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   float64
	}{
		{"NumAcc1", numAcc1, 10000002},
		{"NumAcc2", numAcc2, 1.2},
		{"NumAcc3", numAcc3, 1000000.2},
		{"NumAcc4", numAcc4, 10000000.2},
		{"Michelso", michelso, 299.8524},
		{"Mavro", mavro, 2.001856},
		{"Lew", lew, -177.435},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Tolerance 1e-15: the certified means are exact decimals; the compensated mean is within
			// an ulp or two of the mean of the float64 inputs, which are themselves rounded decimals.
			if got, err := ArithmeticMean(tt.values); err != nil || !near(got, tt.want, 1e-15) {
				t.Errorf("ArithmeticMean(%s) = %v, %v; want %v", tt.name, got, err, tt.want)
			}
		})
	}
}

func TestArithmeticMean_cases(t *testing.T) {
	m := math.MaxFloat64
	tests := []struct {
		name   string
		values []float64
		want   float64
	}{
		{"one element", []float64{-3.5}, -3.5},
		{"two elements", []float64{1, 2}, 1.5},
		{"duplicates", []float64{4, 4, 4}, 4},
		{"negative and zero", []float64{-2, 0, 2, -4}, -1},
		{"cancellation", []float64{1e16, 1, -1e16}, 1.0 / 3},
		{"max of the type", []float64{m, m, m}, m},
		{"min and max cancel", []float64{-m, m}, 0},
		{"overflowing sum", []float64{m, m / 2}, 0.75 * m},
		{"smallest subnormal", []float64{math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64}, math.SmallestNonzeroFloat64},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Tolerance 1e-15: one or two roundings of the sum and the division.
			if got, err := ArithmeticMean(tt.values); err != nil || !near(got, tt.want, 1e-15) {
				t.Errorf("ArithmeticMean(%v) = %v, %v; want %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestArithmeticMean_errors(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   error
	}{
		{"nil", nil, ErrEmpty},
		{"empty", []float64{}, ErrEmpty},
		{"NaN", []float64{1, math.NaN()}, ErrNotFinite},
		{"+Inf", []float64{math.Inf(1)}, ErrNotFinite},
		{"-Inf", []float64{2, math.Inf(-1), 3}, ErrNotFinite},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := ArithmeticMean(tt.values); !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("ArithmeticMean(%v) = %v, %v; want 0, %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestArithmeticMean_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		n := 1 + rng.IntN(500)
		values := make([]float64, n)
		for i := range values {
			values[i] = (rng.Float64() - 0.5) * math.Pow(10, float64(rng.IntN(30)-10))
		}
		want, _ := bigMean(values).Float64()
		// Tolerance: Neumaier's bound, a few ulps of the result plus n·u² of Σ|x| (here 1e-28 of it).
		got, err := ArithmeticMean(values)
		var abs float64
		for _, v := range values {
			abs += math.Abs(v)
		}
		if err != nil || math.Abs(got-want) > 4e-16*math.Abs(want)+1e-28*abs/float64(n) {
			t.Fatalf("trial %d: ArithmeticMean(n=%d) = %v, %v; oracle %v", trial, n, got, err, want)
		}
	}
}

func TestArithmeticMean_doesNotModifyInput(t *testing.T) {
	values := slices.Clone(michelso)
	if _, err := ArithmeticMean(values); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(values, michelso) {
		t.Error("ArithmeticMean modified its input")
	}
}

func TestArithmeticMean_noAllocation(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { benchSink, _ = ArithmeticMean(michelso) }); a != 0 {
		t.Errorf("ArithmeticMean allocates %v times per call, want 0", a)
	}
}

func BenchmarkArithmeticMean(b *testing.B) {
	for _, size := range benchSizes {
		inputs := benchInputs(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; b.Loop(); i++ {
				benchSink, _ = ArithmeticMean(inputs[i&7])
			}
		})
	}
}
