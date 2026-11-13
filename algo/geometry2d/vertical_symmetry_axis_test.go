package geometry2d

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Worked by hand on the specialisation of Atallah 1985 proved in the catalog: the axis must be
// x = (minX + maxX)/2. (1,1),(−1,1) mirror about x = 0; (1,1),(−1,−1) do not; the trapezoid
// (0,0),(3,0),(1,5),(2,5) mirrors about x = 1.5.
func TestHasVerticalSymmetryAxis_example(t *testing.T) {
	tests := []struct {
		points []Point
		want   bool
	}{
		{[]Point{{1, 1}, {-1, 1}}, true},
		{[]Point{{1, 1}, {-1, -1}}, false},
		{[]Point{{0, 0}, {3, 0}, {1, 5}, {2, 5}}, true},
	}
	for _, tt := range tests {
		if got := HasVerticalSymmetryAxis(tt.points); got != tt.want {
			t.Errorf("HasVerticalSymmetryAxis(%v) = %v, want %v", tt.points, got, tt.want)
		}
	}
}

func TestHasVerticalSymmetryAxis_cases(t *testing.T) {
	const maxI, minI = math.MaxInt64, math.MinInt64
	tests := []struct {
		name   string
		points []Point
		want   bool
	}{
		{"empty", nil, true},
		{"one", []Point{{7, -3}}, true},
		{"one column", []Point{{2, 0}, {2, 5}, {2, -9}}, true},
		{"point on the axis", []Point{{0, 0}, {4, 0}, {2, 9}}, true},
		{"off-axis point unmatched", []Point{{0, 0}, {4, 0}, {3, 9}}, false},
		{"duplicates are ignored", []Point{{0, 0}, {0, 0}, {2, 0}}, true},
		{"same xs, different ys", []Point{{0, 0}, {2, 1}}, false},
		{"negative", []Point{{-5, 1}, {-1, 1}, {-3, 4}}, true},
		{"extremes of int64", []Point{{minI, 5}, {maxI, 5}, {-1, 0}, {0, 0}}, true},
		{"extremes, axis at -1/2", []Point{{minI, 5}, {maxI, 5}, {0, 0}}, false},
		{"extremes, one-sided", []Point{{minI, 0}, {maxI, 0}, {maxI, 1}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasVerticalSymmetryAxis(tt.points); got != tt.want {
				t.Errorf("HasVerticalSymmetryAxis(%v) = %v, want %v", tt.points, got, tt.want)
			}
		})
	}
}

func TestHasVerticalSymmetryAxis_doesNotModifyInput(t *testing.T) {
	pts := randomPoints(rand.New(rand.NewPCG(1, 2)), 50, 10)
	orig := slices.Clone(pts)
	_ = HasVerticalSymmetryAxis(pts)
	if !slices.Equal(pts, orig) {
		t.Errorf("HasVerticalSymmetryAxis modified its input")
	}
}

// slowVerticalSymmetry tries every candidate axis 2c ∈ [2·minX, 2·maxX] and checks each point's mirror
// by a linear scan (small coordinates only).
func slowVerticalSymmetry(points []Point) bool {
	if len(points) == 0 {
		return true
	}
	lo, hi := points[0].X, points[0].X
	for _, p := range points {
		lo, hi = min(lo, p.X), max(hi, p.X)
	}
	for twoC := 2 * lo; twoC <= 2*hi; twoC++ {
		ok := true
		for _, p := range points {
			if !slices.Contains(points, Point{twoC - p.X, p.Y}) {
				ok = false
				break
			}
		}
		if ok {
			return true
		}
	}
	return false
}

// Oracle: every candidate axis tried; half the inputs are built symmetric and then perhaps perturbed.
func TestHasVerticalSymmetryAxis_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		pts := randomPoints(rng, rng.IntN(12), 6)
		if trial%2 == 0 {
			twoC := rng.Int64N(9) - 4
			for _, p := range slices.Clone(pts) {
				pts = append(pts, Point{twoC - p.X, p.Y})
			}
			if rng.IntN(3) == 0 && len(pts) > 0 {
				pts[rng.IntN(len(pts))].Y++
			}
			rng.Shuffle(len(pts), func(i, j int) { pts[i], pts[j] = pts[j], pts[i] })
		}
		if got, want := HasVerticalSymmetryAxis(pts), slowVerticalSymmetry(pts); got != want {
			t.Fatalf("trial %d: HasVerticalSymmetryAxis(%v) = %v, want %v", trial, pts, got, want)
		}
	}
}

func BenchmarkHasVerticalSymmetryAxis(b *testing.B) {
	for _, size := range benchSizes {
		// A symmetric set: the worst case, every mirror is looked up.
		half := randomPoints(rand.New(rand.NewPCG(1, 2)), size.n/2, 1<<40)
		pts := slices.Clone(half)
		for _, p := range half {
			pts = append(pts, Point{-p.X, p.Y})
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = HasVerticalSymmetryAxis(pts)
			}
		})
	}
}
