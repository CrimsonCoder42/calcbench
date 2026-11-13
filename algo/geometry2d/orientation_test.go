package geometry2d

import (
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

var benchSizes = []struct {
	name string
	n    int
}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}

// randomPoints returns n points with coordinates uniform in [-span, span].
func randomPoints(rng *rand.Rand, n int, span int64) []Point {
	pts := make([]Point, n)
	for i := range pts {
		pts[i] = Point{rng.Int64N(2*span+1) - span, rng.Int64N(2*span+1) - span}
	}
	return pts
}

// randomCoord returns a coordinate drawn from the extremes of int64 half the time, uniform otherwise.
func randomCoord(rng *rand.Rand) int64 {
	switch rng.IntN(6) {
	case 0:
		return math.MaxInt64 - rng.Int64N(3)
	case 1:
		return math.MinInt64 + rng.Int64N(3)
	case 2:
		return rng.Int64N(7) - 3
	}
	return int64(rng.Uint64())
}

// bigOrientation is the slow obvious version: the cross product in arbitrary precision.
func bigOrientation(a, b, c Point) int {
	bx := new(big.Int).Sub(big.NewInt(b.X), big.NewInt(a.X))
	by := new(big.Int).Sub(big.NewInt(b.Y), big.NewInt(a.Y))
	cx := new(big.Int).Sub(big.NewInt(c.X), big.NewInt(a.X))
	cy := new(big.Int).Sub(big.NewInt(c.Y), big.NewInt(a.Y))
	return new(big.Int).Sub(new(big.Int).Mul(bx, cy), new(big.Int).Mul(by, cx)).Sign()
}

// CLRS 3rd ed. §33.1: the sign of (p1 − p0) × (p2 − p0) gives the turn; worked here on a 3-4-5 triangle.
func TestOrientation_clrsExample(t *testing.T) {
	a, b, c := Point{0, 0}, Point{4, 0}, Point{0, 3}
	if got := Orientation(a, b, c); got != 1 {
		t.Errorf("Orientation(%v, %v, %v) = %d, want +1", a, b, c, got)
	}
	if got := Orientation(a, c, b); got != -1 {
		t.Errorf("Orientation(%v, %v, %v) = %d, want -1", a, c, b, got)
	}
}

func TestOrientation_cases(t *testing.T) {
	const maxI, minI = math.MaxInt64, math.MinInt64
	tests := []struct {
		name    string
		a, b, c Point
		want    int
	}{
		{"collinear diagonal", Point{0, 0}, Point{1, 1}, Point{5, 5}, 0},
		{"collinear, c between", Point{0, 0}, Point{4, 4}, Point{2, 2}, 0},
		{"all equal", Point{3, 3}, Point{3, 3}, Point{3, 3}, 0},
		{"a equals b", Point{1, 2}, Point{1, 2}, Point{7, -4}, 0},
		{"left turn with negatives", Point{-5, -5}, Point{-1, -5}, Point{-1, -2}, 1},
		{"right turn", Point{0, 0}, Point{0, 5}, Point{3, 1}, -1},
		{"extremes, counterclockwise", Point{minI, minI}, Point{maxI, minI}, Point{minI, maxI}, 1},
		{"extremes, clockwise", Point{minI, minI}, Point{minI, maxI}, Point{maxI, minI}, -1},
		{"extremes, collinear", Point{minI, minI}, Point{0, 0}, Point{maxI, maxI}, 0},
		{"near-collinear at the extremes", Point{minI, minI}, Point{maxI, maxI - 1}, Point{maxI - 1, maxI - 2}, -1}, // cross = (N−1)(N−3) − (N−2)² = −1, N = 2⁶⁴
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Orientation(tt.a, tt.b, tt.c); got != tt.want {
				t.Errorf("Orientation(%v, %v, %v) = %d, want %d", tt.a, tt.b, tt.c, got, tt.want)
			}
			if want := bigOrientation(tt.a, tt.b, tt.c); want != tt.want {
				t.Errorf("table entry disagrees with big.Int: %d", want)
			}
		})
	}
}

// Oracle: math/big on coordinates drawn from the whole int64 range, extremes included.
func TestOrientation_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 3000 {
		a := Point{randomCoord(rng), randomCoord(rng)}
		b := Point{randomCoord(rng), randomCoord(rng)}
		c := Point{randomCoord(rng), randomCoord(rng)}
		if got, want := Orientation(a, b, c), bigOrientation(a, b, c); got != want {
			t.Fatalf("trial %d: Orientation(%v, %v, %v) = %d, want %d", trial, a, b, c, got, want)
		}
	}
}

func TestOrientation_noAllocation(t *testing.T) {
	a, b, c := Point{math.MinInt64, 0}, Point{math.MaxInt64, 1}, Point{0, math.MaxInt64}
	if allocs := testing.AllocsPerRun(100, func() { _ = Orientation(a, b, c) }); allocs != 0 {
		t.Errorf("Orientation allocated %v times per call, want 0", allocs)
	}
}

func BenchmarkOrientation(b *testing.B) {
	pts := randomPoints(rand.New(rand.NewPCG(1, 2)), 3, math.MaxInt64/2)
	b.Run("1", func(b *testing.B) {
		b.ReportAllocs()
		s := 0
		for b.Loop() {
			s += Orientation(pts[0], pts[1], pts[2])
		}
		_ = s
	})
}
