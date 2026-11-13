package series

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// theilSenOracle lists every slope with distinct x, sorts them and averages the two middle ones: the
// O(n²) definition (Sen 1968), as the reference computed it.
func theilSenOracle(x, y []float64) float64 {
	var s []float64
	for i := range x {
		for j := i + 1; j < len(x); j++ {
			if x[i] != x[j] {
				s = append(s, (y[j]-y[i])/(x[j]-x[i]))
			}
		}
	}
	slices.Sort(s)
	return mean2(s[(len(s)-1)/2], s[len(s)/2])
}

// closeRel reports |got − want| ≤ tol·max(1, |want|). Tolerance: slopes that tie to within the rounding
// of y − t·x may be ordered differently from their divided values (doc comment), so agreement is to
// rounding, not bit for bit.
func closeRel(got, want, tol float64) bool {
	return math.Abs(got-want) <= tol*math.Max(1, math.Abs(want))
}

func newRand() *rand.Rand { return rand.New(rand.NewPCG(1, 2)) }

func TestTheilSenSlope_example(t *testing.T) {
	// Doc example, hand-checked: slopes 2, 0.5, 2, −1, 2, 5 → median 2.
	got, err := TheilSenSlope([]float64{0, 1, 2, 3}, []float64{1, 3, 2, 7}, newRand())
	if err != nil || got != 2 {
		t.Fatalf("got %v, %v; want 2", got, err)
	}
}

func TestTheilSenSlope_cases(t *testing.T) {
	cases := []struct {
		name string
		x, y []float64
		want float64
	}{
		{"two points", []float64{1, 3}, []float64{2, 8}, 3},
		{"odd count of slopes", []float64{0, 1, 2}, []float64{0, 2, 1}, 0.5},
		{"equal x pairs left out", []float64{0, 0, 1}, []float64{0, 10, 1}, mean2(-9, 1)},
		{"negative slope", []float64{0, 1, 2, 3}, []float64{3, 2, 1, 0}, -1},
		{"flat", []float64{5, 1, 3}, []float64{4, 4, 4}, 0},
		{"one outlier ignored", []float64{0, 1, 2, 3, 4}, []float64{0, 1, 2, 3, 100}, 1},
		{"unsorted x", []float64{3, 0, 2, 1}, []float64{7, 1, 2, 3}, 2},
	}
	for _, c := range cases {
		got, err := TheilSenSlope(c.x, c.y, newRand())
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestTheilSenSlope_errors(t *testing.T) {
	r := newRand()
	cases := []struct {
		name string
		x, y []float64
		rng  *rand.Rand
		want error
	}{
		{"nil rng", []float64{0, 1}, []float64{0, 1}, nil, ErrNilRandom},
		{"length mismatch", []float64{0, 1}, []float64{0}, r, ErrLengthMismatch},
		{"empty", nil, nil, r, ErrTooFew},
		{"one point", []float64{1}, []float64{1}, r, ErrTooFew},
		{"all x equal", []float64{2, 2, 2}, []float64{1, 2, 3}, r, ErrTooFew},
		{"NaN x", []float64{0, math.NaN()}, []float64{0, 1}, r, ErrNotFinite},
		{"Inf y", []float64{0, 1}, []float64{0, math.Inf(1)}, r, ErrNotFinite},
		{"-Inf y", []float64{0, 1}, []float64{math.Inf(-1), 0}, r, ErrNotFinite},
		{"slope overflow", []float64{0, 1e-300}, []float64{0, 1e300}, r, ErrOverflow},
		{"y − t·x overflow", []float64{-1e308, 1e308}, []float64{0, 1}, r, ErrOverflow},
	}
	for _, c := range cases {
		if _, err := TheilSenSlope(c.x, c.y, c.rng); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestTheilSenSlope_doesNotModifyInput(t *testing.T) {
	x, y := []float64{3, 1, 2}, []float64{1, 5, 2}
	xc, yc := slices.Clone(x), slices.Clone(y)
	if _, err := TheilSenSlope(x, y, newRand()); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(x, xc) || !slices.Equal(y, yc) {
		t.Fatal("input modified")
	}
}

// Oracle: random data of sizes on both sides of the listing cutoff (n ≤ 90 lists every slope; above it
// the randomized selection runs), against the O(n²) definition.
func TestTheilSenSlope_oracleRandom(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for trial := 0; trial < 300; trial++ {
		n := 2 + r.IntN(400)
		x, y := make([]float64, n), make([]float64, n)
		slope := r.NormFloat64() * 5
		for i := range x {
			x[i] = r.Float64() * 100
			y[i] = slope*x[i] + r.NormFloat64()*20
			if r.IntN(10) == 0 {
				y[i] += 1000 // outliers
			}
		}
		got, err := TheilSenSlope(x, y, rand.New(rand.NewPCG(uint64(trial), 7)))
		if err != nil {
			t.Fatal(err)
		}
		if want := theilSenOracle(x, y); !closeRel(got, want, 1e-12) {
			t.Fatalf("trial %d n=%d: got %v, want %v", trial, n, got, want)
		}
	}
}

// Oracle: heavy ties (small integer grid, many equal x and many equal slopes) and duplicate points.
func TestTheilSenSlope_oracleTies(t *testing.T) {
	r := rand.New(rand.NewPCG(3, 4))
	for trial := 0; trial < 100; trial++ {
		n := 100 + r.IntN(400)
		x, y := make([]float64, n), make([]float64, n)
		for i := range x {
			x[i] = float64(r.IntN(15))
			y[i] = float64(r.IntN(15))
		}
		if x[0] == x[1] {
			x[1]++
		}
		got, err := TheilSenSlope(x, y, rand.New(rand.NewPCG(uint64(trial), 9)))
		if err != nil {
			t.Fatal(err)
		}
		if want := theilSenOracle(x, y); !closeRel(got, want, 1e-12) {
			t.Fatalf("trial %d n=%d: got %v, want %v", trial, n, got, want)
		}
	}
}

// Oracle: points exactly on a line with a large intercept, where y − t·x cannot separate the tied slopes.
func TestTheilSenSlope_oracleCollinearLargeIntercept(t *testing.T) {
	n := 500
	x, y := make([]float64, n), make([]float64, n)
	for i := range x {
		x[i] = float64(i)
		y[i] = 1e9 + 2*float64(i)
	}
	got, err := TheilSenSlope(x, y, newRand())
	if err != nil {
		t.Fatal(err)
	}
	if want := theilSenOracle(x, y); !closeRel(got, want, 1e-9) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// Oracle at n = 3000 (4.5 million slopes): several contraction rounds.
func TestTheilSenSlope_oracleLarge(t *testing.T) {
	r := rand.New(rand.NewPCG(5, 6))
	n := 3000
	x, y := make([]float64, n), make([]float64, n)
	for i := range x {
		x[i] = r.Float64() * 1e4
		y[i] = 0.3*x[i] + r.NormFloat64()*100
	}
	got, err := TheilSenSlope(x, y, newRand())
	if err != nil {
		t.Fatal(err)
	}
	if want := theilSenOracle(x, y); !closeRel(got, want, 1e-12) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// The result does not depend on the random source.
func TestTheilSenSlope_seedIndependent(t *testing.T) {
	r := rand.New(rand.NewPCG(8, 8))
	n := 1000
	x, y := make([]float64, n), make([]float64, n)
	for i := range x {
		x[i], y[i] = r.Float64(), r.Float64()
	}
	first, _ := TheilSenSlope(x, y, rand.New(rand.NewPCG(1, 1)))
	for s := uint64(2); s < 10; s++ {
		if got, _ := TheilSenSlope(x, y, rand.New(rand.NewPCG(s, s))); !closeRel(got, first, 1e-12) {
			t.Fatalf("seed %d: %v vs %v", s, got, first)
		}
	}
}

func benchData(n int) ([]float64, []float64) {
	r := rand.New(rand.NewPCG(1, 2))
	x, y := make([]float64, n), make([]float64, n)
	for i := range x {
		x[i] = r.Float64() * 1000
		y[i] = 2*x[i] + r.NormFloat64()*50
	}
	return x, y
}

// Old → new: BenchmarkTheilSenSlopeQuadratic is the reference's O(n²) time and space (every slope listed
// and sorted); 1e5 is left out of it (5·10⁹ slopes, 40 GB).
func BenchmarkTheilSenSlope(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		x, y := benchData(n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			r := newRand()
			for b.Loop() {
				if _, err := TheilSenSlope(x, y, r); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkTheilSenSlopeQuadratic(b *testing.B) {
	for _, n := range []int{1_000, 10_000} {
		x, y := benchData(n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = theilSenOracle(x, y)
			}
		})
	}
}
