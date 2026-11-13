package geometry2d

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Worked by hand (the problem as stated by Gajentaan–Overmars 1995; per-anchor counting proved in the
// catalog): four points on y = x beat three on y = 0 and two on x + y = 1.
func TestMaxPointsOnLine_example(t *testing.T) {
	pts := []Point{{0, 0}, {1, 1}, {2, 2}, {3, 3}, {0, 1}, {1, 0}, {2, 0}}
	if got := MaxPointsOnLine(pts); got != 4 {
		t.Errorf("MaxPointsOnLine(example) = %d, want 4", got)
	}
}

func TestMaxPointsOnLine_cases(t *testing.T) {
	const maxI, minI = math.MaxInt64, math.MinInt64
	tests := []struct {
		name   string
		points []Point
		want   int
	}{
		{"empty", nil, 0},
		{"one", []Point{{4, 4}}, 1},
		{"two", []Point{{0, 0}, {5, 9}}, 2},
		{"three in a triangle", []Point{{0, 0}, {1, 0}, {0, 1}}, 2},
		{"all equal", []Point{{2, 3}, {2, 3}, {2, 3}}, 3},
		{"duplicates count on the line", []Point{{0, 0}, {0, 0}, {1, 1}, {5, 0}}, 3},
		{"vertical", []Point{{3, 0}, {3, 7}, {3, -2}, {4, 4}}, 3},
		{"horizontal", []Point{{-1, 2}, {9, 2}, {4, 2}, {4, 5}, {0, 0}}, 3},
		{"slopes of opposite sign differ", []Point{{0, 0}, {1, 2}, {-1, 2}, {2, 4}}, 3},
		{"steep slopes that floats would merge", []Point{{0, 0}, {1, 1 << 60}, {1, 1<<60 + 1}, {2, 1 << 61}}, 3},
		{"negative", []Point{{-1, -1}, {-2, -2}, {-3, -3}, {-1, -3}}, 3},
		{"extremes of int64, collinear", []Point{{minI, minI}, {0, 0}, {maxI, maxI}, {1, 2}}, 3},
		{"extremes of int64, not collinear", []Point{{minI, minI}, {maxI, maxI - 1}, {maxI - 1, maxI - 2}}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxPointsOnLine(tt.points); got != tt.want {
				t.Errorf("MaxPointsOnLine(%v) = %d, want %d", tt.points, got, tt.want)
			}
		})
	}
}

func TestMaxPointsOnLine_doesNotModifyInput(t *testing.T) {
	pts := randomPoints(rand.New(rand.NewPCG(1, 2)), 60, 5)
	orig := slices.Clone(pts)
	_ = MaxPointsOnLine(pts)
	if !slices.Equal(pts, orig) {
		t.Errorf("MaxPointsOnLine modified its input")
	}
}

// slowMaxPointsOnLine is the obvious O(n³) version: for every pair of distinct points count the points
// collinear with them; with no distinct pair every point is the same point.
func slowMaxPointsOnLine(points []Point) int {
	best := min(len(points), 1)
	distinct := false
	for i := range points {
		for j := i + 1; j < len(points); j++ {
			if points[i] == points[j] {
				continue
			}
			distinct = true
			count := 0
			for _, r := range points {
				if Orientation(points[i], points[j], r) == 0 {
					count++
				}
			}
			best = max(best, count)
		}
	}
	if !distinct {
		return len(points)
	}
	return best
}

func TestBinaryGCD(t *testing.T) {
	tests := []struct{ u, v, want uint64 }{
		{0, 0, 0}, {0, 7, 7}, {7, 0, 7}, {12, 18, 6}, {17, 5, 1}, {1 << 63, 1 << 40, 1 << 40},
		{math.MaxUint64, math.MaxUint64, math.MaxUint64}, {math.MaxUint64, 3, 3}, {48, 180, 12},
	}
	for _, tt := range tests {
		if got := binaryGCD(tt.u, tt.v); got != tt.want {
			t.Errorf("binaryGCD(%d, %d) = %d, want %d", tt.u, tt.v, got, tt.want)
		}
	}
}

// Oracle: the O(n³) count on small grids, where collinear triples and duplicates are common.
func TestMaxPointsOnLine_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		pts := randomPoints(rng, rng.IntN(30), int64(1+rng.IntN(5)))
		if got, want := MaxPointsOnLine(pts), slowMaxPointsOnLine(pts); got != want {
			t.Fatalf("trial %d: MaxPointsOnLine(%v) = %d, want %d", trial, pts, got, want)
		}
	}
}

// Θ(n²) work: sizes 1e2, 1e3 and 3e3 keep one iteration under a second.
func BenchmarkMaxPointsOnLine(b *testing.B) {
	for _, size := range []struct {
		name string
		n    int
	}{{"1e2", 100}, {"1e3", 1_000}, {"3e3", 3_000}} {
		pts := randomPoints(rand.New(rand.NewPCG(1, 2)), size.n, 1<<20)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = MaxPointsOnLine(pts)
			}
		})
	}
}
