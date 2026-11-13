package geometry

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// Shared test helpers for the package.

var (
	sink      float64
	sinkPoint Point
	nan       = math.NaN()
	inf       = math.Inf(1)
)

// Relative tolerances, each with its reason.
const (
	// tolFew: a handful of correctly rounded operations, or one call to a library sin, cos, tan, atan2 or
	// sqrt (each within an ulp or two), against a 60-digit reference: up to about 4.5 ulps of 2⁻⁵².
	tolFew = 1e-15
	// tolSome: tens of operations, or a library function's error amplified by a small condition number.
	tolSome = 1e-14
)

func newRNG() *rand.Rand { return rand.New(rand.NewPCG(1, 2)) }

// relNear reports whether got is within tol·|want| of want (within tol when want is 0).
func relNear(got, want, tol float64) bool {
	if want == 0 {
		return math.Abs(got) <= tol
	}
	return math.Abs(got-want) <= tol*math.Abs(want)
}

// valCase is one expected value of a (float64, error) call.
type valCase struct {
	name string
	got  func() (float64, error)
	want float64
	tol  float64 // relative; 0 means exact
}

func runValCases(t *testing.T, cases []valCase) {
	t.Helper()
	for _, c := range cases {
		got, err := c.got()
		ok := err == nil && got == c.want
		if c.tol > 0 {
			ok = err == nil && relNear(got, c.want, c.tol)
		}
		if !ok {
			t.Errorf("%s = (%v, %v), want %v (rel tol %g)", c.name, got, err, c.want, c.tol)
		}
	}
}

// errCase is one documented error.
type errCase struct {
	name string
	call func() error
	want error
}

func runErrCases(t *testing.T, cases []errCase) {
	t.Helper()
	for _, c := range cases {
		if err := c.call(); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// e1 keeps the error of a (float64, error) call; ep the error of a (Point, error) call.
func e1(_ float64, err error) error { return err }
func ep(_ Point, err error) error   { return err }

// assertNoAllocs checks that f does not allocate.
func assertNoAllocs(t *testing.T, name string, f func()) {
	t.Helper()
	if n := testing.AllocsPerRun(100, f); n != 0 {
		t.Errorf("%s allocates %v times per call, want 0", name, n)
	}
}

// Worked example: CRC 31st ed. §4.4.2 eq. (4.4.5); the 3-4-5 right triangle, hand-derived.
func TestPointDistance_example(t *testing.T) {
	runValCases(t, []valCase{
		{"PointDistance((0,0),(3,4))", func() (float64, error) { return PointDistance(Point{0, 0}, Point{3, 4}) }, 5, 0},
	})
}

func TestPointDistance_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"same point", func() (float64, error) { return PointDistance(Point{2, -7}, Point{2, -7}) }, 0, 0},
		{"negative coordinates", func() (float64, error) { return PointDistance(Point{-1, -1}, Point{-4, -5}) }, 5, 0},
		{"no overflow of the squares", func() (float64, error) { return PointDistance(Point{0, 0}, Point{3e200, 4e200}) }, 5e200, tolFew},
		{"no underflow of the squares", func() (float64, error) { return PointDistance(Point{0, 0}, Point{3e-200, 4e-200}) }, 5e-200, tolFew},
	})
}

func TestPointDistance_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN p", func() error { return e1(PointDistance(Point{nan, 0}, Point{})) }, ErrNotFinite},
		{"Inf q", func() error { return e1(PointDistance(Point{}, Point{0, -inf})) }, ErrNotFinite},
		{"overflow", func() error { return e1(PointDistance(Point{-1e308, 0}, Point{1e308, 0})) }, ErrOverflow},
	})
}

func TestPointDistance_oracle(t *testing.T) {
	rng := newRNG()
	for range 400 {
		p := Point{rng.Float64()*2000 - 1000, rng.Float64()*2000 - 1000}
		q := Point{rng.Float64()*2000 - 1000, rng.Float64()*2000 - 1000}
		dx, dy := q.X-p.X, q.Y-p.Y
		want := math.Sqrt(dx*dx + dy*dy)
		if got, err := PointDistance(p, q); err != nil || !relNear(got, want, tolSome) {
			t.Fatalf("PointDistance(%v, %v) = (%v, %v), want %v", p, q, got, err, want)
		}
	}
}

func TestPointDistance_allocs(t *testing.T) {
	assertNoAllocs(t, "PointDistance", func() { sink, _ = PointDistance(Point{1, 2}, Point{4, 6}) })
}

func BenchmarkPointDistance(b *testing.B) {
	b.ReportAllocs()
	p, q := Point{1, 2}, Point{4, 6}
	for b.Loop() {
		sink, _ = PointDistance(p, q)
	}
}
