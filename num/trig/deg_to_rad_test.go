package trig

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// Shared by the package's tests.

// sink and sinkTriangle keep benchmark results observable.
var (
	sink         float64
	sinkTriangle Triangle
)

// Tolerances. ulps4 allows four roundings of a value near 1 (each at most 2^-53 relative); it is used
// where the code does a few correctly rounded operations. Looser tolerances are stated where used.
const ulps4 = 4 * 0x1p-52

// near reports |got − want| ≤ tol·max(1, |want|): relative for |want| ≥ 1, absolute below.
func near(got, want, tol float64) bool {
	return math.Abs(got-want) <= tol*math.Max(1, math.Abs(want))
}

// relNear reports |got − want| ≤ tol·|want|, a purely relative check.
func relNear(got, want, tol float64) bool {
	return math.Abs(got-want) <= tol*math.Abs(want)
}

// radOf and degOf are the test's own conversions, written out so a test does not lean on the code under test.
func radOf(deg float64) float64 { return (deg / 180) * math.Pi }
func degOf(rad float64) float64 { return (rad / math.Pi) * 180 }

// errCase is one call that must fail with want.
type errCase struct {
	name string
	call func() error
	want error
}

func runErrCases(t *testing.T, cases []errCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.call(); !errors.Is(err, tc.want) {
				t.Errorf("err = %v, want %v", err, tc.want)
			}
		})
	}
}

// bigPi is π to 60 digits, for big.Float oracles.
var bigPi, _ = new(big.Float).SetPrec(256).SetString("3.14159265358979323846264338327950288419716939937510582097494")

// newRNG is the fixed-seed generator every oracle test uses.
func newRNG() *rand.Rand { return rand.New(rand.NewPCG(1, 2)) }

// Worked example: NIST SP 811 (2008), §5.1.1, Table 6, 1° = (π/180) rad. 180°, 90° and −45° are exact.
func TestDegToRad_sp811Example(t *testing.T) {
	for _, tc := range []struct{ deg, want float64 }{{180, math.Pi}, {90, math.Pi / 2}, {-45, -math.Pi / 4}, {360, 2 * math.Pi}} {
		if got, err := DegToRad(tc.deg); err != nil || got != tc.want {
			t.Errorf("DegToRad(%v) = (%v, %v), want (%v, nil)", tc.deg, got, err, tc.want)
		}
	}
}

func TestDegToRad_cases(t *testing.T) {
	tests := []struct {
		name     string
		deg      float64
		want     float64
		tol      float64
		exactNeg bool
	}{
		{"zero", 0, 0, 0, false},
		{"negative zero keeps its sign", math.Copysign(0, -1), 0, 0, true},
		{"one degree", 1, 0.017453292519943295, ulps4, false},
		{"not reduced", 720, 4 * math.Pi, 0, false},
		{"largest float64", math.MaxFloat64, math.MaxFloat64 / 180 * math.Pi, ulps4, false},
		{"smallest subnormal", 5e-324, 0, 1e-320, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := DegToRad(tc.deg)
			if err != nil || !near(got, tc.want, tc.tol) || (tc.exactNeg && !math.Signbit(got)) {
				t.Errorf("DegToRad(%v) = (%v, %v), want %v", tc.deg, got, err, tc.want)
			}
		})
	}
}

func TestDegToRad_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { _, err := DegToRad(math.NaN()); return err }, ErrNotFinite},
		{"+Inf", func() error { _, err := DegToRad(math.Inf(1)); return err }, ErrNotFinite},
		{"-Inf", func() error { _, err := DegToRad(math.Inf(-1)); return err }, ErrNotFinite},
	})
}

// Oracle: deg·π/180 in 256-bit big.Float, rounded once. The code rounds twice (the division, the
// product), so it is within one ulp of the oracle: tolerance 2^-52 relative.
func TestDegToRad_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		deg := (rng.Float64() - 0.5) * 1e6
		want, _ := new(big.Float).SetPrec(256).Quo(new(big.Float).SetPrec(256).Mul(big.NewFloat(deg), bigPi), big.NewFloat(180)).Float64()
		if got, err := DegToRad(deg); err != nil || !relNear(got, want, 0x1p-52) {
			t.Fatalf("DegToRad(%v) = (%v, %v), want %v", deg, got, err, want)
		}
	}
}

func TestDegToRad_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = DegToRad(37.5) }); n != 0 {
		t.Errorf("DegToRad allocates %v times, want 0", n)
	}
}

func BenchmarkDegToRad(b *testing.B) {
	b.ReportAllocs()
	x := 37.5
	for b.Loop() {
		sink, _ = DegToRad(x)
	}
}
